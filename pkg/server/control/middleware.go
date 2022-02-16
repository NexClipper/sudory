package control

import (
	"encoding/json"
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/events"
	"github.com/NexClipper/sudory/pkg/server/macro/exceptions"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/google/uuid"

	//lint:ignore ST1001 auto-generated
	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type OperateContext struct {
	Http     echo.Context
	Database database.Context
	Req      interface{} //request value
	TaskId   uint32
}

type (
	// TokenVerifier
	//  토큰 검증
	//  에러: Forbidden
	TokenVerifier func(echo.Context) error

	// Binder
	//  요청 데이터 바인드
	//  에러: BadRequest
	Binder func(echo.Context) (interface{}, error)

	// Operator
	//  요청 처리
	//  에러: InternalServerError
	Operator func(OperateContext) (interface{}, error)

	// HttpResponser
	//  응답
	//  에러: InternalServerError
	HttpResponser func(echo.Context, int, interface{}) error
)
type Option struct {
	TokenVerifier
	Binder
	Operator Operator
	HttpResponser
}

// MakeMiddlewareFunc
//  @param: Option
//  @return: echo.HandlerFunc; func(echo.Context) error
func MakeMiddlewareFunc(opt Option) echo.HandlerFunc {

	exec_token_verifier := func(verifier TokenVerifier, ctx echo.Context) (err error) {
		exceptions.Block{
			Try: func() {
				if verifier == nil {
					return //not error
				}

				err = verifier(ctx) //exec settion-token verify
			},
			Catch: func(ex error) {
				err = ex
			},
		}.Do()
		return
	}

	exec_binder := func(bind Binder, ctx echo.Context) (req interface{}, err error) {
		exceptions.Block{
			Try: func() {
				if bind == nil {
					exceptions.Throw("without binder")
				}
				req, err = bind(ctx) //exec bind
			},
			Catch: func(ex error) {
				err = ex
			},
		}.Do()
		return
	}

	exec_operator := func(tid uint32, operate Operator, ctx echo.Context, v interface{}) (out interface{}, err error) {
		exceptions.Block{
			Try: func() {
				if operate == nil {
					exceptions.Throw("without operator")
				}
				out, err = operate(OperateContext{TaskId: tid, Http: ctx, Req: v}) //exec operate
			},
			Catch: func(ex error) {
				err = ex
			},
		}.Do()
		return
	}

	exec_responser := func(response HttpResponser, ctx echo.Context, status int, v interface{}) (err error) {
		exceptions.Block{
			Try: func() {
				if response == nil {
					exceptions.Throw("without responser")
				}
				err = response(ctx, status, v) //exec response
			},
			Catch: func(ex error) {
				err = ex
			},
		}.Do()
		return
	}

	return func(ctx echo.Context) error {

		var (
			err      error
			req, rsp interface{}
			taskId   uuid.UUID = NewUuid()

			right       = func(b []byte, err error) []byte { return b }
			get_path    = func() string { return ctx.Request().URL.Path }
			get_method  = func() string { return ctx.Request().Method }
			get_status  = func() int { return ctx.Response().Status }
			get_query   = func() string { return ctx.QueryString() }
			get_reqbody = func() []byte { return right(json.Marshal(req)) }
			get_rspbody = func() []byte { return right(json.Marshal(rsp)) }
			get_tid     = func() uint32 { return taskId.ID() }

			path   = get_path()
			method = get_method()
			status = get_status()
			query  = get_query()
			// reqbody = get_reqbody()
			reqbody []byte
			rspbody []byte
		)

		//event invoke
		defer func() {

			args := map[string]interface{}{
				"path":    path,
				"query":   query,
				"method":  method,
				"reqbody": reqbody,
				"rspbody": rspbody,
				"status":  status,
				"error":   err,
			}

			if err == nil {
				delete(args, "error")
			}
			events.Invoke(&events.EventArgs{Sender: path, Args: args})
		}()

		//logging
		defer func() {
			logs.InfoS("C", "tid", get_tid(), "method", method, "path", path, "query", query)
			logs.DebugS("C", "tid", get_tid(), "reqbody", reqbody)

			ErrorWithHandler(err, func(err error) {
				logs.ErrorS(err, "S", "tid", get_tid())
			})
			logs.InfoS("S", "tid", get_tid(), "status", status)
			logs.DebugS("S", "tid", get_tid(), "rspbody", rspbody)
		}()

		err = exec_token_verifier(opt.TokenVerifier, ctx)
		if ErrorWithHandler(err,
			func(err error) { logs.ErrorS(err, "verify token") },
			func(err error) {
				type codecarrier interface {
					Code() int
				}

				status := http.StatusForbidden //기본 http status code

				e, ok := err.(codecarrier) //에러에 코드가 포함되어 있는지 확인
				if ok {
					status = e.Code() //에러에 있는 코드를 가져온다
				}

				//세션-토큰 검증 오류
				if err_ := ctx.String(status, err.Error()); err_ != nil {
					logs.ErrorS(err_, "failed response", "body", err.Error())
				}
			}, //실패 응답
		) {
			return err //return HandlerFunc
		}

		req, err = exec_binder(opt.Binder, ctx)
		reqbody = get_reqbody()
		if ErrorWithHandler(err,
			func(err error) { logs.ErrorS(err, "bind request") },
			func(err error) {
				if err_ := ctx.String(http.StatusBadRequest, err.Error()); err_ != nil {
					logs.ErrorS(err_, "failed response", "body", err.Error())
				}
			}, //실패 응답
		) {
			return err //return HandlerFunc
		}

		rsp, err = exec_operator(get_tid(), opt.Operator, ctx, req)
		rspbody = get_rspbody()
		if ErrorWithHandler(err,
			func(err error) { logs.ErrorS(err, "failed operate") },
			func(err error) {
				//내부작업 오류
				if err_ := ctx.String(http.StatusInternalServerError, err.Error()); err_ != nil {
					logs.ErrorS(err_, "failed response", "body", err.Error())
				}
			}, //실패 응답
		) {
			return err //return HandlerFunc
		}

		err = exec_responser(opt.HttpResponser, ctx, http.StatusOK, rsp)
		if ErrorWithHandler(err,
			func(err error) {
				logs.ErrorS(err, "failed response", "body", err.Error())
			},
		) {
			//응답 오류
			return err //return HandlerFunc
		}
		return nil
	}
}

// func errorToString(format string, err error) (out string) {

// 	if err != nil {
// 		out = fmt.Sprintf(format, err.Error())
// 	}
// 	return out
// }

func HttpResponse(ctx echo.Context, status int, v interface{}) error {
	return ctx.JSON(status, v)
}

func OK() interface{} {
	return "OK"
}

func Lock(engine *xorm.Engine, operate func(OperateContext) (interface{}, error)) func(OperateContext) (interface{}, error) {

	return func(ctx OperateContext) (out interface{}, err error) {
		ctx.Database = database.NewContext(engine) //new database context

		exceptions.Block{
			Try: func() {

				exceptions.Throw(ctx.Database.Tx().Begin()) //begin transaction

				defer func() {
					if err != nil {
						ctx.Database.Tx().Rollback() //rollback
						logs.Debugln("tx rollbacked")
					} else {
						ctx.Database.Tx().Commit() //commit
						logs.Debugln("tx commited")
					}
				}()

				out, err = operate(ctx) //call operate
			},
			Catch: func(ex error) {
				err = errors.Wrap(ex, "catch: exec operate with block")
			},
			Finally: func() {
				ctx.Database.Tx().Close()
				logs.Debugln("tx closed")
			}}.Do()
		return
	}
}

func Nolock(engine *xorm.Engine, operate func(OperateContext) (interface{}, error)) func(OperateContext) (interface{}, error) {

	return func(ctx OperateContext) (out interface{}, err error) {
		ctx.Database = database.NewContext(engine) //new database context

		exceptions.Block{
			Try: func() {

				out, err = operate(ctx) //call operate
			},
			Catch: func(ex error) {
				err = errors.Wrap(ex, "catch: exec operate")
			},
			Finally: func() {
				ctx.Database.Tx().Close()
				logs.Debugln("tx closed")
			}}.Do()
		return
	}
}

package control

import (
	"encoding/json"
	"net/http"
	"sync/atomic"

	"github.com/NexClipper/logger"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/events"
	"github.com/NexClipper/sudory/pkg/server/macro/exceptions"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"

	//lint:ignore ST1001 auto-generated
	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type InterMiddlewarer interface {
	Echo() echo.Context
	Database() database.Context
	Param() interface{} //request value
	Body(i interface{}) error
	Tid() uint64
}

type OperateContext struct {
	Http     echo.Context
	Database database.Context
	Req      interface{} //request value
	TaskId   uint64
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

type Ticker interface {
	Count(uint64) uint64
}

type Ticket struct {
	uint64
}

func (ticker *Ticket) Add(d uint64) uint64 {
	return atomic.AddUint64(&ticker.uint64, d)
}

var ticker = Ticket{}

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

	exec_operator := func(tid uint64, operate Operator, ctx echo.Context, v interface{}) (out interface{}, err error) {
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
			tid      = ticker.Add(1) //get ticket

			right     = func(b []byte, err error) []byte { return b }
			getPath   = func() string { return ctx.Request().URL.Path }
			getMethod = func() string { return ctx.Request().Method }
			getStatus = func() int { return ctx.Response().Status }
			getQuery  = func() string { return ctx.QueryString() }
			getParam  = func() string {

				var s string
				pnames := ctx.ParamNames()
				for n := range pnames {
					if 0 < len(s) {
						s += "&"
					}
					s += pnames[n] + "="
					s += ctx.Param(pnames[n])
				}
				return s
			}
			getReqbody = func() []byte { return right(json.Marshal(req)) }
			getRspbody = func() []byte { return right(json.Marshal(rsp)) }
			getTid     = func() uint64 { return tid }

			path   = getPath()
			method = getMethod()
			status = getStatus()
			query  = getQuery()
			param  = getParam()
			// reqbody = get_reqbody()
			reqbody []byte
			rspbody []byte

			tidSink         = logs.WithId(getTid())
			errVerifyTkSink = tidSink.WithName("failed token verify")
			errRspSink      = tidSink.WithName("failed response")
			errbindSink     = tidSink.WithName("failed bind")
			errOperateSink  = tidSink.WithName("failed operate")
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

		requestSink := tidSink.WithName("C")
		if logs.V(0) {
			requestSink = requestSink.WithValue("method", method, "path", path)
			if logs.V(2) {
				requestSink = requestSink.WithValue("param", param)
				if logs.V(3) {
					requestSink = requestSink.WithValue("query", query)
					if logs.V(5) {
						requestSink = requestSink.WithValue("reqbody", reqbody)
					}
				}
			}
		}
		logger.Info(requestSink.String())

		//logging
		defer func() {
			responseSink := tidSink.WithName("S")
			ErrorWithHandler(err, func(err error) {
				logger.Error(responseSink.WithName("error").WithError(err).String())
			})
			if logs.V(0) {
				responseSink = responseSink.WithValue("status", status)
				if logs.V(5) {
					responseSink = responseSink.WithValue("rspbody", rspbody)
				}
			}
			logger.V(0).Info(responseSink.String())
		}()

		err = exec_token_verifier(opt.TokenVerifier, ctx)
		if ErrorWithHandler(err,
			func(err error) { logger.Error(errVerifyTkSink.WithError(err).String()) },
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
				if erro := ctx.String(status, err.Error()); erro != nil {
					logger.Error(errRspSink.WithError(erro).WithValue("ebody", err).String())
				}
			}, //실패 응답
		) {
			return err //return HandlerFunc
		}

		req, err = exec_binder(opt.Binder, ctx)
		reqbody = getReqbody()
		if ErrorWithHandler(err,
			func(err error) { logger.Error(errbindSink.String()) },
			func(err error) {
				if erro := ctx.String(http.StatusBadRequest, err.Error()); erro != nil {
					logger.Error(errRspSink.WithError(erro).WithValue("ebody", err).String())
				}
			}, //실패 응답
		) {
			return err //return HandlerFunc
		}

		rsp, err = exec_operator(getTid(), opt.Operator, ctx, req)
		rspbody = getRspbody()
		if ErrorWithHandler(err,
			func(err error) { logger.Error(errOperateSink.String()) },
			func(err error) {
				//내부작업 오류
				if erro := ctx.String(http.StatusInternalServerError, err.Error()); erro != nil {
					logger.Error(errRspSink.WithError(erro).WithValue("ebody", err).String())
				}
			}, //실패 응답
		) {
			return err //return HandlerFunc
		}

		err = exec_responser(opt.HttpResponser, ctx, http.StatusOK, rsp)
		if ErrorWithHandler(err,
			func(err error) {
				logger.Error(errRspSink.WithError(err).WithValue("body", rsp).String())
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
					} else {
						ctx.Database.Tx().Commit() //commit
					}
				}()

				out, err = operate(ctx) //call operate
			},
			Catch: func(ex error) {
				err = errors.Wrap(ex, "catch: exec operate with block")
			},
			Finally: func() {
				ctx.Database.Tx().Close()
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
			}}.Do()
		return
	}
}

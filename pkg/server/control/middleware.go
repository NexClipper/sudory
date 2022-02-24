package control

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
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

type (
	Contexter interface {
		TicketId() uint64

		Echo() echo.Context
		Database() database.Context
		SetDatabase(database.Context) Contexter

		Bind(interface{}) error
		Object() interface{}

		Params() map[string]string
		Querys() map[string]string
		Forms() map[string]string
		// Body() []byte
	}

	ParamsHolder interface {
		Params() map[string]string
	}
	QuerysHolder interface {
		Querys() map[string]string
	}
	FormsHolder interface {
		Forms() map[string]string
	}
	BodyHolder interface {
		Body() []byte
	}

	RequestValue struct {
		ticketId uint64

		echo echo.Context
		db   database.Context

		onceParam, onceQuery, onceBody, onceFormParam sync.Once //once
		param, query, formParam                       map[string]string
		body                                          []byte
		object                                        interface{}
	}
)

func (holder RequestValue) TicketId() uint64 {
	return holder.ticketId
}
func (holder *RequestValue) SetTicketId(t uint64) Contexter {
	holder.ticketId = t
	return holder
}

func (holder RequestValue) Echo() echo.Context {
	return holder.echo
}
func (holder *RequestValue) SetEcho(e echo.Context) Contexter {
	holder.echo = e
	return holder
}
func (holder RequestValue) Database() database.Context {
	return holder.db
}
func (holder *RequestValue) SetDatabase(d database.Context) Contexter {
	holder.db = d
	return holder
}

func (holder *RequestValue) Params() map[string]string {
	holder.onceParam.Do(func() {
		holder.param = make(map[string]string)
		for _, name := range holder.echo.ParamNames() {
			holder.param[name] = holder.echo.Param(name)
		}
	})
	return holder.param
}
func (holder *RequestValue) ParamString() string {
	s := make([]string, 0, len(holder.query))
	for key := range holder.param {
		s = append(s, fmt.Sprintf("%s:%s", key, holder.query[key]))
	}
	return strings.Join(s, ",")
}

func (holder *RequestValue) Querys() map[string]string {
	holder.onceQuery.Do(func() {
		holder.query = make(map[string]string)
		for key := range holder.echo.QueryParams() {
			holder.query[key] = holder.echo.QueryParam(key)
		}
	})
	return holder.query
}
func (holder *RequestValue) QueryString() string {
	return holder.echo.QueryString()
}

func (holder *RequestValue) Forms() map[string]string {
	holder.onceFormParam.Do(func() {
		holder.formParam = make(map[string]string)
		formdatas, err := holder.echo.FormParams()
		if err != nil {
			return
		}
		for key := range formdatas {
			holder.formParam[key] = holder.echo.FormValue(key)
		}
	})
	return holder.formParam
}
func (holder *RequestValue) FormString() string {
	s := make([]string, 0, len(holder.formParam))
	for key := range holder.formParam {
		s = append(s, fmt.Sprintf("%s=%s", key, holder.formParam[key]))
	}
	return strings.Join(s, "&")
}

func (holder *RequestValue) Body() []byte {
	holder.onceBody.Do(func() {
		//body read all
		//ranout buffer
		holder.body, _ = ioutil.ReadAll(holder.echo.Request().Body) //read all body
		//restore
		holder.echo.Request().Body = ioutil.NopCloser(bytes.NewBuffer(holder.body))
	})
	return holder.body
}

func (holder *RequestValue) Bind(v interface{}) error {

	if err := holder.echo.Bind(v); err != nil {
		return err
	}

	// if err := json.Unmarshal(holder.Body(), v); err != nil {
	// 	return err
	// }
	holder.object = v
	return nil
}

func (holder *RequestValue) Object() interface{} {
	return holder.object
}

// type OperateContext struct {
// 	Http     echo.Context
// 	Database database.Context
// 	Req      interface{} //request value
// 	TaskId   uint64
// }

type (
	// TokenVerifier
	//  토큰 검증
	//  에러: Forbidden
	TokenVerifier func(Contexter) error

	// Binder
	//  요청 데이터 바인드
	//  에러: BadRequest
	Binder func(Contexter) error

	// Operator
	//  요청 처리
	//  에러: InternalServerError
	Operator func(Contexter) (interface{}, error)

	// HttpResponser
	//  응답
	//  에러: InternalServerError
	HttpResponser func(echo.Context, int, interface{}) error

	Behavior func(Contexter, func(Contexter) (interface{}, error)) (interface{}, error)
)
type Option struct {
	TokenVerifier
	Binder
	Operator
	HttpResponser
	Behavior
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

	exec_token_verifier := func(verifier TokenVerifier, ctx Contexter) (err error) {
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

	exec_binder := func(bind Binder, ctx Contexter) (err error) {
		exceptions.Block{
			Try: func() {
				if bind == nil {
					exceptions.Throw("without binder")
				}
				err = bind(ctx) //exec bind
			},
			Catch: func(ex error) {
				err = ex
			},
		}.Do()
		return
	}

	exec_operator := func(behave Behavior, operate Operator, ctx Contexter) (out interface{}, err error) {
		exceptions.Block{
			Try: func() {
				if operate == nil {
					exceptions.Throw("without operator")
				}
				out, err = behave(ctx, operate) //exec operate
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

		reqval := &RequestValue{}
		reqval.SetEcho(ctx)

		var (
			err error
			rsp interface{}

			onceTicketId sync.Once
			ticketId     uint64
			tid          = func() uint64 {
				onceTicketId.Do(func() {
					ticketId = ticker.Add(1)
				})
				return ticketId
			}
			reqPath   = func() string { return ctx.Request().URL.Path }
			reqMethod = func() string { return ctx.Request().Method }
			reqStatus = func() int { return ctx.Response().Status }
			// reqParam  = func() string {
			// 	var s string
			// 	pnames := ctx.ParamNames()
			// 	for n := range pnames {
			// 		if 0 < len(s) {
			// 			s += ":"
			// 		}
			// 		s += pnames[n] + ","
			// 		s += ctx.Param(pnames[n])
			// 	}
			// 	return s
			// }
			reqForm  = func() string { return reqval.FormString() }
			reqParam = func() string { return reqval.ParamString() }
			reqQuery = func() string { return reqval.QueryString() }
			// onceReqBody sync.Once
			// reqBody     []byte
			// getReqbody  = func() []byte {
			// 	onceReqBody.Do(func() {
			// 		//body read all
			// 		//ranout buffer
			// 		//restore buffer again
			// 		b, err := ioutil.ReadAll(ctx.Request().Body)
			// 		if err != nil {
			// 			reqBody = []byte{}
			// 			return
			// 		}
			// 		//restore
			// 		ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(b))

			// 		reqBody = b
			// 	})
			// 	return reqBody
			// }
			reqBody     = func() []byte { return reqval.Body() }
			onceRspbody sync.Once
			rspbody     []byte
			rspBody     = func() []byte {
				onceRspbody.Do(func() {
					right := func(b []byte, err error) []byte { return b }
					rspbody = right(json.Marshal(rsp))
				})
				return rspbody
			}
		)

		var (
			tidSink         = logs.WithId(tid())
			errVerifyTkSink = tidSink.WithName("failed token verify")
			errRspSink      = tidSink.WithName("failed response")
			errbindSink     = tidSink.WithName("failed bind")
			errOperateSink  = tidSink.WithName("failed operate")
		)

		//event invoke
		defer func() {

			args := map[string]interface{}{}

			args["path"] = reqPath()
			if 0 < len(reqQuery()) {
				args["query"] = reqQuery()
			}
			if 0 < len(reqParam()) {
				args["param"] = reqParam()
			}
			if 0 < len(reqForm()) {
				args["form"] = reqForm()
			}
			args["method"] = reqMethod()
			if 0 < len(reqBody()) {
				args["reqbody"] = reqBody()
			}
			if 0 < len(rspBody()) {
				args["rspbody"] = rspBody()
			}
			args["status"] = reqStatus()
			if err != nil {
				args["error"] = err
			}

			events.Invoke(&events.EventArgs{Sender: reqPath(), Args: args})
		}()

		requestSink := tidSink.WithName("C")
		if logs.V(0) {
			requestSink = requestSink.WithValue("method", reqMethod(), "path", reqPath())
			if logs.V(2) {
				requestSink = requestSink.WithValue("param", reqParam())
				if logs.V(3) {
					requestSink = requestSink.WithValue("query", reqQuery())
					if logs.V(5) {
						requestSink = requestSink.WithValue("reqbody", reqBody())
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
				responseSink = responseSink.WithValue("status", reqStatus())
				if logs.V(5) {
					responseSink = responseSink.WithValue("rspbody", rspBody())
				}
			}
			logger.V(0).Info(responseSink.String())
		}()

		err = exec_token_verifier(opt.TokenVerifier, reqval)
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

		err = exec_binder(opt.Binder, reqval)
		if ErrorWithHandler(err,
			func(err error) { logger.Error(errbindSink.WithError(err).String()) },
			func(err error) {
				if erro := ctx.String(http.StatusBadRequest, err.Error()); erro != nil {
					logger.Error(errRspSink.WithError(erro).WithValue("ebody", err).String())
				}
			}, //실패 응답
		) {
			return err //return HandlerFunc
		}

		rsp, err = exec_operator(opt.Behavior, opt.Operator, reqval)

		if ErrorWithHandler(err,
			func(err error) { logger.Error(errOperateSink.WithError(err).String()) },
			func(err error) {
				//내부작업 오류
				if erro := ctx.String(http.StatusInternalServerError, err.Error()); erro != nil {
					logger.Error(errRspSink.WithError(erro).WithValue("ebody", err).String())
				}
			}, //실패 응답
		) {
			return err //return HandlerFunc
		}

		err = exec_responser(opt.HttpResponser, reqval.Echo(), http.StatusOK, rsp)
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

func Lock(engine *xorm.Engine) func(Contexter, func(Contexter) (interface{}, error)) (interface{}, error) {

	return func(ctx Contexter, operate func(Contexter) (interface{}, error)) (out interface{}, err error) {
		ctx.SetDatabase(database.NewContext(engine)) //new database context

		exceptions.Block{
			Try: func() {

				exceptions.Throw(ctx.Database().Tx().Begin()) //begin transaction

				defer func() {
					if err != nil {
						ctx.Database().Tx().Rollback() //rollback
					} else {
						ctx.Database().Tx().Commit() //commit
					}
				}()

				out, err = operate(ctx) //call operate
			},
			Catch: func(ex error) {
				err = errors.Wrap(ex, "catch: exec operate with block")
			},
			Finally: func() {
				ctx.Database().Tx().Close()
			}}.Do()
		return
	}
}

func Nolock(engine *xorm.Engine) func(Contexter, func(Contexter) (interface{}, error)) (interface{}, error) {

	return func(ctx Contexter, operate func(Contexter) (interface{}, error)) (out interface{}, err error) {
		ctx.SetDatabase(database.NewContext(engine)) //new database context

		exceptions.Block{
			Try: func() {

				out, err = operate(ctx) //call operate
			},
			Catch: func(ex error) {
				err = errors.Wrap(ex, "catch: exec operate")
			},
			Finally: func() {
				ctx.Database().Tx().Close()
			}}.Do()
		return
	}
}

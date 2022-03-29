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
	"github.com/NexClipper/sudory/pkg/server/macro/exceptions"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"

	//lint:ignore ST1001 auto-generated
	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/labstack/echo/v4"
	"xorm.io/xorm"
)

type (
	Context interface {
		//echo
		Echo() echo.Context
		SetEcho(e echo.Context) Context
		Forms() map[string]string
		FormString() string
		Params() map[string]string
		ParamString() string
		Queries() map[string]string
		QueryString() string
		Body() []byte
		Bind(interface{}) error
		Object() interface{}
		//database
		Database() database.Context
		SetDatabase(database.Context) Context

		//ticket
		// TicketId() uint64
	}

	ParamsHolder interface {
		Params() map[string]string
	}
	QueriesHolder interface {
		Queries() map[string]string
	}
	FormsHolder interface {
		Forms() map[string]string
	}
	BodyHolder interface {
		Body() []byte
	}

	RequestValue struct {
		//echo
		echo                    echo.Context
		param, query, formParam OnceMapStringString
		body                    OnceBytes
		object                  interface{}
		//database
		db database.Context
	}
)

type OnceMapStringString struct {
	sync.Once
	v map[string]string
}

type OnceBytes struct {
	sync.Once
	v []byte
}

// func (holder RequestValue) TicketId() uint64 {
// 	return holder.ticketId
// }
// func (holder RequestValue) SetTicketId(t uint64) Contexter {
// 	holder.ticketId = t
// 	return &holder
// }

func (holder RequestValue) Echo() echo.Context {
	return holder.echo
}
func (holder RequestValue) SetEcho(e echo.Context) Context {
	holder.echo = e
	return &holder
}

func (holder *RequestValue) Params() map[string]string {
	holder.param.Do(func() {
		holder.param.v = make(map[string]string)
		for _, name := range holder.echo.ParamNames() {
			holder.param.v[name] = holder.echo.Param(name)
		}
	})
	return holder.param.v
}
func (holder *RequestValue) ParamString() string {
	s := make([]string, 0, len(holder.Params()))
	for key := range holder.Params() {
		s = append(s, fmt.Sprintf("%s:%s", key, holder.Params()[key]))
	}
	return strings.Join(s, ",")
}

func (holder *RequestValue) Queries() map[string]string {
	holder.query.Do(func() {
		holder.query.v = make(map[string]string)
		for key := range holder.echo.QueryParams() {
			holder.query.v[key] = holder.echo.QueryParam(key)
		}
	})
	return holder.query.v
}
func (holder *RequestValue) QueryString() string {
	return holder.echo.QueryString()
}

func (holder *RequestValue) Forms() map[string]string {
	holder.formParam.Do(func() {
		holder.formParam.v = make(map[string]string)
		formdatas, err := holder.echo.FormParams()
		if err != nil {
			return
		}
		for key := range formdatas {
			holder.formParam.v[key] = holder.echo.FormValue(key)
		}
	})
	return holder.formParam.v
}
func (holder *RequestValue) FormString() string {
	s := make([]string, 0, len(holder.Forms()))
	for key := range holder.Forms() {
		s = append(s, fmt.Sprintf("%s=%s", key, holder.Forms()[key]))
	}
	return strings.Join(s, "&")
}

func (holder *RequestValue) Body() []byte {
	holder.body.Do(func() {
		//body read all
		//ranout buffer
		holder.body.v, _ = ioutil.ReadAll(holder.echo.Request().Body) //read all body
		//restore
		holder.echo.Request().Body = ioutil.NopCloser(bytes.NewBuffer(holder.body.v))
	})
	return holder.body.v
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

func (holder RequestValue) Database() database.Context {
	return holder.db
}
func (holder RequestValue) SetDatabase(d database.Context) Context {
	holder.db = d
	return &holder
}

type (
	// TokenVerifier
	//  토큰 검증
	//  에러: Forbidden
	TokenVerifier func(Context) error

	// Binder
	//  요청 데이터 바인드
	//  에러: BadRequest
	Binder func(Context) error

	// Operator
	//  요청 처리
	//  에러: InternalServerError
	Operator func(Context) (interface{}, error)

	// Operate struct {
	// 	status int
	// 	Name   string
	// 	Behavior
	// 	Operator
	// }

	// HttpResponser
	//  응답
	//  에러: InternalServerError
	HttpResponsor func(echo.Context, int, interface{}) error

	Behavior func(Context, func(Context) (interface{}, error)) (interface{}, error)
)
type Option struct {
	TokenVerifier
	Binder
	Operator
	// Operates []Operate
	HttpResponsor
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

	exec_token_verifier := func(verifier TokenVerifier, ctx Context) (err error) {
		exceptions.Block{
			Try: func() {
				if verifier == nil {
					return //do noting
				}

				err = verifier(ctx) //exec settion-token verify
			},
			Catch: func(ex error) {
				err = ex
			},
		}.Do()
		return
	}

	exec_binder := func(bind Binder, ctx Context) (err error) {
		exceptions.Block{
			Try: func() {
				if bind == nil {
					exceptions.Throw("binder is nil")
				}
				err = bind(ctx) //exec bind
			},
			Catch: func(ex error) {
				err = ex
			},
		}.Do()
		return
	}

	exec_operator := func(behave Behavior, operate Operator, ctx Context) (out interface{}, err error) {
		exceptions.Block{
			Try: func() {
				if operate == nil {
					exceptions.Throw("operator is nil")
				}
				out, err = behave(ctx, operate) //exec operate
			},
			Catch: func(ex error) {
				err = ex
			},
		}.Do()
		return
	}

	exec_responser := func(response HttpResponsor, ctx echo.Context, status int, v interface{}) (err error) {
		exceptions.Block{
			Try: func() {
				if response == nil {
					exceptions.Throw("responser is nil")
				}
				err = response(ctx, status, v) //exec response
			},
			Catch: func(ex error) {
				err = ex
			},
		}.Do()
		return
	}

	return func(ctxEcho echo.Context) error {

		var context Context = &RequestValue{}
		context = context.SetEcho(ctxEcho)

		var (
			err error
			rsp interface{}

			reqPath   = func() string { return ctxEcho.Request().URL.Path }
			reqMethod = func() string { return ctxEcho.Request().Method }
			reqStatus = func() int { return ctxEcho.Response().Status }
			// reqForm   = func() string { return context.FormString() }
			reqParam = func() string { return context.ParamString() }
			reqQuery = func() string { return context.QueryString() }
			reqBody  = func() []byte { return context.Body() }
			rspBody  = func() func() []byte {
				var (
					onceRspbody sync.Once
					rspbody     []byte
				)
				return func() []byte {
					onceRspbody.Do(func() {
						right := func(b []byte, err error) []byte { return b }
						rspbody = right(json.Marshal(rsp))
					})
					return rspbody
				}
			}()
		)

		var (
			tidSink    = logs.WithId(ctxEcho.Response().Header().Get(echo.HeaderXRequestID))
			errRspSink = tidSink.WithName("failed to response")
		)

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
				var stack string
				logs.CauseIter(err, func(err error) {
					logs.StackIter(err, func(s string) {
						stack = s
					})
				})

				responseErrorSink := responseSink.WithName("error").WithValue(
					"error", err.Error(),
				)
				if 0 < len(stack) {
					responseErrorSink = responseErrorSink.WithValue(
						"stack", stack,
					)
				}

				logger.Error(responseErrorSink.String())
			})
			if logs.V(0) {
				responseSink = responseSink.WithValue("status", reqStatus())
				if logs.V(5) {
					responseSink = responseSink.WithValue("rspbody", rspBody())
				}
			}
			logger.V(0).Info(responseSink.String())
		}()
		err = exec_token_verifier(opt.TokenVerifier, context)
		if ErrorWithHandler(err,
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
				if err := context.Echo().String(status, err.Error()); err != nil {
					errRspSink := errRspSink.WithValue(
						"error", err.Error(),
					)
					logger.Error(errRspSink.String())
				}
			}, //실패 응답
		) {
			return err //return HandlerFunc
		}
		err = exec_binder(opt.Binder, context)
		if ErrorWithHandler(err,
			func(err error) {
				if erro := context.Echo().String(http.StatusBadRequest, err.Error()); erro != nil {
					errRspSink := errRspSink.WithValue(
						"error", err.Error(),
					)
					logger.Error(errRspSink.String())
				}
			}, //실패 응답
		) {
			return err //return HandlerFunc
		}
		rsp, err = exec_operator(opt.Behavior, opt.Operator, context)
		if ErrorWithHandler(err,
			func(err error) {
				//내부작업 오류
				if err_ := context.Echo().String(http.StatusInternalServerError, err.Error()); err_ != nil {
					errRspSink := errRspSink.WithValue(
						"error", err.Error(),
					)
					logger.Error(errRspSink.String())
				}
			}, //실패 응답
		) {
			return err //return HandlerFunc
		}
		err = exec_responser(opt.HttpResponsor, context.Echo(), http.StatusOK, rsp)
		if ErrorWithHandler(err,
			func(err error) {
				errRspSink := errRspSink.WithValue(
					"error", err.Error(),
				)
				logger.Error(errRspSink.String())
			},
		) {
			//응답 오류
			return err //return HandlerFunc
		}
		return nil
	}
}

func HttpJsonResponsor(ctx echo.Context, status int, v interface{}) error {
	return ctx.JSON(status, v)
}

func OK() interface{} {
	return "OK"
}

func Lock(engine *xorm.Engine) func(Context, func(Context) (interface{}, error)) (interface{}, error) {
	return func(ctx Context, operate func(Context) (interface{}, error)) (interface{}, error) {
		return engine.Transaction(func(s *xorm.Session) (interface{}, error) {
			ctx = ctx.SetDatabase(database.NewXormContext(s)) //new database context

			return operate(ctx)
		})
	}
}

func Nolock(engine *xorm.Engine) func(Context, func(Context) (interface{}, error)) (interface{}, error) {
	return func(ctx Context, operate func(Context) (interface{}, error)) (interface{}, error) {
		ctx = ctx.SetDatabase(database.NewXormContext(engine.NewSession())) //new database context
		defer ctx.Database().Close()
		//close
		return operate(ctx) //call operate
	}
}

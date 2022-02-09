package control

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/events"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/labstack/echo/v4"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"xorm.io/xorm"
)

type OperateContext struct {
	Http     echo.Context
	Database database.Context
	Req      interface{} //request value
}

// TokenVerifier
//  토큰 검증
//  에러: Forbidden
type TokenVerifier func(echo.Context) error

// Binder
//  요청 데이터 바인드
//  에러: BadRequest
type Binder func(echo.Context) (interface{}, error)

// Operator
//  요청 처리
//  에러: InternalServerError
type Operator_experimental func(echo.Context, interface{}) (interface{}, error)

// HttpResponser
//  응답
//  에러: InternalServerError
type HttpResponser func(echo.Context, int, interface{}) error

type Option_experimental struct {
	TokenVerifier
	// TokenAuthorized
	// TokenSetter
	Binder
	Operator Operator_experimental
	HttpResponser
}

// MakeMiddlewareFunc
//  @param: Option
//  @return: echo.HandlerFunc; func(echo.Context) error
func MakeMiddlewareFunc_experimental(opt Option_experimental) echo.HandlerFunc {

	exec_token_verifier := func(verifier TokenVerifier, ctx echo.Context) error {
		if verifier == nil {
			return nil //not error
		}
		err := verifier(ctx) //exec settion-token verify
		if err != nil {
			return err
		}
		return nil
	}

	exec_binder := func(bind Binder, ctx echo.Context) (interface{}, error) {
		if bind == nil {
			return nil, errors.New("without binder")
		}
		req, err := bind(ctx) //exec bind
		if err != nil {
			return nil, err
		}
		return req, nil
	}

	exec_operator := func(operate Operator_experimental, ctx echo.Context, v interface{}) (interface{}, error) {
		if operate == nil {
			return nil, errors.New("without operator")
		}
		rsp, err := operate(ctx, v) //exec operate
		if err != nil {
			return nil, err
		}
		return rsp, nil
	}

	exec_responser := func(response HttpResponser, ctx echo.Context, status int, v interface{}) error {
		if response == nil {
			return errors.New("without responser")
		}

		err := response(ctx, status, v) //exec response
		if err != nil {
			return err
		}
		return nil
	}

	return func(ctx echo.Context) error {

		var err error
		var req, rsp interface{}

		//event invoke
		defer func() {

			path := ctx.Request().URL.Path
			method := ctx.Request().Method
			status := ctx.Response().Status
			query := ctx.QueryString()

			args := map[string]interface{}{
				"path":    path,
				"query":   query,
				"method":  method,
				"reqbody": req,
				"rspbody": rsp,
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
			path := ctx.Request().URL.Path
			method := ctx.Request().Method
			status := ctx.Response().Status
			query := ctx.QueryString()
			reqbody, _ := json.Marshal(req)
			rspbody, _ := json.Marshal(rsp)

			log.Printf("DEBUG: C: time='%v' method='%s' path='%s' query='%s' reqbody='%s'\n", time.Now(), method, path, query, string(reqbody))
			log.Printf("DEBUG: S: time='%v' method='%s' path='%s' query='%s' rspbody='%s' status='%d'%s\n", time.Now(), method, path, query, string(rspbody), status, errorToString("error='%s'", err))
		}()

		err = exec_token_verifier(opt.TokenVerifier, ctx)
		if ErrorWithHandler(err,
			func(err error) { println("token_verify error:", err.Error()) },
			//additional error handler
		) {
			status := http.StatusForbidden

			e, ok := err.(*SessionTokenError)
			if ok {
				status = e.HttpStatus
			}
			//세션-토큰 검증 오류
			if err := ctx.String(status, err.Error()); err != nil {
				//TODO: 응답 실패 처리
			}

			return err //return
		}

		req, err = exec_binder(opt.Binder, ctx)
		if ErrorWithHandler(err,
			func(err error) { println("bind error:", err.Error()) },
			//additional error handler
		) {
			//요청 오류
			if err := ctx.String(http.StatusBadRequest, err.Error()); err != nil {
				//TODO: 응답 실패 처리
			}
			return err //return
		}

		rsp, err = exec_operator(opt.Operator, ctx, req)
		if ErrorWithHandler(err,
			func(err error) { println("operate error:", err.Error()) },
			//additional error handler
		) {
			//내부작업 오류
			if err := ctx.String(http.StatusInternalServerError, err.Error()); err != nil {
				//TODO: 응답 실패 처리
			}
			return err //return
		}

		err = exec_responser(opt.HttpResponser, ctx, http.StatusOK, rsp)
		if ErrorWithHandler(err,
			func(err error) { println("response error:", err.Error()) },
			//additional error handler
		) {
			//응답 오류
			return err //return
		}
		return nil
	}
}

func errorToString(format string, err error) (out string) {

	if err != nil {
		out = fmt.Sprintf(format, err.Error())
	}
	return out
}

func MakeBlockWithLock(engine *xorm.Engine, operate func(OperateContext) (interface{}, error)) func(echo.Context, interface{}) (interface{}, error) {

	return func(http echo.Context, in interface{}) (interface{}, error) {
		var out interface{}
		var err error

		ctx := database.NewContext(engine)
		defer func() { ctx.Close() }()

		tx := ctx.Tx()
		tx.Begin() //begin transaction
		defer func() {
			if err != nil {
				tx.Rollback() //rollback
			} else {
				tx.Commit() //commit
			}
		}()

		out, err = operate(OperateContext{Http: http, Database: ctx, Req: in})
		return out, err
	}
}

func MakeBlockNoLock(engine *xorm.Engine, operate func(OperateContext) (interface{}, error)) func(echo.Context, interface{}) (interface{}, error) {

	return func(http echo.Context, in interface{}) (interface{}, error) {
		var out interface{}
		var err error

		ctx := database.NewContext(engine)
		defer func() { ctx.Close() }()

		out, err = operate(OperateContext{Http: http, Database: ctx, Req: in})
		return out, err
	}
}

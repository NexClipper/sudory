package control

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/events"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/labstack/echo/v4"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"xorm.io/xorm"
)

type Binder func(echo.Context) (interface{}, error)

type Operator func(database.Context, interface{}) (interface{}, error)

type OperationBlockHandler func(engine *xorm.Engine, fn func(ctx database.Context) (interface{}, error)) (interface{}, error)

type HttpResponser func(echo.Context, int, interface{}) error

type CustomReporter func(interface{}) error

type Option struct {
	*xorm.Engine
	Binder
	Operator
	BlockMaker OperationBlockHandler
	HttpResponser
	CustomReporters []CustomReporter
}

// MakeMiddlewareFunc
//  @param:
//            bind Binder | Operator에서 사용하는 데이터 형식에 맞추어 요청 데이터를 변환 및 검증
//       operate Operator | 요청에 대한 처리
//    report HttpReporter | 응답 핸들러
//  @return: echo.HandlerFunc; func(echo.Context) error
func MakeMiddlewareFunc(opt Option) echo.HandlerFunc {
	exec_bind := func(bind Binder, ctx echo.Context) (interface{}, error) {
		if bind == nil {
			return nil, errors.New("without binder")
		}
		req, err := bind(ctx) //exec bind
		if err != nil {
			return nil, err
		}
		return req, nil
	}

	exec_operate := func(operate Operator, v interface{}) (interface{}, error) {
		if operate == nil {
			return nil, errors.New("without operator")
		}
		rsp, err := opt.BlockMaker(opt.Engine, func(ctx database.Context) (interface{}, error) {
			return operate(ctx, v) //exec operate
		})
		if err != nil {
			return nil, err
		}
		return rsp, nil
	}

	exec_response := func(response HttpResponser, ctx echo.Context, status int, v interface{}) error {
		if response == nil {
			return errors.New("without responser")
		}
		err := response(ctx, status, v) //exec report
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
			events.Invoke(ctx, req, rsp, err)
		}()

		req, err = exec_bind(opt.Binder, ctx)
		if ErrorWithHandler(err,
			func(err error) { println(err) },
			//additional error handler
		) {
			//요청 오류
			if err := exec_response(opt.HttpResponser, ctx, http.StatusBadRequest, err.Error()); err != nil {
				//TODO: 응답 실패 처리
			}
			return err //return
		}

		rsp, err = exec_operate(opt.Operator, req)
		if ErrorWithHandler(err,
			func(err error) { println(err) },
			//additional error handler
		) {
			//내부작업 오류
			if err := exec_response(opt.HttpResponser, ctx, http.StatusInternalServerError, err.Error()); err != nil {
				//TODO: 응답 실패 처리
			}
			return err //return
		}
		err = exec_response(opt.HttpResponser, ctx, http.StatusOK, rsp)
		if ErrorWithHandler(err,
			func(err error) { println(err) },
			//additional error handler
		) {
			//응답 오류
			return err //return
		}
		return nil
	}
}

func HttpResponse(ctx echo.Context, status int, v interface{}) error {
	return ctx.JSON(status, v)
}

func Lock(engine *xorm.Engine, operate func(ctx database.Context) (interface{}, error)) (interface{}, error) {
	var v interface{}
	var err error

	ctx := database.NewContext(engine)
	/*
		defer func() { ctx.Close() }()
	*/
	tx := ctx.Tx()
	tx.Begin() //begin transaction
	defer func() {
		if err != nil {
			tx.Rollback() //rollback
		} else {
			tx.Commit() //commit
		}
	}()

	v, err = operate(ctx)

	return v, err
}
func NoLock(engine *xorm.Engine, operate func(ctx database.Context) (interface{}, error)) (interface{}, error) {
	var v interface{}
	var err error

	ctx := database.NewContext(engine)
	/*
		defer func() { ctx.Close() }()
	*/
	v, err = operate(ctx)
	return v, err
}

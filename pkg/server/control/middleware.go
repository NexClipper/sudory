package control

import (
	"net/http"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/labstack/echo/v4"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

type Binder func(ctx echo.Context) (interface{}, error)

type Operator func(v interface{}) (interface{}, error)

type HttpReporter func(echo.Context, int, interface{}) error
type CustomReporter func(interface{}) error

func MakeMiddlewareFunc(bind Binder, operate Operator, report HttpReporter) func(ctx echo.Context) error {

	exec_bind := func(bind Binder, ctx echo.Context) (interface{}, error) {
		if operate == nil {
			return nil, errors.New("without reciver binder")
		}

		req, err := bind(ctx)
		if err != nil {
			err = errors.New(err.Error()) //TODO: define custom error
		}
		return req, err
	}

	exec_operate := func(exec Operator, v interface{}) (interface{}, error) {
		if exec == nil {
			return nil, errors.New("must set operation executor")
		}

		rsp, err := exec(v)
		if err != nil {
			err = errors.New(err.Error()) //TODO: define custom error
		}
		return rsp, err
	}

	exec_response := func(report HttpReporter, ctx echo.Context, status int, v interface{}) error {
		if report == nil {
			return errors.New("must set http responser") //TODO: define custom error
		}

		err := report(ctx, status, v)
		if err != nil {
			err = errors.New(err.Error()) //TODO: define custom error
		}
		return err
	}

	return func(ctx echo.Context) error {

		req, err := exec_bind(bind, ctx)
		if ErrorWithHandler(err,
			func(err error) { println(err) }) {
			//additional error handler
			exec_response(report, ctx, http.StatusBadRequest, err.Error())
			return err //return closure
		}
		rsp, err := exec_operate(operate, req)
		if ErrorWithHandler(err,
			func(err error) { println(err) }) {
			//additional error handler
			exec_response(report, ctx, http.StatusInternalServerError, err.Error())
			return err //return closure
		}
		err = exec_response(report, ctx, http.StatusOK, rsp)
		if ErrorWithHandler(err,
			func(err error) { println("failed to report") }) {
			//additional error handler
			return err //return closure
		}
		return nil
	}
}

func HttpResponse(ctx echo.Context, status int, v interface{}) error {
	return ctx.JSON(status, v)
}

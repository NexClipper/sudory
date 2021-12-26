package operator

import "github.com/labstack/echo/v4"

type ResponseFn func(ctx echo.Context) error

type Operator interface {
	Create(ctx echo.Context) error
}

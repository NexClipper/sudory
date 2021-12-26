package view

import "github.com/labstack/echo/v4"

type Viewer interface {
	Request(ctx echo.Context) error
	Response(ctx echo.Context) error
}

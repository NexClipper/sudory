package view

import (
	"github.com/NexClipper/sudory-prototype-r1/pkg/model"
	"github.com/labstack/echo/v4"
)

type Viewer interface {
	Request(ctx echo.Context) error
	Response(ctx echo.Context, m model.Modeler) error
}

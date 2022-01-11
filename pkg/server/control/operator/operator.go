package operator

import (
	"github.com/NexClipper/sudory/pkg/server/model"
	"github.com/labstack/echo/v4"
)

type ResponseFn func(ctx echo.Context, m model.Modeler) error

type Operator interface {
	Create(ctx echo.Context) error
	Get(ctx echo.Context) error
}

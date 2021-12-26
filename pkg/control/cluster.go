package control

import (
	"github.com/NexClipper/sudory-prototype-r1/pkg/control/operator"
	"github.com/NexClipper/sudory-prototype-r1/pkg/view"
	"github.com/labstack/echo/v4"
)

func (c *Control) CreateCluster(ctx echo.Context) error {
	v := view.NewCreateCluster(operator.NewCluster())
	return v.Request(ctx)
}

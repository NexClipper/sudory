package control

import (
	"github.com/NexClipper/sudory-prototype-r1/pkg/control/operator"
	"github.com/NexClipper/sudory-prototype-r1/pkg/view"
	"github.com/labstack/echo/v4"
)

// CreateService
// @Description Create a Service
// @Accept json
// @Produce json
// @Tags server
// @Router /server/service [post]
// @Param service body model.ReqService true "Service의 정보"
// @Success 200
func (c *Control) CreateService(ctx echo.Context) error {
	v := view.NewCreateService(operator.NewService(c.db))
	return v.Request(ctx)
}

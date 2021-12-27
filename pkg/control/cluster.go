package control

import (
	"github.com/NexClipper/sudory-prototype-r1/pkg/control/operator"
	"github.com/NexClipper/sudory-prototype-r1/pkg/view"
	"github.com/labstack/echo/v4"
)

// CreateCluster godoc
// @Description Create a Cluster
// @Accept json
// @Produce json
// @Router /clusters [post]
// @Param namespace body model.ReqCluster true "Cluster의 정보"
// @Success 200 {object} model.Cluster
func (c *Control) CreateCluster(ctx echo.Context) error {
	v := view.NewCreateCluster(operator.NewCluster(c.db))
	return v.Request(ctx)
}

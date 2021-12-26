package view

import (
	"net/http"

	"github.com/NexClipper/sudory-prototype-r1/pkg/control/operator"
	"github.com/NexClipper/sudory-prototype-r1/pkg/model"
	"github.com/labstack/echo/v4"
)

type CreateCluster struct {
	opr *operator.Cluster
}

func NewCreateCluster(o operator.Operator) Viewer {
	return &CreateCluster{opr: o.(*operator.Cluster)}
}

func (v *CreateCluster) fromModel(m *model.ReqCluster) {
	v.opr.Name = m.Name
	v.opr.Response = v.Response
}

func (v *CreateCluster) Request(ctx echo.Context) error {
	reqModel := &model.ReqCluster{}
	if err := ctx.Bind(reqModel); err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	v.fromModel(reqModel)
	if err := v.opr.Create(ctx); err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	return nil
}

func (v *CreateCluster) Response(ctx echo.Context) error {
	return nil
}

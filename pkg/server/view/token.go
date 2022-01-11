package view

import (
	"net/http"
	"strconv"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/model"
	"github.com/labstack/echo/v4"
)

type CreateToken struct {
	opr *operator.Token
}

func NewCreateToken(o operator.Operator) Viewer {
	return &CreateToken{opr: o.(*operator.Token)}
}

func (v *CreateToken) fromModel(clusterID uint64, m *model.ReqToken) {
	v.opr.ClusterID = clusterID
	v.opr.Key = m.Key
	v.opr.Response = v.Response
}

func (v *CreateToken) Request(ctx echo.Context) error {
	id := ctx.Param("id")

	clusterID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	reqModel := &model.ReqToken{}
	if err := ctx.Bind(reqModel); err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	v.fromModel(clusterID, reqModel)
	if err := v.opr.Create(ctx); err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	return nil
}

func (v *CreateToken) Response(ctx echo.Context, m model.Modeler) error {
	if err := ctx.JSON(http.StatusOK, nil); err != nil {
		return err
	}
	return nil
}

package view

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/model"
	"github.com/labstack/echo/v4"
)

type GetCatalogue struct {
	opr *operator.Catalogue
}

func NewGetCatalogue(o operator.Operator) Viewer {
	return &GetCatalogue{opr: o.(*operator.Catalogue)}
}

func (v *GetCatalogue) fromModel() {
	v.opr.Response = v.Response
}

func (v *GetCatalogue) Request(ctx echo.Context) error {
	v.fromModel()
	if err := v.opr.Get(ctx); err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	return nil
}

func (v *GetCatalogue) Response(ctx echo.Context, m model.Modeler) error {
	catalogues := m.(*model.Catalogues)
	return ctx.JSON(http.StatusOK, catalogues)
}

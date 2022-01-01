package view

import (
	"net/http"

	"github.com/NexClipper/sudory-prototype-r1/pkg/control/operator"
	"github.com/NexClipper/sudory-prototype-r1/pkg/model"
	"github.com/labstack/echo/v4"
)

type CreateService struct {
	opr *operator.Service
}

func NewCreateService(o operator.Operator) Viewer {
	return &CreateService{opr: o.(*operator.Service)}
}

func (v *CreateService) fromModel(m *model.ReqService) {
	v.opr.Name = m.Name
	v.opr.ClusterID = m.ClusterID
	v.opr.StepCount = m.StepCount

	for _, s := range m.Step {
		oprStep := &operator.Step{
			Name:      s.Name,
			Sequence:  s.Sequence,
			Command:   s.Command,
			Parameter: s.Parameter,
		}
		v.opr.Steps = append(v.opr.Steps, oprStep)
	}

	v.opr.Response = v.Response
}

func (v *CreateService) Request(ctx echo.Context) error {
	reqModel := &model.ReqService{}
	if err := ctx.Bind(reqModel); err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	v.fromModel(reqModel)
	if err := v.opr.Create(ctx); err != nil {
		return ctx.JSON(http.StatusInternalServerError, nil)
	}

	return nil
}

func (v *CreateService) Response(ctx echo.Context, m model.Modeler) error {
	if err := ctx.JSON(http.StatusOK, nil); err != nil {
		return err
	}
	return nil
}

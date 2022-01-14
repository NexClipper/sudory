package view

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/model"
	tcommandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
	"github.com/labstack/echo/v4"
)

type CreateTemplateCommand struct {
	opr *operator.CreateTemplateCommand
}

var _ Viewer = (*CreateTemplateCommand)(nil)

func NewCreateTemplateCommand(o operator.Creator) Viewer {
	v := &CreateTemplateCommand{opr: o.(*operator.CreateTemplateCommand)}
	v.opr.Response = v.Response
	return v
}

func (v *CreateTemplateCommand) fromModel(m tcommandv1.HttpReqTemplateCommand) {
	v.opr.TemplateCommand = m.TemplateCommand
	// v.opr.Response = v.Response
}

func (v *CreateTemplateCommand) Request(ctx echo.Context) error {
	reqModel := tcommandv1.HttpReqTemplateCommand{}
	if err := ctx.Bind(&reqModel); err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	v.fromModel(reqModel)
	if err := v.opr.Create(ctx); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return nil
}

func (v *CreateTemplateCommand) Response(ctx echo.Context, m model.Modeler) error {
	if err := ctx.JSON(http.StatusOK, nil); err != nil {
		return err
	}
	return nil
}

type GetTemplateCommand struct {
	opr *operator.GetTemplateCommand
}

var _ Viewer = (*GetTemplateCommand)(nil)

func NewGetTemplateCommand(o operator.Getter) Viewer {
	v := &GetTemplateCommand{opr: o.(*operator.GetTemplateCommand)}
	v.opr.Response = v.Response
	return v
}

func (v *GetTemplateCommand) fromModel(param map[string]string) {
	v.opr.Params = param
	// v.opr.Response = v.Response
}

func (v *GetTemplateCommand) Request(ctx echo.Context) error {

	param := map[string]string{
		"uuid": ctx.Param("uuid"),
	}

	v.fromModel(param)
	if err := v.opr.Get(ctx); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return nil
}

func (v *GetTemplateCommand) Response(ctx echo.Context, m model.Modeler) error {
	rsp, ok := m.(*tcommandv1.HttpRspTemplateCommand)
	if !ok {
		println(ok)
	}
	return ctx.JSON(http.StatusOK, rsp)
}

type UpdateTemplateCommand struct {
	opr *operator.UpdateTemplateCommand
}

var _ Viewer = (*UpdateTemplateCommand)(nil)

func NewUpdateTemplateCommand(o operator.Updater) Viewer {
	v := &UpdateTemplateCommand{opr: o.(*operator.UpdateTemplateCommand)}
	v.opr.Response = v.Response
	return v
}

func (v *UpdateTemplateCommand) fromModel(m *tcommandv1.HttpReqTemplateCommand) {
	v.opr.TemplateCommand = m.TemplateCommand
	// v.opr.Response = v.Response
}

func (v *UpdateTemplateCommand) Request(ctx echo.Context) error {
	reqModel := &tcommandv1.HttpReqTemplateCommand{}
	if err := ctx.Bind(reqModel); err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	v.fromModel(reqModel)
	if err := v.opr.Update(ctx); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return nil
}

func (v *UpdateTemplateCommand) Response(ctx echo.Context, m model.Modeler) error {
	if err := ctx.JSON(http.StatusOK, m); err != nil {
		return err
	}
	return nil
}

type DeleteTemplateCommand struct {
	opr *operator.DeleteTemplateCommand
}

var _ Viewer = (*DeleteTemplate)(nil)

func NewDeleteTemplateCommand(o operator.Remover) Viewer {
	v := &DeleteTemplateCommand{opr: o.(*operator.DeleteTemplateCommand)}
	v.opr.Response = v.Response
	return v
}

func (v *DeleteTemplateCommand) fromModel(param map[string]string) {
	v.opr.Params = param
	// v.opr.Response = v.Response
}

func (v *DeleteTemplateCommand) Request(ctx echo.Context) error {
	param := map[string]string{
		"uuid": ctx.Param("uuid"),
	}

	if len(param["uuid"]) == 0 {
		return ctx.JSON(http.StatusBadRequest, nil)
	}
	v.fromModel(param)
	if err := v.opr.Delete(ctx); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return nil
}

func (v *DeleteTemplateCommand) Response(ctx echo.Context, m model.Modeler) error {
	rsp := m.(*tcommandv1.HttpRspTemplateCommand)
	return ctx.JSON(http.StatusOK, rsp)
}

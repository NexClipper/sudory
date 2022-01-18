package view

// import (
// 	"net/http"

// 	"github.com/NexClipper/sudory/pkg/server/control/operator"
// 	"github.com/NexClipper/sudory/pkg/server/model"
// 	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
// 	"github.com/labstack/echo/v4"
// )

// type CreateTemplate struct {
// 	opr *operator.CreateTemplate
// }

// var _ Viewer = (*CreateTemplate)(nil)

// func NewCreateTemplate(o operator.Creator) Viewer {
// 	v := &CreateTemplate{opr: o.(*operator.CreateTemplate)}
// 	v.opr.Response = v.Response
// 	return v
// }

// func (v *CreateTemplate) fromModel(m templatev1.HttpReqTemplates) {
// 	v.opr.HttpReqTemplates = m
// 	// v.opr.Response = v.Response
// }

// func (v *CreateTemplate) Request(ctx echo.Context) error {
// 	reqModel := &templatev1.HttpReqTemplates{}
// 	if err := ctx.Bind(reqModel); err != nil {
// 		return ctx.JSON(http.StatusBadRequest, nil)
// 	}

// 	v.fromModel(*reqModel)
// 	if err := v.opr.Create(ctx); err != nil {
// 		return ctx.JSON(http.StatusInternalServerError, err)
// 	}

// 	return nil
// }

// func (v *CreateTemplate) Response(ctx echo.Context, m model.Modeler) error {
// 	if err := ctx.JSON(http.StatusOK, nil); err != nil {
// 		return err
// 	}
// 	return nil
// }

// type GetTemplate struct {
// 	opr *operator.GetTemplate
// }

// var _ Viewer = (*GetTemplate)(nil)

// func NewGetTemplate(o operator.Getter) Viewer {
// 	v := &GetTemplate{opr: o.(*operator.GetTemplate)}
// 	v.opr.Response = v.Response
// 	return v
// }

// func (v *GetTemplate) fromModel(param map[string]string) {
// 	v.opr.Params = param
// 	// v.opr.Response = v.Response
// }

// func (v *GetTemplate) Request(ctx echo.Context) error {

// 	param := map[string]string{
// 		"uuid": ctx.Param("uuid"),
// 	}

// 	v.fromModel(param)
// 	if err := v.opr.Get(ctx); err != nil {
// 		return ctx.JSON(http.StatusInternalServerError, err)
// 	}

// 	return nil
// }

// func (v *GetTemplate) Response(ctx echo.Context, m model.Modeler) error {
// 	rsp, ok := m.(*templatev1.HttpRspTemplate)
// 	if !ok {
// 		println(ok)
// 	}
// 	return ctx.JSON(http.StatusOK, rsp)
// }

// type FindTemplate struct {
// 	opr *operator.FindTemplate
// }

// var _ Viewer = (*FindTemplate)(nil)

// func NewFindTemplate(o operator.Getter) Viewer {
// 	v := &FindTemplate{opr: o.(*operator.FindTemplate)}
// 	v.opr.Response = v.Response
// 	return v
// }

// func (v *FindTemplate) fromModel(param map[string]string) {
// 	v.opr.Params = param
// 	// v.opr.Response = v.Response
// }

// func (v *FindTemplate) Request(ctx echo.Context) error {
// 	param := map[string]string{
// 		"uuid":   ctx.QueryParam("uuid"),
// 		"name":   ctx.QueryParam("name"),
// 		"origin": ctx.QueryParam("origin"),
// 	}

// 	v.fromModel(param)
// 	if err := v.opr.Get(ctx); err != nil {
// 		return ctx.JSON(http.StatusInternalServerError, err)
// 	}

// 	return nil
// }

// func (v *FindTemplate) Response(ctx echo.Context, m model.Modeler) error {
// 	rsp := m.(*templatev1.HttpRspTemplates)
// 	return ctx.JSON(http.StatusOK, rsp)
// }

// type UpdateTemplate struct {
// 	opr *operator.UpdateTemplate
// }

// var _ Viewer = (*UpdateTemplate)(nil)

// func NewUpdateTemplate(o operator.Updater) Viewer {
// 	v := &UpdateTemplate{opr: o.(*operator.UpdateTemplate)}
// 	v.opr.Response = v.Response
// 	return v
// }

// func (v *UpdateTemplate) fromModel(m templatev1.HttpReqTemplate) {
// 	v.opr.Template = m.Template
// 	// v.opr.Response = v.Response
// }

// func (v *UpdateTemplate) Request(ctx echo.Context) error {
// 	reqModel := &templatev1.HttpReqTemplate{}
// 	if err := ctx.Bind(reqModel); err != nil {
// 		return ctx.JSON(http.StatusBadRequest, nil)
// 	}

// 	v.fromModel(*reqModel)
// 	if err := v.opr.Update(ctx); err != nil {
// 		return ctx.JSON(http.StatusInternalServerError, err)
// 	}

// 	return nil
// }

// func (v *UpdateTemplate) Response(ctx echo.Context, m model.Modeler) error {
// 	if err := ctx.JSON(http.StatusOK, nil); err != nil {
// 		return err
// 	}
// 	return nil
// }

// type DeleteTemplate struct {
// 	opr *operator.DeleteTemplate
// }

// var _ Viewer = (*DeleteTemplate)(nil)

// func NewDeleteTemplate(o operator.Remover) Viewer {
// 	v := &DeleteTemplate{opr: o.(*operator.DeleteTemplate)}
// 	v.opr.Response = v.Response
// 	return v
// }

// func (v *DeleteTemplate) fromModel(param map[string]string) {
// 	v.opr.Params = param
// 	// v.opr.Response = v.Response
// }

// func (v *DeleteTemplate) Request(ctx echo.Context) error {
// 	param := map[string]string{
// 		"uuid": ctx.Param("uuid"),
// 	}

// 	if len(param["uuid"]) == 0 {
// 		return ctx.JSON(http.StatusBadRequest, nil)
// 	}
// 	v.fromModel(param)
// 	if err := v.opr.Delete(ctx); err != nil {
// 		return ctx.JSON(http.StatusInternalServerError, err)
// 	}

// 	return nil
// }

// func (v *DeleteTemplate) Response(ctx echo.Context, m model.Modeler) error {
// 	rsp := m.(*templatev1.HttpRspTemplate)
// 	return ctx.JSON(http.StatusOK, rsp)
// }

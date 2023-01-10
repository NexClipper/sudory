package control

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/model/template/v3"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @Description Get a template (v2)
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/template
// @Router      /v2/server/template/{uuid} [get]
// @Param       uuid     path string true  "Template Uuid"
// @Success 200 {object} template.HttpRsp_Template
func (ctl ControlVanilla) GetTemplate_v2(ctx echo.Context) (err error) {
	uuid := echoutil.Param(ctx)[__UUID__]

	err = echoutil.WrapHttpError(http.StatusBadRequest,
		func() (err error) {
			if len(echoutil.Param(ctx)[__UUID__]) == 0 {
				err = ErrorInvalidRequestParameter
			}
			return errors.Wrapf(err, "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				))
		},
	)
	if err != nil {
		return
	}

	tmpl := template.Template{}
	tmpl.Uuid = uuid

	tmpl_cond := stmt.And(
		stmt.Equal("uuid", tmpl.Uuid),
		stmt.IsNull("deleted"),
	)

	err = ctl.dialect.QueryRow(tmpl.TableName(), tmpl.ColumnNames(), tmpl_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := tmpl.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusOK, template.HttpRsp_Template(tmpl))
}

// @Description Find []template (v2)
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/template
// @Router      /v2/server/template [get]
// @Param       q            query  string false "query  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       o            query  string false "order  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       p            query  string false "paging github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Success 200 {array} template.HttpRsp_Template
func (ctl ControlVanilla) FindTemplate_v2(ctx echo.Context) error {
	q, err := stmt.ConditionLexer.Parse(echoutil.QueryParam(ctx)["q"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	o, err := stmt.OrderLexer.Parse(echoutil.QueryParam(ctx)["o"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	p, err := stmt.PaginationLexer.Parse(echoutil.QueryParam(ctx)["p"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	// additional conditon
	q = stmt.And(q,
		stmt.IsNull("deleted"),
	)
	// default pagination
	if p == nil {
		p = stmt.Limit(__DEFAULT_DECORATION_LIMIT__)
	}

	rsp := make([]template.HttpRsp_Template, 0, state.ENV__INIT_SLICE_CAPACITY__())
	var tmpl template.Template
	err = ctl.dialect.QueryRows(tmpl.TableName(), tmpl.ColumnNames(), q, o, p)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner, _ int) error {
			err := tmpl.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			rsp = append(rsp, template.HttpRsp_Template(tmpl))

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "failed to find template")
	}

	return ctx.JSON(http.StatusOK, rsp)
}

// @Description List template command (v2)
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/template_command
// @Router      /v2/server/template_command [get]
// @Success 200 {array} template.HttpRsp_TemplateCommand
func (ctl ControlVanilla) FindTemplateCommand_v2(ctx echo.Context) (err error) {
	q, err := stmt.ConditionLexer.Parse(echoutil.QueryParam(ctx)["q"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	o, err := stmt.OrderLexer.Parse(echoutil.QueryParam(ctx)["o"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	p, err := stmt.PaginationLexer.Parse(echoutil.QueryParam(ctx)["p"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	// additional conditon
	q = stmt.And(q,
		stmt.IsNull("deleted"),
	)
	// default pagination
	if p == nil {
		p = stmt.Limit(__DEFAULT_DECORATION_LIMIT__)
	}

	rsp := make([]template.TemplateCommand, 0, state.ENV__INIT_SLICE_CAPACITY__())
	var tmpl template.TemplateCommand
	err = ctl.dialect.QueryRows(tmpl.TableName(), tmpl.ColumnNames(), q, o, p)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner, _ int) error {
			err := tmpl.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			rsp = append(rsp, tmpl)

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "failed to find template_command")
	}

	return ctx.JSON(http.StatusOK, []template.HttpRsp_TemplateCommand(rsp))
}

// @Description Get a template command (v2)
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/template_command
// @Router      /v2/server/template_command/{uuid} [get]
// @Param       name     path   string true  "HttpReqTemplateCommand 의 Uuid"
// @Success 200 {object} template.HttpRsp_TemplateCommand
func (ctl ControlVanilla) GetTemplateCommand_v2(ctx echo.Context) (err error) {
	name := echoutil.Param(ctx)[__NAME__]

	err = echoutil.WrapHttpError(http.StatusBadRequest,
		// func() (err error) {
		// 	if len(echoutil.Param(ctx)[__TEMPLATE_UUID__]) == 0 {
		// 		err = ErrorInvalidRequestParameter
		// 	}
		// 	return errors.Wrapf(err, "valid param%s",
		// 		logs.KVL(
		// 			ParamLog(__TEMPLATE_UUID__, echoutil.Param(ctx)[__TEMPLATE_UUID__])...,
		// 		))
		// },
		func() (err error) {
			if len(echoutil.Param(ctx)[__NAME__]) == 0 {
				err = ErrorInvalidRequestParameter
			}
			return errors.Wrapf(err, "valid param%s",
				logs.KVL(
					ParamLog(__NAME__, echoutil.Param(ctx)[__NAME__])...,
				))
		},
	)
	if err != nil {
		return
	}

	command := template.TemplateCommand{}
	command.Name = name

	eq_uuid := stmt.Equal("name", command.Name)

	err = ctl.dialect.QueryRow(command.TableName(), command.ColumnNames(), eq_uuid, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := command.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusOK, command)
}

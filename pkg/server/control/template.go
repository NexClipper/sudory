package control

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/model/template/v2"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// // @deprecated
// // @Description Create a template
// // @Security    XAuthToken
// // @Accept      json
// // @Produce     json
// // @Tags        server/template
// // @Router      /server/template [post]
// // @Param       template     body   v1.HttpReqTemplate_Create true  "HttpReqTemplate_Create"
// // @Success 200 {object} v1.HttpRspTemplate
// func (ctl Control) CreateTemplate(ctx echo.Context) error {
// 	map_command := func(elems []commandv1.HttpReqTemplateCommand_Create_ByTemplate, mapper func(int, commandv1.HttpReqTemplateCommand_Create_ByTemplate) commandv1.TemplateCommand) []commandv1.TemplateCommand {
// 		rst := make([]commandv1.TemplateCommand, len(elems))
// 		for n := range elems {
// 			rst[n] = mapper(n, elems[n])
// 		}
// 		return rst
// 	}

// 	foreach_command := func(elems []commandv1.TemplateCommand, mapper func(int, commandv1.TemplateCommand) error) error {
// 		for n := range elems {
// 			if err := mapper(n, elems[n]); err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	}

// 	body := new(templatev1.HttpReqTemplate_Create)
// 	if err := echoutil.Bind(ctx, body); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
// 			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
// 				logs.KVL(
// 					"type", TypeName(body),
// 				)))
// 	}
// 	if len(body.Name) == 0 {
// 		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
// 			errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
// 				logs.KVL(
// 					ParamLog(fmt.Sprintf("%s.Name", TypeName(body)), body.Name)...,
// 				)))
// 	}
// 	if len(body.Origin) == 0 {
// 		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
// 			errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
// 				logs.KVL(
// 					ParamLog(fmt.Sprintf("%s.Origin", TypeName(body)), body.Origin)...,
// 				)))
// 	}
// 	if len(body.Commands) == 0 {
// 		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
// 			errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
// 				logs.KVL(
// 					ParamLog(fmt.Sprintf("%s.Commands", TypeName(body)), body.Commands)...,
// 				)))
// 	}
// 	//valid commands
// 	for _, command := range body.Commands {
// 		if len(command.Name) == 0 {
// 			return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
// 				errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
// 					logs.KVL(
// 						ParamLog(fmt.Sprintf("%s.Name", TypeName(command)), command.Name)...,
// 					)))
// 		}
// 		if len(command.Method) == 0 {
// 			return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
// 				errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
// 					logs.KVL(
// 						ParamLog(fmt.Sprintf("%s.Method", TypeName(command)), command.Method)...,
// 					)))
// 		}
// 	}

// 	//property
// 	template := templatev1.Template{}
// 	template.UuidMeta = metav1.NewUuidMeta()
// 	template.LabelMeta = metav1.NewLabelMeta(body.Name, body.Summary)
// 	template.Origin = body.Origin

// 	//create command
// 	commands := map_command(body.Commands, func(i int, iter commandv1.HttpReqTemplateCommand_Create_ByTemplate) commandv1.TemplateCommand {
// 		command := commandv1.TemplateCommand{}
// 		command.UuidMeta = metav1.NewUuidMeta()
// 		command.LabelMeta = metav1.NewLabelMeta(iter.Name, iter.Summary)
// 		command.TemplateUuid = template.Uuid
// 		command.Sequence = newist.Int32(int32(i))
// 		command.Method = iter.Method
// 		command.Args = iter.Args
// 		command.ResultFilter = &iter.ResultFilter

// 		return command
// 	})

// 	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
// 		//create template
// 		template_, err := vault.NewTemplate(db).Create(template)
// 		if err != nil {
// 			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
// 				errors.Wrapf(err, "create template"))
// 		}

// 		if err := foreach_command(commands, func(i int, tc commandv1.TemplateCommand) error {
// 			command, err := vault.NewTemplateCommand(db).Create(tc)
// 			if err != nil {
// 				return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
// 					errors.Wrapf(err, "create template command"))
// 			}

// 			commands[i] = *command

// 			return nil
// 		}); err != nil {
// 			return nil, err
// 		}

// 		return templatev1.HttpRspTemplate{Template: *template_, Commands: commands}, nil
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return ctx.JSON(http.StatusOK, r)
// }

// @Description Get a template
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/template
// @Router      /server/template/{uuid} [get]
// @Param       uuid         path   string true  "Template Uuid"
// @Success 200 {object} template.HttpRsp_Template
func (ctl ControlVanilla) GetTemplate(ctx echo.Context) (err error) {
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

	commands := make([]template.TemplateCommand, 0, state.ENV__INIT_SLICE_CAPACITY__())
	command := template.TemplateCommand{}
	command.TemplateUuid = uuid

	command_cond := stmt.And(
		stmt.Equal("template_uuid", command.TemplateUuid),
		stmt.IsNull("deleted"),
	)

	order := stmt.Asc("sequence")

	err = ctl.dialect.QueryRows(command.TableName(), command.ColumnNames(), command_cond, order, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner, _ int) error {
			err := command.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			commands = append(commands, command)

			return err
		})
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusOK, template.HttpRsp_Template{Template: tmpl, Commands: commands})
}

// @Description Find []template
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/template
// @Router      /server/template [get]
// @Param       q            query  string false "query  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       o            query  string false "order  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       p            query  string false "paging github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Success 200 {array} template.HttpRsp_Template
func (ctl ControlVanilla) FindTemplate(ctx echo.Context) error {
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
	var set_tmpl = map[template.Template]struct{}{}
	var tmpl template.Template
	err = ctl.dialect.QueryRows(tmpl.TableName(), tmpl.ColumnNames(), q, o, p)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner, _ int) error {
			err := tmpl.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			set_tmpl[tmpl] = struct{}{}

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "failed to find template")
	}

	for tmpl := range set_tmpl {
		commands, err := ListTemplateCommand(ctx.Request().Context(), ctl, tmpl.Uuid)
		if err != nil {
			return errors.Wrapf(err, "failed to get list")
		}
		rsp = append(rsp, template.HttpRsp_Template{Template: tmpl, Commands: commands})
	}

	return ctx.JSON(http.StatusOK, rsp)
}

// // @deprecated
// // @Description Update a template
// // @Security    XAuthToken
// // @Accept      json
// // @Produce     json
// // @Tags        server/template
// // @Router      /server/template/{uuid} [put]
// // @Param       uuid         path   string                    true  "Template Uuid"
// // @Param       template     body   v2.HttpReqTemplate_Update true  "HttpReqTemplate_Update"
// // @Success 200 {object} v2.Template
// func (ctl ControlVanilla) UpdateTemplate(ctx echo.Context) error {
// 	body := new(templatev2.HttpReqTemplate_Update)
// 	if err := echoutil.Bind(ctx, body); err != nil {
// 		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
// 			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
// 				logs.KVL(
// 					"type", TypeName(body),
// 				)))
// 	}

// 	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
// 		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
// 			errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
// 				logs.KVL(
// 					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
// 				)))
// 	}

// 	// template := body.Template
// 	uuid := echoutil.Param(ctx)[__UUID__]
// 	time_now := time.Now()

// 	//set uuid from path
// 	template := templatev2.Template{}
// 	template.Uuid = uuid
// 	template.Name = body.Name
// 	template.Summary = body.Summary
// 	template.Origin = body.Origin
// 	template.Updated = *vanilla.NewNullTime(time_now)

// 	updateSet := map[string]interface{}{}
// 	if 0 < len(body.Name) {
// 		updateSet["name"] = template.Name
// 	}
// 	if body.Summary.Valid {
// 		updateSet["summary"] = template.Summary
// 	}
// 	if 0 < len(body.Origin) {
// 		updateSet["naoriginme"] = template.Origin
// 	}
// 	updateSet["updated"] = template.Updated

// 	template_cond := stmt.And(
// 		stmt.Equal("uuid", template.Uuid),
// 		stmt.IsNull("deleted"),
// 	)

// 	err := func() error {
// 		//upate template
// 		affected, err := stmtex.Update(template.TableName(), updateSet, template_cond).
// 			ExecContext(ctx.Request().Context(), ctl, ctl.Dialect())
// 		if err != nil {
// 			return errors.Wrapf(err, "update template")
// 		}
// 		if affected == 0 {
// 			return errors.WithStack(ErrorNoAffected)
// 		}
// 		return nil
// 	}()
// 	if err != nil {
// 		return err
// 	}

// 	return ctx.JSON(http.StatusOK, template)
// }

// // @deprecated
// // @Description Delete a template
// // @Security    XAuthToken
// // @Accept      json
// // @Produce     json
// // @Tags        server/template
// // @Router      /server/template/{uuid} [delete]
// // @Param       uuid         path   string true  "Template 의 Uuid"
// // @Success 200
// func (ctl Control) DeleteTemplate(ctx echo.Context) error {
// 	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
// 		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
// 			errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
// 				logs.KVL(
// 					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
// 				)))
// 	}
// 	uuid := echoutil.Param(ctx)[__UUID__]

// 	_, err := ctl.Scope(func(db database.Context) (interface{}, error) {
// 		if err := vault.NewTemplateCommand(db).Delete_ByTemplate(uuid); err != nil {
// 			return nil, errors.Wrapf(err, "delete template command")
// 		}
// 		if err := vault.NewTemplate(db).Delete(uuid); err != nil {
// 			return nil, errors.Wrapf(err, "delete template")
// 		}

// 		return nil, nil
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return ctx.JSON(http.StatusOK, OK())
// }

package control

import (
	"fmt"
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
	templatev2 "github.com/NexClipper/sudory/pkg/server/model/template/v2"
	commandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @deprecated
// @Description Create a template
// @Accept      json
// @Produce     json
// @Tags        server/template
// @Router      /server/template [post]
// @Param       x_auth_token header string                    false "client session token"
// @Param       template     body   v1.HttpReqTemplate_Create true  "HttpReqTemplate_Create"
// @Success 200 {object} v1.HttpRspTemplate
func (ctl Control) CreateTemplate(ctx echo.Context) error {
	map_command := func(elems []commandv1.HttpReqTemplateCommand_Create_ByTemplate, mapper func(int, commandv1.HttpReqTemplateCommand_Create_ByTemplate) commandv1.TemplateCommand) []commandv1.TemplateCommand {
		rst := make([]commandv1.TemplateCommand, len(elems))
		for n := range elems {
			rst[n] = mapper(n, elems[n])
		}
		return rst
	}

	foreach_command := func(elems []commandv1.TemplateCommand, mapper func(int, commandv1.TemplateCommand) error) error {
		for n := range elems {
			if err := mapper(n, elems[n]); err != nil {
				return err
			}
		}
		return nil
	}

	body := new(templatev1.HttpReqTemplate_Create)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}
	if len(body.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.Name", TypeName(body)), body.Name)...,
				)))
	}
	if len(body.Origin) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.Origin", TypeName(body)), body.Origin)...,
				)))
	}
	if len(body.Commands) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.Commands", TypeName(body)), body.Commands)...,
				)))
	}
	//valied commands
	for _, command := range body.Commands {
		if len(command.Name) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
				errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
					logs.KVL(
						ParamLog(fmt.Sprintf("%s.Name", TypeName(command)), command.Name)...,
					)))
		}
		if len(command.Method) == 0 {
			return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
				errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
					logs.KVL(
						ParamLog(fmt.Sprintf("%s.Method", TypeName(command)), command.Method)...,
					)))
		}
	}

	//property
	template := templatev1.Template{}
	template.UuidMeta = metav1.NewUuidMeta()
	template.LabelMeta = metav1.NewLabelMeta(body.Name, body.Summary)
	template.Origin = body.Origin

	//create command
	commands := map_command(body.Commands, func(i int, iter commandv1.HttpReqTemplateCommand_Create_ByTemplate) commandv1.TemplateCommand {
		command := commandv1.TemplateCommand{}
		command.UuidMeta = metav1.NewUuidMeta()
		command.LabelMeta = metav1.NewLabelMeta(iter.Name, iter.Summary)
		command.TemplateUuid = template.Uuid
		command.Sequence = newist.Int32(int32(i))
		command.Method = iter.Method
		command.Args = iter.Args
		command.ResultFilter = &iter.ResultFilter

		return command
	})

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		//create template
		template_, err := vault.NewTemplate(db).Create(template)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "create template"))
		}

		if err := foreach_command(commands, func(i int, tc commandv1.TemplateCommand) error {
			command, err := vault.NewTemplateCommand(db).Create(tc)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
					errors.Wrapf(err, "create template command"))
			}

			commands[i] = *command

			return nil
		}); err != nil {
			return nil, err
		}

		return templatev1.HttpRspTemplate{Template: *template_, Commands: commands}, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description Get a template
// @Accept      json
// @Produce     json
// @Tags        server/template
// @Router      /server/template/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Template Uuid"
// @Success 200 {object} v2.HttpRsp_Template
func (ctl ControlVanilla) GetTemplate(ctx echo.Context) (err error) {
	uuid := echoutil.Param(ctx)[__UUID__]

	err = echoutil.WrapHttpError(http.StatusBadRequest,
		func() (err error) {
			if len(echoutil.Param(ctx)[__UUID__]) == 0 {
				err = ErrorInvalidRequestParameter()
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

	template := templatev2.Template{}
	template.Uuid = uuid

	eq_uuid := vanilla.Equal("uuid", template.Uuid)

	err = vanilla.Stmt.Select(template.TableName(), template.ColumnNames(), eq_uuid.Parse(), nil, nil).
		QueryRow(ctl)(func(scan vanilla.Scanner) error {
		err := template.Scan(scan)
		return errors.Wrapf(err, "failed to scan")
	})
	if err != nil {
		return
	}

	commands := make([]templatev2.TemplateCommand, 0, state.ENV__INIT_SLICE_CAPACITY__())
	command := templatev2.TemplateCommand{}
	command.TemplateUuid = uuid

	eq_uuid = vanilla.Equal("template_uuid", command.TemplateUuid)
	order := vanilla.Asc("sequence")

	err = vanilla.Stmt.Select(command.TableName(), command.ColumnNames(), eq_uuid.Parse(), order.Parse(), nil).
		QueryRows(ctl)(func(scan vanilla.Scanner, _ int) error {
		err := command.Scan(scan)
		if err == nil {
			commands = append(commands, command)
		}
		return errors.Wrapf(err, "failed to scan")
	})
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusOK, templatev2.HttpRsp_Template{Template: template, Commands: commands})
}

// @Description Find []template
// @Accept      json
// @Produce     json
// @Tags        server/template
// @Router      /server/template [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success 200 {array} v2.HttpRsp_Template
func (ctl ControlVanilla) FindTemplate(ctx echo.Context) error {
	q, o, p, err := ParseDecoration(echoutil.QueryParam(ctx))
	err = errors.Wrapf(err, "ParseDecoration%v", logs.KVL(
		"query", echoutil.QueryParamString(ctx),
	))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	rsp := make([]templatev2.HttpRsp_Template, 0, state.ENV__INIT_SLICE_CAPACITY__())

	template := templatev2.Template{}
	err = vanilla.Stmt.Select(template.TableName(), template.ColumnNames(), q, o, p).
		QueryRows(ctl)(func(scan vanilla.Scanner, _ int) error {
		if err = template.Scan(scan); err != nil {
			return errors.Wrapf(err, "failed to scan")
		}

		commands, err := ListTemplateCommand(ctl, template.Uuid)
		rsp = append(rsp, templatev2.HttpRsp_Template{Template: template, Commands: commands})
		return errors.Wrapf(err, "failed to get list")
	})
	if err != nil {
		return errors.Wrapf(err, "failed to find template")
	}

	return ctx.JSON(http.StatusOK, rsp)
}

// @deprecated
// @Description Update a template
// @Accept      json
// @Produce     json
// @Tags        server/template
// @Router      /server/template/{uuid} [put]
// @Param       x_auth_token header string                    false "client session token"
// @Param       uuid         path   string                    true  "Template Uuid"
// @Param       template     body   v1.HttpReqTemplate_Update true  "HttpReqTemplate_Update"
// @Success 200 {object} v1.Template
func (ctl Control) UpdateTemplate(ctx echo.Context) error {
	body := new(templatev1.HttpReqTemplate_Update)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	// template := body.Template
	uuid := echoutil.Param(ctx)[__UUID__]

	//set uuid from path
	template := templatev1.Template{}
	template.Uuid = uuid
	template.LabelMeta = metav1.NewLabelMeta(body.Name, body.Summary)
	template.Origin = body.Origin

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		//upate template
		tempalte_, err := vault.NewTemplate(db).Update(template)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "update template"))
		}

		return tempalte_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// @deprecated
// @Description Delete a template
// @Accept      json
// @Produce     json
// @Tags        server/template
// @Router      /server/template/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Template 의 Uuid"
// @Success 200
func (ctl Control) DeleteTemplate(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}
	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		if err := vault.NewTemplateCommand(db).Delete_ByTemplate(uuid); err != nil {
			return nil, errors.Wrapf(err, "delete template command")
		}
		if err := vault.NewTemplate(db).Delete(uuid); err != nil {
			return nil, errors.Wrapf(err, "delete template")
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, OK())
}

package control

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
	commandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Create Template
// @Description Create a template
// @Accept      json
// @Produce     json
// @Tags        server/template
// @Router      /server/template [post]
// @Param       template body v1.HttpReqTemplateWithCommands true "HttpReqTemplateWithCommands"
// @Success 200 {object} v1.HttpReqTemplateWithCommands
func (ctl Control) CreateTemplate(ctx echo.Context) error {
	map_command := func(elems []commandv1.TemplateCommand, mapper func(commandv1.TemplateCommand) commandv1.TemplateCommand) []commandv1.TemplateCommand {
		rst := make([]commandv1.TemplateCommand, len(elems))
		for n := range elems {
			rst[n] = mapper(elems[n])
		}
		return rst
	}

	body := new(templatev1.HttpReqTemplateWithCommands)
	if err := ctx.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	template := body.Template
	commmands := body.Commands

	//property
	template.UuidMeta = NewUuidMeta()
	template.LabelMeta = NewLabelMeta(template.Name, template.Summary)

	//create command
	seq := int32(0)
	commmands = map_command(commmands, func(tc commandv1.TemplateCommand) commandv1.TemplateCommand {
		//LabelMeta
		tc.UuidMeta = NewUuidMeta()
		tc.LabelMeta = NewLabelMeta(tc.Name, tc.Summary)
		//TemplateUuid
		tc.TemplateUuid = template.Uuid
		//Sequence
		tc.Sequence = newist.Int32(seq)
		seq++

		return tc
	})

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		//create template
		template_, err := vault.NewTemplate(db).
			Create(templatev1.TemplateWithCommands{Template: template, Commands: commmands})
		if err != nil {
			return nil, errors.Wrapf(err, "create template")
		}

		return templatev1.HttpRspTemplateWithCommands{DbSchemaTemplateWithCommands: *template_}, nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, r)
}

// Get Template
// @Description Get a template
// @Accept      json
// @Produce     json
// @Tags        server/template
// @Router      /server/template/{uuid} [get]
// @Param       uuid path string true "Template 의 Uuid"
// @Success 200 {object} v1.HttpReqTemplateWithCommands
func (ctl Control) GetTemplate(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	template, err := vault.NewTemplate(ctl.NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			errors.Wrapf(err, "NewTemplate Get"))
	}

	return ctx.JSON(http.StatusOK, templatev1.HttpRspTemplateWithCommands{DbSchemaTemplateWithCommands: *template})
}

// Find []Template
// @Description Find []template
// @Accept      json
// @Produce     json
// @Tags        server/template
// @Router      /server/template [get]
// @Param       q query string false "query  pkg/server/database/prepared/README.md"
// @Param       o query string false "order  pkg/server/database/prepared/README.md"
// @Param       p query string false "paging pkg/server/database/prepared/README.md"
// @Success 200 {array} v1.HttpReqTemplateWithCommands
func (ctl Control) FindTemplate(ctx echo.Context) error {
	templates, err := vault.NewTemplate(ctl.NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			errors.Wrapf(err, "find template"))
	}

	return ctx.JSON(http.StatusOK, templatev1.TransToHttpRspTemplateWithCommands(templates))
}

// Update Template
// @Description Update a template
// @Accept      json
// @Produce     json
// @Tags        server/template
// @Router      /server/template/{uuid} [put]
// @Param       uuid     path string false "Template 의 Uuid"
// @Param       template body v1.HttpReqTemplate true "HttpReqTemplate"
// @Success 200 {object} v1.HttpRspTemplate
func (ctl Control) UpdateTemplate(ctx echo.Context) error {

	body := new(templatev1.HttpReqTemplate)
	if err := ctx.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	tempalte := body.Template
	uuid := echoutil.Param(ctx)[__UUID__]

	//set uuid from path
	tempalte.Uuid = uuid

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		//upate template
		tempalte_, err := vault.NewTemplate(db).Update(tempalte)
		if err != nil {
			return nil, errors.Wrapf(err, "update template")
		}

		return templatev1.HttpRspTemplate{DbSchema: *tempalte_}, nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, r)
}

// Delete Template
// @Description Delete a template
// @Accept      json
// @Produce     json
// @Tags        server/template
// @Router      /server/template/{uuid} [delete]
// @Param       uuid path string true "Template 의 Uuid"
// @Success 200
func (ctl Control) DeleteTemplate(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}
	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		//template 삭제
		if err := vault.NewTemplate(db).Delete(uuid); err != nil {
			return nil, errors.Wrapf(err, "NewTemplate Delete")
		}

		return nil, nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, OK())
}

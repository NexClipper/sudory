package control

import (
	"github.com/NexClipper/sudory/pkg/server/control/vault"
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
func (c *Control) CreateTemplate() func(ctx echo.Context) error {
	binder := func(ctx Context) error {
		body := new(templatev1.HttpReqTemplateWithCommands)
		if err := ctx.Bind(body); err != nil {
			return ErrorBindRequestObject(err)
		}
		return nil
	}
	operator := func(ctx Context) (interface{}, error) {
		map_command := func(elems []commandv1.TemplateCommand, mapper func(commandv1.TemplateCommand) commandv1.TemplateCommand) []commandv1.TemplateCommand {
			rst := make([]commandv1.TemplateCommand, len(elems))
			for n := range elems {
				rst[n] = mapper(elems[n])
			}
			return rst
		}

		body, ok := ctx.Object().(*templatev1.HttpReqTemplateWithCommands)
		if !ok {
			return nil, ErrorFailedCast()
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

		//create template
		template_, err := vault.NewTemplate(ctx.Database()).
			Create(templatev1.TemplateWithCommands{Template: template, Commands: commmands})
		if err != nil {
			return nil, errors.Wrapf(err, "NewTemplate Create")
		}

		return templatev1.HttpRspTemplateWithCommands{DbSchemaTemplateWithCommands: *template_}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Context) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "CreateTemplate binder")
			}
			return nil
		},
		Operator: func(ctx Context) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "CreateTemplate operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Lock(c.db.Engine()),
	})
}

// Get Template
// @Description Get a template
// @Accept      json
// @Produce     json
// @Tags        server/template
// @Router      /server/template/{uuid} [get]
// @Param       uuid path string true "Template 의 Uuid"
// @Success 200 {object} v1.HttpReqTemplateWithCommands
func (c *Control) GetTemplate() func(ctx echo.Context) error {
	binder := func(ctx Context) error {
		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}

		return nil
	}
	operator := func(ctx Context) (interface{}, error) {
		uuid := ctx.Params()[__UUID__]

		template, err := vault.NewTemplate(ctx.Database()).Get(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "NewTemplate Get")
		}

		return templatev1.HttpRspTemplateWithCommands{DbSchemaTemplateWithCommands: *template}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Context) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "GetTemplate binder")
			}
			return nil
		},
		Operator: func(ctx Context) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "GetTemplate operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Nolock(c.db.Engine()),
	})
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
func (c *Control) FindTemplate() func(ctx echo.Context) error {
	binder := func(ctx Context) error {
		return nil
	}
	operator := func(ctx Context) (interface{}, error) {
		templates, err := vault.NewTemplate(ctx.Database()).Query(ctx.Queries())
		if err != nil {
			return nil, errors.Wrapf(err, "NewTemplate Query")
		}

		templates_ := make([]templatev1.HttpRspTemplateWithCommands, len(templates))
		for n := range templates {
			templates_[n].DbSchemaTemplateWithCommands = templates[n]
		}
		return templates_, nil //pop
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Context) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "FindTemplate binder")
			}
			return nil
		},
		Operator: func(ctx Context) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "FindTemplate operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Nolock(c.db.Engine()),
	})
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
func (c *Control) UpdateTemplate() func(ctx echo.Context) error {

	binder := func(ctx Context) error {
		body := new(templatev1.HttpReqTemplate)
		if err := ctx.Bind(body); err != nil {
			return ErrorBindRequestObject(err)
		}

		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}

		return nil
	}
	operator := func(ctx Context) (interface{}, error) {
		body, ok := ctx.Object().(*templatev1.HttpReqTemplate)
		if !ok {
			return nil, ErrorFailedCast()
		}

		tempalte := body.Template

		uuid := ctx.Params()[__UUID__]

		//set uuid from path
		tempalte.Uuid = uuid

		//upate template
		tempalte_, err := vault.NewTemplate(ctx.Database()).Update(tempalte)
		if err != nil {
			return nil, errors.Wrapf(err, "NewTemplate Update")
		}

		return templatev1.HttpRspTemplate{DbSchema: *tempalte_}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Context) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "UpdateTemplate binder")
			}
			return nil
		},
		Operator: func(ctx Context) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "UpdateTemplate operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Lock(c.db.Engine()),
	})
}

// Delete Template
// @Description Delete a template
// @Accept      json
// @Produce     json
// @Tags        server/template
// @Router      /server/template/{uuid} [delete]
// @Param       uuid path string true "Template 의 Uuid"
// @Success 200
func (c *Control) DeleteTemplate() func(ctx echo.Context) error {

	binder := func(ctx Context) error {
		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}

		return nil
	}
	operator := func(ctx Context) (interface{}, error) {
		uuid := ctx.Params()[__UUID__]

		//template 삭제
		if err := vault.NewTemplate(ctx.Database()).Delete(uuid); err != nil {
			return nil, errors.Wrapf(err, "NewTemplate Delete")
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Context) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "DeleteTemplate binder")
			}
			return nil
		},
		Operator: func(ctx Context) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "DeleteTemplate operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Lock(c.db.Engine()),
	})
}

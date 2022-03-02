package control

import (
	"fmt"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
	commandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
	"github.com/labstack/echo/v4"
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
	binder := func(ctx Contexter) error {
		body := new(templatev1.HttpReqTemplateWithCommands)
		if err := ctx.Bind(body); err != nil {
			return ErrorBindRequestObject(err)
		}
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		body, ok := ctx.Object().(*templatev1.HttpReqTemplateWithCommands)
		if !ok {
			return nil, ErrorFailedCast()
		}

		template := body.Template
		commmands := body.Commands

		//property
		template.UuidMeta = NewUuidMeta()
		template.LabelMeta = NewLabelMeta(template.Name, template.Summary)

		//create template
		err := operator.NewTemplate(ctx.Database()).
			Create(template)
		if err != nil {
			return nil, err
		}
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
		err = foreach_command(commmands, func(command commandv1.TemplateCommand) error {
			if err := operator.NewTemplateCommand(ctx.Database()).
				Create(command); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		return templatev1.HttpReqTemplateWithCommands{Template: template, Commands: commmands}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
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
	binder := func(ctx Contexter) error {
		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}

		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		uuid := ctx.Params()[__UUID__]

		//get template
		template, err := operator.NewTemplate(ctx.Database()).
			Get(uuid)
		if err != nil {
			return nil, err
		}
		//find command
		where := "template_uuid = ?"
		template_uuid := uuid
		commands, err := operator.NewTemplateCommand(ctx.Database()).
			Find(where, template_uuid)
		if err != nil {
			return nil, err
		}

		return templatev1.HttpRspTemplateWithCommands{Template: *template, Commands: commands}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
	})
}

// Find []Template
// @Description Find []template
// @Accept      json
// @Produce     json
// @Tags        server/template
// @Router      /server/template [get]
// @Param       uuid   query string false "Template 의 Uuid"
// @Param       name   query string false "Template 의 Name"
// @Param       origin query string false "Template 의 Origin"
// @Success 200 {array} v1.HttpReqTemplateWithCommands
func (c *Control) FindTemplate() func(ctx echo.Context) error {
	binder := func(ctx Contexter) error {
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		//make condition
		args := make([]interface{}, 0)
		add, build := StringBuilder()

		for key, val := range ctx.Querys() {
			switch key {
			case "uuid":
				args = append(args, fmt.Sprintf("%s%%", val)) //앞 부분 부터 일치 해야함
			default:
				args = append(args, fmt.Sprintf("%%%s%%", val))
			}
			add(fmt.Sprintf("%s LIKE ?", key)) //조건문 만들기
		}
		where := build(" AND ")

		//find template
		templates, err := operator.NewTemplate(ctx.Database()).
			Find(where, args...)
		if err != nil {
			return nil, err
		}
		//make response
		rspadd, rspbuild := templatev1.HttpRspBuilder(len(templates))
		err = foreach_template(templates, func(template templatev1.Template) error {
			template_uuid := template.Uuid
			where := "template_uuid = ?"
			//find commands
			commands, err := operator.NewTemplateCommand(ctx.Database()).
				Find(where, template_uuid)
			if err != nil {
				return err
			}
			rspadd(template, commands) //넣
			return nil
		})
		if err != nil {
			return nil, err
		}
		return rspbuild(), nil //pop
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
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

	binder := func(ctx Contexter) error {
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
	operator := func(ctx Contexter) (interface{}, error) {
		body, ok := ctx.Object().(*templatev1.HttpReqTemplate)
		if !ok {
			return nil, ErrorFailedCast()
		}

		tempalte := body.Template

		uuid := ctx.Params()[__UUID__]

		//set uuid from path
		tempalte.Uuid = uuid

		//upate template
		err := operator.NewTemplate(ctx.Database()).
			Update(tempalte)
		if err != nil {
			return nil, err
		}

		return templatev1.HttpRspTemplate{Template: tempalte}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
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

	binder := func(ctx Contexter) error {
		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}

		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		uuid := ctx.Params()[__UUID__]

		//command 테이블에 데이터 있는 경우 삭제 방지
		where := "template_uuid = ?"
		command, err := operator.NewTemplateCommand(ctx.Database()).Find(where, uuid)
		if err != nil {
			return nil, err
		}
		if len(command) == 0 {
			return nil, fmt.Errorf("commands not empty")
		}

		//template 삭제
		if err := operator.NewTemplate(ctx.Database()).Delete(uuid); err != nil {
			return nil, err
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Lock(c.db.Engine()),
	})
}

func foreach_command(elems []commandv1.TemplateCommand, fn func(commandv1.TemplateCommand) error) error {
	for _, it := range elems {
		if err := fn(it); err != nil {
			return err
		}
	}
	return nil
}

func map_command(elems []commandv1.TemplateCommand, mapper func(commandv1.TemplateCommand) commandv1.TemplateCommand) []commandv1.TemplateCommand {
	rst := make([]commandv1.TemplateCommand, len(elems))
	for n := range elems {
		rst[n] = mapper(elems[n])
	}
	return rst
}

func foreach_template(elems []templatev1.Template, fn func(templatev1.Template) error) error {
	for _, it := range elems {
		if err := fn(it); err != nil {
			return err
		}
	}
	return nil
}

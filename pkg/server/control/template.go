package control

import (
	"fmt"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
	commandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
	"github.com/labstack/echo/v4"
)

// Create Template
// @Description Create a template
// @Accept json
// @Produce json
// @Tags server/template
// @Router /server/template [post]
// @Param template body v1.HttpReqTemplateCreate true "HttpReqTemplateCreate"
// @Success 200
func (c *Control) CreateTemplate() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := new(templatev1.HttpReqTemplateCreate)
		err := ctx.Bind(req)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}
		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(*templatev1.HttpReqTemplateCreate)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//create template
		err := operator.NewTemplate(ctx.Database).
			Create(req.Template)
		if err != nil {
			return nil, err
		}
		//create command
		err = foreach_command(req.Commands, func(command commandv1.TemplateCommand) error {
			if err := operator.NewTemplateCommand(ctx.Database).
				Create(command); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Lock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// Get Template
// @Description Get a template
// @Accept json
// @Produce json
// @Tags server/template
// @Router /server/template/{uuid} [get]
// @Param uuid path string true "Template 의 Uuid"
// @Success 200 {object} v1.HttpRspTemplate
func (c *Control) GetTemplate() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}

		//check request params
		if len(req[__UUID__]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}

		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		//get template
		uuid := req[__UUID__]
		template, err := operator.NewTemplate(ctx.Database).
			Get(uuid)
		if err != nil {
			return nil, err
		}
		//find command
		where := "template_uuid = ?"
		template_uuid := uuid
		commands, err := operator.NewTemplateCommand(ctx.Database).
			Find(where, template_uuid)
		if err != nil {
			return nil, err
		}

		return templatev1.HttpRspTemplate{Template: *template, Commands: commands}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Nolock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// Find []Template
// @Description Find []template
// @Accept json
// @Produce json
// @Tags server/template
// @Router /server/template [get]
// @Param uuid   query string false "Template 의 Uuid"
// @Param name   query string false "Template 의 Name"
// @Param origin query string false "Template 의 Origin"
// @Success 200 {array} v1.HttpRspTemplate
func (c *Control) FindTemplate() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for key := range ctx.QueryParams() {
			req[key] = ctx.QueryParam(key)
		}
		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		//make condition
		args := make([]interface{}, 0)
		add, build := StringBuilder()

		for key, val := range req {
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
		templates, err := operator.NewTemplate(ctx.Database).
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
			commands, err := operator.NewTemplateCommand(ctx.Database).
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
		Operator:      Nolock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// Update Template
// @Description Update a template
// @Accept json
// @Produce json
// @Tags server/template
// @Router /server/template/{uuid} [put]
// @Param uuid     path string false "Template 의 Uuid"
// @Param template body v1.HttpReqTemplate true "HttpReqTemplate"
// @Success 200
func (c *Control) UpdateTemplate() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req[__UUID__].(string)) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		body := new(templatev1.HttpReqTemplate)
		err := ctx.Bind(body)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}
		req[__BODY__] = body
		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}
		uuid, ok := req[__UUID__].(string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		body, ok := req[__BODY__].(*templatev1.HttpReqTemplate)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//set uuid from path
		body.Template.Uuid = uuid

		//upate template
		err := operator.NewTemplate(ctx.Database).
			Update(body.Template)
		if err != nil {
			return nil, err
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Lock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// Delete Template
// @Description Delete a template
// @Accept json
// @Produce json
// @Tags server/template
// @Router /server/template/{uuid} [delete]
// @Param uuid path string true "Template 의 Uuid"
// @Success 200
func (c *Control) DeleteTemplate() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req[__UUID__]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//delete template
		uuid := req[__UUID__]
		err := operator.NewTemplate(ctx.Database).
			Delete(uuid)
		if err != nil {
			return nil, err
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Lock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
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

func foreach_template(elems []templatev1.Template, fn func(templatev1.Template) error) error {
	for _, it := range elems {
		if err := fn(it); err != nil {
			return err
		}
	}
	return nil
}

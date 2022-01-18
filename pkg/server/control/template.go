package control

import (
	"fmt"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
	"github.com/labstack/echo/v4"
)

// Create Template
// @Description Create a template
// @Accept json
// @Produce json
// @Tags server
// @Router /server/template [post]
// @Param template body v1.HttpReqTemplate true "HttpReqTemplate"
// @Success 200
func (c *Control) CreateTemplate() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := new(templatev1.HttpReqTemplate)
		err := ctx.Bind(req)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}
		return req, nil
	}
	operator := func(v interface{}) (interface{}, error) {
		req, ok := v.(*templatev1.HttpReqTemplate)
		if !ok {
			return nil, ErrorFailedCast()
		}

		err := operator.NewTemplate(c.db).
			Create(req.Template)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	return MakeMiddlewareFunc(binder, operator, HttpResponse)
}

// Get Template
// @Description Get a template
// @Accept json
// @Produce json
// @Tags server
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
		if len(req["uuid"]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}

		return req, nil
	}
	operator := func(v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		uuid := req["uuid"]

		record, err := operator.NewTemplate(c.db).
			Get(uuid)
		if err != nil {
			return nil, err
		}
		return templatev1.HttpRspTemplate{Template: *record}, nil
	}

	return MakeMiddlewareFunc(binder, operator, HttpResponse)
}

// Find []Template
// @Description Find []template
// @Accept json
// @Produce json
// @Tags server
// @Router /server/template [get]
// @Param uuid   query string false "Template 의 Uuid"
// @Param name   query string false "Template 의 Name"
// @Param origin query string false "Template 의 Origin"
// @Success 200 {array} v1.HttpRspTemplate
func (c *Control) FindTemplate() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for key, _ := range ctx.QueryParams() {
			req[key] = ctx.QueryParam(key)
		}
		return req, nil
	}
	operator := func(v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		args := make([]interface{}, 0)
		join, build := StringJoin()

		for key, val := range req {
			switch key {
			case "uuid":
				args = append(args, fmt.Sprintf("%s%%", val)) //앞 부분 부터 일치 해야함
			default:
				args = append(args, fmt.Sprintf("%%%s%%", val))
			}
			join(fmt.Sprintf("%s LIKE ?", key)) //조건문 만들기
		}

		// where := "uuid LIKE ? AND name LIKE ? AND origin LIKE ?"
		// uuid := fmt.Sprintf("%s%%", req["uuid"])
		// name := fmt.Sprintf("%%%s%%", req["name"])
		// origin := fmt.Sprintf("%%%s%%", req["origin"])
		where := build(" AND ")

		rst, err := operator.NewTemplate(c.db).
			Find(where, args...)
		if err != nil {
			return nil, err
		}
		return templatev1.TransToHttpRsp(rst), nil
	}

	return MakeMiddlewareFunc(binder, operator, HttpResponse)
}

// Update Template
// @Description Update a template
// @Accept json
// @Produce json
// @Tags server
// @Router /server/template/{uuid} [put]
// @Param template body v1.HttpReqTemplate true "HttpReqTemplate"
// @Success 200
func (c *Control) UpdateTemplate() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req["uuid"].(string)) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		body := new(templatev1.HttpReqTemplate)
		err := ctx.Bind(body)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}
		req["_"] = body
		return req, nil
	}
	operator := func(v interface{}) (interface{}, error) {
		req, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}
		uuid, ok := req["uuid"].(string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		body, ok := req["_"].(*templatev1.HttpReqTemplate)
		if !ok {
			return nil, ErrorFailedCast()
		}

		body.Template.Uuid = uuid

		err := operator.NewTemplate(c.db).
			Update(body.Template)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	return MakeMiddlewareFunc(binder, operator, HttpResponse)
}

// Delete Template
// @Description Delete a template
// @Accept json
// @Produce json
// @Tags server
// @Router /server/template/{uuid} [delete]
// @Param uuid path string true "Template 의 Uuid"
// @Success 200
func (c *Control) DeleteTemplate() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req["uuid"]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		uuid := req["uuid"]

		err := operator.NewTemplate(c.db).
			Delete(uuid)
		if err != nil {
			return nil, err
		}

		where := "template_uuid = ?"
		template_uuid := req["uuid"]

		tcommand, err := operator.NewTemplateCommand(c.db).Find(where, template_uuid)
		if err != nil {
			return nil, err
		}
		for _, it := range tcommand {
			err := operator.NewTemplateCommand(c.db).Delete(it.Uuid)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	}

	return MakeMiddlewareFunc(binder, operator, HttpResponse)
}

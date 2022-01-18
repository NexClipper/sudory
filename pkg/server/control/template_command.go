package control

import (
	"github.com/NexClipper/sudory/pkg/server/control/operator"
	commandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
	"github.com/labstack/echo/v4"
)

// Create Template Command
// @Description Create a template command
// @Accept json
// @Produce json
// @Tags server
// @Router /server/template/{template_uuid}/command [post]
// @Param template_uuid path string true "HttpReqTemplate 의 Uuid"
// @Param command       body v1.HttpReqTemplateCommand true "HttpReqTemplateCommand"
// @Success 200
func (c *Control) CreateTemplateCommand() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req["template_uuid"].(string)) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}

		body := new(commandv1.HttpReqTemplateCommand)
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
		template_uuid, ok := req["template_uuid"].(string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		body, ok := req["_"].(*commandv1.HttpReqTemplateCommand)
		if !ok {
			return nil, ErrorFailedCast()
		}

		body.TemplateCommand.TemplateUuid = template_uuid

		err := operator.NewTemplateCommand(c.db).
			Create(body.TemplateCommand)

		return nil, err
	}

	return MakeMiddlewareFunc(binder, operator, HttpResponse)
}

// Get Template Command
// @Description Get a template command
// @Accept json
// @Produce json
// @Tags server
// @Router /server/template/{template_uuid}/command [get]
// @Param template_uuid path string true "HttpReqTemplate 의 Uuid"
// @Success 200 {array} v1.HttpRspTemplate
func (c *Control) GetTemplateCommands() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req["template_uuid"]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		where := "template_uuid = ?"
		template_uuid := req["template_uuid"]

		rst, err := operator.NewTemplateCommand(c.db).
			Find(where, template_uuid)

		return commandv1.TransToHttpRsp(rst), err
	}

	return MakeMiddlewareFunc(binder, operator, HttpResponse)
}

// Get Template Command
// @Description Get a template command
// @Accept json
// @Produce json
// @Tags server
// @Router /server/template/{template_uuid}/command/{uuid} [get]
// @Param template_uuid path string true "HttpReqTemplate 의 Uuid"
// @Param uuid          path string true "HttpReqTemplateCommand 의 Uuid"
// @Success 200 {object} v1.HttpRspTemplate
func (c *Control) GetTemplateCommand() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req["template_uuid"]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
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

		_ = req["template_uuid"]
		uuid := req["uuid"]
		rst, err := operator.NewTemplateCommand(c.db).
			Get(uuid)
		if err != nil {
			return nil, err
		}
		return commandv1.HttpRspTemplateCommand{TemplateCommand: *rst}, nil
	}

	return MakeMiddlewareFunc(binder, operator, HttpResponse)
}

// Update Template Command
// @Description Update a template command
// @Accept json
// @Produce json
// @Tags server
// @Router /server/template/{template_uuid}/command/{uuid} [put]
// @Param template_uuid path string true "HttpReqTemplateCommand 의 TemplateUuid"
// @Param uuid          path string true "HttpReqTemplateCommand 의 Uuid"
// @Param command       body v1.HttpReqTemplateCommand true "HttpReqTemplateCommand"
// @Success 200
func (c *Control) UpdateTemplateCommand() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req["template_uuid"].(string)) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		if len(req["uuid"].(string)) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}

		body := new(commandv1.HttpReqTemplateCommand)
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
		template_uuid, ok := req["template_uuid"].(string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		uuid, ok := req["uuid"].(string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		body, ok := req["_"].(*commandv1.HttpReqTemplateCommand)
		if !ok {
			return nil, ErrorFailedCast()
		}

		body.TemplateCommand.TemplateUuid = template_uuid
		body.TemplateCommand.Uuid = uuid

		err := operator.NewTemplateCommand(c.db).
			Update(body.TemplateCommand)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	return MakeMiddlewareFunc(binder, operator, HttpResponse)
}

// Delete Template Command
// @Description Delete a template command
// @Accept json
// @Produce json
// @Tags server
// @Router /server/template/{template_uuid}/command/{uuid} [delete]
// @Param template_uuid path string true "HttpReqTemplate 의 Uuid"
// @Param uuid          path string true "HttpReqTemplateCommand 의 Uuid"
// @Success 200
func (c *Control) DeleteTemplateCommand() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req["template_uuid"]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
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

		_ = req["template_uuid"]
		uuid := req["uuid"]
		err := operator.NewTemplateCommand(c.db).
			Delete(uuid)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	return MakeMiddlewareFunc(binder, operator, HttpResponse)
}

package control

import (
	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/database"
	commandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
	"github.com/labstack/echo/v4"
)

// Create Template Command
// @Description Create a template command
// @Accept json
// @Produce json
// @Tags server/template_command
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
		if len(req[__TEMPLATE_UUID__].(string)) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}

		body := new(commandv1.HttpReqTemplateCommand)
		err := ctx.Bind(body)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}
		req[__BODY__] = body
		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}
		template_uuid, ok := req[__TEMPLATE_UUID__].(string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		body, ok := req[__BODY__].(*commandv1.HttpReqTemplateCommand)
		if !ok {
			return nil, ErrorFailedCast()
		}

		body.TemplateCommand.TemplateUuid = template_uuid

		err := operator.NewTemplateCommand(ctx).
			Create(body.TemplateCommand)
		if err != nil {
			return nil, err
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    Lock,
	})
}

// Get Template Commands
// @Description Get template commands
// @Accept json
// @Produce json
// @Tags server/template_command
// @Router /server/template/{template_uuid}/command [get]
// @Param template_uuid path string true "HttpReqTemplate 의 Uuid"
// @Success 200 {array} v1.HttpRspTemplate
func (c *Control) GetTemplateCommands() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req[__TEMPLATE_UUID__]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		where := "template_uuid = ?"
		template_uuid := req[__TEMPLATE_UUID__]

		rst, err := operator.NewTemplateCommand(ctx).
			Find(where, template_uuid)

		return commandv1.TransToHttpRsp(rst), err
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    NoLock,
	})
}

// Get Template Command
// @Description Get a template command
// @Accept json
// @Produce json
// @Tags server/template_command
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
		if len(req[__TEMPLATE_UUID__]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		if len(req[__UUID__]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		_ = req[__TEMPLATE_UUID__]
		uuid := req[__UUID__]
		rst, err := operator.NewTemplateCommand(ctx).
			Get(uuid)
		if err != nil {
			return nil, err
		}
		return commandv1.HttpRspTemplateCommand{TemplateCommand: *rst}, nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    NoLock,
	})
}

// Update Template Command
// @Description Update a template command
// @Accept json
// @Produce json
// @Tags server/template_command
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
		if len(req[__TEMPLATE_UUID__].(string)) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		if len(req[__UUID__].(string)) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}

		body := new(commandv1.HttpReqTemplateCommand)
		err := ctx.Bind(body)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}
		req[__BODY__] = body

		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}
		template_uuid, ok := req[__TEMPLATE_UUID__].(string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		uuid, ok := req[__UUID__].(string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		body, ok := req[__BODY__].(*commandv1.HttpReqTemplateCommand)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//set template uuid from path
		body.TemplateCommand.TemplateUuid = template_uuid
		//set uuid from path
		body.TemplateCommand.Uuid = uuid

		err := operator.NewTemplateCommand(ctx).
			Update(body.TemplateCommand)
		if err != nil {
			return nil, err
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    Lock,
	})
}

// Delete Template Command
// @Description Delete a template command
// @Accept json
// @Produce json
// @Tags server/template_command
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
		if len(req[__TEMPLATE_UUID__]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		if len(req[__UUID__]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		_ = req[__TEMPLATE_UUID__]
		uuid := req[__UUID__]
		err := operator.NewTemplateCommand(ctx).
			Delete(uuid)
		if err != nil {
			return nil, err
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    Lock,
	})
}

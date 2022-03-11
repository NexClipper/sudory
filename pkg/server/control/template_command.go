package control

import (
	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	commandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Create Template Command
// @Description Create a template command
// @Accept      json
// @Produce     json
// @Tags        server/template_command
// @Router      /server/template/{template_uuid}/command [post]
// @Param       template_uuid path string true "HttpReqTemplate 의 Uuid"
// @Param       command       body v1.HttpReqTemplateCommand true "HttpReqTemplateCommand"
// @Success 200 {object} v1.HttpRspTemplateCommand
func (c *Control) CreateTemplateCommand() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {
		body := new(commandv1.HttpReqTemplateCommand)
		if err := ctx.Bind(body); err != nil {
			return ErrorBindRequestObject(err)
		}
		if body.Name == nil {
			return ErrorInvaliedRequestParameterName("Name")
		}
		if body.Method == nil {
			return ErrorInvaliedRequestParameterName("Method")
		}

		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__TEMPLATE_UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__TEMPLATE_UUID__)
		}

		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		body, ok := ctx.Object().(*commandv1.HttpReqTemplateCommand)
		if !ok {
			return nil, ErrorFailedCast()
		}
		command := body.TemplateCommand

		template_uuid := ctx.Params()[__TEMPLATE_UUID__]

		//property
		command.UuidMeta = NewUuidMeta()
		command.LabelMeta = NewLabelMeta(command.Name, command.Summary)
		command.TemplateUuid = template_uuid
		if command.Sequence == nil {
			//마직막 순서를 지정하기 위해서 스텝을 가져온다
			where := "template_uuid = ?"
			commands, err := vault.NewTemplateCommand(ctx.Database()).
				Find(where, template_uuid)
			if err != nil {
				return nil, errors.Wrapf(err, "NewTemplateCommand Find")
			}
			//스탭 순서 지정
			command.Sequence = newist.Int32(int32(len(commands)))
		}
		command_, err := vault.NewTemplateCommand(ctx.Database()).Create(command)
		if err != nil {
			return nil, errors.Wrapf(err, "NewTemplateCommand Create")
		}

		//ChainingSequence
		if err := vault.NewTemplateCommand(ctx.Database()).ChainingSequence(template_uuid, command.Uuid); err != nil {
			return nil, errors.Wrapf(err, "NewTemplateCommand ChainingSequence")
		}

		return commandv1.HttpRspTemplateCommand{DbSchema: *command_}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Contexter) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "CreateTemplateCommand binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "CreateTemplateCommand operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Lock(c.db.Engine()),
	})
}

// FindTemplateCommand
// @Description Find template command
// @Accept      json
// @Produce     json
// @Tags        server/template_command
// @Router      /server/template/{template_uuid}/command [get]
// @Param       template_uuid path string true "HttpReqTemplate 의 Uuid"
// @Success 200 {array} v1.HttpRspTemplateCommand
func (c *Control) FindTemplateCommand() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {
		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__TEMPLATE_UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__TEMPLATE_UUID__)
		}

		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		where := "template_uuid = ?"
		template_uuid := ctx.Params()[__TEMPLATE_UUID__]

		rst, err := vault.NewTemplateCommand(ctx.Database()).
			Find(where, template_uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "NewTemplateCommand Find")
		}
		return commandv1.TransToHttpRsp(rst), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Contexter) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "FindTemplateCommand binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "FindTemplateCommand operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Nolock(c.db.Engine()),
	})
}

// Get Template Command
// @Description Get a template command
// @Accept      json
// @Produce     json
// @Tags        server/template_command
// @Router      /server/template/{template_uuid}/command/{uuid} [get]
// @Param       template_uuid path string true "HttpReqTemplate 의 Uuid"
// @Param       uuid          path string true "HttpReqTemplateCommand 의 Uuid"
// @Success 200 {object} v1.HttpRspTemplateCommand
func (c *Control) GetTemplateCommand() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {
		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__TEMPLATE_UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__TEMPLATE_UUID__)
		}
		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}

		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		_ = ctx.Params()[__TEMPLATE_UUID__]
		uuid := ctx.Params()[__UUID__]

		rst, err := vault.NewTemplateCommand(ctx.Database()).Get(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "NewTemplateCommand Get")
		}
		return commandv1.HttpRspTemplateCommand{DbSchema: *rst}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Contexter) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "GetTemplateCommand binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "GetTemplateCommand operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Nolock(c.db.Engine()),
	})
}

// Update Template Command
// @Description Update a template command
// @Accept      json
// @Produce     json
// @Tags        server/template_command
// @Router      /server/template/{template_uuid}/command/{uuid} [put]
// @Param       template_uuid path string true "HttpReqTemplateCommand 의 TemplateUuid"
// @Param       uuid          path string true "HttpReqTemplateCommand 의 Uuid"
// @Param       command       body v1.HttpReqTemplateCommand true "HttpReqTemplateCommand"
// @Success 200 {object} v1.HttpRspTemplateCommand
func (c *Control) UpdateTemplateCommand() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {
		body := new(commandv1.HttpReqTemplateCommand)
		if err := ctx.Bind(body); err != nil {
			return ErrorBindRequestObject(err)
		}

		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__TEMPLATE_UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__TEMPLATE_UUID__)
		}
		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}

		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		body, ok := ctx.Object().(*commandv1.HttpReqTemplateCommand)
		if !ok {
			return nil, ErrorFailedCast()
		}

		template_uuid := ctx.Params()[__TEMPLATE_UUID__]
		uuid := ctx.Params()[__UUID__]

		command := body.TemplateCommand

		//set template uuid from path
		command.TemplateUuid = template_uuid
		//set uuid from path
		command.Uuid = uuid

		command_, err := vault.NewTemplateCommand(ctx.Database()).
			Update(body.TemplateCommand)
		if err != nil {
			return nil, errors.Wrapf(err, "NewTemplateCommand Update")
		}

		//ChainingSequence
		if err := vault.NewTemplateCommand(ctx.Database()).ChainingSequence(template_uuid, command.Uuid); err != nil {
			return nil, errors.Wrapf(err, "NewTemplateCommand ChainingSequence")
		}

		return commandv1.HttpRspTemplateCommand{DbSchema: *command_}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Contexter) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "UpdateTemplateCommand binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "UpdateTemplateCommand operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Lock(c.db.Engine()),
	})
}

// Delete Template Command
// @Description Delete a template command
// @Accept      json
// @Produce     json
// @Tags        server/template_command
// @Router      /server/template/{template_uuid}/command/{uuid} [delete]
// @Param       template_uuid path string true "HttpReqTemplate 의 Uuid"
// @Param       uuid          path string true "HttpReqTemplateCommand 의 Uuid"
// @Success 200
func (c *Control) DeleteTemplateCommand() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {
		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__TEMPLATE_UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__TEMPLATE_UUID__)
		}
		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}

		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		_ = ctx.Params()[__TEMPLATE_UUID__]
		uuid := ctx.Params()[__UUID__]

		if err := vault.NewTemplateCommand(ctx.Database()).Delete(uuid); err != nil {
			return nil, errors.Wrapf(err, "NewTemplateCommand Delete")
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Contexter) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "DeleteTemplateCommand binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "DeleteTemplateCommand operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Lock(c.db.Engine()),
	})
}

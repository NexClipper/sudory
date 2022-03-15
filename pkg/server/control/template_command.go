package control

import (
	"sort"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
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
func (c Control) CreateTemplateCommand() func(ctx echo.Context) error {

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
		chaining_sequence := func(command commandv1.TemplateCommand) error {
			where := "template_uuid = ? AND uuid <> ?"
			args := []interface{}{
				command.TemplateUuid,
				command.Uuid,
			}
			commands, err := vault.NewTemplateCommand(ctx.Database()).Find(where, args...)
			if err != nil {
				return errors.Wrapf(err, "NewTemplateCommand Find")
			}

			//sort -> Sequence ASC
			sort.Slice(commands, func(i, j int) bool {
				return nullable.Int32(commands[i].Sequence).Value() < nullable.Int32(commands[j].Sequence).Value()
			})

			commands_ := make([]commandv1.DbSchema, 0, len(commands)+1)
			for i := range commands {
				itor := commands[i]

				//순서 중간에 비교해서 적용하려는 command 우선으로 적용
				sequence := nullable.Int32(command.Sequence)
				if sequence.Has() && int32(i) == sequence.Value() {
					commands_ = append(commands_, commandv1.DbSchema{TemplateCommand: command})
				}
				commands_ = append(commands_, itor)
			}
			//마지막에 비교해서 빠져있으면 넣는다
			//command.Sequence 중간에 껴넣는게 아니라면 마지막에 위치 시킨다
			if len(commands) == len(commands_) {
				commands_ = append(commands_, commandv1.DbSchema{TemplateCommand: command})
			}

		LOOP:
			for i := range commands_ {
				itor := commands_[i]

				//Sequence 동일함
				sequence := nullable.Int32(command.Sequence)
				if sequence.Has() && int32(i) == sequence.Value() {
					continue LOOP //change nothing
				}

				itor.Sequence = newist.Int32(int32(i))
				if _, err := vault.NewTemplateCommand(ctx.Database()).Update(itor.TemplateCommand); err != nil {
					return errors.Wrapf(err, "NewTemplateCommand Update")
				}
			}
			return nil
		}

		body, ok := ctx.Object().(*commandv1.HttpReqTemplateCommand)
		if !ok {
			return nil, ErrorFailedCast()
		}
		command := body.TemplateCommand

		template_uuid := ctx.Params()[__TEMPLATE_UUID__]

		if _, err := vault.NewTemplate(ctx.Database()).Get(template_uuid); err != nil {
			return nil, errors.Wrapf(err, "NewTemplate Get")
		}
		//property
		command.UuidMeta = NewUuidMeta()
		command.LabelMeta = NewLabelMeta(command.Name, command.Summary)
		command.TemplateUuid = template_uuid

		command_, err := vault.NewTemplateCommand(ctx.Database()).Create(command)
		if err != nil {
			return nil, errors.Wrapf(err, "NewTemplateCommand Create")
		}

		// Chaining
		if err := chaining_sequence(command); err != nil {
			return nil, errors.Wrapf(err, "chaining sequence")
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
func (c Control) FindTemplateCommand() func(ctx echo.Context) error {

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
func (c Control) GetTemplateCommand() func(ctx echo.Context) error {

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
func (c Control) UpdateTemplateCommand() func(ctx echo.Context) error {
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
		chaining_sequence := func(command commandv1.TemplateCommand) error {
			//valid
			if command.Sequence == nil {
				return nil //Sequence 값이 있어야 처리
			}

			where := "template_uuid = ? AND uuid <> ?"
			args := []interface{}{
				command.TemplateUuid,
				command.Uuid,
			}
			commands, err := vault.NewTemplateCommand(ctx.Database()).Find(where, args...)
			if err != nil {
				return errors.Wrapf(err, "NewTemplateCommand Find")
			}

			//sort -> Sequence ASC
			sort.Slice(commands, func(i, j int) bool {
				return nullable.Int32(commands[i].Sequence).Value() < nullable.Int32(commands[j].Sequence).Value()
			})

			commands_ := make([]commandv1.DbSchema, 0, len(commands)+1)
			for i := range commands {
				itor := commands[i]

				//순서 중간에 비교해서 적용하려는 command 우선으로 적용
				sequence := nullable.Int32(command.Sequence)
				if sequence.Has() && int32(i) == sequence.Value() {
					commands_ = append(commands_, commandv1.DbSchema{TemplateCommand: command})
				}
				commands_ = append(commands_, itor)
			}
			//마지막에 비교해서 빠져있으면 넣는다
			//command.Sequence 중간에 껴넣는게 아니라면 마지막에 위치 시킨다
			if len(commands) == len(commands_) {
				commands_ = append(commands_, commandv1.DbSchema{TemplateCommand: command})
			}

		LOOP:
			for i := range commands_ {
				itor := commands_[i]

				//Sequence 동일함
				sequence := nullable.Int32(command.Sequence)
				if sequence.Has() && int32(i) == sequence.Value() {
					continue LOOP //change nothing
				}

				itor.Sequence = newist.Int32(int32(i))
				if _, err := vault.NewTemplateCommand(ctx.Database()).Update(itor.TemplateCommand); err != nil {
					return errors.Wrapf(err, "NewTemplateCommand Update")
				}
			}
			return nil
		}

		body, ok := ctx.Object().(*commandv1.HttpReqTemplateCommand)
		if !ok {
			return nil, ErrorFailedCast()
		}

		template_uuid := ctx.Params()[__TEMPLATE_UUID__]
		uuid := ctx.Params()[__UUID__]
		command := body.TemplateCommand

		command.TemplateUuid = template_uuid //set template uuid from path
		command.Uuid = uuid                  //set uuid from path

		command_, err := vault.NewTemplateCommand(ctx.Database()).
			Update(command)
		if err != nil && !macro.Eqaul(database.ErrorNoAffected(), errors.Cause(err)) {
			return nil, errors.Wrapf(err, "NewTemplateCommand Update")
		}

		if err := chaining_sequence(command_.TemplateCommand); err != nil {
			return nil, errors.Wrapf(err, "chaining sequence")
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
func (c Control) DeleteTemplateCommand() func(ctx echo.Context) error {

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
		chaining_sequence := func(template_uuid, uuid string) error {
			vault_command := vault.NewTemplateCommand(ctx.Database())
			where := "template_uuid = ?"
			args := []interface{}{
				template_uuid,
			}
			commands, err := vault_command.Find(where, args...)
			if err != nil {
				return errors.Wrapf(err, "NewTemplateCommand Find")
			}

			//sort -> Sequence ASC
			sort.Slice(commands, func(i, j int) bool {
				return nullable.Int32(commands[i].Sequence).Value() < nullable.Int32(commands[j].Sequence).Value()
			})

		LOOP:
			for i := range commands {
				command := commands[i]

				//Sequence 동일함
				sequence := nullable.Int32(command.Sequence)
				if sequence.Has() && sequence.Value() == int32(i) {
					continue LOOP //change nothing
				}

				command.Sequence = newist.Int32(int32(i))
				if _, err := vault_command.Update(command.TemplateCommand); err != nil {
					return errors.Wrapf(err, "NewTemplateCommand Update")
				}
			}
			return nil
		}

		template_uuid := ctx.Params()[__TEMPLATE_UUID__]
		uuid := ctx.Params()[__UUID__]

		if err := vault.NewTemplateCommand(ctx.Database()).Delete(uuid); err != nil {
			return nil, errors.Wrapf(err, "NewTemplateCommand Delete")
		}

		if err := chaining_sequence(template_uuid, uuid); err != nil {
			return nil, errors.Wrapf(err, "chaining sequence")
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

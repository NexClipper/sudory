package control

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	templatev2 "github.com/NexClipper/sudory/pkg/server/model/template/v2"
	commandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @deprecated
// @Description Create a template command
// @Accept      json
// @Produce     json
// @Tags        server/template_command
// @Router      /server/template/{template_uuid}/command [post]
// @Param       x_auth_token  header string                           false "client session token"
// @Param       template_uuid path   string                           true  "HttpReqTemplate Uuid"
// @Param       command       body   v1.HttpReqTemplateCommand_Create true  "HttpReqTemplateCommand_Create"
// @Success 200 {object} v1.TemplateCommand
func (ctl Control) CreateTemplateCommand(ctx echo.Context) error {
	chaining_sequence := func(db database.Context, command commandv1.TemplateCommand) error {
		where := "template_uuid = ? AND uuid <> ?"
		args := []interface{}{
			command.TemplateUuid,
			command.Uuid,
		}
		commands, err := vault.NewTemplateCommand(db).Find(where, args...)
		if err != nil {
			return errors.Wrapf(err, "find template command")
		}

		//sort -> Sequence ASC
		sort.Slice(commands, func(i, j int) bool {
			return nullable.Int32(commands[i].Sequence).Value() < nullable.Int32(commands[j].Sequence).Value()
		})

		commands_ := make([]commandv1.TemplateCommand, 0, len(commands)+1)
		for i := range commands {
			itor := commands[i]

			//순서 중간에 비교해서 적용하려는 command 우선으로 적용
			sequence := nullable.Int32(command.Sequence)
			if sequence.Has() && int32(i) == sequence.Value() {
				commands_ = append(commands_, command)
			}
			commands_ = append(commands_, itor)
		}
		//마지막에 비교해서 빠져있으면 넣는다
		//command.Sequence 중간에 껴넣는게 아니라면 마지막에 위치 시킨다
		if len(commands) == len(commands_) {
			commands_ = append(commands_, command)
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
			if _, err := vault.NewTemplateCommand(db).Update(itor); err != nil {
				return errors.Wrapf(err, "update template command")
			}
		}
		return nil
	}

	body := new(commandv1.HttpReqTemplateCommand_Create)
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
	if len(body.Method) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.Method", TypeName(body)), body.Method)...,
				)))
	}
	if len(echoutil.Param(ctx)[__TEMPLATE_UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__TEMPLATE_UUID__, echoutil.Param(ctx)[__TEMPLATE_UUID__])...,
				)))
	}

	// command := body.TemplateCommand
	template_uuid := echoutil.Param(ctx)[__TEMPLATE_UUID__]

	//vailed template
	_, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		if _, err := vault.NewTemplate(db).Get(template_uuid); err != nil {
			return nil, errors.Wrapf(err, "valid%s",
				logs.KVL(
					"template", template_uuid,
				))
		}

		return nil, nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	//property
	command := commandv1.TemplateCommand{}
	command.UuidMeta = metav1.NewUuidMeta()
	command.LabelMeta = metav1.NewLabelMeta(body.Name, body.Summary)
	command.TemplateUuid = template_uuid
	command.Sequence = &body.Sequence
	command.Method = body.Method
	command.Args = body.Args
	command.ResultFilter = &body.ResultFilter

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		command_, err := vault.NewTemplateCommand(db).Create(command)
		if err != nil {
			return nil, errors.Wrapf(err, "create template command")
		}

		// Chaining
		if err := chaining_sequence(db, command); err != nil {
			return nil, errors.Wrapf(err, "chaining sequence")
		}

		return command_, nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return ctx.JSON(http.StatusOK, r)
}

// @Description List template command
// @Accept      json
// @Produce     json
// @Tags        server/template_command
// @Router      /server/template/{template_uuid}/command [get]
// @Param       x_auth_token  header string false "client session token"
// @Param       template_uuid path   string true  "HttpReqTemplate Uuid"
// @Success 200 {array} v2.HttpRsp_TemplateCommand
func (ctl ControlVanilla) ListTemplateCommand(ctx echo.Context) (err error) {
	template_uuid := echoutil.Param(ctx)[__TEMPLATE_UUID__]

	err = echoutil.WrapHttpError(http.StatusBadRequest,
		func() (err error) {
			if len(echoutil.Param(ctx)[__TEMPLATE_UUID__]) == 0 {
				err = ErrorInvalidRequestParameter()
			}
			return errors.Wrapf(err, "valid param%s",
				logs.KVL(
					ParamLog(__TEMPLATE_UUID__, echoutil.Param(ctx)[__TEMPLATE_UUID__])...,
				))
		},
	)
	if err != nil {
		return
	}

	commands, err := ListTemplateCommand(ctl, template_uuid)
	if err != nil {
		return
	}

	rsp := make([]templatev2.HttpRsp_TemplateCommand, len(commands))
	for i := range commands {
		rsp[i] = templatev2.HttpRsp_TemplateCommand(commands[i])
	}

	return ctx.JSON(http.StatusOK, rsp)

}

func ListTemplateCommand(ctl ControlVanilla, template_uuid string) ([]templatev2.TemplateCommand, error) {
	rsp := make([]templatev2.TemplateCommand, 0, state.ENV__INIT_SLICE_CAPACITY__())

	command := templatev2.TemplateCommand{}
	command.TemplateUuid = template_uuid

	eq_uuid := vanilla.Equal("template_uuid", command.TemplateUuid)
	order := vanilla.Asc("sequence")

	err := vanilla.Stmt.Select(command.TableName(), command.ColumnNames(), eq_uuid.Parse(), order.Parse(), nil).
		QueryRows(ctl)(func(scan vanilla.Scanner, _ int) error {
		err := command.Scan(scan)
		rsp = append(rsp, command)
		return errors.Wrapf(err, "failed to scan")
	})

	return rsp, err
}

// @Description Get a template command
// @Accept      json
// @Produce     json
// @Tags        server/template_command
// @Router      /server/template/{template_uuid}/command/{uuid} [get]
// @Param       x_auth_token  header string false "client session token"
// @Param       template_uuid path   string true  "HttpReqTemplate Uuid"
// @Param       uuid          path   string true  "HttpReqTemplateCommand 의 Uuid"
// @Success 200 {object} v2.HttpRsp_TemplateCommand
func (ctl ControlVanilla) GetTemplateCommand(ctx echo.Context) (err error) {
	_ = echoutil.Param(ctx)[__TEMPLATE_UUID__]
	uuid := echoutil.Param(ctx)[__UUID__]

	err = echoutil.WrapHttpError(http.StatusBadRequest,
		func() (err error) {
			if len(echoutil.Param(ctx)[__TEMPLATE_UUID__]) == 0 {
				err = ErrorInvalidRequestParameter()
			}
			return errors.Wrapf(err, "valid param%s",
				logs.KVL(
					ParamLog(__TEMPLATE_UUID__, echoutil.Param(ctx)[__TEMPLATE_UUID__])...,
				))
		},
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

	command := templatev2.TemplateCommand{}
	command.Uuid = uuid

	eq_uuid := vanilla.Equal("uuid", command.Uuid)

	err = vanilla.Stmt.Select(command.TableName(), command.ColumnNames(), eq_uuid.Parse(), nil, nil).
		QueryRow(ctl)(func(scan vanilla.Scanner) error {
		err := command.Scan(scan)
		return errors.Wrapf(err, "failed to scan")
	})
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusOK, command)
}

// @deprecated
// @Description Update a template command
// @Accept      json
// @Produce     json
// @Tags        server/template_command
// @Router      /server/template/{template_uuid}/command/{uuid} [put]
// @Param       x_auth_token  header string                           false "client session token"
// @Param       template_uuid path   string                           true  "HttpReqTemplateCommand TemplateUuid"
// @Param       uuid          path   string                           true  "HttpReqTemplateCommand Uuid"
// @Param       command       body   v1.HttpReqTemplateCommand_Update true  "HttpReqTemplateCommand_Update"
// @Success 200 {object} v1.TemplateCommand
func (ctl Control) UpdateTemplateCommand(ctx echo.Context) error {
	chaining_sequence := func(db database.Context, command commandv1.TemplateCommand) error {
		//valid
		if command.Sequence == nil {
			return nil //Sequence 값이 있어야 처리
		}

		where := "template_uuid = ? AND uuid <> ?"
		args := []interface{}{
			command.TemplateUuid,
			command.Uuid,
		}
		commands, err := vault.NewTemplateCommand(db).Find(where, args...)
		if err != nil {
			return errors.Wrapf(err, "find template command")
		}

		//sort -> Sequence ASC
		sort.Slice(commands, func(i, j int) bool {
			return nullable.Int32(commands[i].Sequence).Value() < nullable.Int32(commands[j].Sequence).Value()
		})

		commands_ := make([]commandv1.TemplateCommand, 0, len(commands)+1)
		for i := range commands {
			itor := commands[i]

			//순서 중간에 비교해서 적용하려는 command 우선으로 적용
			sequence := nullable.Int32(command.Sequence)
			if sequence.Has() && int32(i) == sequence.Value() {
				commands_ = append(commands_, command)
			}
			commands_ = append(commands_, itor)
		}
		//마지막에 비교해서 빠져있으면 넣는다
		//command.Sequence 중간에 껴넣는게 아니라면 마지막에 위치 시킨다
		if len(commands) == len(commands_) {
			commands_ = append(commands_, command)
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
			if _, err := vault.NewTemplateCommand(db).Update(itor); err != nil {
				return errors.Wrapf(err, "update template command")
			}
		}
		return nil
	}

	body := new(commandv1.HttpReqTemplateCommand_Update)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if len(echoutil.Param(ctx)[__TEMPLATE_UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__TEMPLATE_UUID__, echoutil.Param(ctx)[__TEMPLATE_UUID__])...,
				)))
	}
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	template_uuid := echoutil.Param(ctx)[__TEMPLATE_UUID__]
	uuid := echoutil.Param(ctx)[__UUID__]
	// command := body.TemplateCommand

	//property
	command := commandv1.TemplateCommand{}
	command.Uuid = uuid                  //set uuid from path
	command.TemplateUuid = template_uuid //set template uuid from path

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		command_, err := vault.NewTemplateCommand(db).Update(command)
		if err != nil && !macro.Eqaul(database.ErrorNoAffected(), errors.Cause(err)) {
			return nil, errors.Wrapf(err, "update template command")
		}

		if err := chaining_sequence(db, *command_); err != nil {
			return nil, errors.Wrapf(err, "chaining sequence")
		}

		return command_, nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return ctx.JSON(http.StatusOK, r)
}

// @deprecated
// @Description Delete a template command
// @Accept      json
// @Produce     json
// @Tags        server/template_command
// @Router      /server/template/{template_uuid}/command/{uuid} [delete]
// @Param       x_auth_token  header string false "client session token"
// @Param       template_uuid path   string true  "HttpReqTemplate Uuid"
// @Param       uuid          path   string true  "HttpReqTemplateCommand 의 Uuid"
// @Success 200
func (ctl Control) DeleteTemplateCommand(ctx echo.Context) error {
	chaining_sequence := func(db database.Context, template_uuid, uuid string) error {
		vault_command := vault.NewTemplateCommand(db)
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
			if _, err := vault_command.Update(command); err != nil {
				return errors.Wrapf(err, "NewTemplateCommand Update")
			}
		}
		return nil
	}

	if len(echoutil.Param(ctx)[__TEMPLATE_UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__TEMPLATE_UUID__, echoutil.Param(ctx)[__TEMPLATE_UUID__])...,
				)))
	}
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	template_uuid := echoutil.Param(ctx)[__TEMPLATE_UUID__]
	uuid := echoutil.Param(ctx)[__UUID__]
	_, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		if err := vault.NewTemplateCommand(db).Delete(uuid); err != nil {
			return nil, errors.Wrapf(err, "delete template command")
		}

		if err := chaining_sequence(db, template_uuid, uuid); err != nil {
			return nil, errors.Wrapf(err, "chaining sequence")
		}

		return nil, nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return ctx.JSON(http.StatusOK, OK())
}

package operator

import (
	"github.com/NexClipper/sudory/pkg/server/database"

	tcommandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
	"github.com/labstack/echo/v4"
)

type CreateTemplateCommand struct {
	OperateContext
	tcommandv1.HttpReqTemplateCommand
}

var _ Creator = (*CreateTemplateCommand)(nil)

func NewCreateTemplateCommand(d *database.DBManipulator) Creator {
	return &CreateTemplateCommand{OperateContext: OperateContext{Db: d}}
}

func (o *CreateTemplateCommand) toModel() tcommandv1.TemplateCommand {
	return o.HttpReqTemplateCommand.TemplateCommand
}

func (o *CreateTemplateCommand) Create(ctx echo.Context) error {
	model := o.toModel()

	_, err := o.Db.CreateTemplateCommand(tcommandv1.DbSchemaTemplateCommand{TemplateCommand: model})
	if err != nil {
		return err
	}

	if o.Response != nil {
		o.Response(ctx, nil)
	}

	return nil
}

type GetTemplateCommand KeyValueParam

func NewGetTemplateCommand(d *database.DBManipulator) Getter {
	return (*GetTemplateCommand)(NewKeyValueParam(OperateContext{
		Db: d,
	}))
}

func (o *GetTemplateCommand) toModel() map[string]string {
	return o.Params
}

func (o *GetTemplateCommand) Get(ctx echo.Context) error {
	model := o.toModel()

	r, err := o.Db.GetTemplateCommand(model["uuid"])
	if err != nil {
		return err
	}

	if o.Response != nil {
		//convert http response object
		m := tcommandv1.HttpRspTemplateCommand{TemplateCommand: r.TemplateCommand}
		o.Response(ctx, &m)
	}

	return nil
}

type UpdateTemplateCommand struct {
	OperateContext
	tcommandv1.HttpReqTemplateCommand
}

var _ Updater = (*UpdateTemplateCommand)(nil)

func NewUpdateTemplateCommand(d *database.DBManipulator) Updater {
	return &UpdateTemplateCommand{OperateContext: OperateContext{Db: d}}
}

func (o *UpdateTemplateCommand) toModel() *tcommandv1.TemplateCommand {
	return &o.TemplateCommand
}

func (o *UpdateTemplateCommand) Update(ctx echo.Context) error {
	model := o.toModel()

	m := tcommandv1.DbSchemaTemplateCommand{TemplateCommand: *model}
	_, err := o.Db.UpdateTemplateCommand(m)
	if err != nil {
		return err
	}

	if o.Response != nil {
		o.Response(ctx, nil)
	}

	return nil
}

type DeleteTemplateCommand KeyValueParam

func NewDeleteTemplateCommand(d *database.DBManipulator) Remover {
	return (*DeleteTemplateCommand)(NewKeyValueParam(OperateContext{
		Db: d,
	}))
}

func (o *DeleteTemplateCommand) toModel() map[string]string {
	return o.Params
}

func (o *DeleteTemplateCommand) Delete(ctx echo.Context) error {
	model := o.toModel()

	_, err := o.Db.DeleteTemplateCommand(model["uuid"])
	if err != nil {
		return err
	}

	if o.Response != nil {
		o.Response(ctx, nil)
	}

	return nil
}

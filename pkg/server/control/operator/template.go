package operator

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	. "github.com/NexClipper/sudory/pkg/server/macro"

	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
	"github.com/labstack/echo/v4"
)

type CreateTemplate struct {
	OperateContext
	templatev1.HttpReqTemplates
}

var _ Creator = (*CreateTemplate)(nil)

func NewCreateTemplate(d *database.DBManipulator) Creator {
	return &CreateTemplate{OperateContext: OperateContext{Db: d}}
}

func (o *CreateTemplate) toModel() []templatev1.Template {
	return templatev1.TransFormHttpReqTemplate([]templatev1.HttpReqTemplate(o.HttpReqTemplates))
}

func (o *CreateTemplate) Create(ctx echo.Context) error {
	model := o.toModel()

	vaild := func(m *templatev1.Template) {
		if len(m.Uuid) == 0 {
			m.Uuid = UuidNewString()
		}
	}

	for n := range model {
		vaild(&model[n])
	}

	m := templatev1.TransToDbSchema(model)

	_, err := o.Db.CreateTemplate(m)
	if err != nil {
		return err
	}

	if o.Response != nil {
		o.Response(ctx, nil)
	}

	return nil
}

type GetTemplate KeyValueParam

func NewGetTemplate(d *database.DBManipulator) Getter {
	return (*GetTemplate)(NewKeyValueParam(OperateContext{
		Db: d,
	}))
}

func (o *GetTemplate) toModel() map[string]string {
	return o.Params
}

func (o *GetTemplate) Get(ctx echo.Context) error {
	model := o.toModel()

	r, err := o.Db.GetTemplate(model)
	if err != nil {
		return err
	}

	if o.Response != nil {
		//convert http response object
		m := templatev1.HttpRspTemplate{Template: r.Template}
		o.Response(ctx, &m)
	}

	return nil
}

type SearchTemplate KeyValueParam

var _ Getter = (*SearchTemplate)(nil)

func NewSearchTemplate(d *database.DBManipulator) Getter {
	return (*SearchTemplate)(NewKeyValueParam(OperateContext{
		Db: d,
	}))
}

func (o *SearchTemplate) toModel() map[string]string {
	return o.Params
}

func (o *SearchTemplate) Get(ctx echo.Context) error {
	param := o.toModel()

	r, err := o.Db.SearchTemplate(param)
	if err != nil {
		return err
	}

	if o.Response != nil {
		//convert http response object
		m := make([]templatev1.HttpRspTemplate, len(r))
		for n, it := range r {
			m[n].Template = it.Template
		}
		mm := templatev1.HttpRspTemplates(m)
		o.Response(ctx, &mm)
	}

	return nil
}

type UpdateTemplate struct {
	OperateContext
	templatev1.HttpReqTemplate
}

var _ Updater = (*UpdateTemplate)(nil)

func NewUpdateTemplate(d *database.DBManipulator) Updater {
	return &UpdateTemplate{OperateContext: OperateContext{Db: d}}
}

func (o *UpdateTemplate) toModel() *templatev1.Template {
	return &o.Template
}

func (o *UpdateTemplate) Update(ctx echo.Context) error {
	model := o.toModel()

	m := &templatev1.DbSchemaTemplate{Template: *model}
	_, err := o.Db.UpdateTemplate(m)
	if err != nil {
		return err
	}

	if o.Response != nil {
		o.Response(ctx, nil)
	}

	return nil
}

type DeleteTemplate KeyValueParam

func NewDeleteTemplate(d *database.DBManipulator) Remover {
	return (*DeleteTemplate)(NewKeyValueParam(OperateContext{
		Db: d,
	}))
}

func (o *DeleteTemplate) toModel() map[string]string {
	return o.Params
}

func (o *DeleteTemplate) Delete(ctx echo.Context) error {
	model := o.toModel()

	_, err := o.Db.DeleteTemplate(model["uuid"])
	if err != nil {
		return err
	}

	if o.Response != nil {
		o.Response(ctx, nil)
	}

	return nil
}

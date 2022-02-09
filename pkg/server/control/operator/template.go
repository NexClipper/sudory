package operator

import (
	"github.com/NexClipper/sudory/pkg/server/database"

	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
)

//Template
type Template struct {
	// ctx *database.DBManipulator
	ctx database.Context
}

func NewTemplate(ctx database.Context) *Template {
	return &Template{ctx: ctx}
}

func (o *Template) Create(model templatev1.Template) error {
	err := o.ctx.CreateTemplate(templatev1.DbSchemaTemplate{Template: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *Template) Get(uuid string) (*templatev1.Template, error) {

	record, err := o.ctx.GetTemplate(uuid)
	if err != nil {
		return nil, err
	}

	return &record.Template, nil
}

func (o *Template) Find(where string, args ...interface{}) ([]templatev1.Template, error) {
	r, err := o.ctx.FindTemplate(where, args...)
	if err != nil {
		return nil, err
	}

	records := templatev1.TransFormDbSchema(r)

	return records, nil
}

func (o *Template) Update(model templatev1.Template) error {

	err := o.ctx.UpdateTemplate(templatev1.DbSchemaTemplate{Template: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *Template) Delete(uuid string) error {

	//Template Command 삭제
	where := "template_uuid = ?"
	record, err := o.ctx.FindTemplateCommand(where, uuid)
	if err != nil {
		return err
	}

	for _, it := range record {
		err := o.ctx.DeleteTemplateCommand(it.Uuid)
		if err != nil {
			return err
		}
	}
	//Template 삭제
	err = o.ctx.DeleteTemplate(uuid)
	if err != nil {
		return err
	}

	return nil
}

package operator

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	. "github.com/NexClipper/sudory/pkg/server/macro"

	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
)

//Template
type Template struct {
	db *database.DBManipulator
}

func NewTemplate(d *database.DBManipulator) *Template {
	return &Template{db: d}
}

func (o *Template) Create(model templatev1.Template) error {

	//test only
	//uuid 없는 경우 자동 생성
	if false {
		vaild := func(m *templatev1.Template) {
			if len(m.Uuid) == 0 {
				m.Uuid = UuidNewString()
			}
		}

		//uuid 설정 안된 것 강제 생성
		vaild(&model)
	}

	err := o.db.CreateTemplate(templatev1.DbSchemaTemplate{Template: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *Template) Get(uuid string) (*templatev1.Template, error) {

	record, err := o.db.GetTemplate(uuid)
	if err != nil {
		return nil, err
	}

	return &record.Template, nil
}

func (o *Template) Find(where string, args ...interface{}) ([]templatev1.Template, error) {

	r, err := o.db.FindTemplate(where, args...)
	if err != nil {
		return nil, err
	}

	records := templatev1.TransFormDbSchema(r)

	return records, nil
}

func (o *Template) Update(model templatev1.Template) error {

	err := o.db.UpdateTemplate(templatev1.DbSchemaTemplate{Template: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *Template) Delete(uuid string) error {

	err := o.db.DeleteTemplate(uuid)
	if err != nil {
		return err
	}

	//Template Command 레코드 삭제
	where := "template_uuid = ?"
	record, err := o.db.FindTemplateCommand(where, uuid)
	if err != nil {
		return err
	}

	for _, it := range record {
		err := o.db.DeleteTemplateCommand(it.Uuid)
		if err != nil {
			return err
		}
	}

	return nil
}

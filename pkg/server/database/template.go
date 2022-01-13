package database

import (
	"errors"
	"fmt"

	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
)

//create template
func (d *DBManipulator) CreateTemplate(m []templatev1.DbSchemaTemplate) (int64, error) {
	var err error
	tx := d.session()
	tx.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	cnt, err := tx.Insert(m)

	return cnt, err
}

//get template
func (d *DBManipulator) GetTemplate(params map[string]string) (*templatev1.DbSchemaTemplate, error) {
	tx := d.session()

	uuid := params["uuid"]

	//SELECT * FROM template WHERE uuid = ? LIMIT 1
	var model templatev1.DbSchemaTemplate
	has, err := tx.Where("uuid = ?", uuid).
		Get(&model)

	if !has {
		return nil, errors.New(ErrorRecordWasNotFound)
	}

	return &model, err
}

//search []template
func (d *DBManipulator) SearchTemplate(params map[string]string) ([]templatev1.DbSchemaTemplate, error) {
	tx := d.session()

	uuid := fmt.Sprintf("%s%%", params["uuid"])
	name := fmt.Sprintf("%%%s%%", params["name"])
	origin := fmt.Sprintf("%%%s%%", params["origin"])

	//SELECT * FROM template WHERE uuid LIKE ? AND name LIKE ? AND origin LIKE ?
	var model []templatev1.DbSchemaTemplate = make([]templatev1.DbSchemaTemplate, 0)
	err := tx.Where("uuid LIKE ? AND name LIKE ? AND origin LIKE ?", uuid, name, origin).
		Find(&model)

	return model, err
}

//update template
func (d *DBManipulator) UpdateTemplate(m *templatev1.DbSchemaTemplate) (int64, error) {
	var err error
	tx := d.session()
	tx.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	cnt, err := tx.Where("uuid = ?", m.Uuid).
		Update(m)

	return cnt, err
}

//delete template
func (d *DBManipulator) DeleteTemplate(uuid string) (int64, error) {
	var err error
	tx := d.session()
	tx.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	record := new(templatev1.DbSchemaTemplate)

	//DELETE FROM template WHERE uuid = ?
	cnt, err := tx.Where("uuid = ?", uuid).
		Delete(record)

	return cnt, err
}

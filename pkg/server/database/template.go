package database

import (
	"errors"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
)

/* CreateTemplate
   @return int64, error
   @method insert
   @from Template
   @condition []templatev1.DbSchemaTemplate
*/
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

	affect, err := tx.Insert(m)

	return affect, err
}

/* GetTemplate
   @return DbSchemaTemplate, error
   @method get
   @from Template
   @condition uuid
*/
func (d *DBManipulator) GetTemplate(uuid string) (*templatev1.DbSchemaTemplate, error) {
	tx := d.session()

	var model templatev1.DbSchemaTemplate
	//SELECT * FROM template WHERE uuid = ? LIMIT 1
	has, err := tx.Where("uuid = ?", uuid).
		Get(&model)

	if !has {
		return nil, errors.New(ErrorRecordWasNotFound)
	}

	return &model, err
}

/* FindTemplate
   @return []templatev1.DbSchemaTemplate, error
   @method find
   @from Template
   @condition cond, args
*/
func (d *DBManipulator) FindTemplate(where string, args ...interface{}) ([]templatev1.DbSchemaTemplate, error) {
	tx := d.session()

	//SELECT * FROM template WHERE [cond]
	var model []templatev1.DbSchemaTemplate = make([]templatev1.DbSchemaTemplate, 0)
	err := tx.Where(where, args...).
		Find(&model)

	return model, err
}

/* UpdateTemplate
   @return int64, error
   @method update
   @from Template
   @condition DbSchemaTemplate
*/
func (d *DBManipulator) UpdateTemplate(m templatev1.DbSchemaTemplate) (int64, error) {
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

	//아이디를 가져오기
	var record = new(templatev1.DbSchemaTemplate)
	//SELECT * FROM template WHERE uuid = ? LIMIT 1
	has, err := tx.Where("uuid = ?", m.Uuid).
		Get(record)

	if HasError(err) {
		return -1, err
	}
	if !has {
		return -2, errors.New(ErrorRecordWasNotFound)
	}

	affect, err := tx.Where("uuid = ?", m.Uuid).
		Update(m)

	return affect, err
}

/* DeleteTemplate
   @return int64, error
   @method delete
   @from Template
   @condition uuid
*/
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

	//아이디를 가져오기
	record := new(templatev1.DbSchemaTemplate)
	//SELECT * FROM template WHERE uuid = ? LIMIT 1
	has, err := tx.Where("uuid = ?", uuid).
		Get(record)

	if HasError(err) {
		return -1, err
	}
	if !has {
		return -2, errors.New(ErrorRecordWasNotFound)
	}

	//DELETE FROM template WHERE uuid = ?
	affect, err := tx.Where("uuid = ?", uuid).
		Delete(record)

	return affect, err
}

package database

import (
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
)

/* CreateServiceStep
   @return error
   @method insert
   @from ServiceStep
   @condition []stepv1.DbSchemaServiceStep
*/
func (d *DBManipulator) CreateServiceStep(m stepv1.DbSchemaServiceStep) error {
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

	affect, err := tx.Insert(&m)
	if err != nil {
		return err
	}
	if !(0 < affect) {
		return ErrorNoAffecte()
	}
	return nil
}

/* GetServiceStep
   @return DbSchemaServiceStep, error
   @method get
   @from ServiceStep
   @condition uuid
*/
func (d *DBManipulator) GetServiceStep(uuid string) (*stepv1.DbSchemaServiceStep, error) {
	tx := d.session()

	record := new(stepv1.DbSchemaServiceStep)
	has, err := tx.Where("uuid = ?", uuid).
		Get(record)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrorRecordWasNotFound()
	}

	return record, err
}

/* FindServiceStep
   @return []stepv1.DbSchemaServiceStep, error
   @method find
   @from ServiceStep
   @condition where, args
*/
func (d *DBManipulator) FindServiceStep(where string, args ...interface{}) ([]stepv1.DbSchemaServiceStep, error) {
	tx := d.session()

	//SELECT * FROM {table} WHERE [cond]
	model := make([]stepv1.DbSchemaServiceStep, 0)
	err := tx.Where(where, args...).
		Find(&model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

/* UpdateServiceStep
   @return error
   @method update
   @from ServiceStep
   @condition DbSchemaServiceStep
*/
func (d *DBManipulator) UpdateServiceStep(m stepv1.DbSchemaServiceStep) error {
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

	affect, err := tx.Where("uuid = ?", m.Uuid).
		Update(&m)
	if err != nil {
		return err
	}
	if !(0 < affect) {
		return ErrorNoAffecte()
	}
	return nil
}

/* DeleteServiceStep
   @return error
   @method delete
   @from ServiceStep
   @condition uuid
*/
func (d *DBManipulator) DeleteServiceStep(uuid string) error {
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

	record := new(stepv1.DbSchemaServiceStep)
	//DELETE FROM {table} WHERE uuid = ?
	affect, err := tx.Where("uuid = ?", uuid).
		Delete(record)
	if err != nil {
		return err
	}
	if !(0 < affect) {
		return ErrorNoAffecte()
	}
	return nil
}

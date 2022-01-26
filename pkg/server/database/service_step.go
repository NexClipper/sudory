package database

import (
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
)

// CreateServiceStep
//  @return error
//  @method insert
//  @from ServiceStep
//  @condition []DbSchemaServiceStep
func (ctx Session) CreateServiceStep(m stepv1.DbSchemaServiceStep) error {
	tx := ctx.Tx()

	affect, err := tx.
		AllCols().Insert(&m)
	if err != nil {
		return err
	}
	if !(0 < affect) {
		return ErrorNoAffecte()
	}
	return nil
}

// GetServiceStep
//  @return DbSchemaServiceStep, error
//  @method get
//  @from ServiceStep
//  @condition uuid
func (ctx Session) GetServiceStep(uuid string) (*stepv1.DbSchemaServiceStep, error) {
	tx := ctx.Tx()

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

// FindServiceStep
//  @return []stepv1.DbSchemaServiceStep, error
//  @method find
//  @from ServiceStep
//  @condition where, args
func (ctx Session) FindServiceStep(where string, args ...interface{}) ([]stepv1.DbSchemaServiceStep, error) {
	tx := ctx.Tx()

	//SELECT * FROM {table} WHERE [cond]
	model := make([]stepv1.DbSchemaServiceStep, 0)
	err := tx.Where(where, args...).
		Find(&model)
	if err != nil {
		return nil, err
	}
	return model, nil
}

// UpdateServiceStep
//  @return error
//  @method update
//  @from ServiceStep
//  @condition DbSchemaServiceStep
func (ctx Session) UpdateServiceStep(m stepv1.DbSchemaServiceStep) error {
	tx := ctx.Tx()

	affect, err := tx.Where("uuid = ?", m.Uuid).
		AllCols().Update(&m)
	if err != nil {
		return err
	}
	if !(0 < affect) {
		return ErrorNoAffecte()
	}
	return nil
}

// DeleteServiceStep
//  @return error
//  @method delete
//  @from ServiceStep
//  @condition uuid
func (ctx Session) DeleteServiceStep(uuid string) error {
	tx := ctx.Tx()

	record := new(stepv1.DbSchemaServiceStep)
	//DELETE FROM {table} WHERE uuid = ?
	affect, err := tx.Where("uuid = ?", uuid).
		Delete(record)
	if err != nil {
		return err
	}
	if !(0 < affect) {
		return nil //idempotent
	}
	return nil
}

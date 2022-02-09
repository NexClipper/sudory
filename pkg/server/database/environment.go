package database

import (
	envv1 "github.com/NexClipper/sudory/pkg/server/model/environment/v1"
)

// CreateEnvironment
//  @return error
//  @method insert
//  @from Environment
//  @condition DbSchemaEnvironment
func (ctx Session) CreateEnvironment(m envv1.DbSchemaEnvironment) error {
	tx := ctx.Tx()

	affect, err := tx.Insert(&m)
	if err != nil {
		return err
	}
	if !(0 < affect) {
		return ErrorNoAffecte()
	}
	return nil
}

// GetEnvironment
//  @return DbSchemaEnvironment, error
//  @method get
//  @from Environment
//  @condition uuid
func (ctx Session) GetEnvironment(uuid string) (*envv1.DbSchemaEnvironment, error) {
	tx := ctx.Tx()

	var record = new(envv1.DbSchemaEnvironment)
	has, err := tx.Where("uuid = ?", uuid).
		Get(record)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrorRecordWasNotFound()
	}

	return record, nil
}

// FindEnvironment
//  @return []DbSchemaEnvironment, error
//  @method find
//  @from Environment
//  @condition where, args
func (ctx Session) FindEnvironment(where string, args ...interface{}) ([]envv1.DbSchemaEnvironment, error) {
	tx := ctx.Tx()

	model := make([]envv1.DbSchemaEnvironment, 0)
	err := tx.Where(where, args...).
		Find(&model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// UpdateEnvironment
//  @return error
//  @method update
//  @from Environment
//  @condition DbSchemaEnvironment
func (ctx Session) UpdateEnvironment(m envv1.DbSchemaEnvironment) error {
	tx := ctx.Tx()

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

// DeleteEnvironment
//  @return error
//  @method delete
//  @from Environment
//  @condition uuid
func (ctx Session) DeleteEnvironment(uuid string) error {
	tx := ctx.Tx()

	record := new(envv1.DbSchemaEnvironment)
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

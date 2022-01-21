package database

import (
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
)

// CreateService
//  @return error
//  @method insert
//  @from Service
//  @condition DbSchemaService
func (ctx Session) CreateService(m servicev1.DbSchemaService) error {
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

// GetService
//  @return DbSchemaService, error
//  @method get
//  @from Service
//  @condition uuid
func (ctx Session) GetService(uuid string) (*servicev1.DbSchemaService, error) {
	tx := ctx.Tx()

	record := new(servicev1.DbSchemaService)
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

// FindService
//  @return []DbSchemaService, error
//  @method find
//  @from Service
//  @condition where, args
func (ctx Session) FindService(where string, args ...interface{}) ([]servicev1.DbSchemaService, error) {
	tx := ctx.Tx()

	//SELECT * FROM {table} WHERE [cond]
	var model = make([]servicev1.DbSchemaService, 0)
	err := tx.Where(where, args...).
		Find(&model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// UpdateService
//  @return error
//  @method update
//  @from Service
//  @condition DbSchemaService
func (ctx Session) UpdateService(m servicev1.DbSchemaService) error {
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

// DeleteService
//  @return error
//  @method delete
//  @from Service
//  @condition uuid
func (ctx Session) DeleteService(uuid string) error {
	tx := ctx.Tx()

	record := new(servicev1.DbSchemaService)
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

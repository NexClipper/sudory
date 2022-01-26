package database

import (
	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
)

// CreateTemplate
//  @return error
//  @method insert
//  @from Template
//  @condition templatev1.DbSchemaTemplate
func (ctx Session) CreateTemplate(m templatev1.DbSchemaTemplate) error {
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

// GetTemplate
//  @return DbSchemaTemplate, error
//  @method get
//  @from Template
//  @condition uuid
func (ctx Session) GetTemplate(uuid string) (*templatev1.DbSchemaTemplate, error) {
	tx := ctx.Tx()

	var record = new(templatev1.DbSchemaTemplate)
	//SELECT * FROM {table} WHERE uuid = ? LIMIT 1
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

// FindTemplate
//  @return []templatev1.DbSchemaTemplate, error
//  @method find
//  @from Template
//  @condition where, args
func (ctx Session) FindTemplate(where string, args ...interface{}) ([]templatev1.DbSchemaTemplate, error) {
	tx := ctx.Tx()

	//SELECT * FROM template WHERE [cond]
	model := make([]templatev1.DbSchemaTemplate, 0)
	err := tx.Where(where, args...).
		Find(&model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// UpdateTemplate
//  @return error
//  @method update
//  @from Template
//  @condition DbSchemaTemplate
func (ctx Session) UpdateTemplate(m templatev1.DbSchemaTemplate) error {
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

// DeleteTemplate
//  @return error
//  @method delete
//  @from Template
//  @condition uuid
func (ctx Session) DeleteTemplate(uuid string) error {
	tx := ctx.Tx()

	record := new(templatev1.DbSchemaTemplate)
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

package database

import (
	clientv1 "github.com/NexClipper/sudory/pkg/server/model/client/v1"
)

// CreateClient
//  @return error
//  @method insert
//  @from Client
//  @condition clientv1.DbSchemaClient
func (ctx Session) CreateClient(m clientv1.DbSchemaClient) error {
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

// GetClient
//  @return DbSchemaClient, error
//  @method get
//  @from Client
//  @condition uuid
func (ctx Session) GetClient(uuid string) (*clientv1.DbSchemaClient, error) {
	tx := ctx.Tx()

	var record = new(clientv1.DbSchemaClient)
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

// FindClient
//  @return []clientv1.DbSchemaClient, error
//  @method find
//  @from Client
//  @condition where, args
func (ctx Session) FindClient(where string, args ...interface{}) ([]clientv1.DbSchemaClient, error) {
	tx := ctx.Tx()

	model := make([]clientv1.DbSchemaClient, 0)
	err := tx.Where(where, args...).
		Find(&model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// UpdateClient
//  @return error
//  @method update
//  @from Client
//  @condition DbSchemaClient
func (ctx Session) UpdateClient(m clientv1.DbSchemaClient) error {
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

// DeleteClient
//  @return error
//  @method delete
//  @from Client
//  @condition uuid
func (ctx Session) DeleteClient(uuid string) error {
	tx := ctx.Tx()

	record := new(clientv1.DbSchemaClient)
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

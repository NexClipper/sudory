package database

import (
	"github.com/NexClipper/sudory/pkg/server/database/query_parser"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
)

// CreateSession
//  @return error
//  @method insert
//  @from Session
//  @condition DbSchemaSession
func (ctx Session) CreateSession(m sessionv1.DbSchemaSession) error {
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

// GetSession
//  @return DbSchemaSession, error
//  @method get
//  @from Session
//  @condition uuid
func (ctx Session) GetSession(uuid string) (*sessionv1.DbSchemaSession, error) {
	tx := ctx.Tx()

	var record = new(sessionv1.DbSchemaSession)
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

// FindSession
//  @return []DbSchemaSession, error
//  @method find
//  @from Session
//  @condition where, args
func (ctx Session) FindSession(where string, args ...interface{}) ([]sessionv1.DbSchemaSession, error) {
	tx := ctx.Tx()

	//SELECT * FROM {table} WHERE [cond]
	var model = make([]sessionv1.DbSchemaSession, 0)
	err := tx.Where(where, args...).
		Find(&model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// QuerySession
//  @return []DbSchemaSession, error
//  @method find
//  @from Session
//  @condition where, args
func (ctx Session) QuerySession(query *query_parser.QueryParser) ([]sessionv1.DbSchemaSession, error) {
	tx := ctx.Tx()

	model := make([]sessionv1.DbSchemaSession, 0)
	err := query.Prepare(tx).
		Find(&model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// UpdateSession
//  @return error
//  @method update
//  @from Session
//  @condition DbSchemaSession
func (ctx Session) UpdateSession(m sessionv1.DbSchemaSession) error {
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

// DeleteSession
//  @return error
//  @method delete
//  @from Session
//  @condition uuid
func (ctx Session) DeleteSession(uuid string) error {
	tx := ctx.Tx()

	record := new(sessionv1.DbSchemaSession)
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

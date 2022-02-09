package database

import (
	"github.com/NexClipper/sudory/pkg/server/database/query_parser"
	tokenv1 "github.com/NexClipper/sudory/pkg/server/model/token/v1"
)

// CreateToken
//  @return error
//  @method insert
//  @from Token
//  @condition DbSchemaToken
func (ctx Session) CreateToken(m tokenv1.DbSchemaToken) error {
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

// GetToken
//  @return DbSchemaToken, error
//  @method get
//  @from Token
//  @condition uuid
func (ctx Session) GetToken(uuid string) (*tokenv1.DbSchemaToken, error) {
	tx := ctx.Tx()

	var record = new(tokenv1.DbSchemaToken)
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

// FindToken
//  @return []DbSchemaToken, error
//  @method find
//  @from Token
//  @condition where, args
func (ctx Session) FindToken(where string, args ...interface{}) ([]tokenv1.DbSchemaToken, error) {
	tx := ctx.Tx()

	//SELECT * FROM {table} WHERE [cond]
	var model = make([]tokenv1.DbSchemaToken, 0)
	err := tx.Where(where, args...).
		Find(&model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// QueryToken
//  @return []DbSchemaToken, error
//  @method find
//  @from Token
//  @condition where, args
func (ctx Session) QueryToken(query *query_parser.QueryParser) ([]tokenv1.DbSchemaToken, error) {
	tx := ctx.Tx()

	model := make([]tokenv1.DbSchemaToken, 0)
	err := query.Prepare(tx).
		Find(&model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// UpdateToken
//  @return error
//  @method update
//  @from Token
//  @condition DbSchemaToken
func (ctx Session) UpdateToken(m tokenv1.DbSchemaToken) error {
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

// DeleteToken
//  @return error
//  @method delete
//  @from Token
//  @condition uuid
func (ctx Session) DeleteToken(uuid string) error {
	tx := ctx.Tx()

	record := new(tokenv1.DbSchemaToken)
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

package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	tokenv1 "github.com/NexClipper/sudory/pkg/server/model/token/v1"
	"github.com/pkg/errors"
)

type Token struct {
	ctx database.Context
}

func NewToken(ctx database.Context) *Token {
	return &Token{ctx: ctx}
}

func (vault Token) CreateClusterToken(token tokenv1.Token) (*tokenv1.DbSchema, error) {
	//create
	record := &tokenv1.DbSchema{Token: token}
	if err := vault.ctx.Create(record); err != nil {
		return nil, errors.Wrapf(err, "database create")
	}

	return record, nil
}

func (vault Token) Get(uuid string) (*tokenv1.DbSchema, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	record := &tokenv1.DbSchema{}
	if err := vault.ctx.Where(where, args...).Get(record); err != nil {
		return nil, errors.Wrapf(err, "database get where=%s args=%+v", where, args)
	}

	return record, nil
}

func (vault Token) Find(where string, args ...interface{}) ([]tokenv1.DbSchema, error) {
	records := make([]tokenv1.DbSchema, 0)
	if err := vault.ctx.Where(where, args...).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find where=%s args=%+v", where, args)
	}

	return records, nil
}

func (vault Token) Query(query map[string]string) ([]tokenv1.DbSchema, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser query=%+v", query)
	}

	//find service
	records := make([]tokenv1.DbSchema, 0)
	if err := vault.ctx.Prepared(preparer).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find query=%+v", query)
	}

	return records, nil
}

func (vault Token) Update(model tokenv1.Token) (*tokenv1.DbSchema, error) {
	where := "uuid = ?"
	args := []interface{}{
		model.Uuid,
	}
	record := &tokenv1.DbSchema{Token: model}
	if err := vault.ctx.Where(where, args...).Update(record); err != nil {
		return nil, errors.Wrapf(err, "database update where=%s args=%+v", where, args)
	}

	//make result
	record_, err := vault.Get(record.Uuid)
	if err != nil {
		return nil, errors.Wrapf(err, "make update result")
	}

	return record_, nil
}

func (vault Token) Delete(uuid string) error {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &tokenv1.DbSchema{}
	if err := vault.ctx.Where(where, args...).Delete(model); err != nil {
		return errors.Wrapf(err, "database delete where=%s args=%+v", where, args)
	}

	return nil
}

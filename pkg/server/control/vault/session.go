package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
	"github.com/pkg/errors"
)

type Session struct {
	ctx database.Context
}

func NewSession(ctx database.Context) *Session {
	return &Session{ctx: ctx}
}

func (vault Session) Create(model sessionv1.Session) (*sessionv1.DbSchema, error) {
	record := &sessionv1.DbSchema{Session: model}
	if err := vault.ctx.Create(record); err != nil {
		return nil, errors.Wrapf(err, "database create")
	}

	return record, nil
}

func (vault Session) Get(uuid string) (*sessionv1.DbSchema, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	record := &sessionv1.DbSchema{}
	if err := vault.ctx.Where(where, args...).Get(record); err != nil {
		return nil, errors.Wrapf(err, "database get%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return record, nil
}

func (vault Session) Find(where string, args ...interface{}) ([]sessionv1.DbSchema, error) {
	records := make([]sessionv1.DbSchema, 0)
	if err := vault.ctx.Where(where, args...).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return records, nil
}

func (vault Session) Query(query map[string]string) ([]sessionv1.DbSchema, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser%v",
			logs.KVL(
				"query", query,
			))
	}

	//find service
	records := make([]sessionv1.DbSchema, 0)
	if err := vault.ctx.Prepared(preparer).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"query", query,
			))
	}

	return records, nil
}

func (vault Session) Update(model sessionv1.Session) (*sessionv1.DbSchema, error) {
	where := "uuid = ?"
	args := []interface{}{
		model.Uuid,
	}
	record := &sessionv1.DbSchema{Session: model}
	if err := vault.ctx.Where(where, args...).Update(record); err != nil {
		return nil, errors.Wrapf(err, "database update%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return record, nil
}

func (vault Session) Delete(uuid string) error {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}

	record := &sessionv1.DbSchema{}
	if err := vault.ctx.Where(where, args...).Delete(record); err != nil {
		return errors.Wrapf(err, "database delete%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return nil
}

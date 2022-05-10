package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type Session struct {
	tx *xorm.Session
}

func NewSession(tx *xorm.Session) *Session {
	return &Session{tx: tx}
}

func (vault Session) Create(model sessionv1.Session) (*sessionv1.Session, error) {
	if err := database.XormCreate(vault.tx, &model); err != nil {
		return nil, errors.Wrapf(err, "create %v", model.TableName())
	}

	return &model, nil
}

func (vault Session) Get(uuid string) (*sessionv1.Session, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &sessionv1.Session{}
	if err := database.XormGet(
		vault.tx.Where(where, args...), model); err != nil {
		return nil, errors.Wrapf(err, "get %v", model.TableName())
	}

	return model, nil
}

func (vault Session) Find(where string, args ...interface{}) ([]sessionv1.Session, error) {
	models := make([]sessionv1.Session, 0)
	if err := database.XormFind(
		vault.tx.Where(where, args...), &models); err != nil {
		return nil, errors.Wrapf(err, "find %v", new(sessionv1.Session).TableName())
	}

	return models, nil
}

func (vault Session) Query(query map[string]string) ([]sessionv1.Session, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser%v",
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]sessionv1.Session, 0)
	if err := database.XormFind(
		preparer.Prepared(vault.tx), &models); err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(sessionv1.Session).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	return models, nil
}

func (vault Session) Update(model sessionv1.Session) (*sessionv1.Session, error) {
	where := "uuid = ?"
	args := []interface{}{
		model.Uuid,
	}
	if err := database.XormUpdate(
		vault.tx.Where(where, args...), &model); err != nil {
		return nil, errors.Wrapf(err, "update %v", model.TableName())
	}

	return &model, nil
}

func (vault Session) Delete(uuid string) error {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}

	model := &sessionv1.Session{}
	if err := database.XormDelete(
		vault.tx.Where(where, args...), model); err != nil {
		return errors.Wrapf(err, "delete %v", model.TableName())
	}

	return nil
}

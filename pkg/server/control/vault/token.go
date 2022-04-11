package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	tokenv1 "github.com/NexClipper/sudory/pkg/server/model/token/v1"
	"github.com/pkg/errors"
)

type Token struct {
	ctx database.Context
}

func NewToken(ctx database.Context) *Token {
	return &Token{ctx: ctx}
}

func (vault Token) CreateToken(model tokenv1.Token) (*tokenv1.Token, error) {
	//create
	if err := vault.ctx.Create(&model); err != nil {
		return nil, errors.Wrapf(err, "database create")
	}
	return &model, nil
}

func (vault Token) Get(uuid string) (*tokenv1.Token, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &tokenv1.Token{}
	if err := vault.ctx.Where(where, args...).Get(model); err != nil {
		return nil, errors.Wrapf(err, "database get%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return model, nil
}

func (vault Token) Find(where string, args ...interface{}) ([]tokenv1.Token, error) {
	models := make([]tokenv1.Token, 0)
	if err := vault.ctx.Where(where, args...).Find(&models); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return models, nil
}

func (vault Token) Query(query map[string]string) ([]tokenv1.Token, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser%v",
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]tokenv1.Token, 0)
	if err := vault.ctx.Prepared(preparer).Find(&models); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"query", query,
			))
	}

	return models, nil
}

func (vault Token) Update(model tokenv1.Token) (*tokenv1.Token, error) {
	where := "uuid = ?"
	args := []interface{}{
		model.Uuid,
	}
	if err := vault.ctx.Where(where, args...).Update(&model); err != nil {
		return nil, errors.Wrapf(err, "database update%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return &model, nil
}

func (vault Token) Delete(uuid string) error {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &tokenv1.Token{}
	if err := vault.ctx.Where(where, args...).Delete(model); err != nil {
		return errors.Wrapf(err, "database delete%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return nil
}

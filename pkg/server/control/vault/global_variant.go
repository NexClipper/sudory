package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	globvarv1 "github.com/NexClipper/sudory/pkg/server/model/global_variant/v1"
	"github.com/pkg/errors"
)

type GlobalVariant struct {
	ctx database.Context
}

func NewEnvironment(ctx database.Context) *GlobalVariant {
	return &GlobalVariant{ctx: ctx}
}

func (vault GlobalVariant) Create(model globvarv1.GlobalVariant) (*globvarv1.GlobalVariant, error) {
	if err := vault.ctx.Create(&model); err != nil {
		return nil, errors.Wrapf(err, "database create")
	}

	return &model, nil
}

func (vault GlobalVariant) Get(uuid string) (*globvarv1.GlobalVariant, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &globvarv1.GlobalVariant{}
	if err := vault.ctx.Where(where, args...).Get(model); err != nil {
		return nil, errors.Wrapf(err, "database get%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return model, nil
}

func (vault GlobalVariant) Find(where string, args ...interface{}) ([]globvarv1.GlobalVariant, error) {
	models := make([]globvarv1.GlobalVariant, 0)
	if err := vault.ctx.Where(where, args...).Find(&models); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return models, nil
}

func (vault GlobalVariant) Query(query map[string]string) ([]globvarv1.GlobalVariant, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser%v",
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]globvarv1.GlobalVariant, 0)
	if err := vault.ctx.Prepared(preparer).Find(&models); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"query", query,
			))
	}

	return models, nil
}

func (vault GlobalVariant) Update(model globvarv1.GlobalVariant) (*globvarv1.GlobalVariant, error) {
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

func (vault GlobalVariant) Delete(uuid string) error {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &globvarv1.GlobalVariant{}
	if err := vault.ctx.Where(where, args...).Delete(model); err != nil {
		return errors.Wrapf(err, "database delete%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return nil
}

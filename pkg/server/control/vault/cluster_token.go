package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	clustertokenv1 "github.com/NexClipper/sudory/pkg/server/model/cluster_token/v1"
	"github.com/pkg/errors"
)

type ClusterToken struct {
	ctx database.Context
}

func NewClusterToken(ctx database.Context) *ClusterToken {
	return &ClusterToken{ctx: ctx}
}

func (vault ClusterToken) CreateToken(model clustertokenv1.ClusterToken) (*clustertokenv1.ClusterToken, error) {
	//create
	if err := vault.ctx.Create(&model); err != nil {
		return nil, errors.Wrapf(err, "database create")
	}
	return &model, nil
}

func (vault ClusterToken) Get(uuid string) (*clustertokenv1.ClusterToken, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &clustertokenv1.ClusterToken{}
	if err := vault.ctx.Where(where, args...).Get(model); err != nil {
		return nil, errors.Wrapf(err, "database get%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return model, nil
}

func (vault ClusterToken) Find(where string, args ...interface{}) ([]clustertokenv1.ClusterToken, error) {
	models := make([]clustertokenv1.ClusterToken, 0)
	if err := vault.ctx.Where(where, args...).Find(&models); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return models, nil
}

func (vault ClusterToken) Query(query map[string]string) ([]clustertokenv1.ClusterToken, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser%v",
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]clustertokenv1.ClusterToken, 0)
	if err := vault.ctx.Prepared(preparer).Find(&models); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"query", query,
			))
	}

	return models, nil
}

func (vault ClusterToken) Update(model clustertokenv1.ClusterToken) (*clustertokenv1.ClusterToken, error) {
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

func (vault ClusterToken) Delete(uuid string) error {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &clustertokenv1.ClusterToken{}
	if err := vault.ctx.Where(where, args...).Delete(model); err != nil {
		return errors.Wrapf(err, "database delete%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return nil
}

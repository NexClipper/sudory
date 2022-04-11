package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	clusterv1 "github.com/NexClipper/sudory/pkg/server/model/cluster/v1"
	"github.com/pkg/errors"
)

type Cluster struct {
	ctx database.Context
}

func NewCluster(ctx database.Context) *Cluster {
	return &Cluster{ctx: ctx}
}

func (vault Cluster) Create(model clusterv1.Cluster) (*clusterv1.Cluster, error) {
	if err := vault.ctx.Create(&model); err != nil {
		return nil, errors.Wrapf(err, "database create")
	}

	return &model, nil
}

func (vault Cluster) Get(uuid string) (*clusterv1.Cluster, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &clusterv1.Cluster{}
	if err := vault.ctx.Where(where, args...).Get(model); err != nil {
		return nil, errors.Wrapf(err, "database get%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return model, nil
}

func (vault Cluster) Find(where string, args ...interface{}) ([]clusterv1.Cluster, error) {
	models := make([]clusterv1.Cluster, 0)
	if err := vault.ctx.Where(where, args...).Find(&models); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return models, nil
}

func (vault Cluster) Query(query map[string]string) ([]clusterv1.Cluster, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser%v",
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]clusterv1.Cluster, 0)
	if err := vault.ctx.Prepared(preparer).Find(&models); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"query", query,
			))
	}

	return models, nil
}

func (vault Cluster) Update(model clusterv1.Cluster) (*clusterv1.Cluster, error) {
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

func (vault Cluster) Delete(uuid string) error {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &clusterv1.Cluster{}
	if err := vault.ctx.Where(where, args...).Delete(model); err != nil {
		return errors.Wrapf(err, "database delete%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return nil
}

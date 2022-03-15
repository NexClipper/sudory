package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	clusterv1 "github.com/NexClipper/sudory/pkg/server/model/cluster/v1"
	"github.com/pkg/errors"
)

type Cluster struct {
	ctx database.Context
}

func NewCluster(ctx database.Context) *Cluster {
	return &Cluster{ctx: ctx}
}

func (vault Cluster) Create(model clusterv1.Cluster) (*clusterv1.DbSchema, error) {
	record := &clusterv1.DbSchema{Cluster: model}
	if err := vault.ctx.Create(record); err != nil {
		return nil, errors.Wrapf(err, "database create")
	}

	return record, nil
}

func (vault Cluster) Get(uuid string) (*clusterv1.DbSchema, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	record := &clusterv1.DbSchema{}
	if err := vault.ctx.Where(where, args...).Get(record); err != nil {
		return nil, errors.Wrapf(err, "database get where=%s args=%+v", where, args)
	}

	return record, nil
}

func (vault Cluster) Find(where string, args ...interface{}) ([]clusterv1.DbSchema, error) {
	records := make([]clusterv1.DbSchema, 0)
	if err := vault.ctx.Where(where, args...).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find where=%s args=%+v", where, args)
	}

	return records, nil
}

func (vault Cluster) Query(query map[string]string) ([]clusterv1.DbSchema, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser query=%+v", query)
	}

	//find service
	records := make([]clusterv1.DbSchema, 0)
	if err := vault.ctx.Prepared(preparer).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find query=%+v", query)
	}

	return records, nil
}

func (vault Cluster) Update(model clusterv1.Cluster) (*clusterv1.DbSchema, error) {
	where := "uuid = ?"
	args := []interface{}{
		model.Uuid,
	}
	record := &clusterv1.DbSchema{Cluster: model}
	if err := vault.ctx.Where(where, args...).Update(record); err != nil {
		return nil, errors.Wrapf(err, "database update where=%s args=%+v", where, args)
	}

	return record, nil
}

func (vault Cluster) Delete(uuid string) error {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	record := &clusterv1.DbSchema{}
	if err := vault.ctx.Where(where, args...).Delete(record); err != nil {
		return errors.Wrapf(err, "database delete where=%s args=%+v", where, args)
	}

	return nil
}

package operator

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	clusterv1 "github.com/NexClipper/sudory/pkg/server/model/cluster/v1"
)

type Cluster struct {
	ctx database.Context
}

func NewCluster(ctx database.Context) *Cluster {
	return &Cluster{ctx: ctx}
}

func (o *Cluster) Create(model clusterv1.Cluster) error {
	err := o.ctx.CreateCluster(clusterv1.DbSchemaCluster{Cluster: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *Cluster) Get(uuid string) (*clusterv1.Cluster, error) {

	record, err := o.ctx.GetCluster(uuid)
	if err != nil {
		return nil, err
	}

	return &record.Cluster, nil
}

func (o *Cluster) Find(where string, args ...interface{}) ([]clusterv1.Cluster, error) {
	r, err := o.ctx.FindCluster(where, args...)
	if err != nil {
		return nil, err
	}

	records := clusterv1.TransFormDbSchema(r)

	return records, nil
}

func (o *Cluster) Update(model clusterv1.Cluster) error {

	err := o.ctx.UpdateCluster(clusterv1.DbSchemaCluster{Cluster: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *Cluster) Delete(uuid string) error {

	err := o.ctx.DeleteCluster(uuid)
	if err != nil {
		return err
	}

	return nil
}

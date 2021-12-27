package operator

import (
	"github.com/NexClipper/sudory-prototype-r1/pkg/database"
	"github.com/NexClipper/sudory-prototype-r1/pkg/model"
	"github.com/labstack/echo/v4"
)

type Cluster struct {
	db *database.DBManipulator

	Name string

	Response ResponseFn
}

func NewCluster(d *database.DBManipulator) Operator {
	return &Cluster{db: d}
}

func (o *Cluster) toModel() *model.Cluster {
	m := &model.Cluster{
		Name: o.Name,
	}

	return m
}

func (o *Cluster) Create(ctx echo.Context) error {
	cluster := o.toModel()

	_, err := o.db.CreateCluster(cluster)
	if err != nil {
		return err
	}

	if o.Response != nil {
		o.Response(ctx)
	}

	return nil
}

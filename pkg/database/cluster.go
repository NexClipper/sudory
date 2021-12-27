package database

import "github.com/NexClipper/sudory-prototype-r1/pkg/model"

func (d *DBManipulator) CreateCluster(m *model.Cluster) (int64, error) {
	tx := d.session()

	return tx.Insert(m)
}

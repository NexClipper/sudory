package database

import "github.com/NexClipper/sudory-prototype-r1/pkg/server/model"

func (d *DBManipulator) CreateCluster(m *model.Cluster) (int64, error) {
	tx := d.session()
	tx.Begin()

	cnt, err := tx.Insert(m)

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	return cnt, err
}

func (d *DBManipulator) GetCluster(m *model.Cluster) (*model.Cluster, error) {
	tx := d.session()

	var cluster model.Cluster
	_, err := tx.ID(m.ID).Get(&cluster)

	return &cluster, err
}

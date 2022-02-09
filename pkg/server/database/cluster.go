package database

import (
	clusterv1 "github.com/NexClipper/sudory/pkg/server/model/cluster/v1"
)

// CreateCluster
//  @return error
//  @method insert
//  @from Cluster
//  @condition DbSchemaCluster
func (ctx Session) CreateCluster(m clusterv1.DbSchemaCluster) error {
	tx := ctx.Tx()

	affect, err := tx.Insert(&m)
	if err != nil {
		return err
	}
	if !(0 < affect) {
		return ErrorNoAffecte()
	}
	return nil
}

// GetCluster
//  @return DbSchemaCluster, error
//  @method get
//  @from Cluster
//  @condition uuid
func (ctx Session) GetCluster(uuid string) (*clusterv1.DbSchemaCluster, error) {
	tx := ctx.Tx()

	var record = new(clusterv1.DbSchemaCluster)
	has, err := tx.Where("uuid = ?", uuid).
		Get(record)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrorRecordWasNotFound()
	}

	return record, nil
}

// FindCluster
//  @return []DbSchemaCluster, error
//  @method find
//  @from Cluster
//  @condition where, args
func (ctx Session) FindCluster(where string, args ...interface{}) ([]clusterv1.DbSchemaCluster, error) {
	tx := ctx.Tx()

	model := make([]clusterv1.DbSchemaCluster, 0)
	err := tx.Where(where, args...).
		Find(&model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// UpdateCluster
//  @return error
//  @method update
//  @from Cluster
//  @condition DbSchemaCluster
func (ctx Session) UpdateCluster(m clusterv1.DbSchemaCluster) error {
	tx := ctx.Tx()

	affect, err := tx.Where("uuid = ?", m.Uuid).
		Update(&m)
	if err != nil {
		return err
	}
	if !(0 < affect) {
		return ErrorNoAffecte()
	}
	return nil
}

// DeleteCluster
//  @return error
//  @method delete
//  @from Cluster
//  @condition uuid
func (ctx Session) DeleteCluster(uuid string) error {
	tx := ctx.Tx()

	record := new(clusterv1.DbSchemaCluster)
	affect, err := tx.Where("uuid = ?", uuid).
		Delete(record)
	if err != nil {
		return err
	}
	if !(0 < affect) {
		return nil //idempotent
	}
	return nil
}

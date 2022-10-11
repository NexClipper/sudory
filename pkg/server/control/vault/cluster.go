package vault

// import (
// 	"github.com/NexClipper/sudory/pkg/server/database"
// 	"github.com/NexClipper/sudory/pkg/server/database/prepare"
// 	"github.com/NexClipper/sudory/pkg/server/macro/logs"
// 	clusterv1 "github.com/NexClipper/sudory/pkg/server/model/cluster/v1"
// 	"github.com/pkg/errors"
// 	"xorm.io/xorm"
// )

// type Cluster struct {
// 	tx *xorm.Session
// }

// func NewCluster(tx *xorm.Session) *Cluster {
// 	return &Cluster{tx: tx}
// }

// func (vault Cluster) Create(model clusterv1.Cluster) (*clusterv1.Cluster, error) {
// 	if err := database.XormCreate(vault.tx, &model); err != nil {
// 		return nil, errors.Wrapf(err, "create %v", model.TableName())
// 	}

// 	return &model, nil
// }

// func (vault Cluster) Get(uuid string) (*clusterv1.Cluster, error) {
// 	where := "uuid = ?"
// 	args := []interface{}{
// 		uuid,
// 	}
// 	model := &clusterv1.Cluster{}
// 	if err := database.XormGet(
// 		vault.tx.Where(where, args...), model); err != nil {
// 		return nil, errors.Wrapf(err, "get %v", model.TableName())
// 	}

// 	return model, nil
// }

// func (vault Cluster) Find(where string, args ...interface{}) ([]clusterv1.Cluster, error) {
// 	models := make([]clusterv1.Cluster, 0)
// 	if err := database.XormFind(
// 		vault.tx.Where(where, args...), &models); err != nil {
// 		return nil, errors.Wrapf(err, "find %v", new(clusterv1.Cluster).TableName())
// 	}

// 	return models, nil
// }

// func (vault Cluster) Query(query map[string]string) ([]clusterv1.Cluster, error) {
// 	//parse query
// 	preparer, err := prepare.NewParser(query)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "query %v%v", new(clusterv1.Cluster).TableName(),
// 			logs.KVL(
// 				"query", query,
// 			))
// 	}

// 	//find service
// 	models := make([]clusterv1.Cluster, 0)
// 	if err := database.XormFind(
// 		preparer.Prepared(vault.tx), &models); err != nil {
// 		return nil, errors.Wrapf(err, "query %v%v", new(clusterv1.Cluster).TableName(),
// 			logs.KVL(
// 				"query", query,
// 			))
// 	}

// 	return models, nil
// }

// func (vault Cluster) Update(model clusterv1.Cluster) (*clusterv1.Cluster, error) {
// 	where := "uuid = ?"
// 	args := []interface{}{
// 		model.Uuid,
// 	}

// 	if err := database.XormUpdate(
// 		vault.tx.Where(where, args...), &model); err != nil {
// 		return nil, errors.Wrapf(err, "update %v", model.TableName())
// 	}

// 	return &model, nil
// }

// func (vault Cluster) Delete(uuid string) error {
// 	where := "uuid = ?"
// 	args := []interface{}{
// 		uuid,
// 	}
// 	model := &clusterv1.Cluster{}
// 	if err := database.XormDelete(
// 		vault.tx.Where(where, args...), model); err != nil {
// 		return errors.Wrapf(err, "delete %v", model.TableName())
// 	}

// 	return nil
// }

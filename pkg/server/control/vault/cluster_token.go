package vault

// import (
// 	"github.com/NexClipper/sudory/pkg/server/database"
// 	"github.com/NexClipper/sudory/pkg/server/database/prepare"
// 	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
// 	"github.com/NexClipper/sudory/pkg/server/macro/logs"
// 	clustertokenv1 "github.com/NexClipper/sudory/pkg/server/model/cluster_token/v1"
// 	clustertokenv2 "github.com/NexClipper/sudory/pkg/server/model/cluster_token/v2"
// 	"github.com/pkg/errors"
// 	"xorm.io/xorm"
// )

// type ClusterToken struct {
// 	tx *xorm.Session
// }

// func NewClusterToken(tx *xorm.Session) *ClusterToken {
// 	return &ClusterToken{tx: tx}
// }

// func (vault ClusterToken) CreateToken(model clustertokenv1.ClusterToken) (*clustertokenv1.ClusterToken, error) {
// 	if err := database.XormCreate(vault.tx, &model); err != nil {
// 		return nil, errors.Wrapf(err, "create %v", model.TableName())
// 	}
// 	return &model, nil
// }

// func (vault ClusterToken) Get(uuid string) (*clustertokenv1.ClusterToken, error) {
// 	where := "uuid = ?"
// 	args := []interface{}{
// 		uuid,
// 	}
// 	model := &clustertokenv1.ClusterToken{}
// 	if err := database.XormGet(
// 		vault.tx.Where(where, args...), model); err != nil {
// 		return nil, errors.Wrapf(err, "get %v", model.TableName())
// 	}

// 	return model, nil
// }

// func (vault ClusterToken) Find(where string, args ...interface{}) ([]clustertokenv1.ClusterToken, error) {
// 	models := make([]clustertokenv1.ClusterToken, 0)
// 	if err := database.XormFind(
// 		vault.tx.Where(where, args...), &models); err != nil {
// 		return nil, errors.Wrapf(err, "find %v", new(clustertokenv1.ClusterToken).TableName())
// 	}

// 	return models, nil
// }

// func (vault ClusterToken) Query(query map[string]string) ([]clustertokenv1.ClusterToken, error) {
// 	//parse query
// 	preparer, err := prepare.NewParser(query)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "query %v%v", new(clustertokenv1.ClusterToken).TableName(),
// 			logs.KVL(
// 				"query", query,
// 			))
// 	}

// 	//find service
// 	models := make([]clustertokenv1.ClusterToken, 0)
// 	if err := database.XormFind(
// 		preparer.Prepared(vault.tx), &models); err != nil {
// 		return nil, errors.Wrapf(err, "query %v%v", new(clustertokenv1.ClusterToken).TableName(),
// 			logs.KVL(
// 				"query", query,
// 			))
// 	}

// 	return models, nil
// }

// func (vault ClusterToken) Update(model clustertokenv1.ClusterToken) (*clustertokenv1.ClusterToken, error) {
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

// func (vault ClusterToken) Delete(uuid string) error {
// 	where := "uuid = ?"
// 	args := []interface{}{
// 		uuid,
// 	}
// 	model := &clustertokenv1.ClusterToken{}
// 	if err := database.XormDelete(
// 		vault.tx.Where(where, args...), model); err != nil {
// 		return errors.Wrapf(err, "delete %v", model.TableName())
// 	}

// 	return nil
// }

// func GetClusterToken(tx vanilla.Preparer, uuid string) (token clustertokenv2.ClusterToken, err error) {
// 	token.Uuid = uuid
// 	eq_uuid := vanilla.Equal("uuid", token.Uuid)
// 	// cond := vanilla.And(
// 	// 	eq_uuid,
// 	// 	vanilla.Equal("deleted", token.Deleted),
// 	// )

// 	// token = new(clustertokenv2.ClusterToken)
// 	err = vanilla.Stmt.Select(token.TableName(), token.ColumnNames(), eq_uuid.Parse(), nil, nil).
// 		QueryRow(tx)(func(scan vanilla.Scanner) (err error) {
// 		err = token.Scan(scan)
// 		return
// 	})

// 	return
// }

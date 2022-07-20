package vault

// import (
// 	"github.com/NexClipper/sudory/pkg/server/database"
// 	"github.com/NexClipper/sudory/pkg/server/database/prepare"
// 	"github.com/NexClipper/sudory/pkg/server/macro/logs"
// 	globvarv1 "github.com/NexClipper/sudory/pkg/server/model/global_variables/v1"
// 	"github.com/pkg/errors"
// 	"xorm.io/xorm"
// )

// type GlobalVariables struct {
// 	tx *xorm.Session
// }

// func NewGlobalVariables(ctx *xorm.Session) *GlobalVariables {
// 	return &GlobalVariables{tx: ctx}
// }

// func (vault GlobalVariables) Create(record globvarv1.GlobalVariables) (*globvarv1.GlobalVariables, error) {
// 	if err := database.XormCreate(
// 		vault.tx, &record); err != nil {
// 		return nil, errors.Wrapf(err, "create %v", record.TableName())
// 	}

// 	return &record, nil
// }

// func (vault GlobalVariables) Get(uuid string) (*globvarv1.GlobalVariables, error) {
// 	where := "uuid = ?"
// 	args := []interface{}{
// 		uuid,
// 	}
// 	record := &globvarv1.GlobalVariables{}
// 	if err := database.XormGet(
// 		vault.tx.Where(where, args...), record); err != nil {
// 		return nil, errors.Wrapf(err, "get %v", record.TableName())
// 	}

// 	return record, nil
// }

// func (vault GlobalVariables) Find(where string, args ...interface{}) ([]globvarv1.GlobalVariables, error) {
// 	models := make([]globvarv1.GlobalVariables, 0)
// 	if err := database.XormFind(
// 		vault.tx.Where(where, args...), &models); err != nil {
// 		return nil, errors.Wrapf(err, "find %v", new(globvarv1.GlobalVariables).TableName())
// 	}

// 	return models, nil
// }

// func (vault GlobalVariables) Query(query map[string]string) ([]globvarv1.GlobalVariables, error) {
// 	//parse query
// 	preparer, err := prepare.NewParser(query)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "prepare newParser%v",
// 			logs.KVL(
// 				"query", query,
// 			))
// 	}

// 	//find
// 	records := make([]globvarv1.GlobalVariables, 0)
// 	if err := database.XormFind(
// 		preparer.Prepared(vault.tx), &records); err != nil {
// 		return nil, errors.Wrapf(err, "query %v", new(globvarv1.GlobalVariables).TableName())
// 	}

// 	return records, nil
// }

// func (vault GlobalVariables) Update(record globvarv1.GlobalVariables) (*globvarv1.GlobalVariables, error) {
// 	where := "uuid = ?"
// 	args := []interface{}{
// 		record.Uuid,
// 	}
// 	if err := database.XormUpdate(
// 		vault.tx.Where(where, args...), &record); err != nil {
// 		return nil, errors.Wrapf(err, "update %v", record.TableName())
// 	}

// 	return &record, nil
// }

// func (vault GlobalVariables) Delete(uuid string) error {
// 	where := "uuid = ?"
// 	args := []interface{}{
// 		uuid,
// 	}
// 	record := &globvarv1.GlobalVariables{}
// 	if err := database.XormDelete(
// 		vault.tx.Where(where, args...), record); err != nil {
// 		return errors.Wrapf(err, "delete %v", record.TableName())
// 	}

// 	return nil
// }

package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	globvarv1 "github.com/NexClipper/sudory/pkg/server/model/global_variant/v1"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type GlobalVariant struct {
	tx *xorm.Session
}

func NewGlobalVariant(ctx *xorm.Session) *GlobalVariant {
	return &GlobalVariant{tx: ctx}
}

func (vault GlobalVariant) Create(record globvarv1.GlobalVariant) (*globvarv1.GlobalVariant, error) {
	if err := database.XormCreate(
		vault.tx, &record); err != nil {
		return nil, errors.Wrapf(err, "create %v", record.TableName())
	}

	return &record, nil
}

func (vault GlobalVariant) Get(uuid string) (*globvarv1.GlobalVariant, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	record := &globvarv1.GlobalVariant{}
	if err := database.XormGet(
		vault.tx.Where(where, args...), record); err != nil {
		return nil, errors.Wrapf(err, "get %v", record.TableName())
	}

	return record, nil
}

func (vault GlobalVariant) Find(where string, args ...interface{}) ([]globvarv1.GlobalVariant, error) {
	models := make([]globvarv1.GlobalVariant, 0)
	if err := database.XormFind(
		vault.tx.Where(where, args...), &models); err != nil {
		return nil, errors.Wrapf(err, "find %v", new(globvarv1.GlobalVariant).TableName())
	}

	return models, nil
}

func (vault GlobalVariant) Query(query map[string]string) ([]globvarv1.GlobalVariant, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser%v",
			logs.KVL(
				"query", query,
			))
	}

	//find
	records := make([]globvarv1.GlobalVariant, 0)
	if err := database.XormFind(
		preparer.Prepared(vault.tx), &records); err != nil {
		return nil, errors.Wrapf(err, "query %v", new(globvarv1.GlobalVariant).TableName())
	}

	return records, nil
}

func (vault GlobalVariant) Update(record globvarv1.GlobalVariant) (*globvarv1.GlobalVariant, error) {
	where := "uuid = ?"
	args := []interface{}{
		record.Uuid,
	}
	if err := database.XormUpdate(
		vault.tx.Where(where, args...), &record); err != nil {
		return nil, errors.Wrapf(err, "update %v", record.TableName())
	}

	return &record, nil
}

func (vault GlobalVariant) Delete(uuid string) error {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	record := &globvarv1.GlobalVariant{}
	if err := database.XormDelete(
		vault.tx.Where(where, args...), record); err != nil {
		return errors.Wrapf(err, "delete %v", record.TableName())
	}

	return nil
}

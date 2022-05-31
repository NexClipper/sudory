package vault

import (
	"fmt"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	channelv1 "github.com/NexClipper/sudory/pkg/server/model/channel/v1"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type NotifierStatus struct {
	tx *xorm.Session
}

func NewNotifierStatus(tx *xorm.Session) *NotifierStatus {
	return &NotifierStatus{tx: tx}
}

func (vault NotifierStatus) Create(model channelv1.NotifierStatus) (*channelv1.NotifierStatus, error) {
	if err := database.XormCreate(vault.tx, &model); err != nil {
		return nil, errors.Wrapf(err, "create %v", model.TableName())
	}

	return &model, nil
}

func (vault NotifierStatus) Get(uuid string) (*channelv1.NotifierStatus, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &channelv1.NotifierStatus{}
	if err := database.XormGet(
		vault.tx.Where(where, args...), model); err != nil {
		return nil, errors.Wrapf(err, "get %v", model.TableName())
	}

	return model, nil
}

func (vault NotifierStatus) Find(where string, args ...interface{}) ([]channelv1.NotifierStatus, error) {
	models := make([]channelv1.NotifierStatus, 0)
	if err := database.XormFind(
		vault.tx.Where(where, args...), &models); err != nil {
		return nil, errors.Wrapf(err, "find %v", new(channelv1.NotifierStatus).TableName())
	}

	return models, nil
}

func (vault NotifierStatus) Query(query map[string]string) ([]channelv1.NotifierStatus, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(channelv1.NotifierStatus).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]channelv1.NotifierStatus, 0)
	if err := database.XormFind(
		preparer.Prepared(vault.tx), &models); err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(channelv1.NotifierStatus).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	return models, nil
}

// func (vault EventNotifierStatus) Update(model eventv1.EventNotifierStatus) (*eventv1.EventNotifierStatus, error) {
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

func (vault NotifierStatus) Delete(uuid string) error {
	//delete event notifier status
	model := new(channelv1.NotifierStatus)
	if err := database.XormDelete(
		vault.tx.Where("uuid = ?", uuid), model); err != nil {
		return errors.Wrapf(err, "delete %v", model.TableName())
	}

	return nil
}

func (vault NotifierStatus) Rotate(notifier_uuid string, limit int) error {
	smt := func(s string, n int) []string {
		r := make([]string, n)
		for i := 0; i < n; i++ {
			r[i] = s
		}
		return r
	}

	model := new(channelv1.NotifierStatus)
	rows, err := vault.tx.Where("notifier_uuid = ?", notifier_uuid).Desc("created").Limit(2*limit, limit).Cols("uuid").Rows(model)
	if err != nil {
		return errors.Wrapf(err, "find rotate records %v", model.TableName())
	}
	defer rows.Close()

	// if true {

	args := make([]interface{}, 0, limit)
	for rows.Next() {
		model := channelv1.NotifierStatus{}
		// var uuid string
		if err := rows.Scan(&model); err != nil {
			return errors.Wrapf(err, "scan a record %v", model.TableName())
		}
		args = append(args, model.Uuid)
	}

	if 0 < len(args) {
		if _, err := vault.tx.Where(fmt.Sprintf("uuid IN (%s)", strings.Join(smt("?", len(args)), ",")), args...).Delete(model); err != nil {
			return errors.Wrapf(err, "delete margin %v", model.TableName())
		}
	}
	// }

	// if true {
	// 	args := make([]interface{}, 0, limit)
	// 	for rows.Next() {
	// 		var uuid string
	// 		if err := rows.Scan(&uuid); err != nil {
	// 			return errors.Wrapf(err, "scan a record %v", model.TableName())
	// 		}

	// 		r := eventv1.EventNotifierStatus{}
	// 		r.Uuid = uuid

	// 		args = append(args, &r)
	// 	}

	// 	if 0 < len(args) {
	// 		if _, err := vault.tx.Where(fmt.Sprintf("uuid IN (%s)", strings.Join(smt("?", len(args)), ",")), args...).Delete(model); err != nil {
	// 			return errors.Wrapf(err, "delete margin %v", model.TableName())
	// 		}
	// 	}
	// }

	return nil
}

func (vault NotifierStatus) CreateAndRotate(status channelv1.NotifierStatus, limit int) error {
	if _, err := vault.Create(status); err != nil {
		return errors.Wrapf(err, "create channel notifier status")
	}

	if err := vault.Rotate(status.NotifierUuid, limit); err != nil {
		return errors.Wrapf(err, "create channel notifier status")
	}

	return nil
}

package vault

import (
	"fmt"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	eventv1 "github.com/NexClipper/sudory/pkg/server/model/event/v1"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type EventNotifierStatus struct {
	tx *xorm.Session
}

func NewEventNotifierStatus(tx *xorm.Session) *EventNotifierStatus {
	return &EventNotifierStatus{tx: tx}
}

func (vault EventNotifierStatus) Create(model eventv1.EventNotifierStatus) (*eventv1.EventNotifierStatus, error) {
	if err := database.XormCreate(vault.tx, &model); err != nil {
		return nil, errors.Wrapf(err, "create %v", model.TableName())
	}

	return &model, nil
}

func (vault EventNotifierStatus) Get(uuid string) (*eventv1.EventNotifierStatus, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &eventv1.EventNotifierStatus{}
	if err := database.XormGet(
		vault.tx.Where(where, args...), model); err != nil {
		return nil, errors.Wrapf(err, "get %v", model.TableName())
	}

	return model, nil
}

func (vault EventNotifierStatus) Find(where string, args ...interface{}) ([]eventv1.EventNotifierStatus, error) {
	models := make([]eventv1.EventNotifierStatus, 0)
	if err := database.XormFind(
		vault.tx.Where(where, args...), &models); err != nil {
		return nil, errors.Wrapf(err, "find %v", new(eventv1.EventNotifierStatus).TableName())
	}

	return models, nil
}

func (vault EventNotifierStatus) Query(query map[string]string) ([]eventv1.EventNotifierStatus, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(eventv1.EventNotifierStatus).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]eventv1.EventNotifierStatus, 0)
	if err := database.XormFind(
		preparer.Prepared(vault.tx), &models); err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(eventv1.EventNotifierStatus).TableName(),
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

func (vault EventNotifierStatus) Delete(uuid string) error {
	//delete event notifier status
	model := new(eventv1.EventNotifierStatus)
	if err := database.XormDelete(
		vault.tx.Where("uuid = ?", uuid), model); err != nil {
		return errors.Wrapf(err, "delete %v", model.TableName())
	}

	return nil
}

func (vault EventNotifierStatus) Rotate(notifier_uuid string, limit int) error {
	model := new(eventv1.EventNotifierStatus)
	rows, err := vault.tx.Where("notifier_uuid = ?").OrderBy("created").Limit(2*limit, limit).Cols("uuid").Rows(model)
	if err != nil {
		return errors.Wrapf(err, "find rotate records %v", model.TableName())
	}

	if false {
		smt := func(s string, n int) []string {
			r := make([]string, n)
			for i := 0; i < n; i++ {
				r[i] = s
			}
			return r
		}
		args := make([]interface{}, 0, limit)
		for i := 0; i < limit && rows.Next(); i++ {
			var uuid string
			if err := rows.Scan(&uuid); err != nil {
				return errors.Wrapf(err, "scan a record %v", model.TableName())
			}
			args = append(args, uuid)
		}

		if 0 < len(args) {
			if _, err := vault.tx.Where(fmt.Sprintf("uuid IN (%s)", strings.Join(smt("?", len(args)), ",")), args...).Delete(model); err != nil {
				return errors.Wrapf(err, "delete margin %v", model.TableName())
			}
		}
	}

	if true {
		args := make([]interface{}, 0, limit)
		for i := 0; i < limit && rows.Next(); i++ {
			var uuid string
			if err := rows.Scan(&uuid); err != nil {
				return errors.Wrapf(err, "scan a record %v", model.TableName())
			}

			r := eventv1.EventNotifierStatus{}
			r.Uuid = uuid

			args = append(args, &r)
		}

		if 0 < len(args) {
			if _, err := vault.tx.Delete(args...); err != nil {
				return errors.Wrapf(err, "delete margin %v", model.TableName())
			}
		}
	}

	return nil
}

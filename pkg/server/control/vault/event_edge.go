package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	eventv1 "github.com/NexClipper/sudory/pkg/server/model/event/v1"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type EventNotifierEdge struct {
	tx *xorm.Session
}

func NewEventNotifierEdge(tx *xorm.Session) *EventNotifierEdge {
	return &EventNotifierEdge{tx: tx}
}

func (vault EventNotifierEdge) Create(model eventv1.EventNotifierEdge) (*eventv1.EventNotifierEdge, error) {
	if err := database.XormCreate(vault.tx, &model); err != nil {
		return nil, errors.Wrapf(err, "create %v", model.TableName())
	}

	return &model, nil
}

func (vault EventNotifierEdge) Get(event_uuid, notifier_uuid string) (*eventv1.EventNotifierEdge, error) {
	where := "event_uuid = ? AND notifier_uuid = ?"
	args := []interface{}{
		event_uuid, notifier_uuid,
	}
	model := &eventv1.EventNotifierEdge{}
	if err := database.XormGet(
		vault.tx.Where(where, args...), model); err != nil {
		return nil, errors.Wrapf(err, "get %v", model.TableName())
	}

	return model, nil
}

func (vault EventNotifierEdge) Find(where string, args ...interface{}) ([]eventv1.EventNotifierEdge, error) {
	models := make([]eventv1.EventNotifierEdge, 0)
	if err := database.XormFind(
		vault.tx.Where(where, args...), &models); err != nil {
		return nil, errors.Wrapf(err, "find %v", new(eventv1.EventNotifierEdge).TableName())
	}

	return models, nil
}

// func (vault EventNotifierEdge) Query(query map[string]string) ([]eventv1.EventNotifierEdge, error) {
// 	//parse query
// 	preparer, err := prepare.NewParser(query)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "query %v%v", new(eventv1.EventNotifierEdge).TableName(),
// 			logs.KVL(
// 				"query", query,
// 			))
// 	}

// 	//find service
// 	models := make([]eventv1.EventNotifierEdge, 0)
// 	if err := database.XormFind(
// 		preparer.Prepared(vault.tx), &models); err != nil {
// 		return nil, errors.Wrapf(err, "query %v%v", new(eventv1.EventNotifierEdge).TableName(),
// 			logs.KVL(
// 				"query", query,
// 			))
// 	}

// 	return models, nil
// }

func (vault EventNotifierEdge) Delete(event_uuid, notifier_uuid string) error {
	where := "event_uuid = ? AND notifier_uuid = ?"
	args := []interface{}{
		event_uuid, notifier_uuid,
	}
	model := &eventv1.EventNotifierEdge{}
	if err := database.XormDelete(
		vault.tx.Where(where, args...), model); err != nil {
		return errors.Wrapf(err, "delete %v", model.TableName())
	}

	return nil
}

package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	eventv1 "github.com/NexClipper/sudory/pkg/server/model/event/v1"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type EventNotifierWebhook struct {
	tx *xorm.Session
}

func NewEventNotifierWebhook(tx *xorm.Session) *EventNotifierWebhook {
	return &EventNotifierWebhook{tx: tx}
}

func (vault EventNotifierWebhook) Create(model eventv1.EventNotifierWebhook) (*eventv1.EventNotifierWebhook, error) {
	if err := database.XormCreate(vault.tx, &model); err != nil {
		return nil, errors.Wrapf(err, "create %v", model.TableName())
	}

	return &model, nil
}

func (vault EventNotifierWebhook) Get(uuid string) (*eventv1.EventNotifierWebhook, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &eventv1.EventNotifierWebhook{}
	if err := database.XormGet(
		vault.tx.Where(where, args...), model); err != nil {
		return nil, errors.Wrapf(err, "get %v", model.TableName())
	}

	return model, nil
}

func (vault EventNotifierWebhook) Find(where string, args ...interface{}) ([]eventv1.EventNotifierWebhook, error) {
	models := make([]eventv1.EventNotifierWebhook, 0)
	if err := database.XormFind(
		vault.tx.Where(where, args...), &models); err != nil {
		return nil, errors.Wrapf(err, "find %v", new(eventv1.EventNotifierWebhook).TableName())
	}

	return models, nil
}

func (vault EventNotifierWebhook) Query(query map[string]string) ([]eventv1.Event, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(eventv1.EventNotifierWebhook).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]eventv1.Event, 0)
	if err := database.XormFind(
		preparer.Prepared(vault.tx), &models); err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(eventv1.EventNotifierWebhook).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	return models, nil
}

func (vault EventNotifierWebhook) Update(model eventv1.EventNotifierWebhook) (*eventv1.EventNotifierWebhook, error) {
	where := "uuid = ?"
	args := []interface{}{
		model.Uuid,
	}

	if err := database.XormUpdate(
		vault.tx.Where(where, args...), &model); err != nil {
		return nil, errors.Wrapf(err, "update %v", model.TableName())
	}

	return &model, nil
}

func (vault EventNotifierWebhook) Delete(uuid string) error {
	//delete event notifier edge
	edge := new(eventv1.EventNotifierEdge)
	if err := database.XormDelete(
		vault.tx.Where("notifier_uuid = ?", uuid), edge); err != nil {
		return errors.Wrapf(err, "delete %v", edge.TableName())
	}

	//delete event notifier webhook
	model := new(eventv1.EventNotifierWebhook)
	if err := database.XormDelete(
		vault.tx.Where("uuid = ?", uuid), model); err != nil {
		return errors.Wrapf(err, "delete %v", model.TableName())
	}

	return nil
}

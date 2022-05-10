package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	eventv1 "github.com/NexClipper/sudory/pkg/server/model/event/v1"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type EventNotifierRabbitMq struct {
	tx *xorm.Session
}

func NewEventNotifierRabbitMq(tx *xorm.Session) *EventNotifierRabbitMq {
	return &EventNotifierRabbitMq{tx: tx}
}

func (vault EventNotifierRabbitMq) Create(model eventv1.EventNotifierRabbitMq) (*eventv1.EventNotifierRabbitMq, error) {
	if err := database.XormCreate(vault.tx, &model); err != nil {
		return nil, errors.Wrapf(err, "create %v", model.TableName())
	}

	return &model, nil
}

func (vault EventNotifierRabbitMq) Get(uuid string) (*eventv1.EventNotifierRabbitMq, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &eventv1.EventNotifierRabbitMq{}
	if err := database.XormGet(
		vault.tx.Where(where, args...), model); err != nil {
		return nil, errors.Wrapf(err, "get %v", model.TableName())
	}

	return model, nil
}

func (vault EventNotifierRabbitMq) Find(where string, args ...interface{}) ([]eventv1.EventNotifierRabbitMq, error) {
	models := make([]eventv1.EventNotifierRabbitMq, 0)
	if err := database.XormFind(
		vault.tx.Where(where, args...), &models); err != nil {
		return nil, errors.Wrapf(err, "find %v", new(eventv1.EventNotifierRabbitMq).TableName())
	}

	return models, nil
}

func (vault EventNotifierRabbitMq) Query(query map[string]string) ([]eventv1.EventNotifierRabbitMq, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(eventv1.EventNotifierRabbitMq).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]eventv1.EventNotifierRabbitMq, 0)
	if err := database.XormFind(
		preparer.Prepared(vault.tx), &models); err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(eventv1.EventNotifierRabbitMq).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	return models, nil
}

func (vault EventNotifierRabbitMq) Update(model eventv1.EventNotifierRabbitMq) (*eventv1.EventNotifierRabbitMq, error) {
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

func (vault EventNotifierRabbitMq) Delete(uuid string) error {
	//delete event notifier edge
	edge := new(eventv1.EventNotifierEdge)
	if err := database.XormDelete(
		vault.tx.Where("notifier_uuid = ?", uuid), edge); err != nil {
		return errors.Wrapf(err, "delete %v", edge.TableName())
	}

	//delete event notifier rabbitmq
	model := &eventv1.EventNotifierRabbitMq{}
	if err := database.XormDelete(
		vault.tx.Where("uuid = ?", uuid), model); err != nil {
		return errors.Wrapf(err, "delete %v", model.TableName())
	}

	return nil
}

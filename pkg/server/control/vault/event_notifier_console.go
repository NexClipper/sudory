package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	eventv1 "github.com/NexClipper/sudory/pkg/server/model/event/v1"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type EventNotifierConsole struct {
	tx *xorm.Session
}

func NewEventNotifierConsole(tx *xorm.Session) *EventNotifierConsole {
	return &EventNotifierConsole{tx: tx}
}

func (vault EventNotifierConsole) Create(model eventv1.EventNotifierConsole) (*eventv1.EventNotifierConsole, error) {
	if err := database.XormCreate(vault.tx, &model); err != nil {
		return nil, errors.Wrapf(err, "create %v", model.TableName())
	}

	return &model, nil
}

func (vault EventNotifierConsole) Get(uuid string) (*eventv1.EventNotifierConsole, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &eventv1.EventNotifierConsole{}
	if err := database.XormGet(
		vault.tx.Where(where, args...), model); err != nil {
		return nil, errors.Wrapf(err, "get %v", model.TableName())
	}

	return model, nil
}

func (vault EventNotifierConsole) Find(where string, args ...interface{}) ([]eventv1.EventNotifierConsole, error) {
	models := make([]eventv1.EventNotifierConsole, 0)
	if err := database.XormFind(
		vault.tx.Where(where, args...), &models); err != nil {
		return nil, errors.Wrapf(err, "find %v", new(eventv1.EventNotifierRabbitMq).TableName())
	}

	return models, nil
}

func (vault EventNotifierConsole) Query(query map[string]string) ([]eventv1.EventNotifierConsole, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(eventv1.EventNotifierConsole).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]eventv1.EventNotifierConsole, 0)
	if err := database.XormFind(
		preparer.Prepared(vault.tx), &models); err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(eventv1.EventNotifierConsole).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	return models, nil
}

func (vault EventNotifierConsole) Update(model eventv1.EventNotifierConsole) (*eventv1.EventNotifierConsole, error) {
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

func (vault EventNotifierConsole) Delete(uuid string) error {
	//delete event notifier edge
	edge := new(eventv1.EventNotifierEdge)
	if err := database.XormDelete(
		vault.tx.Where("notifier_uuid = ?", uuid), edge); err != nil {
		return errors.Wrapf(err, "delete %v", edge.TableName())
	}

	//delete event notifier console
	model := &eventv1.EventNotifierConsole{}
	if err := database.XormDelete(
		vault.tx.Where("uuid = ?", uuid), model); err != nil {
		return errors.Wrapf(err, "delete %v", model.TableName())
	}

	return nil
}

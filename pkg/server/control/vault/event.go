package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	eventv1 "github.com/NexClipper/sudory/pkg/server/model/event/v1"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type Event struct {
	tx *xorm.Session
}

func NewEvent(tx *xorm.Session) *Event {
	return &Event{tx: tx}
}

func (vault Event) Create(model eventv1.Event) (*eventv1.Event, error) {
	if err := database.XormCreate(vault.tx, &model); err != nil {
		return nil, errors.Wrapf(err, "create %v", model.TableName())
	}

	return &model, nil
}

func (vault Event) Get(uuid string) (*eventv1.Event, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &eventv1.Event{}
	if err := database.XormGet(
		vault.tx.Where(where, args...), model); err != nil {
		return nil, errors.Wrapf(err, "get %v", model.TableName())
	}

	return model, nil
}

func (vault Event) Find(where string, args ...interface{}) ([]eventv1.Event, error) {
	models := make([]eventv1.Event, 0)
	if err := database.XormFind(
		vault.tx.Where(where, args...), &models); err != nil {
		return nil, errors.Wrapf(err, "find %v", new(eventv1.Event).TableName())
	}

	return models, nil
}

func (vault Event) Query(query map[string]string) ([]eventv1.Event, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(eventv1.Event).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]eventv1.Event, 0)
	if err := database.XormFind(
		preparer.Prepared(vault.tx), &models); err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(eventv1.Event).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	return models, nil
}

func (vault Event) Update(model eventv1.Event) (*eventv1.Event, error) {
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

func (vault Event) Delete(uuid string) error {
	//delete event notifier edge
	edge := new(eventv1.EventNotifierEdge)
	if err := database.XormDelete(
		vault.tx.Where("event_uuid = ?", uuid), edge); err != nil {
		return errors.Wrapf(err, "delete %v", edge.TableName())
	}

	//delete event
	event := new(eventv1.Event)
	if err := database.XormDelete(
		vault.tx.Where("uuid = ?", uuid), event); err != nil {
		return errors.Wrapf(err, "delete %v", event.TableName())
	}

	return nil
}

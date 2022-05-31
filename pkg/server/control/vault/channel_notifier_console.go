package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	channelv1 "github.com/NexClipper/sudory/pkg/server/model/channel/v1"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type NotifierConsole struct {
	tx *xorm.Session
}

func NewNotifierConsole(tx *xorm.Session) *NotifierConsole {
	return &NotifierConsole{tx: tx}
}

func (vault NotifierConsole) Create(model channelv1.NotifierConsole) (*channelv1.NotifierConsole, error) {
	if err := database.XormCreate(vault.tx, &model); err != nil {
		return nil, errors.Wrapf(err, "create %v", model.TableName())
	}

	return &model, nil
}

func (vault NotifierConsole) Get(uuid string) (*channelv1.NotifierConsole, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &channelv1.NotifierConsole{}
	if err := database.XormGet(
		vault.tx.Where(where, args...), model); err != nil {
		return nil, errors.Wrapf(err, "get %v", model.TableName())
	}

	return model, nil
}

func (vault NotifierConsole) Find(where string, args ...interface{}) ([]channelv1.NotifierConsole, error) {
	models := make([]channelv1.NotifierConsole, 0)
	if err := database.XormFind(
		vault.tx.Where(where, args...), &models); err != nil {
		return nil, errors.Wrapf(err, "find %v", new(channelv1.NotifierRabbitMq).TableName())
	}

	return models, nil
}

func (vault NotifierConsole) Query(query map[string]string) ([]channelv1.NotifierConsole, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(channelv1.NotifierConsole).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]channelv1.NotifierConsole, 0)
	if err := database.XormFind(
		preparer.Prepared(vault.tx), &models); err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(channelv1.NotifierConsole).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	return models, nil
}

func (vault NotifierConsole) Update(model channelv1.NotifierConsole) (*channelv1.NotifierConsole, error) {
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

func (vault NotifierConsole) Delete(uuid string) error {
	//delete event notifier edge
	edge := new(channelv1.ChannelNotifierEdge)
	if err := database.XormDelete(
		vault.tx.Where("notifier_type = ? AND notifier_uuid = ?", channelv1.NotifierTypeConsole.String(), uuid), edge); err != nil {
		return errors.Wrapf(err, "delete %v", edge.TableName())
	}

	//delete event notifier console
	model := &channelv1.NotifierConsole{}
	if err := database.XormDelete(
		vault.tx.Where("uuid = ?", uuid), model); err != nil {
		return errors.Wrapf(err, "delete %v", model.TableName())
	}

	return nil
}

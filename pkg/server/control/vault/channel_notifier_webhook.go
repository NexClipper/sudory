package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	channelv1 "github.com/NexClipper/sudory/pkg/server/model/channel/v1"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type NotifierWebhook struct {
	tx *xorm.Session
}

func NewNotifierWebhook(tx *xorm.Session) *NotifierWebhook {
	return &NotifierWebhook{tx: tx}
}

func (vault NotifierWebhook) Create(model channelv1.NotifierWebhook) (*channelv1.NotifierWebhook, error) {
	if err := database.XormCreate(vault.tx, &model); err != nil {
		return nil, errors.Wrapf(err, "create %v", model.TableName())
	}

	return &model, nil
}

func (vault NotifierWebhook) Get(uuid string) (*channelv1.NotifierWebhook, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &channelv1.NotifierWebhook{}
	if err := database.XormGet(
		vault.tx.Where(where, args...), model); err != nil {
		return nil, errors.Wrapf(err, "get %v", model.TableName())
	}

	return model, nil
}

func (vault NotifierWebhook) Find(where string, args ...interface{}) ([]channelv1.NotifierWebhook, error) {
	models := make([]channelv1.NotifierWebhook, 0)
	if err := database.XormFind(
		vault.tx.Where(where, args...), &models); err != nil {
		return nil, errors.Wrapf(err, "find %v", new(channelv1.NotifierWebhook).TableName())
	}

	return models, nil
}

func (vault NotifierWebhook) Query(query map[string]string) ([]channelv1.NotifierWebhook, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(channelv1.NotifierWebhook).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]channelv1.NotifierWebhook, 0)
	if err := database.XormFind(
		preparer.Prepared(vault.tx), &models); err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(channelv1.NotifierWebhook).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	return models, nil
}

func (vault NotifierWebhook) Update(model channelv1.NotifierWebhook) (*channelv1.NotifierWebhook, error) {
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

func (vault NotifierWebhook) Delete(uuid string) error {
	//delete event notifier edge
	edge := new(channelv1.ChannelNotifierEdge)
	if err := database.XormDelete(
		vault.tx.Where("notifier_type = ? AND notifier_uuid = ?", channelv1.NotifierTypeWebhook.String(), uuid), edge); err != nil {
		return errors.Wrapf(err, "delete %v", edge.TableName())
	}

	//delete event notifier webhook
	model := new(channelv1.NotifierWebhook)
	if err := database.XormDelete(
		vault.tx.Where("uuid = ?", uuid), model); err != nil {
		return errors.Wrapf(err, "delete %v", model.TableName())
	}

	return nil
}

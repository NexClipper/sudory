package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	channelv1 "github.com/NexClipper/sudory/pkg/server/model/channel/v1"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type Channel struct {
	tx *xorm.Session
}

func NewChannel(tx *xorm.Session) *Channel {
	return &Channel{tx: tx}
}

func (vault Channel) Create(model channelv1.Channel) (*channelv1.Channel, error) {
	if err := database.XormCreate(vault.tx, &model); err != nil {
		return nil, errors.Wrapf(err, "create %v", model.TableName())
	}

	return &model, nil
}

func (vault Channel) Get(uuid string) (*channelv1.Channel, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &channelv1.Channel{}
	if err := database.XormGet(
		vault.tx.Where(where, args...), model); err != nil {
		return nil, errors.Wrapf(err, "get %v", model.TableName())
	}

	return model, nil
}

func (vault Channel) Find(where string, args ...interface{}) ([]channelv1.Channel, error) {
	models := make([]channelv1.Channel, 0)
	if err := database.XormFind(
		vault.tx.Where(where, args...), &models); err != nil {
		return nil, errors.Wrapf(err, "find %v", new(channelv1.Channel).TableName())
	}

	return models, nil
}

func (vault Channel) FindAll() ([]channelv1.Channel, error) {
	models := make([]channelv1.Channel, 0)
	if err := database.XormFind(
		vault.tx, &models); err != nil {
		return nil, errors.Wrapf(err, "find all %v", new(channelv1.Channel).TableName())
	}

	return models, nil
}

func (vault Channel) Query(query map[string]string) ([]channelv1.Channel, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(channelv1.Channel).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]channelv1.Channel, 0)
	if err := database.XormFind(
		preparer.Prepared(vault.tx), &models); err != nil {
		return nil, errors.Wrapf(err, "query %v%v", new(channelv1.Channel).TableName(),
			logs.KVL(
				"query", query,
			))
	}

	return models, nil
}

func (vault Channel) Update(model channelv1.Channel) (*channelv1.Channel, error) {
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

func (vault Channel) Delete(uuid string) error {
	//delete event notifier edge
	edge := new(channelv1.ChannelNotifierEdge)
	if err := database.XormDelete(
		vault.tx.Where("channel_uuid = ?", uuid), edge); err != nil {
		return errors.Wrapf(err, "delete %v", edge.TableName())
	}

	//delete event
	event := new(channelv1.Channel)
	if err := database.XormDelete(
		vault.tx.Where("uuid = ?", uuid), event); err != nil {
		return errors.Wrapf(err, "delete %v", event.TableName())
	}

	return nil
}

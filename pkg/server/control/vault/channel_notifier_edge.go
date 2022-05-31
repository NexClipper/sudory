package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	channelv1 "github.com/NexClipper/sudory/pkg/server/model/channel/v1"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type ChannelNotifierEdge struct {
	tx *xorm.Session
}

func NewChannelNotifierEdge(tx *xorm.Session) *ChannelNotifierEdge {
	return &ChannelNotifierEdge{tx: tx}
}

func (vault ChannelNotifierEdge) Create(model channelv1.ChannelNotifierEdge) (*channelv1.ChannelNotifierEdge, error) {
	if err := database.XormCreate(vault.tx, &model); err != nil {
		return nil, errors.Wrapf(err, "create %v", model.TableName())
	}

	return &model, nil
}

func (vault ChannelNotifierEdge) Get(channel_uuid, notifier_uuid string) (*channelv1.ChannelNotifierEdge, error) {
	where := "channel_uuid = ? AND notifier_uuid = ?"
	args := []interface{}{
		channel_uuid, notifier_uuid,
	}
	model := &channelv1.ChannelNotifierEdge{}
	if err := database.XormGet(
		vault.tx.Where(where, args...), model); err != nil {
		return nil, errors.Wrapf(err, "get %v", model.TableName())
	}

	return model, nil
}

func (vault ChannelNotifierEdge) Find(where string, args ...interface{}) ([]channelv1.ChannelNotifierEdge, error) {
	models := make([]channelv1.ChannelNotifierEdge, 0)
	if err := database.XormFind(
		vault.tx.Where(where, args...), &models); err != nil {
		return nil, errors.Wrapf(err, "find %v", new(channelv1.ChannelNotifierEdge).TableName())
	}

	return models, nil
}

func (vault ChannelNotifierEdge) Delete(channel_uuid, notifier_uuid string) error {
	where := "channel_uuid = ? AND notifier_uuid = ?"
	args := []interface{}{
		channel_uuid, notifier_uuid,
	}
	model := &channelv1.ChannelNotifierEdge{}
	if err := database.XormDelete(
		vault.tx.Where(where, args...), model); err != nil {
		return errors.Wrapf(err, "delete %v", model.TableName())
	}

	return nil
}

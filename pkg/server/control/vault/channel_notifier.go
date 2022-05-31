package vault

import (
	channelv1 "github.com/NexClipper/sudory/pkg/server/model/channel/v1"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type ChannelNotifier struct {
	tx *xorm.Session
}

func NewChannelNotifier(tx *xorm.Session) *ChannelNotifier {
	return &ChannelNotifier{tx: tx}
}

func (vault ChannelNotifier) Get(notifier_type, notifier_uuid string) (interface{}, error) {
	type_, err := channelv1.ParseNotifierType(notifier_type)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid notifier types")
	}

	switch type_ {
	case channelv1.NotifierTypeConsole:
		notifier, err := NewNotifierConsole(vault.tx).Get(notifier_uuid)
		if err != nil {
			err = errors.Wrapf(err, "get console notifier")
		}
		return notifier, err
	case channelv1.NotifierTypeWebhook:
		notifier, err := NewNotifierWebhook(vault.tx).Get(notifier_uuid)
		if err != nil {
			err = errors.Wrapf(err, "get webhook notifier")
		}
		return notifier, err
	case channelv1.NotifierTypeRabbitmq:
		notifier, err := NewNotifierRabbitMq(vault.tx).Get(notifier_uuid)
		if err != nil {
			err = errors.Wrapf(err, "get rabbitmq notifier")
		}
		return notifier, err
	}

	return nil, errors.Errorf("invalid notifier types")
}

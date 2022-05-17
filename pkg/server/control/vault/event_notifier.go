package vault

import (
	eventv1 "github.com/NexClipper/sudory/pkg/server/model/event/v1"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type EventNotifier struct {
	tx *xorm.Session
}

func NewEventNotifier(tx *xorm.Session) *EventNotifier {
	return &EventNotifier{tx: tx}
}

func (vault EventNotifier) Get(notifier_type, notifier_uuid string) (map[string]interface{}, error) {
	type_, err := eventv1.ParseEventNotifierType(notifier_type)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid notifier types")
	}

	switch type_ {
	case eventv1.EventNotifierTypeConsole:
		notifier, err := NewEventNotifierConsole(vault.tx).Get(notifier_uuid)
		if err != nil {
			err = errors.Wrapf(err, "get console notifier")
		}
		return map[string]interface{}{notifier_type: notifier}, err
	case eventv1.EventNotifierTypeWebhook:
		notifier, err := NewEventNotifierWebhook(vault.tx).Get(notifier_uuid)
		if err != nil {
			err = errors.Wrapf(err, "get webhook notifier")
		}
		return map[string]interface{}{notifier_type: notifier}, err
	case eventv1.EventNotifierTypeRabbitmq:
		notifier, err := NewEventNotifierRabbitMq(vault.tx).Get(notifier_uuid)
		if err != nil {
			err = errors.Wrapf(err, "get rabbitmq notifier")
		}
		return map[string]interface{}{notifier_type: notifier}, err
	}

	return nil, errors.Errorf("invalid notifier types")
}

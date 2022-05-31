package managed_event

import (
	"fmt"
	"time"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
	channelv1 "github.com/NexClipper/sudory/pkg/server/model/channel/v1"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

var _ Notifier = (*RabbitMQNotifier)(nil)

type RabbitMQNotifier struct {
	opt *channelv1.NotifierRabbitMq

	connection *amqp.Connection //RabbitMQ //amqp.Connection
	channel    *amqp.Channel    //RabbitMQ //amqp.Channel
}

func NewRabbitMqNotifier(opt *channelv1.NotifierRabbitMq) (*RabbitMQNotifier, error) {
	notifier := &RabbitMQNotifier{}
	notifier.opt = opt

	return notifier, nil
}
func (notifier RabbitMQNotifier) Type() fmt.Stringer {
	return notifier.opt.Type()
}
func (notifier RabbitMQNotifier) Uuid() string {
	return notifier.opt.Uuid
}

func (notifier RabbitMQNotifier) Property() map[string]string {
	return map[string]string{
		"type":        notifier.opt.Type().String(),
		"uuid":        notifier.opt.Uuid,
		"url":         notifier.opt.Url,
		"exchange":    nullable.String(notifier.opt.Exchange).Value(),
		"routing-key": nullable.String(notifier.opt.RoutingKey).Value(),
	}
}

func (notifier *RabbitMQNotifier) Close() {
	//disconnect rabbitmq
	var established bool = !(notifier.connection == nil || notifier.connection.IsClosed())
	if established {
		notifier.connection.Close()
	}
}

func (notifier *RabbitMQNotifier) OnNotify(factory MarshalFactoryResult) error {
	var established bool = !(notifier.connection == nil || notifier.connection.IsClosed())
	if !established {
		conn, ch, err := notifier.Dial(notifier.opt.Url)
		if err != nil {
			return errors.Wrapf(err, "dial to rabbimq%s",
				logs.KVL(
					"url", notifier.opt.Url,
				))
		}
		notifier.connection = conn
		notifier.channel = ch
	}

	opt := notifier.opt
	ch := notifier.channel

	opt.MessageContentType = newist.String(notifier.opt.ContentType)
	b, err := factory(notifier.opt.ContentType)
	if err != nil {
		return errors.Wrapf(err, "marshal factory")
	}
	for _, b := range b {
		if err := notifier.Publish(*opt, ch, b); err != nil {
			return errors.Wrapf(err, "publish to rabbimq%s",
				logs.KVL(
					"opt", opt,
				))
		}
	}

	return nil
}

func (RabbitMQNotifier) Dial(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "dial to amqp url=%s", url)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "create rabbitmq channel")
	}

	return conn, ch, nil
}

func (RabbitMQNotifier) Publish(opt channelv1.NotifierRabbitMq, ch *amqp.Channel, b []byte) error {
	publishing := amqp.Publishing{}
	publishing.ContentType = nullable.String(opt.MessageContentType).Value()
	publishing.ContentEncoding = nullable.String(opt.MessageContentEncoding).Value()
	publishing.Headers = opt.MessageHeaders
	publishing.DeliveryMode = nullable.Uint8(opt.MessageDeliveryMode).Value()
	publishing.Priority = nullable.Uint8(opt.MessagePriority).Value()
	publishing.CorrelationId = nullable.String(opt.MessageCorrelationId).Value()
	publishing.ReplyTo = nullable.String(opt.MessageReplyTo).Value()
	publishing.Expiration = nullable.String(opt.MessageExpiration).Value()
	publishing.MessageId = nullable.String(opt.MessageMessageId).Value()
	if nullable.Bool(opt.MessageTimestamp).Value() {
		publishing.Timestamp = time.Now()
	}
	publishing.Type = nullable.String(opt.MessageType).Value()
	publishing.UserId = nullable.String(opt.MessageUserId).Value()
	publishing.AppId = nullable.String(opt.MessageAppId).Value()
	publishing.Body = b

	if err := ch.Publish(
		nullable.String(opt.Exchange).Value(),
		nullable.String(opt.RoutingKey).Value(),
		nullable.Bool(opt.Mandatory).Value(),
		nullable.Bool(opt.Immediate).Value(),
		publishing,
	); err != nil {
		return errors.Wrapf(err, "publish to rabbitmq%s",
			logs.KVL("exchange", opt.Exchange,
				"routing_key", opt.RoutingKey,
				"mandatory", opt.Mandatory,
				"immediate", opt.Immediate,
				"publishing", publishing,
			))
	}

	return nil
}

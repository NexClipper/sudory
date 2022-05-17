package event

import (
	"fmt"
	"time"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type RabbitMQNotifier struct {
	opt RabbitMQNotifierConfig
	sub EventNotifierMultiplexer

	connection *amqp.Connection //RabbitMQ //amqp.Connection
	channel    *amqp.Channel    //RabbitMQ //amqp.Channel
}

func NewRabbitMqNotifier(opt RabbitMQNotifierConfig) (*RabbitMQNotifier, error) {
	conn, ch, err := new(RabbitMQNotifier).Dial(opt.Url)
	if err != nil {
		return nil, errors.Wrapf(err, "dial to rabbimq%s",
			logs.KVL(
				"url", opt.Url,
			))
	}

	notifier := &RabbitMQNotifier{}
	notifier.opt = opt
	notifier.connection = conn
	notifier.channel = ch

	return notifier, nil
}
func (notifier RabbitMQNotifier) Type() fmt.Stringer {
	return NotifierTypeRabbitMQ
}

func (notifier RabbitMQNotifier) Property() map[string]string {
	return map[string]string{
		"name":        notifier.sub.(EventNotifiMuxConfigHolder).Config().Name,
		"type":        notifier.Type().String(),
		"url":         notifier.opt.Url,
		"exchange":    notifier.opt.Exchange,
		"routing-key": notifier.opt.RoutingKey,
	}
}

func (notifier *RabbitMQNotifier) Regist(sub EventNotifierMultiplexer) {
	//Subscribe
	if !(sub == nil && notifier.sub != nil) {
		notifier.sub = sub
		notifier.sub.Notifiers().Add(notifier)
	}
}

func (notifier *RabbitMQNotifier) Close() {
	//Unsubscribe
	if notifier.sub != nil {
		notifier.sub.Notifiers().Remove(notifier)
		notifier.sub = nil
	}

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

	opt.MessageContentType = "application/json"
	b, err := factory("application/json")
	if err != nil {
		return errors.Wrapf(err, "marshal factory")
	}
	for _, b := range b {
		if err := notifier.Publish(opt, ch, b); err != nil {
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

func (RabbitMQNotifier) Publish(opt RabbitMQNotifierConfig, ch *amqp.Channel, b []byte) error {
	publishing := amqp.Publishing{}
	publishing.ContentType = opt.MessageContentType
	publishing.ContentEncoding = opt.MessageContentEncoding
	publishing.Headers = opt.MessageHeaders
	publishing.DeliveryMode = opt.MessageDeliveryMode
	publishing.Priority = opt.MessagePriority
	publishing.CorrelationId = opt.MessageCorrelationId
	publishing.ReplyTo = opt.MessageReplyTo
	publishing.Expiration = opt.MessageExpiration
	publishing.MessageId = opt.MessageMessageId
	if opt.MessageTimestamp {
		publishing.Timestamp = time.Now()
	}
	publishing.Type = opt.Type
	publishing.UserId = opt.MessageUserId
	publishing.AppId = opt.MessageAppId
	publishing.Body = b

	if err := ch.Publish(opt.Exchange, opt.RoutingKey, opt.Mandatory, opt.Immediate, publishing); err != nil {
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

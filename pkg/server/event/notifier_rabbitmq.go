package event

import (
	"bytes"
	"strconv"
	"time"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type rabbitMQNotifier struct {
	opt RabbitMQNotifierConfig
	sub EventSubscriber

	connection *amqp.Connection //RabbitMQ //amqp.Connection
	channel    *amqp.Channel    //RabbitMQ //amqp.Channel
}

func NewRabbitMqNotifier(opt RabbitMQNotifierConfig) (*rabbitMQNotifier, error) {
	conn, ch, err := new(rabbitMQNotifier).dial(opt.Url)
	if err != nil {
		return nil, errors.Wrapf(err, "dial to rabbimq%s",
			logs.KVL(
				"url", opt.Url,
			))
	}

	notifier := &rabbitMQNotifier{}
	notifier.opt = opt
	notifier.connection = conn
	notifier.channel = ch

	return notifier, nil
}
func (notifier rabbitMQNotifier) Type() string {
	return NotifierTypeRabbitMQ.String()
}

func (notifier rabbitMQNotifier) Property() map[string]string {
	return map[string]string{
		"name":        notifier.sub.Config().Name,
		"type":        notifier.Type(),
		"url":         notifier.opt.Url,
		"exchange":    notifier.opt.Exchange,
		"routing-key": notifier.opt.RoutingKey,
	}
}

func (notifier rabbitMQNotifier) PropertyString() string {
	buff := bytes.Buffer{}
	for key, value := range notifier.Property() {
		if 0 < buff.Len() {
			buff.WriteString(" ")
		}
		buff.WriteString(key)
		buff.WriteString("=")
		buff.WriteString(strconv.Quote(value))
	}
	return buff.String()
}

func (notifier *rabbitMQNotifier) Regist(sub EventSubscriber) {
	//Subscribe
	if !(sub == nil && notifier.sub != nil) {
		notifier.sub = sub
		notifier.sub.Notifiers().Add(notifier)
	}
}

func (notifier *rabbitMQNotifier) Close() {
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

func (notifier *rabbitMQNotifier) OnNotify(factory MarshalFactory) error {
	var established bool = !(notifier.connection == nil || notifier.connection.IsClosed())
	if !established {
		conn, ch, err := notifier.dial(notifier.opt.Url)
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
		if err := notifier.publish(opt, ch, b); err != nil {
			return errors.Wrapf(err, "publish to rabbimq%s",
				logs.KVL(
					"opt", opt,
				))
		}
	}

	return nil
}

func (notifier rabbitMQNotifier) OnNotifyAsync(factory MarshalFactory) <-chan NotifierFuture {
	future := make(chan NotifierFuture)
	go func() {
		defer close(future)

		future <- NotifierFuture{&notifier, notifier.OnNotify(factory)}
	}()

	return future
}

func (rabbitMQNotifier) dial(url string) (*amqp.Connection, *amqp.Channel, error) {
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

func (rabbitMQNotifier) publish(opt RabbitMQNotifierConfig, ch *amqp.Channel, b []byte) error {
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
				"pouting_key", opt.RoutingKey,
				"mandatory", opt.Mandatory,
				"immediate", opt.Immediate,
				"publishing", publishing,
			))
	}

	return nil
}

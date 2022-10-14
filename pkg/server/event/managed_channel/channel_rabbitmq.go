package managed_channel

import (
	"fmt"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	channelv2 "github.com/NexClipper/sudory/pkg/server/model/channel/v3"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type ChannelRabbitMQ struct {
	uuid string
	opt  *channelv2.RabbitMqConfig

	connection *amqp.Connection //RabbitMQ //amqp.Connection
	channel    *amqp.Channel    //RabbitMQ //amqp.Channel
}

func NewChannelRabbitMQ(uuid string, opt *channelv2.RabbitMqConfig) *ChannelRabbitMQ {
	notifier := &ChannelRabbitMQ{}
	notifier.uuid = uuid
	notifier.opt = opt

	return notifier
}
func (channel ChannelRabbitMQ) Type() fmt.Stringer {
	return channel.opt.Type()
}
func (channel ChannelRabbitMQ) Uuid() string {
	return channel.uuid
}

func (channel ChannelRabbitMQ) Property() map[string]string {
	return map[string]string{
		"type":        channel.opt.Type().String(),
		"uuid":        channel.uuid,
		"url":         channel.opt.Url,
		"exchange":    channel.opt.ChannelPublish.Exchange.String,
		"routing-key": channel.opt.ChannelPublish.RoutingKey.String,
	}
}

func (channel *ChannelRabbitMQ) Close() {
	//disconnect rabbitmq
	var established bool = !(channel.connection == nil || channel.connection.IsClosed())
	if established {
		channel.connection.Close()
	}
}

func (channel *ChannelRabbitMQ) OnNotify(factory *MarshalFactory) error {
	const content_type = "application/json"

	var established bool = !(channel.connection == nil || channel.connection.IsClosed())
	if !established {
		conn, ch, err := channel.Dial(channel.opt.Url)
		if err != nil {
			return errors.Wrapf(err, "dial to rabbimq%s",
				logs.KVL(
					"url", channel.opt.Url,
				))
		}
		channel.connection = conn
		channel.channel = ch
	}
	opt := channel.opt
	ch := channel.channel

	opt.Publishing.MessageContentType = *vanilla.NewNullString(content_type)
	b, err := factory.Marshal(content_type)
	if err != nil {
		return errors.Wrapf(err, "marshal factory")
	}

	if err := channel.Publish(opt, ch, b); err != nil {
		return errors.Wrapf(err, "publish to rabbimq%s",
			logs.KVL(
				"opt", opt,
			))
	}

	return nil
}

func (ChannelRabbitMQ) Dial(url string) (*amqp.Connection, *amqp.Channel, error) {
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

func (ChannelRabbitMQ) Publish(opt *channelv2.RabbitMqConfig, ch *amqp.Channel, b []byte) error {
	publishing := amqp.Publishing{}
	publishing.ContentType = opt.Publishing.MessageContentType.String
	publishing.ContentEncoding = opt.Publishing.MessageContentEncoding.String
	publishing.Headers = opt.Publishing.MessageHeaders.Object
	publishing.DeliveryMode = opt.Publishing.MessageDeliveryMode.Byte
	publishing.Priority = opt.Publishing.MessagePriority.Byte
	publishing.CorrelationId = opt.Publishing.MessageCorrelationId.String
	publishing.ReplyTo = opt.Publishing.MessageReplyTo.String
	publishing.Expiration = opt.Publishing.MessageExpiration.String
	publishing.MessageId = opt.Publishing.MessageMessageId.String
	if opt.Publishing.MessageTimestamp.Bool {
		publishing.Timestamp = time.Now()
	}
	publishing.Type = opt.Publishing.MessageType.String
	publishing.UserId = opt.Publishing.MessageUserId.String
	publishing.AppId = opt.Publishing.MessageAppId.String
	publishing.Body = b

	if err := ch.Publish(
		opt.ChannelPublish.Exchange.String,
		opt.ChannelPublish.RoutingKey.String,
		opt.ChannelPublish.Mandatory.Bool,
		opt.ChannelPublish.Immediate.Bool,
		publishing,
	); err != nil {
		return errors.Wrapf(err, "publish to rabbitmq%s",
			logs.KVL(
				"options", opt,
			))
	}

	return nil
}

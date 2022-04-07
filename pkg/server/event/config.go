package event

import (
	"time"
)

type EventConfig struct {
	EventSubscribeConfigs []EventSubscribeConfig `yaml:"events,omitempty"`
}

type EventSubscribeConfig struct {
	Name            string           `yaml:"name"`
	UpdateInterval  time.Duration    `yaml:"update-interval"`
	NotifierConfigs []NotifierConfig `yaml:"notifiers,omitempty"`
}

type NotifierConfig map[string]interface{}

type ConsoleNotifierConfig struct{}

type FileNotifierConfig struct {
	Type string `yaml:"type"` //Notifier Type

	//file
	Path string `yaml:"path"`
}

type WebhookNotifierConfig struct {
	Type string `yaml:"type"` //Notifier Type

	//http
	Method         string            `yaml:"method"`
	Url            string            `yaml:"url"`
	RequestHeaders map[string]string `yaml:"request-headers,omitempty"`

	//for timeout context
	RequestTimeout time.Duration `yaml:"request-timeout"`
}

type RabbitMQNotifierConfig struct {
	Type string `yaml:"type"` //Notifier Type

	//amqp Dial
	Url string `yaml:"url"`

	//amqp.Channel.Publish
	Exchange   string `yaml:"exchange"`
	RoutingKey string `yaml:"routing-key"`
	Mandatory  bool   `yaml:"mandatory"`
	Immediate  bool   `yaml:"immediate"`

	//amqp.Publishing
	MessageHeaders         map[string]interface{} `yaml:"message-headers,omitempty"` // Application or header exchange table
	MessageContentType     string                 `yaml:"message-content-type"`      // MIME content type
	MessageContentEncoding string                 `yaml:"message-content-encoding"`  // MIME content encoding
	MessageDeliveryMode    uint8                  `yaml:"message-delivery-mode"`     // queue implementation use - Transient (1) or Persistent (2)
	MessagePriority        uint8                  `yaml:"message-priority"`          // queue implementation use - 0 to 9
	MessageCorrelationId   string                 `yaml:"message-correlation-id"`    // application use - correlation identifier
	MessageReplyTo         string                 `yaml:"message-reply-to"`          // application use - address to to reply to (ex: RPC)
	MessageExpiration      string                 `yaml:"message-expiration"`        // implementation use - message expiration spec
	MessageMessageId       string                 `yaml:"message-message-id"`        // application use - message identifier
	MessageTimestamp       bool                   `yaml:"message-timestamp"`         // application use - message timestamp
	MessageType            string                 `yaml:"message-type"`              // application use - message type name
	MessageUserId          string                 `yaml:"message-user-id"`           // application use - creating user id
	MessageAppId           string                 `yaml:"message-app-id"`            // application use - creating application
}

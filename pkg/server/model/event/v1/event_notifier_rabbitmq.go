package v1

import metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"

type EventNotifierRabbitMqProperty struct {
	//amqp Dial
	Url string `json:"url" xorm:"'url' notnull"`

	//amqp.Channel.Publish
	Exchange   *string `json:"exchange"    xorm:"'exchange'    null"`
	RoutingKey *string `json:"routing_key" xorm:"'routing_key' null"`
	Mandatory  *bool   `json:"mandatory"   xorm:"'mandatory'   null"`
	Immediate  *bool   `json:"immediate"   xorm:"'immediate'   null"`

	//amqp.Publishing
	MessageHeaders         map[string]interface{} `json:"message_headers,omitempty" xorm:"'message_headers'          null"` // Application or header exchange table
	MessageContentType     *string                `json:"message_content_type"      xorm:"'message_content_type'     null"` // MIME content type
	MessageContentEncoding *string                `json:"message_content_encoding"  xorm:"'message_content_encoding' null"` // MIME content encoding
	MessageDeliveryMode    *uint8                 `json:"message_delivery_mode"     xorm:"'message_delivery_mode'    null"` // queue implementation use - Transient (1) or Persistent (2)
	MessagePriority        *uint8                 `json:"message_priority"          xorm:"'message_priority'         null"` // queue implementation use - 0 to 9
	MessageCorrelationId   *string                `json:"message_correlation_id"    xorm:"'message_correlation_id'   null"` // application use - correlation identifier
	MessageReplyTo         *string                `json:"message_reply_to"          xorm:"'message_reply_to'         null"` // application use - address to to reply to (ex: RPC)
	MessageExpiration      *string                `json:"message_expiration"        xorm:"'message_expiration'       null"` // implementation use - message expiration spec
	MessageMessageId       *string                `json:"message_message_id"        xorm:"'message_message_id'       null"` // application use - message identifier
	MessageTimestamp       *bool                  `json:"message_timestamp"         xorm:"'message_timestamp'        null"` // application use - message timestamp
	MessageType            *string                `json:"message_type"              xorm:"'message_type'             null"` // application use - message type name
	MessageUserId          *string                `json:"message_user_id"           xorm:"'message_user_id'          null"` // application use - creating user id
	MessageAppId           *string                `json:"message_app_id"            xorm:"'message_app_id'           null"` // application use - creating application
}

type EventNotifierRabbitMq struct {
	metav1.DbMeta                 `json:",inline" xorm:"extends"` //inline dbmeta
	metav1.UuidMeta               `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta              `json:",inline" xorm:"extends"` //inline labelmeta
	EventNotifierRabbitMqProperty `json:",inline" xorm:"extends"` //inline property
}

func (EventNotifierRabbitMq) TableName() string {
	return "event_notifier_rabbitmq"
}

type EventNotifierRabbitMq_create struct {
	metav1.LabelMeta              `json:",inline" xorm:"extends"` //inline labelmeta
	EventNotifierRabbitMqProperty `json:",inline" xorm:"extends"` //inline property
}

package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

type EventNotifierWebhookProperty struct {
	//http
	Method         string            `json:"method"                    xorm:"'method'          varchar(255) notnull"`                                    //
	Url            string            `json:"url"                       xorm:"'url'                          notnull"`                                    //
	RequestHeaders map[string]string `json:"request_headers,omitempty" xorm:"'request_headers'              null"   `                                    //
	RequestTimeout string            `json:"request_timeout"           xorm:"'request_timeout' varchar(16)  null    comment('fmt(time.ParseDuration)')"` //for timeout context
}

func (EventNotifierWebhookProperty) Type() EventNotifierType {
	return EventNotifierTypeWebhook
}

type EventNotifierWebhook struct {
	metav1.DbMeta                `json:",inline" xorm:"extends"` //inline dbmeta
	metav1.UuidMeta              `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta             `json:",inline" xorm:"extends"` //inline labelmeta
	EventNotifierWebhookProperty `json:",inline" xorm:"extends"` //inline property
	MIME                         `json:",inline" xorm:"extends"` //inline MIME
}

func (EventNotifierWebhook) TableName() string {
	return "event_notifier_webhook"
}

type EventNotifierWebhook_create struct {
	metav1.LabelMeta             `json:",inline" xorm:"extends"` //inline labelmeta
	EventNotifierWebhookProperty `json:",inline" xorm:"extends"` //inline property
	MIME                         `json:",inline" xorm:"extends"` //inline MIME
}

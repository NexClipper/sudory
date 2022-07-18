package v2

import (
	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type NotifierWebhook_essential struct {
	//http
	Method         string               `column:"method,default('')"         json:"method"`
	Url            string               `column:"url,default('')"            json:"url"`
	RequestHeaders vanilla.NullKeyValue `column:"request_headers"            json:"request_headers,omitempty" swaggertype:"object"`
	RequestTimeout uint                 `column:"request_timeout,default(0)" json:"request_timeout"` // second
}
type NotifierWebhook_property struct {
	NotifierWebhook_essential `json:",inline"`

	Created vanilla.NullTime `column:"created" json:"created,omitempty" swaggertype:"string"`
	Updated vanilla.NullTime `column:"updated" json:"updated,omitempty" swaggertype:"string"`
}

func (NotifierWebhook_property) Type() NotifierType {
	return NotifierTypeWebhook
}

func (NotifierWebhook_property) TableName() string {
	return "managed_channel_notifier_webhook"
}

type NotifierWebhook struct {
	NotifierWebhook_property `json:",inline"`

	Uuid string `column:"uuid"    json:"uuid,omitempty"` // pk
}

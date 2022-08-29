package v3

import (
	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type WebhookConfig struct {
	//http
	Method         string               `column:"method,default('')"         json:"method"`
	Url            string               `column:"url,default('')"            json:"url"`
	RequestHeaders vanilla.NullKeyValue `column:"request_headers"            json:"request_headers,omitempty" swaggertype:"object"`
	RequestTimeout uint                 `column:"request_timeout,default(0)" json:"request_timeout"` // second
}

func (WebhookConfig) Type() NotifierType {
	return NotifierTypeWebhook
}

func (cfg WebhookConfig) Valid() bool {
	return true
}

type NotifierWebhook_update = WebhookConfig

type NotifierWebhook_property = WebhookConfig

type NotifierWebhook struct {
	Uuid string `column:"uuid" json:"uuid,omitempty"` // pk

	NotifierWebhook_property `json:",inline"`

	// Created vanilla.NullTime `column:"created" json:"created,omitempty" swaggertype:"string"`
	// Updated vanilla.NullTime `column:"updated" json:"updated,omitempty" swaggertype:"string"`
}

func (NotifierWebhook) TableName() string {
	return "managed_channel_notifier_webhook"
}

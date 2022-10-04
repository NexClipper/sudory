package v3

import (
	"fmt"
	"strings"

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

func (cfg WebhookConfig) Valid() error {
	const http = "http://"
	const https = "https://"
	if strings.Index(cfg.Url, http) != 0 && strings.Index(cfg.Url, https) != 0 {
		return fmt.Errorf("url is not an expression of the Webhook protocol")
	}
	if len(cfg.Method) == 0 {
		return fmt.Errorf("missing method")
	}
	return nil
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

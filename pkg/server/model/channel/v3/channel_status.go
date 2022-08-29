package v3

import (
	"time"
)

type ChannelStatus struct {
	Uuid    string    `column:"uuid"    json:"uuid,omitempty"`    // pk
	Created time.Time `column:"created" json:"created,omitempty"` // pk
	Message string    `column:"message" json:"message,omitempty"`
}

func (ChannelStatus) TableName() string {
	return "managed_channel_status"
}

type ChannelStatusOption_update = ChannelStatusOption_property

type ChannelStatusOption_property struct {
	StatusMaxCount uint `column:"status_max_count,default(0)" json:"status_max_count"`
}

type ChannelStatusOption struct {
	Uuid string `column:"uuid"                        json:"uuid,omitempty"` // pk

	ChannelStatusOption_property `json:",inline"`

	// Created vanilla.NullTime `column:"created"                     json:"created,omitempty" swaggertype:"string"`
	// Updated vanilla.NullTime `column:"updated"                     json:"updated,omitempty" swaggertype:"string"`
}

func (ChannelStatusOption) TableName() string {
	return "managed_channel_status_option"
}

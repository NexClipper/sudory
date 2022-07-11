package v2

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type ChannelStatus struct {
	Uuid    string    `column:"uuid"    json:"uuid,omitempty"`    // pk
	Created time.Time `column:"created" json:"created,omitempty"` // pk
	Message string    `column:"message" json:"message,omitempty"`
}

func (ChannelStatus) TableName() string {
	return "managed_channel_status"
}

type ChannelStatusOption_essential struct {
	StatusMaxCount uint `column:"status_max_count,default(0)" json:"status_max_count"`
}

type ChannelStatusOption_property struct {
	ChannelStatusOption_essential `json:",inline"`

	Created vanilla.NullTime `column:"created" json:"created,omitempty" swaggertype:"string"`
	Updated vanilla.NullTime `column:"updated" json:"updated,omitempty" swaggertype:"string"`
}

func (ChannelStatusOption_property) TableName() string {
	return "managed_channel_status_option"
}

type ChannelStatusOption struct {
	ChannelStatusOption_property `json:",inline"`

	Uuid string `column:"uuid" json:"uuid,omitempty"` // pk
}

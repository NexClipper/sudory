package v2

import (
	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type NotifierSlackhook_essential struct {
	Url            string `column:"url,default('')"            json:"url"`
	RequestTimeout uint   `column:"request_timeout,default(0)" json:"request_timeout"` // second
}
type NotifierSlackhook_property struct {
	NotifierSlackhook_essential `json:",inline"`

	Created vanilla.NullTime `column:"created" json:"created,omitempty" swaggertype:"string"`
	Updated vanilla.NullTime `column:"updated" json:"updated,omitempty" swaggertype:"string"`
}

func (NotifierSlackhook_property) Type() NotifierType {
	return NotifierTypeSlackhook
}

func (NotifierSlackhook_property) TableName() string {
	return "managed_channel_notifier_slackhook"
}

type NotifierSlackhook struct {
	NotifierSlackhook_property `json:",inline"`

	Uuid string `column:"uuid"    json:"uuid,omitempty"` // pk
}

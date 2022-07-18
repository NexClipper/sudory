package v2

import (
	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type NotifierConsole_essential struct{}

type NotifierConsole_property struct {
	Created vanilla.NullTime `column:"created" json:"created,omitempty" swaggertype:"string"`
	Updated vanilla.NullTime `column:"updated" json:"updated,omitempty" swaggertype:"string"`

	NotifierConsole_essential `json:",inline"`
}

func (NotifierConsole_property) Type() NotifierType {
	return NotifierTypeConsole
}

func (NotifierConsole_property) TableName() string {
	return "managed_channel_notifier_console"
}

type NotifierConsole struct {
	Uuid string `column:"uuid"    json:"uuid,omitempty"` // pk

	NotifierConsole_property `json:",inline"` //inline property
}

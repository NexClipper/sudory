package v2

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

//GlobalVariables Property
type GlobalVariables_property struct {
	Value vanilla.NullString `column:"value" json:"value,omitempty" swaggertype:"string"`
}

func (GlobalVariables_property) TableName() string {
	return "global_variables"
}

//GlobalVariables
type GlobalVariables struct {
	ID   int64  `column:"id"   json:"id,omitempty"`   // pk
	Uuid string `column:"uuid" json:"uuid,omitempty"` // uuid

	GlobalVariables_property `json:",inline"`

	Created time.Time        `column:"created" json:"created,omitempty"`
	Updated vanilla.NullTime `column:"updated" json:"updated,omitempty" swaggertype:"string"`
	Deleted vanilla.NullTime `column:"deleted" json:"deleted,omitempty" swaggertype:"string"`
}

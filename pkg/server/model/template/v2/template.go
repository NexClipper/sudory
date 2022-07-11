package v2

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type Template_essential struct {
	Name    string             `column:"name"    json:"name,omitempty"`
	Summary vanilla.NullString `column:"summary" json:"summary,omitempty" swaggertype:"string"`
	Origin  string             `column:"origin"  json:"origin,omitempty"`
}

func (Template_essential) TableName() string {
	return "template"
}

type Template struct {
	// Id   uint64 `column:"id"   json:"id,omitempty"` //pk
	Uuid string `column:"uuid" json:"uuid,omitempty"` //pk

	Template_essential `json:",inline"`

	Created time.Time        `column:"created" json:"created,omitempty"`
	Updated vanilla.NullTime `column:"updated" json:"updated,omitempty" swaggertype:"string"`
	Deleted vanilla.NullTime `column:"deleted" json:"deleted,omitempty" swaggertype:"string"`
}

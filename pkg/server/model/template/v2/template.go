package v2

import (
	"time"

	noxorm "github.com/NexClipper/sudory/pkg/server/model/noxorm/v2"
)

type Template_essential struct {
	Name    string            `column:"name"    json:"name,omitempty"`
	Summary noxorm.NullString `column:"summary" json:"summary,omitempty"`
	Origin  string            `column:"origin"  json:"origin,omitempty"`
}

func (Template_essential) TableName() string {
	return "template"
}

type Template struct {
	// Id   uint64 `column:"id"   json:"id,omitempty"` //pk
	Uuid string `column:"uuid" json:"uuid,omitempty"` //pk

	Template_essential `json:",inline"`

	Created time.Time       `column:"created" json:"created,omitempty"`
	Updated noxorm.NullTime `column:"updated" json:"updated,omitempty"`
	Deleted noxorm.NullTime `column:"deleted" json:"deleted,omitempty"`
}

package v2

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

//GlobalVariables
type GlobalVariables struct {
	ID      int64              `column:"id"      json:"id,omitempty"`   // pk
	Uuid    string             `column:"uuid"    json:"uuid,omitempty"` // uuid
	Name    string             `column:"name"    json:"name,omitempty"`
	Summary vanilla.NullString `column:"summary" json:"summary,omitempty" swaggertype:"string"`
	Value   vanilla.NullString `column:"value"   json:"value,omitempty"   swaggertype:"string"`
	Created time.Time          `column:"created" json:"created,omitempty"`
	Updated vanilla.NullTime   `column:"updated" json:"updated,omitempty" swaggertype:"string"`
	Deleted vanilla.NullTime   `column:"deleted" json:"deleted,omitempty" swaggertype:"string"`
}

func (GlobalVariables) TableName() string {
	return "global_variables"
}

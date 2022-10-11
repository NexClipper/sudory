package v3

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type Cluster_create struct {
	Uuid          string             `json:"uuid,omitempty"` // uuid
	Name          string             `column:"name"           json:"name,omitempty"`
	Summary       vanilla.NullString `column:"summary"        json:"summary,omitempty"        swaggertype:"string"`
	PollingOption vanilla.NullObject `column:"polling_option" json:"polling_option,omitempty" swaggertype:"object"`
	PoliingLimit  vanilla.NullInt    `column:"polling_limit"  json:"polling_limit,omitempty"  swaggertype:"integer"`
}

type Cluster_update struct {
	Name          string             `column:"name"           json:"name,omitempty"`
	Summary       vanilla.NullString `column:"summary"        json:"summary,omitempty"        swaggertype:"string"`
	PollingOption vanilla.NullObject `column:"polling_option" json:"polling_option,omitempty" swaggertype:"object"`
	PoliingLimit  vanilla.NullInt    `column:"polling_limit"  json:"polling_limit,omitempty"  swaggertype:"integer"`
}

type Cluster struct {
	ID            int64              `column:"id"             json:"id,omitempty"`   // pk
	Uuid          string             `column:"uuid"           json:"uuid,omitempty"` // uuid
	Name          string             `column:"name"           json:"name,omitempty"`
	Summary       vanilla.NullString `column:"summary"        json:"summary,omitempty"        swaggertype:"string"`
	PollingOption vanilla.NullObject `column:"polling_option" json:"polling_option,omitempty" swaggertype:"object"`
	PoliingLimit  int                `column:"polling_limit"  json:"polling_limit,omitempty"`
	Created       time.Time          `column:"created"        json:"created,omitempty"`
	Updated       vanilla.NullTime   `column:"updated"        json:"updated,omitempty" swaggertype:"string"`
	Deleted       vanilla.NullTime   `column:"deleted"        json:"deleted,omitempty" swaggertype:"string"`
}

func (Cluster) TableName() string {
	return "cluster"
}

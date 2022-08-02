package v2

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type Cluster_essential struct {
	Name          string             `column:"name"           json:"name,omitempty"`
	Summary       vanilla.NullString `column:"summary"        json:"summary,omitempty"        swaggertype:"string"`
	PollingOption vanilla.NullObject `column:"polling_option" json:"polling_option,omitempty" swaggertype:"object"`
	PoliingLimit  int                `column:"polling_limit"  json:"polling_limit,omitempty"`
}

func (Cluster_essential) TableName() string {
	return "cluster"
}

type Cluster struct {
	Cluster_essential `json:",inline"`

	ID      int64            `column:"id"      json:"id,omitempty"`   // pk
	Uuid    string           `column:"uuid"    json:"uuid,omitempty"` // uuid
	Created time.Time        `column:"created" json:"created,omitempty"`
	Updated vanilla.NullTime `column:"updated" json:"updated,omitempty" swaggertype:"string"`
	Deleted vanilla.NullTime `column:"deleted" json:"deleted,omitempty" swaggertype:"string"`
}

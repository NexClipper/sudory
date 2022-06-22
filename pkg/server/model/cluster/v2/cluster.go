package v2

import (
	"time"

	noxorm "github.com/NexClipper/sudory/pkg/server/model/noxorm/v2"
)

type Cluster_essential struct {
	Name          string            `column:"name"           json:"name,omitempty"`
	Summary       noxorm.NullString `column:"summary"        json:"summary,omitempty"`
	PollingOption noxorm.NullJson   `column:"polling_option" json:"polling_option,omitempty"`
	PoliingLimit  int               `column:"polling_limit"  json:"polling_limit,omitempty"`
}

func (Cluster_essential) TableName() string {
	return "cluster"
}

type Cluster struct {
	ID   int64  `column:"id"   json:"id,omitempty"`   // pk
	Uuid string `column:"uuid" json:"uuid,omitempty"` // uuid

	Cluster_essential `json:",inline"`

	Created time.Time       `column:"created" json:"created,omitempty"`
	Updated noxorm.NullTime `column:"updated" json:"updated,omitempty"`
	Deleted noxorm.NullTime `column:"deleted" json:"deleted,omitempty"`
}

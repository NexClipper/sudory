package clusterinfos

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type ClusterInformation struct {
	ID            uint             `column:"id"             json:"id,omitempty"`           // pk
	ClusterUuid   string           `column:"cluster_uuid"   json:"cluster_uuid,omitempty"` // uuid
	PollingCount  vanilla.NullInt  `column:"polling_count"  json:"polling_count,omitempty"`
	PollingOffset vanilla.NullTime `column:"polling_offset" json:"polling_offset,omitempty"`
	Created       time.Time        `column:"created"        json:"created,omitempty"`
	Updated       vanilla.NullTime `column:"updated"        json:"updated,omitempty" swaggertype:"string"`
	Deleted       vanilla.NullTime `column:"deleted"        json:"deleted,omitempty" swaggertype:"string"`
}

func (ClusterInformation) TableName() string {
	return "cluster_information"
}

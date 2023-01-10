package service

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type servicestatusTableName struct{}

func (servicestatusTableName) TableName() string {
	return "service_status_v2"
}

type pkServiceStatus struct {
	PartitionDate time.Time `column:"pdate"        json:"partition_date"` // pk date
	ClusterUuid   string    `column:"cluster_uuid" json:"cluster_uuid"`   // pk char(32) cluster.uuid
	Uuid          string    `column:"uuid"         json:"uuid"`           // pk char(32) service.uuid
	Created       time.Time `column:"created"      json:"created"`        // pk datetime(6)
}

// type ServiceStatus_create struct {
// 	servicestatusTableName `json:"-"`

// 	pkServiceStatus `json:",inline"`

// 	Status  StepStatus
// 	Started vanilla.NullTime
// 	Ended   vanilla.NullTime
// }

// type ServiceStatus_update struct {
// 	servicestatusTableName `json:"-"`

// 	Status    StepStatus
// 	Started   vanilla.NullTime
// 	Ended     vanilla.NullTime
// 	Timestamp time.Time
// }

type ServiceStatus struct {
	servicestatusTableName `json:"-"`

	pkServiceStatus `json:",inline"`
	StepMax         int                `column:"step_max" json:"step_max,omitempty"`
	StepSeq         int                `column:"step_seq" json:"step_seq,omitempty"`
	Status          StepStatus         `column:"status"   json:"status,omitempty"`
	Started         vanilla.NullTime   `column:"started"  json:"started,omitempty"       swaggertype:"string"`
	Ended           vanilla.NullTime   `column:"ended"    json:"ended,omitempty"         swaggertype:"string"`
	Message         vanilla.NullString `column:"message"  json:"message,omitempty"       swaggertype:"string"`
}

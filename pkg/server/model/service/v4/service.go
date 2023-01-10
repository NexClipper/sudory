package service

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	cryptov2 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v2"
)

type serviceTableName struct{}

func (serviceTableName) TableName() string {
	return "service_v2"
}

type pkService struct {
	PartitionDate time.Time `column:"pdate"        json:"partition_date"` // pk date
	ClusterUuid   string    `column:"cluster_uuid" json:"cluster_uuid"`   // pk char(32) cluster.uuid
	Uuid          string    `column:"uuid"         json:"uuid"`           // pk char(32) service.uuid
	// Created       time.Time `column:"created"      json:"created"`        // pk datetime(6)
}

// type Service_create struct {
// 	serviceTableName `json:"-"`

// 	pkService         `json:",inline"`
// 	Name              string             `column:"name"               json:"name,omitempty"`
// 	Summary           vanilla.NullString `column:"summary"            json:"summary,omitempty"            swaggertype:"string"`
// 	TemplateUuid      string             `column:"template_uuid"      json:"template_uuid,omitempty"`
// 	StepCount         int                `column:"step_count"         json:"step_count,omitempty"`
// 	Priority          Priority           `column:"priority"           json:"priority,omitempty"`
// 	SubscribedChannel vanilla.NullString `column:"subscribed_channel" json:"subscribed_channel,omitempty" swaggertype:"string"`
// 	StepPosition      int                `column:"step_position"      json:"step_position,omitempty"`
// 	Status            StepStatus         `column:"status"             json:"status,omitempty"`
// 	Created           time.Time          `column:"created"            json:"created,omitempty"`
// }

// type Service_update struct {
// 	serviceTableName `json:"-"`

// 	AssignedClientUuid vanilla.NullString
// 	StepPosition       int
// 	Status             StepStatus
// 	Message            vanilla.NullString
// 	Timestamp          time.Time
// }

type Service struct {
	serviceTableName `json:"-"`

	pkService         `json:",inline"`
	Name              string                `column:"name"               json:"name,omitempty"`
	Summary           vanilla.NullString    `column:"summary"            json:"summary,omitempty"            swaggertype:"string"`
	TemplateUuid      string                `column:"template_uuid"      json:"template_uuid,omitempty"`
	Flow              string                `column:"flow"               json:"flow,omitempty"`
	Inputs            cryptov2.CryptoObject `column:"inputs"             json:"inputs,omitempty"             swaggertype:"object"`
	StepMax           int                   `column:"step_max"           json:"step_max,omitempty"`
	SubscribedChannel vanilla.NullString    `column:"subscribed_channel" json:"subscribed_channel,omitempty" swaggertype:"string"`
	Priority          Priority              `column:"priority"           json:"priority,omitempty"`
	Created           time.Time             `column:"created"            json:"created,omitempty"`
}

type Service_polling struct {
	serviceTableName `json:"-"`

	pkService `json:",inline"`
	Created   time.Time `column:"created"`
	Priority  Priority  `column:"priority"`
}

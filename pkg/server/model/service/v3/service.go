package v3

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type serviceTableName struct{}

func (serviceTableName) TableName() string {
	return "service"
}

type pkService struct {
	ClusterUuid   string    `column:"cluster_uuid" json:"cluster_uuid"`   // pk char(32) cluster.uuid
	Uuid          string    `column:"uuid"         json:"uuid"`           // pk char(32) service.uuid
	Timestamp     time.Time `column:"timestamp"    json:"timestamp"`      // pk datetime(6)
	PartitionDate time.Time `column:"pdate"        json:"partition_date"` // pk date
}

type Service_create struct {
	serviceTableName `json:"-"`

	pkService         `json:",inline"`
	Name              string             `column:"name"               json:"name,omitempty"`
	Summary           vanilla.NullString `column:"summary"            json:"summary,omitempty"            swaggertype:"string"`
	TemplateUuid      string             `column:"template_uuid"      json:"template_uuid,omitempty"`
	StepCount         int                `column:"step_count"         json:"step_count,omitempty"`
	Priority          Priority           `column:"priority"           json:"priority,omitempty"`
	SubscribedChannel vanilla.NullString `column:"subscribed_channel" json:"subscribed_channel,omitempty" swaggertype:"string"`
	StepPosition      int                `column:"step_position"      json:"step_position,omitempty"`
	Status            StepStatus         `column:"status"             json:"status,omitempty"`
	Created           time.Time          `column:"created"            json:"created,omitempty"`
}

type Service_update struct {
	serviceTableName `json:"-"`

	AssignedClientUuid vanilla.NullString
	StepPosition       int
	Status             StepStatus
	Message            vanilla.NullString
	Timestamp          time.Time
}

type Service struct {
	serviceTableName `json:"-"`

	pkService          `json:",inline"`
	Name               string             `column:"name"                 json:"name,omitempty"`
	Summary            vanilla.NullString `column:"summary"              json:"summary,omitempty"              swaggertype:"string"`
	TemplateUuid       string             `column:"template_uuid"        json:"template_uuid,omitempty"`
	StepCount          int                `column:"step_count"           json:"step_count,omitempty"`
	Priority           Priority           `column:"priority"             json:"priority,omitempty"`
	SubscribedChannel  vanilla.NullString `column:"subscribed_channel"   json:"subscribed_channel,omitempty"   swaggertype:"string"`
	AssignedClientUuid vanilla.NullString `column:"assigned_client_uuid" json:"assigned_client_uuid,omitempty" swaggertype:"string"`
	StepPosition       int                `column:"step_position"        json:"step_position,omitempty"`
	Status             StepStatus         `column:"status"               json:"status,omitempty"`
	Message            vanilla.NullString `column:"message"              json:"message,omitempty"              swaggertype:"string"`
	Created            time.Time          `column:"created"              json:"created,omitempty"`
}

type Service_polling struct {
	serviceTableName `json:"-"`

	Uuid      string     `column:"uuid"`
	Timestamp time.Time  `column:"timestamp"`
	Priority  Priority   `column:"priority"`
	Status    StepStatus `column:"status"`
	Created   time.Time  `column:"created"`
}

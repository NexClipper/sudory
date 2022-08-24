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
	PartitionDate time.Time `column:"pdate"        json:"partition_date"` // pk patition hash
	ClusterUuid   string    `column:"cluster_uuid" json:"cluster_uuid"`   // pk cluster uuid
	Uuid          string    `column:"uuid"         json:"uuid"`           // pk service uuid
	Revision      uint8     `column:"revision"     json:"revision"`       // pk record revision
}

type Service_create struct {
	serviceTableName `json:"-"`

	PK                pkService          `json:",inline"`
	Name              string             `column:"name"               json:"name,omitempty"`
	Summary           vanilla.NullString `column:"summary"            json:"summary,omitempty"            swaggertype:"string"`
	TemplateUuid      string             `column:"template_uuid"      json:"template_uuid,omitempty"`
	StepCount         uint8              `column:"step_count"         json:"step_count,omitempty"`
	SubscribedChannel vanilla.NullString `column:"subscribed_channel" json:"subscribed_channel,omitempty" swaggertype:"string"`
	StepPosition      uint8              `column:"step_position"      json:"step_position,omitempty"`
	Status            StepStatus         `column:"status"             json:"status,omitempty"`
	Created           time.Time          `column:"created"            json:"created,omitempty"`
}

type Service_update struct {
	serviceTableName `json:"-"`

	AssignedClientUuid vanilla.NullString
	StepPosition       uint8
	Status             StepStatus
	Message            vanilla.NullString
	Updated            vanilla.NullTime
}

type Service struct {
	serviceTableName `json:"-"`

	PK                 pkService          `json:",inline"`
	Name               string             `column:"name"                 json:"name,omitempty"`
	Summary            vanilla.NullString `column:"summary"              json:"summary,omitempty"              swaggertype:"string"`
	TemplateUuid       string             `column:"template_uuid"        json:"template_uuid,omitempty"`
	StepCount          uint8              `column:"step_count"           json:"step_count,omitempty"`
	SubscribedChannel  vanilla.NullString `column:"subscribed_channel"   json:"subscribed_channel,omitempty"   swaggertype:"string"`
	AssignedClientUuid vanilla.NullString `column:"assigned_client_uuid" json:"assigned_client_uuid,omitempty" swaggertype:"string"`
	StepPosition       uint8              `column:"step_position"        json:"step_position,omitempty"`
	Status             StepStatus         `column:"status"               json:"status,omitempty"`
	Message            vanilla.NullString `column:"message"              json:"message,omitempty"              swaggertype:"string"`
	Created            time.Time          `column:"created"              json:"created,omitempty"`
	Updated            vanilla.NullTime   `column:"updated"              json:"updated,omitempty"              swaggertype:"string"`
}

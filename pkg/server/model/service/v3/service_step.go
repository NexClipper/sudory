package v3

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	cryptov2 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v2"
)

type serviceStepTableName struct{}

func (serviceStepTableName) TableName() string {
	return "service_step"
}

type pkServiceStep struct {
	PartitionDate time.Time `column:"pdate"        json:"partition_date"` // pk patition hash
	ClusterUuid   string    `column:"cluster_uuid" json:"cluster_uuid"`   // pk cluster uuid
	Uuid          string    `column:"uuid"         json:"uuid"`           // pk step uuid
	Sequence      uint8     `column:"seq"          json:"seq"`            // pk step sequence
	Revision      uint8     `column:"revision"     json:"revision"`       // pk record revision
}

type ServiceStep_create struct {
	serviceStepTableName `json:"-"`

	PK           pkServiceStep         `json:",inline"`
	Name         string                `column:"name"          json:"name,omitempty"`
	Summary      vanilla.NullString    `column:"summary"       json:"summary,omitempty"       swaggertype:"string"`
	Method       string                `column:"method"        json:"method,omitempty"`
	Args         cryptov2.CryptoObject `column:"args"          json:"args,omitempty"          swaggertype:"object"`
	ResultFilter vanilla.NullString    `column:"result_filter" json:"result_filter,omitempty" swaggertype:"string"`
	Status       StepStatus            `column:"status"        json:"status,omitempty"`
	Created      time.Time             `column:"created"       json:"created,omitempty"`
}

type ServiceStep_update struct {
	serviceStepTableName `json:"-"`

	Status  StepStatus
	Started vanilla.NullTime
	Ended   vanilla.NullTime
	Updated vanilla.NullTime
}

type ServiceStep struct {
	serviceStepTableName `json:"-"`

	PK           pkServiceStep         `json:",inline"`
	Name         string                `column:"name"          json:"name,omitempty"`
	Summary      vanilla.NullString    `column:"summary"       json:"summary,omitempty"       swaggertype:"string"`
	Method       string                `column:"method"        json:"method,omitempty"`
	Args         cryptov2.CryptoObject `column:"args"          json:"args,omitempty"          swaggertype:"object"`
	ResultFilter vanilla.NullString    `column:"result_filter" json:"result_filter,omitempty" swaggertype:"string"`
	Status       StepStatus            `column:"status"        json:"status,omitempty"`
	Started      vanilla.NullTime      `column:"started"       json:"started,omitempty"       swaggertype:"string"`
	Ended        vanilla.NullTime      `column:"ended"         json:"ended,omitempty"         swaggertype:"string"`
	Created      time.Time             `column:"created"       json:"created,omitempty"`
	Updated      vanilla.NullTime      `column:"updated"       json:"updated,omitempty"       swaggertype:"string"`
}

func (ServiceStep) TableName() string {
	return "service_step"
}

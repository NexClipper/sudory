package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/model/orm"
)

type Status int32

const (
	StatusDone  Status = 0
	StatusInit  Status = 1 << 0
	StatusRun   Status = 1 << 1
	StatusError Status = 1 << 2
)

//ServiceProperty Property
type ServiceProperty struct {
	//클러스터 UUID
	ClusterUuid string `json:"cluster_uuid" xorm:"char(32) notnull index 'cluster_uuid' comment('cluster's uuid')"`
	// //스탭 카운트
	// StepCount int32 `json:"step_count,omitempty" xorm:"int notnull default(0) 'step_count' comment('step_count')"`
	// // 스탭 Position
	// StepPosition int32 `json:"step_position,omitempty" xorm:"int notnull default(0) 'step_position' comment('step_position')"`
	//Status
	Status int32 `json:"status,omitempty" xorm:"int notnull default(1) index 'status' comment('status')"`
}

//MODEL: SERVICE
type Service struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	ServiceProperty  `json:",inline" xorm:"extends"` //inline property
}

//DATABASE SCHEMA: SERVICE
type DbSchemaService struct {
	metav1.DbMeta `xorm:"extends"`
	Service       `xorm:"extends"`
}

var _ orm.TableName = (*DbSchemaService)(nil)

func (DbSchemaService) TableName() string {
	return "service"
}

//HTTP REQUEST BODY: SERVICE
type HttpReqService struct {
	Service `json:",inline"`
}

//HTTP REQUEST BODY: SERVICE (with steps)
type HttpReqServiceWithSteps struct {
	Service `json:",inline"`
	Steps   []ServiceStep `json:"steps,inline"`
}

//HTTP REQUEST BODY: SERVICE (client)
type HttpReqServiceClient struct {
	Service `json:",inline"`
}

//HTTP RESPONSE BODY: SERVICE
type HttpRspService struct {
	Service `json:",inline"`
}

//HTTP RESPONSE BODY: SERVICE (client)
type HttpRspServiceClient struct {
	Service `json:",inline"`
	Steps   []ServiceStep `json:"steps,inline"`
}

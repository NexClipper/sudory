package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/model/orm"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
)

type Status int32

const (
	StatusRegiste Status = 0
	StatusSend    Status = 1 << 0
	StatusRun     Status = 1 << 1
	StatusDone    Status = 1 << 2
	StatusError   Status = 1 << 3
)

//ServiceProperty Property
type ServiceProperty struct {
	//클러스터 UUID
	ClusterUuid string `json:"cluster_uuid" xorm:"char(32) notnull index 'cluster_uuid' comment('cluster's uuid')"`
	//스탭 카운트
	StepCount int32 `json:"step_count,omitempty" xorm:"int notnull default(0) 'step_count' comment('step_count')"`
	//스탭 Position
	StepPosition int32 `json:"step_position,omitempty" xorm:"int notnull default(0) 'step_position' comment('step_position')"`
	//Status
	Status int32 `json:"status,omitempty" xorm:"int notnull default(0) index 'status' comment('status')"`
	//Type; 0: Once, 1: repeat(epoch, interval)
	Type int32 `json:"type,omitempty" xorm:"int notnull default(0) index 'type' comment('type')"`
	//Epoch -1: infinite, 0 :
	Epoch int32 `json:"epoch,omitempty" xorm:"int notnull default(0) index 'epoch' comment('epoch (times)')"`
	//Interval
	Interval int32 `json:"interval,omitempty" xorm:"int notnull default(0) index 'interval' comment('interval (sec)')"`
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
	Steps   []stepv1.ServiceStep `json:"steps,inline"`
}

//HTTP REQUEST BODY: SERVICE (client)
type HttpReqServiceClient struct {
	ClusterUuid string               `json:"cluster_uuid"`
	Steps       []stepv1.ServiceStep `json:",inline"`
}

//HTTP RESPONSE BODY: SERVICE
type HttpRspService struct {
	Service `json:",inline"`
	Steps   []stepv1.ServiceStep `json:"steps,inline"`
}

//HTTP RESPONSE BODY: SERVICIES
type HttpRspServicies struct {
	Service []HttpRspService
}

//변환 DbSchema -> Service
func TransFormDbSchema(s []DbSchemaService) []Service {
	var out = make([]Service, len(s))
	for n, it := range s {
		out[n] = it.Service
	}
	return out
}

//변환 Service -> HtppRsp
func TransToHttpRsp(s []Service) []HttpRspService {
	var out = make([]HttpRspService, len(s))
	for n, it := range s {
		out[n].Service = it
	}
	return out
}

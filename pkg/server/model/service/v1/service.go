package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/model/orm"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
)

type Status int32

const (
	StatusRegist     Status = 0
	StatusSend       Status = 1 << 0
	StatusProcessing Status = 1 << 1
	StatusSuccess    Status = 1 << 2
	StatusFail       Status = 1 << 3
)

//ServiceProperty Property
type ServiceProperty struct {
	//텔플릿 UUID
	TemplateUuid *string `json:"template_uuid,omitempty" xorm:"char(32) notnull index 'template_uuid' comment('template's uuid')"`
	//오리진 종류
	OriginKind *string `json:"origin_kind,omitempty" xorm:"varchar(255) notnull index 'origin_kind' comment('service origin kind')"`
	//오리진 UUID
	OriginUuid *string `json:"origin_uuid,omitempty" xorm:"char(32) notnull index 'origin_uuid' comment('service origin uuid')"`
	//클러스터 UUID
	ClusterUuid string `json:"cluster_uuid" xorm:"char(32) notnull index 'cluster_uuid' comment('cluster's uuid')"`
	//할당된 클라이언트 UUID
	AssignedClientUuid string `json:"assigned_client_uuid" xorm:"char(32) notnull index 'assigned_client_uuid' comment('client's uuid when service assigned')"`
	//스탭 카운트
	StepCount *int32 `json:"step_count,omitempty" xorm:"int null default(0) 'step_count' comment('step_count')"`
	//스탭 Position
	StepPosition *int32 `json:"step_position,omitempty" xorm:"int null default(0) 'step_position' comment('step_position')"`
	//Type; 0: Once, 1: repeat(epoch, interval)
	Type *int32 `json:"type,omitempty" xorm:"int null default(0) 'type' comment('type')"`
	//Epoch -1: infinite, 0 :
	Epoch *int32 `json:"epoch,omitempty" xorm:"int null default(0) 'epoch' comment('epoch (times)')"`
	//Interval
	Interval *int32 `json:"interval,omitempty" xorm:"int null default(0) 'interval' comment('interval (sec)')"`
	//Status
	Status *int32 `json:"status,omitempty" xorm:"int null default(0) index 'status' comment('status')"`
	//Result 스탭 실행 결과(정상:'결과', 오류:'오류 메시지')
	Result *string `json:"result,omitempty" xorm:"longtext null 'result' comment('result')"`
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

//HTTP RESPONSE BODY: SERVICE
type HttpRspService struct {
	Service `json:",inline"`
}

//HTTP REQUEST BODY: SERVICE (with steps)
type HttpReqServiceWithSteps struct {
	Service `json:",inline"`
	Steps   []stepv1.ServiceStep `json:"steps"`
}

//HTTP RESPONSE BODY: SERVICE
type HttpRspServiceWithSteps struct {
	Service `json:",inline"`
	Steps   []stepv1.ServiceStep `json:"steps"`
}

//HTTP REQUEST BODY: SERVICE (client)
type HttpReqClientSideService struct {
	Service `json:",inline"`
	Steps   []stepv1.ServiceStep `json:"steps"`
}

//HTTP RESPONSE BODY: SERVICE (client)
type HttpRspClientSideService struct {
	Service `json:",inline"`
	Steps   []stepv1.ServiceStep `json:"steps"`
}

//변환 DbSchema -> Service
func TransFormDbSchema(s []DbSchemaService) []Service {
	var out = make([]Service, len(s))
	for n, it := range s {
		out[n] = it.Service
	}
	return out
}

// //변환 Service -> HtppRsp
// func TransToHttpRsp(s []Service) []HttpRspService {
// 	var out = make([]HttpRspService, len(s))
// 	for n, it := range s {
// 		out[n].Service = it
// 	}
// 	return out
// }

//Build Template -> HttpRsp
func HttpRspBuilder(length int) (func(a Service, b []stepv1.ServiceStep), func() []HttpRspServiceWithSteps) {
	var pos int = 0
	queue := make([]HttpRspServiceWithSteps, length)
	pusher := func(a Service, b []stepv1.ServiceStep) {
		queue[pos] = HttpRspServiceWithSteps{Service: a, Steps: b}
		pos++
	}
	poper := func() []HttpRspServiceWithSteps {
		return queue
	}
	return pusher, poper
}

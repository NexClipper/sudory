package v1

import (
	"time"

	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/model/orm"
)

type ServiceStepProperty struct {
	//서비스 Uuid
	ServiceUuid *string `json:"service_uuid,omitempty" xorm:"char(32) notnull index 'service_uuid' comment('services uuid')"`
	//순서
	Sequence *int32 `json:"sequence,omitempty" xorm:"int null default(0) 'sequence' comment('sequence')"`
	//메소드
	Method *string `json:"method,omitempty" xorm:"varchar(255) null 'method' comment('method')"`
	//arguments
	Args map[string]interface{} `json:"args,omitempty" xorm:"text null 'args' comment('args')"`
	//Status 상태
	Status *int32 `json:"status,omitempty" xorm:"int null index default(0) 'status' comment('status')"`
	//Result 스탭 실행 결과(정상:'결과', 오류:'오류 메시지')
	Result *string `json:"result,omitempty" xorm:"longtext null 'result' comment('result')"`
	//Started 스탭 시작 시간
	Started *time.Time `json:"srated,omitempty" xorm:"datetime null comment('step start time')"`
	//Started 스탭 완료 시간
	Ended *time.Time `json:"ended,omitempty" xorm:"datetime null comment(step end time)"`
}

type ServiceStep struct {
	metav1.UuidMeta     `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta    `json:",inline" xorm:"extends"` //inline labelmeta
	ServiceStepProperty `json:",inline" xorm:"extends"` //inline property
}

//DATABASE SCHEMA: SERVICE
type DbSchemaServiceStep struct {
	metav1.DbMeta `xorm:"extends"`
	ServiceStep   `xorm:"extends"`
}

var _ orm.TableName = (*DbSchemaServiceStep)(nil)

func (DbSchemaServiceStep) TableName() string {
	return "service_step"
}

type ServiceStepPropertyEssential struct {
	//arguments
	Args map[string]interface{} `json:"args,omitempty"`
}

type ServiceStepEssential struct {
	ServiceStepPropertyEssential `json:",inline"` //inline property
}

//HTTP REQUEST BODY: SERVICE_STEP
type HttpReqServiceStep struct {
	ServiceStep `json:",inline"`
}

//HTTP RESPONSE BODY: SERVICE_STEP
type HttpRspServiceStep struct {
	ServiceStep `json:",inline"`
}

//변환 DbSchema -> Step
func TransFormDbSchema(s []DbSchemaServiceStep) []ServiceStep {
	var out = make([]ServiceStep, len(s))
	for n, it := range s {
		out[n] = it.ServiceStep
	}
	return out
}

//변환 Step -> HttpRsp
func TransToHttpRsp(s []ServiceStep) []HttpRspServiceStep {
	var out = make([]HttpRspServiceStep, len(s))
	for n, it := range s {
		out[n].ServiceStep = it
	}
	return out
}

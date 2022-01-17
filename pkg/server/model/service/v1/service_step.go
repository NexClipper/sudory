package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/model/orm"
)

type ServiceStepProperty struct {
	//서비스 Uuid
	ServiceUuid string `json:"service_uuid" xorm:"char(32) notnull index 'service_uuid' comment('services uuid')"`
	//순서
	Sequence int32 `json:"sequence,omitempty" xorm:"int notnull 'sequence' comment('sequence')"`
	//메소드
	Method string `json:"method,omitempty" xorm:"varchar(255) notnull 'method' comment('method')"`
	//arguments
	Args map[string]string `json:"args,omitempty" xorm:"text null 'args' comment('args')"`
	//Status 상태
	Status int32 `json:"status,omitempty" xorm:"int notnull default(1) 'status' comment('status')"`
	//Result 스탭 실행 결과(정상:'결과', 오류:'오류 메시지')
	Result string `json:"result,omitempty" xorm:"varchar(255) notnull default('') 'result' comment('result')"`
}

type ServiceStep struct {
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

//HTTP REQUEST BODY: SERVICE_STEP
type HttpReqServiceStep struct {
	ServiceStep `json:",inline"`
}

//HTTP RESPONSE BODY: SERVICE_STEP
type HttpRspServiceStep struct {
	ServiceStep `json:",inline"`
}
package v1

import (
	"time"

	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

type ServiceStepProperty struct {
	ServiceUuid  string                 `json:"service_uuid"            xorm:"'service_uuid'  char(32)      notnull index comment('services uuid')"`
	Sequence     *int32                 `json:"sequence,omitempty"      xorm:"'sequence'      int           notnull       comment('sequence')"`
	Method       string                 `json:"method"                  xorm:"'method'        varchar(255)  notnull       comment('method')"`
	Args         map[string]interface{} `json:"args,omitempty"          xorm:"'args'          text          null          comment('args')"`
	ResultFilter *string                `json:"result_filter,omitempty" xorm:"'result_filter' varchar(4096) null          comment('result_filter')"`
	Status       *int32                 `json:"status,omitempty"        xorm:"'status'        int           notnull index comment('status')"`
	Result       *string                `json:"result,omitempty"        xorm:"'result'        longtext      null          comment('result')"`
	Started      *time.Time             `json:"started,omitempty"       xorm:"'started'       datetime      null          comment('step start time')"`
	Ended        *time.Time             `json:"ended,omitempty"         xorm:"'ended'         datetime      null          comment('step end time)'"`
}

//DATABASE SCHEMA: SERVICE
type ServiceStep struct {
	metav1.DbMeta       `json:",inline" xorm:"extends"`
	metav1.UuidMeta     `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta    `json:",inline" xorm:"extends"` //inline labelmeta
	ServiceStepProperty `json:",inline" xorm:"extends"` //inline property
}

func (ServiceStep) TableName() string {
	return "service_step"
}

type HttpReqServiceStep_Create_ByService struct {
	Args map[string]interface{} `json:"args,omitempty"`
}

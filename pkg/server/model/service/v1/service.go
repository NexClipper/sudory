package v1

import (
	"math"
	"sort"

	cryptov1 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v1"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
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
	TemplateUuid       string           `json:"template_uuid,omitempty"        xorm:"'template_uuid'        char(32)     notnull index comment('template_uuid')"`        //
	ClusterUuid        string           `json:"cluster_uuid,omitempty"         xorm:"'cluster_uuid'         char(32)     notnull index comment('cluster_uuid')"`         //
	AssignedClientUuid string           `json:"assigned_client_uuid,omitempty" xorm:"'assigned_client_uuid' char(32)     notnull index comment('assigned client_uuid')"` //
	StepCount          *int32           `json:"step_count,omitempty"           xorm:"'step_count'           int          notnull       comment('step_count')"`           //
	StepPosition       *int32           `json:"step_position,omitempty"        xorm:"'step_position'        int          notnull       comment('step_position')"`        //
	Type               *int32           `json:"type,omitempty"                 xorm:"'type'                 int          notnull       comment('type')"`                 //0: once, 1: repeat(epoch, interval)
	Epoch              *int32           `json:"epoch,omitempty"                xorm:"'epoch'                int          notnull       comment('epoch (times)')"`        //-1: infinite
	Interval           *int32           `json:"interval,omitempty"             xorm:"'interval'             int          notnull       comment('interval (sec)')"`       //
	Status             *int32           `json:"status,omitempty"               xorm:"'status'               int          notnull index comment('status')"`               //
	Result             *cryptov1.String `json:"result,omitempty"               xorm:"'result'               longtext     null          comment('result')"`               //실행 결과(정상:'결과', 오류:'오류 메시지')
	SubscribeEvent     string           `json:"subscribe_event,omitempty"      xorm:"'subscribe_event'      varchar(255) null          comment('subscribe event')"`      //서비스 POLL 결과 전달 이벤트 이름
	//  필드 타입에 포인터를 사용하는 이유:
	//    xorm을 사용하면서 초기값을 갖는 타입들은
	//    레코드를 수정할 때 해당 컬럼을 무시하기 때문에
	//    초기값으로 수정이 필요한 필드는 포인터를 사용한다
}

func (property ServiceProperty) ChaniningStep(steps []stepv1.ServiceStep) ServiceProperty {
	ptrInt32 := func(n int32) *int32 {
		return &n
	}

	property.StepCount = ptrInt32(int32(len(steps)))
	property.StepPosition = ptrInt32(0)
	property.Status = ptrInt32(int32(StatusRegist))

	//sort steps by sequence
	sort.Slice(steps, func(i, j int) bool {
		var a, b int32 = math.MaxInt32, math.MaxInt32
		if steps[i].Sequence != nil {
			a = *steps[i].Sequence
		}
		if steps[j].Sequence != nil {
			b = *steps[j].Sequence
		}
		return a < b
	})

	step_status := int32(StatusRegist)
	for i, step := range steps {

		if step.Status == nil {
			step.Status = ptrInt32(int32(StatusRegist))
		}

		if step_status <= *step.Status {
			step_status = *step.Status

			//step이 성공 상태이지만 마지막이 아니면
			//service의 상태를 진행중으로 표시
			if i+1 < len(steps) && int32(StatusSuccess) == step_status {
				step_status = int32(StatusProcessing)
			}

			*property.StepPosition = int32(i) + 1 //step position
			*property.Status = step_status        //step status
		}
	}

	return property
}

//DATABASE SCHEMA: SERVICE
type Service struct {
	metav1.DbMeta    `json:",inline" xorm:"extends"`
	metav1.UuidMeta  `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	ServiceProperty  `json:",inline" xorm:"extends"` //inline property
}

func (Service) TableName() string {
	return "service"
}

type HttpReqService_Create struct {
	metav1.LabelMeta `json:",inline"`                             //inline labelmeta
	TemplateUuid     string                                       `json:"template_uuid"`
	ClusterUuid      string                                       `json:"cluster_uuid"`
	SubscribeEvent   string                                       `json:"subscribe_event,omitempty"`
	Steps            []stepv1.HttpReqServiceStep_Create_ByService `json:"steps,omitempty"`
}

type HttpRspService struct {
	Service `json:",inline"`
	Steps   []stepv1.ServiceStep `json:"steps,omitempty"`
}

// func (object HttpRspService) MarshalJSON() ([]byte, error) {
// 	object.ServiceProperty = object.ChaniningStep(object.Steps)
// 	v := struct {
// 		Service `json:",inline"`
// 		Steps   []stepv1.ServiceStep `json:",inline"`
// 	}{
// 		Service: object.Service,
// 		Steps:   object.Steps,
// 	}

// 	return json.Marshal(v)
// }

type HttpReqService_ClientSide HttpRspService

type HttpRspService_ClientSide HttpRspService

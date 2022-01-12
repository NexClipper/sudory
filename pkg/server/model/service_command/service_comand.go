package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/model/orm"
)

type ServiceCommand struct {
	metav1.DbMeta    `json:",inline" xorm:"extends"`
	metav1.LabelMeta `json:",inline" xorm:"extends"`
	//템플릿 UUID
	TemplateUuid string `json:"template_uuid" orm:"template_uuid"`
	//메소드
	//@example: "kubernetes.deployment.get.v1", "kubernetes.pod.list.v1"
	Method string `json:"methods,omitempty" orm:"methods"`
	//arguments
	Args map[string]string `json:"args,omitempty" orm:"args"`
}

var _ orm.TableName = (*ServiceCommand)(nil)

func (ServiceCommand) TableName() string {
	return "service_command_v1"
}

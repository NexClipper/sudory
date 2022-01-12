package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/model/orm"
)

type ServiceTemplate struct {
	metav1.DbMeta    `json:",inline" xorm:"extends"`
	metav1.LabelMeta `json:",inline" xorm:"extends"`
	//origin
	//@example: predefined, userdefined
	Origin string `json:",omitempty"`
}

var _ orm.TableName = (*ServiceTemplate)(nil)

func (ServiceTemplate) TableName() string {
	return "service_template_v1"
}

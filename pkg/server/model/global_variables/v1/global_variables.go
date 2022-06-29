package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

//GlobalVariables Property
type GlobalVariablesProperty struct {
	Value *string `json:"value" xorm:"'value' text null comment('value')"`
}

//GlobalVariables
type GlobalVariables struct {
	metav1.DbMeta           `json:",inline" xorm:"extends"` //inline dbmeta
	metav1.UuidMeta         `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta        `json:",inline" xorm:"extends"` //inline labelmeta
	GlobalVariablesProperty `json:",inline" xorm:"extends"` //inline property
}

func (GlobalVariables) TableName() string {
	return "global_variables"
}

//HTTP REQUEST BODY: GlobalVariables Update
type HttpReqGlobalVariables_update struct {
	GlobalVariablesProperty `json:",inline" xorm:"extends"` //inline property
}

package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

//GlobalVariant Property
type GlobalVariantProperty struct {
	Value *string `json:"value" xorm:"'value' text null comment('value')"`
}

//GlobalVariant
type GlobalVariant struct {
	metav1.DbMeta         `json:",inline" xorm:"extends"` //inline dbmeta
	metav1.UuidMeta       `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta      `json:",inline" xorm:"extends"` //inline labelmeta
	GlobalVariantProperty `json:",inline" xorm:"extends"` //inline property
}

func (GlobalVariant) TableName() string {
	return "global_variant"
}

//HTTP REQUEST BODY: GlobalVariant Update
type HttpReqGlobalVariant_Update struct {
	GlobalVariantProperty `json:",inline" xorm:"extends"` //inline property
}

package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

//Environment Property
type EnvironmentProperty struct {
	Value *string `json:"value" xorm:"'value' text null comment('value')"`
}

//Environment
type Environment struct {
	metav1.DbMeta       `json:",inline" xorm:"extends"` //inline dbmeta
	metav1.UuidMeta     `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta    `json:",inline" xorm:"extends"` //inline labelmeta
	EnvironmentProperty `json:",inline" xorm:"extends"` //inline property
}

func (Environment) TableName() string {
	return "environment"
}

//HTTP REQUEST BODY: Environment Update
type HttpReqEnvironment_Update struct {
	EnvironmentProperty `json:",inline" xorm:"extends"` //inline property
}

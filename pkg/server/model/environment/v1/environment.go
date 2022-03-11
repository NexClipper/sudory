package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

//Environment Property
type EnvironmentProperty struct {
	Value *string `json:"value,omitempty" xorm:"text null 'value' comment('value')"`
}

//Environment
type Environment struct {
	metav1.UuidMeta     `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta    `json:",inline" xorm:"extends"` //inline labelmeta
	EnvironmentProperty `json:",inline" xorm:"extends"` //inline property
}

//DATABASE SCHEMA: Environment
type DbSchema struct {
	metav1.DbMeta `xorm:"extends"`
	Environment   `xorm:"extends"`
}

func (DbSchema) TableName() string {
	return "environment"
}

//HTTP RESPONSE BODY: Environment
type HttpRspEnvironment struct {
	DbSchema `json:",inline"`
}

// //변환 DbSchema -> Environment
// func TransFromDbSchema(s []DbSchema) []Environment {
// 	var out = make([]Environment, len(s))
// 	for n := range s {
// 		out[n] = s[n].Environment
// 	}
// 	return out
// }

//변환 Environment -> HttpRsp
func TransToHttpRsp(s []DbSchema) []HttpRspEnvironment {
	var out = make([]HttpRspEnvironment, len(s))
	for n := range s {
		out[n].DbSchema = s[n]
	}
	return out
}

type EnvironmentUpdate struct {
	EnvironmentProperty `json:",inline" xorm:"extends"` //inline property
}

//HTTP REQUEST BODY: Environment
type HttpReqEnvironmentUpdate struct {
	EnvironmentUpdate `json:",inline"`
}

package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/model/orm"
)

//Environment Property
type EnvironmentProperty struct {
	Value *string `json:"value,omitempty" xorm:"text null 'value' comment('value')"`
}

//Environment
type Environment struct {
	metav1.LabelMeta    `json:",inline" xorm:"extends"` //inline labelmeta
	EnvironmentProperty `json:",inline" xorm:"extends"` //inline property
}

//DATABASE SCHEMA: Environment
type DbSchemaEnvironment struct {
	metav1.DbMeta `xorm:"extends"`
	Environment   `xorm:"extends"`
}

var _ orm.TableName = (*DbSchemaEnvironment)(nil)

func (DbSchemaEnvironment) TableName() string {
	return "environment"
}

//HTTP REQUEST BODY: Environment
type HttpReqEnvironment struct {
	Environment `json:",inline"`
}

//HTTP RESPONSE BODY: Environment
type HttpRspEnvironment struct {
	Environment `json:",inline"`
}

//변환 DbSchema -> Environment
func TransFormDbSchema(s []DbSchemaEnvironment) []Environment {
	var out = make([]Environment, len(s))
	for n := range s {
		out[n] = s[n].Environment
	}
	return out
}

//변환 Environment -> HttpRsp
func TransToHttpRsp(s []Environment) []HttpRspEnvironment {
	var out = make([]HttpRspEnvironment, len(s))
	for n := range s {
		out[n].Environment = s[n]
	}
	return out
}

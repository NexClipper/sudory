package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/model/orm"
)

//TemplateCommand Property
type TemplateCommandProperty struct {
	//템플릿 UUID
	TemplateUuid string `json:"template_uuid" xorm:"char(32) notnull index 'template_uuid' comment('templates uuid')"`
	//순서
	Sequence *int32 `json:"sequence,omitempty" xorm:"int null default(0) 'sequence' comment('sequence')"`
	//메소드
	//@example: "kubernetes.deployment.get.v1", "kubernetes.pod.list.v1"
	Method *string `json:"method,omitempty" xorm:"varchar(255) null 'method' comment('method')"`
	//arguments
	Args map[string]interface{} `json:"args,omitempty" xorm:"text null 'args' comment('args')"`
}

//MODEL: TEMPLATE_COMMAND
type TemplateCommand struct {
	metav1.UuidMeta         `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta        `json:",inline" xorm:"extends"` //inline labelmeta
	TemplateCommandProperty `json:",inline" xorm:"extends"` //inline property
}

//DATABASE SCHEMA: TEMPLATE_COMMAND
type DbSchema struct {
	metav1.DbMeta   `xorm:"extends"`
	TemplateCommand `xorm:"extends"`
}

var _ orm.TableName = (*DbSchema)(nil)

func (DbSchema) TableName() string {
	return "template_command"
}

//HTTP REQUEST BODY: TEMPLATE
type HttpReqTemplateCommand struct {
	TemplateCommand `json:",inline"`
}

//HTTP RESPONSE BODY: TEMPLATE
type HttpRspTemplateCommand struct {
	DbSchema `json:",inline"`
}

// //변환 TemplateCommand -> DbSchema
// func TransToDbSchema(s []TemplateCommand) []DbSchema {
// 	var out = make([]DbSchema, len(s))
// 	for n, it := range s {
// 		out[n] = DbSchema{TemplateCommand: it}
// 	}
// 	return out
// }

//변환 DbSchema -> TemplateCommand
func TransFromDbSchema(s []DbSchema) []TemplateCommand {
	var out = make([]TemplateCommand, len(s))
	for n, it := range s {
		out[n] = it.TemplateCommand
	}
	return out
}

// //변환 HttpReq -> TemplateCommand
// func TransFormHttpReq(s []HttpReqTemplateCommand) []TemplateCommand {
// 	var out = make([]TemplateCommand, len(s))
// 	for n, it := range s {
// 		out[n] = it.TemplateCommand
// 	}
// 	return out
// }

//변환 TemplateCommand -> HttpRsp
func TransToHttpRsp(s []DbSchema) []HttpRspTemplateCommand {
	var out = make([]HttpRspTemplateCommand, len(s))
	for n, it := range s {
		out[n].DbSchema = it
	}
	return out
}

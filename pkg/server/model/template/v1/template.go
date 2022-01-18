package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/model/orm"
)

//Template Property
type TemplateProperty struct {
	//origin
	//@example: predefined, userdefined
	Origin string `json:"origin,omitempty" xorm:"varchar(255) null 'origin' comment('origin')"`
}

//Template
type Template struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	TemplateProperty `json:",inline" xorm:"extends"` //inline property
}

//DATABASE SCHEMA: TEMPLATE
type DbSchemaTemplate struct {
	metav1.DbMeta `xorm:"extends"`
	Template      `xorm:"extends"`
}

var _ orm.TableName = (*DbSchemaTemplate)(nil)

func (DbSchemaTemplate) TableName() string {
	return "template"
}

//HTTP REQUEST BODY: TEMPLATE
type HttpReqTemplate struct {
	Template `json:",inline"`
}

//HTTP REQUEST BODY: TEMPLATE
type HttpRspTemplate struct {
	Template `json:",inline"`
}

//변환 DbSchema -> Template
func TransFormDbSchema(s []DbSchemaTemplate) []Template {
	var out = make([]Template, len(s))
	for n, it := range s {
		out[n] = it.Template
	}
	return out
}

//변환 Template -> DbSchema
func TransToDbSchema(s []Template) []DbSchemaTemplate {
	var out = make([]DbSchemaTemplate, len(s))
	for n, it := range s {
		out[n] = DbSchemaTemplate{Template: it}
	}
	return out
}

//변환 HttpReq -> Template
func TransFormHttpReq(s []HttpReqTemplate) []Template {
	var out = make([]Template, len(s))
	for n, it := range s {
		out[n] = it.Template
	}
	return out
}

//변환 Template -> HttpRsp
func TransToHttpRsp(s []Template) []HttpRspTemplate {
	var out = make([]HttpRspTemplate, len(s))
	for n, it := range s {
		out[n] = HttpRspTemplate{Template: it}
	}
	return out
}

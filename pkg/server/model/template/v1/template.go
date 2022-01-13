package v1

import (
	"github.com/NexClipper/sudory/pkg/server/model"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/model/orm"
)

//Template Property
type TemplateProperty struct {
	//origin
	//@example: predefined, userdefined
	Origin string `json:"origin,omitempty" orm:"origin"`
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
	return "service_template_v1"
}

//HTTP REQUEST BODY: TEMPLATE
type HttpReqTemplate struct {
	Template `json:",inline"`
}

type HttpReqTemplates []HttpReqTemplate

//HTTP REQUEST BODY: TEMPLATE
type HttpRspTemplate struct {
	Template `json:",inline"`
}

var _ model.Modeler = (*HttpRspTemplate)(nil)

func (HttpRspTemplate) GetType() string {
	return "HTTP RSP TEMPLATE"
}

//HTTP REQUEST BODY: MANY TEMPLATE
type HttpRspTemplates []HttpRspTemplate

var _ model.Modeler = (*HttpRspTemplates)(nil)

func (HttpRspTemplates) GetType() string {
	return "HTTP RSP []TEMPLATE"
}

//변환 DbSchema* -> Template
func TransFormDbSchemaTemplate(s []DbSchemaTemplate) []Template {
	var out = make([]Template, len(s))
	for n, it := range s {
		out[n] = it.Template
	}
	return out
}

//변환 Template -> DbSchema*
func TransToDbSchema(s []Template) []DbSchemaTemplate {
	var out = make([]DbSchemaTemplate, len(s))
	for n, it := range s {
		out[n] = DbSchemaTemplate{Template: it}
	}
	return out
}

//변환 ttpReq* -> Template
func TransFormHttpReqTemplate(s []HttpReqTemplate) []Template {
	var out = make([]Template, len(s))
	for n, it := range s {
		out[n] = it.Template
	}
	return out
}

package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/model/orm"
	commandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
)

//Template Property
type TemplateProperty struct {
	//origin
	//@example: predefined, userdefined
	Origin *string `json:"origin,omitempty" xorm:"varchar(255) null 'origin' comment('origin')"`
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

//HTTP REQUEST BODY: TEMPLATE with template_command
type HttpReqTemplateCreate struct {
	Template `json:",inline"`
	Commands []commandv1.TemplateCommand `json:"commands"`
}

//HTTP RESPONSE BODY: TEMPLATE with template_command
type HttpRspTemplate struct {
	Template `json:",inline"`
	Commands []commandv1.TemplateCommand `json:"commands"`
}

//변환 DbSchema -> Template
func TransFormDbSchema(s []DbSchemaTemplate) []Template {
	var out = make([]Template, len(s))
	for n, it := range s {
		out[n] = it.Template
	}
	return out
}

//Build Template -> HttpRsp
func HttpRspBuilder(length int) (func(a Template, b []commandv1.TemplateCommand), func() []HttpRspTemplate) {
	var pos int = 0
	queue := make([]HttpRspTemplate, length)
	pusher := func(a Template, b []commandv1.TemplateCommand) {
		queue[pos] = HttpRspTemplate{Template: a, Commands: b}
		pos++
	}
	poper := func() []HttpRspTemplate {
		return queue
	}
	return pusher, poper
}

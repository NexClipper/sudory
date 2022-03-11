package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
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
	metav1.UuidMeta  `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	TemplateProperty `json:",inline" xorm:"extends"` //inline property
}

//DATABASE SCHEMA: TEMPLATE
type DbSchema struct {
	metav1.DbMeta `xorm:"extends"`
	Template      `xorm:"extends"`
}

func (DbSchema) TableName() string {
	return "template"
}

//HTTP REQUEST BODY: TEMPLATE
type HttpReqTemplate struct {
	Template `json:",inline"`
}

//HTTP RESPONSE BODY: TEMPLATE
type HttpRspTemplate struct {
	DbSchema `json:",inline"`
}

type TemplateWithCommands struct {
	Template `json:",inline"`
	Commands []commandv1.TemplateCommand `json:"commands"`
}

type DbSchemaTemplateWithCommands struct {
	DbSchema `json:",inline"`
	Commands []commandv1.DbSchema `json:"commands"`
}

//HTTP REQUEST BODY: TEMPLATE with template_command
type HttpReqTemplateWithCommands struct {
	TemplateWithCommands `json:",inline"`
}

//HTTP RESPONSE BODY: TEMPLATE with template_command
type HttpRspTemplateWithCommands struct {
	DbSchemaTemplateWithCommands `json:",inline"`
}

// //변환 DbSchema -> Template
// func TransFormDbSchema(s []DbSchema) []Template {
// 	var out = make([]Template, len(s))
// 	for n, it := range s {
// 		out[n] = it.Template
// 	}
// 	return out
// }

// //Build Template -> HttpRsp
// func HttpRspBuilder(length int) (func(a Template, b []commandv1.TemplateCommand), func() []HttpRspTemplateWithCommands) {
// 	var pos int = 0
// 	queue := make([]HttpRspTemplateWithCommands, length)
// 	pusher := func(a Template, b []commandv1.TemplateCommand) {
// 		queue[pos] = HttpRspTemplateWithCommands{Template: a, Commands: b}
// 		pos++
// 	}
// 	poper := func() []HttpRspTemplateWithCommands {
// 		return queue
// 	}
// 	return pusher, poper
// }

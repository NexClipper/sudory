package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

// TemplateCommandProperty
type TemplateCommandProperty struct {
	TemplateUuid string                 `json:"template_uuid"           xorm:"'template_uuid' char(32)      notnull index comment('templates uuid')"`
	Sequence     *int32                 `json:"sequence,omitempty"      xorm:"'sequence'      int           notnull       comment('sequence')"`
	Method       string                 `json:"method,omitempty"        xorm:"'method'        varchar(255)  notnull       comment('method')"`
	Args         map[string]interface{} `json:"args,omitempty"          xorm:"'args'          text          null          comment('args')"`
	ResultFilter *string                `json:"result_filter,omitempty" xorm:"'result_filter' varchar(4096) null          comment('result_filter')"`
}

// DATABASE SCHEMA: TEMPLATE_COMMAND
type TemplateCommand struct {
	metav1.DbMeta           `json:",inline" xorm:"extends"` //inline dbmeta
	metav1.UuidMeta         `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta        `json:",inline" xorm:"extends"` //inline labelmeta
	TemplateCommandProperty `json:",inline" xorm:"extends"` //inline property
}

func (TemplateCommand) TableName() string {
	return "template_command"
}

type HttpReqTemplateCommand_Create struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	TemplateUuid     string                          `json:"template_uuid"`
	Sequence         int32                           `json:"sequence,omitempty"`
	Method           string                          `json:"method,omitempty"`
	Args             map[string]interface{}          `json:"args,omitempty"`
	ResultFilter     string                          `json:"result_filter,omitempty"`
}

type HttpReqTemplateCommand_Create_ByTemplate struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	Method           string                          `json:"method,omitempty"`
	Args             map[string]interface{}          `json:"args,omitempty"`
	ResultFilter     string                          `json:"result_filter,omitempty"`
}

type HttpReqTemplateCommand_Update struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	Sequence         int32                           `json:"sequence,omitempty"`
	Method           string                          `json:"method,omitempty"`
	Args             map[string]interface{}          `json:"args,omitempty"`
	ResultFilter     string                          `json:"result_filter,omitempty"`
}

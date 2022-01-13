package v1

import (
	"github.com/NexClipper/sudory/pkg/server/model"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/model/orm"
)

//TemplateCommand Property
type TemplateCommandProperty struct {
	//템플릿 UUID
	TemplateUuid string `json:"template_uuid" orm:"template_uuid"`
	//메소드
	//@example: "kubernetes.deployment.get.v1", "kubernetes.pod.list.v1"
	Method string `json:"methods,omitempty" orm:"methods"`
	//arguments
	Args map[string]string `json:"args,omitempty" orm:"args"`
}

//MODEL: TEMPLATE_COMMAND
type TemplateCommand struct {
	metav1.LabelMeta        `json:",inline" xorm:"extends"` //inline labelmeta
	TemplateCommandProperty `json:",inline" xorm:"extends"` //inline property
}

//DATABASE SCHEMA: TEMPLATE_COMMAND
type DbSchemaTemplateCommand struct {
	metav1.DbMeta   `xorm:"extends"`
	TemplateCommand `xorm:"extends"`
}

var _ orm.TableName = (*DbSchemaTemplateCommand)(nil)

func (DbSchemaTemplateCommand) TableName() string {
	return "service_command_v1"
}

//HTTP REQUEST BODY: TEMPLATE
type HttpReqTemplateCommand struct {
	TemplateCommand `json:",inline"`
}

//HTTP RESPONSE BODY: TEMPLATE
type HttpRspTemplateCommand struct {
	TemplateCommand `json:",inline"`
}

var _ model.Modeler = (*HttpRspTemplateCommand)(nil)

func (HttpRspTemplateCommand) GetType() string {
	return "HTTP RSP TEMPLATE_COMMAND"
}

//HTTP RESPONSE BODY: MANY TEMPLATE
type HttpRspTemplateCommands []HttpRspTemplateCommand

var _ model.Modeler = (*HttpRspTemplateCommands)(nil)

func (HttpRspTemplateCommands) GetType() string {
	return "HTTP RSP []TEMPLATE_COMMAND"
}

package database

import (
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
	tcommandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
	"xorm.io/xorm"
)

// Context
//  데이터베이스 컨텍스트
type Context interface {
	Tx() *xorm.Session //get a session
	Close() error      //close session

	//template
	CreateTemplate(m templatev1.DbSchemaTemplate) error
	GetTemplate(uuid string) (*templatev1.DbSchemaTemplate, error)
	FindTemplate(where string, args ...interface{}) ([]templatev1.DbSchemaTemplate, error)
	UpdateTemplate(m templatev1.DbSchemaTemplate) error
	DeleteTemplate(uuid string) error
	//template command
	CreateTemplateCommand(model tcommandv1.DbSchemaTemplateCommand) error
	GetTemplateCommand(uuid string) (*tcommandv1.DbSchemaTemplateCommand, error)
	FindTemplateCommand(where string, args ...interface{}) ([]tcommandv1.DbSchemaTemplateCommand, error)
	UpdateTemplateCommand(model tcommandv1.DbSchemaTemplateCommand) error
	DeleteTemplateCommand(uuid string) error
	//service
	CreateService(m servicev1.DbSchemaService) error
	GetService(uuid string) (*servicev1.DbSchemaService, error)
	FindService(where string, args ...interface{}) ([]servicev1.DbSchemaService, error)
	UpdateService(m servicev1.DbSchemaService) error
	DeleteService(uuid string) error
	//service step
	CreateServiceStep(m stepv1.DbSchemaServiceStep) error
	GetServiceStep(uuid string) (*stepv1.DbSchemaServiceStep, error)
	FindServiceStep(where string, args ...interface{}) ([]stepv1.DbSchemaServiceStep, error)
	UpdateServiceStep(m stepv1.DbSchemaServiceStep) error
	DeleteServiceStep(uuid string) error
}

// Session
type Session struct {
	tx *xorm.Session
}

func NewContext(engine *xorm.Engine) Context {
	return &Session{tx: engine.NewSession()}
}

// Tx
//  get a session
func (me Session) Tx() *xorm.Session {
	return me.tx
}

// Close
//  close session
func (me Session) Close() error {
	return me.tx.Close()
}

// implementation
var _ Context = (*Session)(nil)

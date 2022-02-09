package database

import (
	"github.com/NexClipper/sudory/pkg/server/database/query_parser"
	clientv1 "github.com/NexClipper/sudory/pkg/server/model/client/v1"
	clusterv1 "github.com/NexClipper/sudory/pkg/server/model/cluster/v1"
	envv1 "github.com/NexClipper/sudory/pkg/server/model/environment/v1"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
	commandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
	tokenv1 "github.com/NexClipper/sudory/pkg/server/model/token/v1"
	"xorm.io/xorm"
)

// Context
//  데이터베이스 컨텍스트
type Context interface {
	Tx() *xorm.Session //get a session
	Close() error      //close session

	//client
	CreateClient(m clientv1.DbSchemaClient) error
	GetClient(uuid string) (*clientv1.DbSchemaClient, error)
	FindClient(where string, args ...interface{}) ([]clientv1.DbSchemaClient, error)
	UpdateClient(m clientv1.DbSchemaClient) error
	DeleteClient(uuid string) error
	//cluster
	CreateCluster(m clusterv1.DbSchemaCluster) error
	GetCluster(uuid string) (*clusterv1.DbSchemaCluster, error)
	FindCluster(where string, args ...interface{}) ([]clusterv1.DbSchemaCluster, error)
	UpdateCluster(m clusterv1.DbSchemaCluster) error
	DeleteCluster(uuid string) error
	//template
	CreateTemplate(m templatev1.DbSchemaTemplate) error
	GetTemplate(uuid string) (*templatev1.DbSchemaTemplate, error)
	FindTemplate(where string, args ...interface{}) ([]templatev1.DbSchemaTemplate, error)
	UpdateTemplate(m templatev1.DbSchemaTemplate) error
	DeleteTemplate(uuid string) error
	//template command
	CreateTemplateCommand(model commandv1.DbSchemaTemplateCommand) error
	GetTemplateCommand(uuid string) (*commandv1.DbSchemaTemplateCommand, error)
	FindTemplateCommand(where string, args ...interface{}) ([]commandv1.DbSchemaTemplateCommand, error)
	UpdateTemplateCommand(model commandv1.DbSchemaTemplateCommand) error
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
	//environment
	CreateEnvironment(m envv1.DbSchemaEnvironment) error
	GetEnvironment(uuid string) (*envv1.DbSchemaEnvironment, error)
	FindEnvironment(where string, args ...interface{}) ([]envv1.DbSchemaEnvironment, error)
	UpdateEnvironment(m envv1.DbSchemaEnvironment) error
	DeleteEnvironment(uuid string) error
	//session
	CreateSession(m sessionv1.DbSchemaSession) error
	GetSession(uuid string) (*sessionv1.DbSchemaSession, error)
	FindSession(where string, args ...interface{}) ([]sessionv1.DbSchemaSession, error)
	QuerySession(query *query_parser.QueryParser) ([]sessionv1.DbSchemaSession, error)
	UpdateSession(m sessionv1.DbSchemaSession) error
	DeleteSession(uuid string) error
	//token
	CreateToken(m tokenv1.DbSchemaToken) error
	GetToken(uuid string) (*tokenv1.DbSchemaToken, error)
	FindToken(where string, args ...interface{}) ([]tokenv1.DbSchemaToken, error)
	QueryToken(query *query_parser.QueryParser) ([]tokenv1.DbSchemaToken, error)
	UpdateToken(m tokenv1.DbSchemaToken) error
	DeleteToken(uuid string) error
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

// @title SUDORY
// @version 0.0.1
// @description this is a sudory server.
// @contact.url https://nexclipper.io
// @contact.email jaehoon@nexclipper.io
// @securityDefinitions.apikey ClientAuth
// @in header
// @name x-sudory-client-token
// @description token for client api
// @securityDefinitions.apikey XAuthToken
// @in header
// @name x_auth_token
// @description limit for access sudory api
package route

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/NexClipper/sudory/pkg/server/config"
	"github.com/NexClipper/sudory/pkg/server/control"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	mysqlFlavor "github.com/NexClipper/sudory/pkg/server/database/vanilla/excute/dialects/mysql"
	pprof "github.com/NexClipper/sudory/pkg/server/route/pprof"
	"github.com/NexClipper/sudory/pkg/version"
	"github.com/pkg/errors"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/NexClipper/sudory/pkg/server/route/docs"
)

func init() {
	//swago docs version
	docs.SwaggerInfo.Version = version.Version
}

type Route struct {
	e *echo.Echo

	Port                   int32
	UseTls                 bool
	TlsCertificateFilename string
	TlsPrivateKeyFilename  string
}

func New(cfg *config.Config, db *sql.DB) *Route {

	e := echo.New()
	ctl := control.NewVanilla(db, excute.GetSqlExcutor(mysqlFlavor.Dialect()))

	//echo cors config
	e.Use(echoCORSConfig(cfg))

	if true {
		//request id generator
		e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Generator: func() func() string {
				var (
					id uint64
				)
				return func() string {
					id := atomic.AddUint64(&id, 1)
					return fmt.Sprintf("%d", id)
				}
			}(),
		}))
	}
	//logger
	if true {
		e.Use(echoServiceLogger(config.LoggerInfoOutput))
	}

	// enable request 'Content-Encoding' header handler
	if true {
		e.Use(middleware.Decompress())
	}

	// enable response 'Content-Encoding' header handler
	if true {
		e.Use(middleware.Gzip())
	}

	//echo error handler
	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		echoErrorResponder(err, ctx)
		echoErrorLogger(err, ctx)
	}
	//echo recover
	e.Use(echoRecover())

	//swago
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	// pprof
	pprof.Wrap("/debug/pprof", e)

	{
		// /client/auth*
		e.POST("/client/auth", ctl.AuthClient)
		e.POST("/sudory/client/auth", ctl.AuthClient)

		group := e.Group("")
		// @Security ClientSessionToken
		group.Use(ClientSessionToken(db, excute.GetSqlExcutor(mysqlFlavor.Dialect())))
		// /client/service*
		group.GET("/client/service", ctl.PollingService)
		group.PUT("/client/service", ctl.UpdateService)
		group.GET("/sudory/client/service", ctl.PollingService)
		group.PUT("/sudory/client/service", ctl.UpdateService)
	}

	{
		group := e.Group("")
		group.Use(XAuthToken(cfg)) // @Security XAuthToken

		groupV2 := e.Group("/v2")
		groupV2.Use(XAuthToken(cfg)) // @Security XAuthToken

		// /server/auth*
		group.POST("/server/tenant", ctl.Tenant)
		group.POST("/sudory/server/tenant", ctl.Tenant)

		// /server/template*
		group.GET("/server/template", ctl.FindTemplate)
		group.GET("/server/template/:uuid", ctl.GetTemplate)
		group.GET("/sudory/server/template", ctl.FindTemplate)
		group.GET("/sudory/server/template/:uuid", ctl.GetTemplate)

		// /server/template/:template_uuid/command*
		group.GET("/server/template/:template_uuid/command", ctl.ListTemplateCommand)
		group.GET("/server/template/:template_uuid/command/:uuid", ctl.GetTemplateCommand)
		group.GET("/sudory/server/template/:template_uuid/command", ctl.ListTemplateCommand)
		group.GET("/sudory/server/template/:template_uuid/command/:uuid", ctl.GetTemplateCommand)

		// /server/template_recipe*
		group.GET("/server/template_recipe", ctl.FindTemplateRecipe)
		group.GET("/sudory/server/template_recipe", ctl.FindTemplateRecipe)

		// /v2/server/template*
		groupV2.GET("/server/template", ctl.FindTemplate_v2)
		groupV2.GET("/server/template/:uuid", ctl.GetTemplate_v2)
		groupV2.GET("/sudory/server/template", ctl.FindTemplate_v2)
		groupV2.GET("/sudory/server/template/:uuid", ctl.GetTemplate_v2)

		// /v2/server/template_command*
		groupV2.GET("/server/template_command", ctl.FindTemplateCommand_v2)
		groupV2.GET("/server/template_command/:uuid", ctl.GetTemplateCommand_v2)
		groupV2.GET("/sudory/server/template_command", ctl.FindTemplateCommand_v2)
		groupV2.GET("/sudory/server/template_command/:uuid", ctl.GetTemplateCommand_v2)

		// /server/global_variables*
		group.GET("/server/global_variables", ctl.FindGlobalVariables)
		group.GET("/server/global_variables/:uuid", ctl.GetGlobalVariables)
		group.PUT("/server/global_variables/:uuid", ctl.UpdateGlobalVariablesValue)
		group.GET("/sudory/server/global_variables", ctl.FindGlobalVariables)
		group.GET("/sudory/server/global_variables/:uuid", ctl.GetGlobalVariables)
		group.PUT("/sudory/server/global_variables/:uuid", ctl.UpdateGlobalVariablesValue)

	}

	{
		group := e.Group("")
		group.Use(ServiceAuthorizationBearerToken()) // @Security ServiceAuthorizationBearerToken

		groupV2 := e.Group("/v2")
		groupV2.Use(ServiceAuthorizationBearerToken()) // @Security ServiceAuthorizationBearerToken

		// /server/cluster*
		group.GET("/server/cluster", ctl.FindCluster)
		group.GET("/server/cluster/:uuid", ctl.GetCluster)
		group.POST("/server/cluster", ctl.CreateCluster)
		group.PUT("/server/cluster/:uuid", ctl.UpdateCluster)
		group.PUT("/server/cluster/:uuid/polling/regular", ctl.UpdateClusterPollingRegular)
		group.PUT("/server/cluster/:uuid/polling/smart", ctl.UpdateClusterPollingSmart)
		group.DELETE("/server/cluster/:uuid", ctl.DeleteCluster)
		group.GET("/sudory/server/cluster", ctl.FindCluster)
		group.GET("/sudory/server/cluster/:uuid", ctl.GetCluster)
		group.POST("/sudory/server/cluster", ctl.CreateCluster)
		group.PUT("/sudory/server/cluster/:uuid", ctl.UpdateCluster)
		group.PUT("/sudory/server/cluster/:uuid/polling/regular", ctl.UpdateClusterPollingRegular)
		group.PUT("/sudory/server/cluster/:uuid/polling/smart", ctl.UpdateClusterPollingSmart)
		group.DELETE("/sudory/server/cluster/:uuid", ctl.DeleteCluster)

		// /server/service*
		group.GET("/server/service", ctl.FindService)
		group.GET("/server/service/:uuid", ctl.GetService)
		group.POST("/server/service", ctl.CreateService)
		group.GET("/server/service/:uuid/result", ctl.GetServiceResult)
		group.GET("/sudory/server/service", ctl.FindService)
		group.GET("/sudory/server/service/:uuid", ctl.GetService)
		group.POST("/sudory/server/service", ctl.CreateService)
		group.GET("/sudory/server/service/:uuid/result", ctl.GetServiceResult)

		// /server/service_step*
		group.GET("/server/service/step", ctl.FindServiceStep)
		group.GET("/server/service/:uuid/step", ctl.GetServiceSteps)
		group.GET("/server/service/:uuid/step/:sequence", ctl.GetServiceStep)
		group.GET("/sudory/server/service/step", ctl.FindServiceStep)
		group.GET("/sudory/server/service/:uuid/step", ctl.GetServiceSteps)
		group.GET("/sudory/server/service/:uuid/step/:sequence", ctl.GetServiceStep)

		// /v2/server/service*
		groupV2.GET("/server/service", ctl.FindService_v2)
		groupV2.GET("/server/service/:uuid", ctl.GetService_v2)
		groupV2.POST("/server/service", ctl.CreateService_v2)
		groupV2.GET("/server/service/:uuid/result", ctl.GetServiceResult_v2)
		groupV2.GET("/sudory/server/service", ctl.FindService_v2)
		groupV2.GET("/sudory/server/service/:uuid", ctl.GetService_v2)
		groupV2.POST("/sudory/server/service", ctl.CreateService_v2)
		groupV2.GET("/sudory/server/service/:uuid/result", ctl.GetServiceResult_v2)

		// /server/session*
		group.GET("/server/session", ctl.FindSession)
		group.GET("/server/session/:uuid", ctl.GetSession)
		group.DELETE("/server/session/:uuid", ctl.DeleteSession)
		group.GET("/server/session/cluster/:cluster_uuid/alive", ctl.AliveClusterSession)
		group.GET("/sudory/server/session", ctl.FindSession)
		group.GET("/sudory/server/session/:uuid", ctl.GetSession)
		group.DELETE("/sudory/server/session/:uuid", ctl.DeleteSession)
		group.GET("/sudory/server/session/cluster/:cluster_uuid/alive", ctl.AliveClusterSession)

		// /server/cluster_token*
		group.GET("/server/cluster_token", ctl.FindClusterToken)
		group.GET("/server/cluster_token/:uuid", ctl.GetClusterToken)
		group.PUT("/server/cluster_token/:uuid/label", ctl.UpdateClusterTokenLabel)
		group.DELETE("/server/cluster_token/:uuid", ctl.DeleteClusterToken)
		group.POST("/server/cluster_token", ctl.CreateClusterToken)
		group.PUT("/server/cluster_token/:uuid/refresh", ctl.RefreshClusterTokenTime)
		group.PUT("/server/cluster_token/:uuid/expire", ctl.ExpireClusterToken)
		group.GET("/sudory/server/cluster_token", ctl.FindClusterToken)
		group.GET("/sudory/server/cluster_token/:uuid", ctl.GetClusterToken)
		group.PUT("/sudory/server/cluster_token/:uuid/label", ctl.UpdateClusterTokenLabel)
		group.DELETE("/sudory/server/cluster_token/:uuid", ctl.DeleteClusterToken)
		group.POST("/sudory/server/cluster_token", ctl.CreateClusterToken)
		group.PUT("/sudory/server/cluster_token/:uuid/refresh", ctl.RefreshClusterTokenTime)
		group.PUT("/sudory/server/cluster_token/:uuid/expire", ctl.ExpireClusterToken)

		// /server/channels*
		group.POST("/server/channels", ctl.CreateChannel)
		group.GET("/server/channels", ctl.FindChannel)
		group.GET("/server/channels/:uuid", ctl.GetChannel)
		group.PUT("/server/channels/:uuid", ctl.UpdateChannel)
		group.DELETE("/server/channels/:uuid", ctl.DeleteChannel)
		group.POST("/sudory/server/channels", ctl.CreateChannel)
		group.GET("/sudory/server/channels", ctl.FindChannel)
		group.GET("/sudory/server/channels/:uuid", ctl.GetChannel)
		group.PUT("/sudory/server/channels/:uuid", ctl.UpdateChannel)
		group.DELETE("/sudory/server/channels/:uuid", ctl.DeleteChannel)

		// /server/channels/:uuid/notifiers/*
		group.GET("/server/channels/:uuid/notifiers/edge", ctl.GetChannelNotifierEdge)
		group.PUT("/server/channels/:uuid/notifiers/console", ctl.UpdateChannelNotifierConsole)
		group.PUT("/server/channels/:uuid/notifiers/rabbitmq", ctl.UpdateChannelNotifierRabbitMq)
		group.PUT("/server/channels/:uuid/notifiers/webhook", ctl.UpdateChannelNotifierWebhook)
		group.PUT("/server/channels/:uuid/notifiers/slackhook", ctl.UpdateChannelNotifierSlackhook)
		group.GET("/sudory/server/channels/:uuid/notifiers/edge", ctl.GetChannelNotifierEdge)
		group.PUT("/sudory/server/channels/:uuid/notifiers/console", ctl.UpdateChannelNotifierConsole)
		group.PUT("/sudory/server/channels/:uuid/notifiers/rabbitmq", ctl.UpdateChannelNotifierRabbitMq)
		group.PUT("/sudory/server/channels/:uuid/notifiers/webhook", ctl.UpdateChannelNotifierWebhook)
		group.PUT("/sudory/server/channels/:uuid/notifiers/slackhook", ctl.UpdateChannelNotifierSlackhook)

		// /server/channels/status
		group.GET("/server/channels/status", ctl.FindChannelStatus)
		group.GET("/sudory/server/channels/status", ctl.FindChannelStatus)

		// /server/channels/:uuid/status*
		group.GET("/server/channels/:uuid/status", ctl.ListChannelStatus)
		group.DELETE("/server/channels/:uuid/status/purge", ctl.PurgeChannelStatus)
		group.PUT("/server/channels/:uuid/status/option", ctl.UpdateChannelStatusOption)
		group.GET("/server/channels/:uuid/status/option", ctl.GetChannelStatusOption)
		group.GET("/sudory/server/channels/:uuid/status", ctl.ListChannelStatus)
		group.DELETE("/sudory/server/channels/:uuid/status/purge", ctl.PurgeChannelStatus)
		group.PUT("/sudory/server/channels/:uuid/status/option", ctl.UpdateChannelStatusOption)
		group.GET("/sudory/server/channels/:uuid/status/option", ctl.GetChannelStatusOption)

		// /server/channels/:uuid/format*
		group.GET("/server/channels/:uuid/format", ctl.GetChannelFormat)
		group.PUT("/server/channels/:uuid/format", ctl.UpdateChannelFormat)
		group.GET("/sudory/server/channels/:uuid/format", ctl.GetChannelFormat)
		group.PUT("/sudory/server/channels/:uuid/format", ctl.UpdateChannelFormat)

	}

	return &Route{
		e:                      e,
		Port:                   cfg.Host.Port,
		UseTls:                 cfg.Host.TlsEnable,
		TlsCertificateFilename: cfg.Host.TlsCertificateFilename,
		TlsPrivateKeyFilename:  cfg.Host.TlsPrivateKeyFilename,
	}
}

func (r *Route) Start() error {
	if r.UseTls {
		crt, err := os.ReadFile(r.TlsCertificateFilename)
		if err != nil {
			err = errors.Wrapf(err, "faild to read tls certificate file name=%v", r.TlsCertificateFilename)
			return err
		}
		key, err := os.ReadFile(r.TlsPrivateKeyFilename)
		if err != nil {
			err = errors.Wrapf(err, "faild to read tls privateKey file name=%v", r.TlsPrivateKeyFilename)
			return err
		}

		return StartTLS(r.e, r.Port, crt, key)
	}
	return Start(r.e, r.Port)
}

func Start(e *echo.Echo, port int32) error {
	go func() {
		address := fmt.Sprintf(":%d", port)
		if err := e.Start(address); err != nil {
			e.Logger.Info("shut down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func StartTLS(e *echo.Echo, port int32, crt, key []byte) error {
	go func() {
		address := fmt.Sprintf(":%d", port)
		if err := e.StartTLS(address, crt, key); err != nil {
			e.Logger.Info("shut down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

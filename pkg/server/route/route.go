// @title SUDORY
// @version 0.0.1
// @description this is a sudory server.
// @contact.url https://nexclipper.io
// @contact.email jaehoon@nexclipper.io
package route

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"time"

	"github.com/NexClipper/logger"
	"github.com/NexClipper/sudory/pkg/server/config"
	"github.com/NexClipper/sudory/pkg/server/control"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/version"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/NexClipper/sudory/pkg/server/route/docs"
)

type Route struct {
	e *echo.Echo
}

func New(cfg *config.Config, db *database.DBManipulator) *Route {
	e := echo.New()
	controller := control.New(db)

	//echo cors config
	echoCORSConfig(e, cfg)

	if true {
		e.Logger.SetOutput(config.LoggerInfoOutput)

		//echo logger
		format := `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}",` +
			`"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n"
		logconfig := middleware.DefaultLoggerConfig
		logconfig.Output = config.LoggerInfoOutput
		logconfig.Format = format
		e.Use(middleware.LoggerWithConfig(logconfig))
	}

	if true {
		//request id generator
		e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Generator: func() func() string {
				var (
					ticketId uint64
				)
				return func() string {
					atomic.AddUint64(&ticketId, 1)
					return fmt.Sprintf("%d", ticketId)
				}
			}(),
		}))
	}

	//echo error handler
	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		echoErrorHandler(err, ctx)
		echoErrorHandlerLogger(err, ctx)
	}
	// defaultHeader := `{"time":"${time_rfc3339_nano}","id":"${id}","level":"${level}","prefix":"${prefix}"}`
	// e.Logger.SetHeader(defaultHeader)
	// e.Logger.SetOutput(config.LoggerInfoOutput)

	//logger

	//swago
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	//swago docs version
	docs.SwaggerInfo.Version = version.Version

	client_group := e.Group("/client")
	{
		//route /client/service*
		client_group.PUT("/service", controller.PollService, func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) (err error) {
				if err := controller.VerifyClientSessionToken(c); err != nil {
					return err
				}

				if err := next(c); err != nil {
					return err
				}

				return nil
			}
		})
		//route /client/auth*
		client_group.POST("/auth", controller.AuthClient)
	}

	server_group := e.Group("/server")
	{
		if strings.Compare(version.Version, "dev") != 0 {
			server_group.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) (err error) {
					const key = "x_auth_token"
					header := c.Request().Header.Get(key)

					if len(header) == 0 {
						return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
							fmt.Errorf("not found request header%s",
								logs.KVL(
									"key", key,
								)))
					}

					if strings.Compare(header, "SUDORY") != 0 {
						return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
							fmt.Errorf("not found request header%s",
								logs.KVL(
									"key", key,
								)))
					}

					if err := next(c); err != nil {
						return err
					}

					return nil
				}
			})
		}

		//route /server/client*
		server_group.GET("/client", controller.FindClient)
		server_group.GET("/client/:uuid", controller.GetClient)
		server_group.DELETE("/client/:uuid", controller.DeleteClient)
		//route /server/cluster*
		server_group.GET("/cluster", controller.FindCluster)
		server_group.GET("/cluster/:uuid", controller.GetCluster)
		server_group.POST("/cluster", controller.CreateCluster)
		server_group.PUT("/cluster/:uuid", controller.UpdateCluster)
		server_group.PUT("/cluster/:uuid/polling/raguler", controller.UpdateClusterPollingRaguler)
		server_group.PUT("/cluster/:uuid/polling/smart", controller.UpdateClusterPollingSmart)
		server_group.DELETE("/cluster/:uuid", controller.DeleteCluster)
		//route /server/template*
		server_group.GET("/template", controller.FindTemplate)
		server_group.GET("/template/:uuid", controller.GetTemplate)
		server_group.POST("/template", controller.CreateTemplate)
		server_group.PUT("/template/:uuid", controller.UpdateTemplate)
		server_group.DELETE("/template/:uuid", controller.DeleteTemplate)
		//route /server/template/:template_uuid/command*
		server_group.GET("/template/:template_uuid/command", controller.FindTemplateCommand)
		server_group.GET("/template/:template_uuid/command/:uuid", controller.GetTemplateCommand)
		server_group.POST("/template/:template_uuid/command", controller.CreateTemplateCommand)
		server_group.PUT("/template/:template_uuid/command/:uuid", controller.UpdateTemplateCommand)
		server_group.DELETE("/template/:template_uuid/command/:uuid", controller.DeleteTemplateCommand)
		//route /server/template_recipe*
		server_group.GET("/template_recipe", controller.FindTemplateRecipe)
		//route /server/service*
		server_group.GET("/service", controller.FindService)
		server_group.GET("/service/:uuid", controller.GetService)
		server_group.GET("/service/:uuid/result", controller.GetServiceResult)
		server_group.POST("/service", controller.CreateService)
		// router.e.PUT("/service/:uuid", controller.UpdateService)
		server_group.DELETE("/service/:uuid", controller.DeleteService)
		//route /server/service_step*
		server_group.GET("/service/:service_uuid/step", controller.FindServiceStep)
		server_group.GET("/service/:service_uuid/step/:uuid", controller.GetServiceStep)
		//route /server/environment*
		server_group.GET("/environment", controller.FindEnvironment)
		server_group.GET("/environment/:uuid", controller.GetEnvironment)
		server_group.PUT("/environment/:uuid", controller.UpdateEnvironmentValue)
		//route /server/session*
		server_group.GET("/session", controller.FindSession)
		server_group.GET("/session/:uuid", controller.GetSession)
		server_group.DELETE("/session/:uuid", controller.DeleteSession)
		//route /server/token*
		server_group.GET("/token", controller.FindToken)
		server_group.GET("/token/:uuid", controller.GetToken)
		server_group.PUT("/token/:uuid/label", controller.UpdateTokenLabel)
		server_group.DELETE("/token/:uuid", controller.DeleteToken)
		//route /server/token/cluster/*
		server_group.POST("/token/cluster", controller.CreateClusterToken)
		server_group.PUT("/token/cluster/:uuid/refresh", controller.RefreshClusterTokenTime)
		server_group.PUT("/token/cluster/:uuid/expire", controller.ExpireClusterToken)
	}
	/*TODO: 라우트 연결 기능 구현
	e.POST("/cluster/:id/token", controller.CreateToken)

	//route /server/catalogue
	e.GET("/catalogue", controller.GetCatalogue)

	e.POST("/client/regist", controller.CreateClient)
	*/

	return &Route{e: e}
}

func (r *Route) Start(port int32) error {
	go func() {
		address := fmt.Sprintf(":%d", port)
		if err := r.e.Start(address); err != nil {
			r.e.Logger.Info("shut down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := r.e.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func echoCORSConfig(_echo *echo.Echo, _config *config.Config) {
	CORSConfig := middleware.DefaultCORSConfig //use default cors config
	//cors allow orign
	if 0 < len(_config.CORSConfig.AllowOrigins) {
		origins := strings.Split(_config.CORSConfig.AllowOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}

		CORSConfig.AllowOrigins = origins
	}
	//cors allow method
	if 0 < len(_config.CORSConfig.AllowMethods) {
		methods := strings.Split(_config.CORSConfig.AllowMethods, ",")
		for i := range methods {
			methods[i] = strings.TrimSpace(methods[i]) //trim space
			methods[i] = strings.ToUpper(methods[i])   //to upper
		}

		CORSConfig.AllowMethods = methods
	}

	_echo.Use(middleware.CORSWithConfig(CORSConfig))

	fmt.Fprintf(os.Stdout, "ECHO CORS Config:\n")
	fmt.Fprintf(os.Stdout, "- allow-origins: %v\n", strings.Join(CORSConfig.AllowOrigins, ", "))
	fmt.Fprintf(os.Stdout, "- allow-methods: %v\n", strings.Join(CORSConfig.AllowMethods, ", "))
	fmt.Fprintf(os.Stdout, "%s\n", strings.Repeat("_", 40))
}

func echoErrorHandler(err error, ctx echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if he.Internal != nil {
			err = he.Internal
		}
	}

	ctx.JSON(code, map[string]interface{}{
		"code": code,
		// "status":     http.StatusText(code),
		"message": err.Error(),
	})
}

func echoErrorHandlerLogger(err error, ctx echo.Context) {
	code := -1
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		err = he.Internal
	}

	var stack string = "none"
	logs.CauseIter(err, func(err error) {
		logs.StackIter(err, func(s string) {
			stack = s
		})
	})

	id := ctx.Response().Header().Get(echo.HeaderXRequestID)

	reqbody := echoutil.Body(ctx)

	logger.Error(fmt.Errorf("%w%v", err, logs.KVL(
		"id", id,
		"code", code,
		"reqbody", reqbody,
		"stack", stack,
	)))
}

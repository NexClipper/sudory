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

	"github.com/NexClipper/sudory/pkg/server/config"
	"github.com/NexClipper/sudory/pkg/server/control"
	"github.com/NexClipper/sudory/pkg/server/database"
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
		// echoErrorHandlerLogger(err, ctx)
	}
	// defaultHeader := `{"time":"${time_rfc3339_nano}","id":"${id}","level":"${level}","prefix":"${prefix}"}`
	// e.Logger.SetHeader(defaultHeader)
	// e.Logger.SetOutput(config.LoggerInfoOutput)

	//logger

	//swago
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	//swago docs version
	docs.SwaggerInfo.Version = version.Version

	//route /client/service*
	e.PUT("/client/service", controller.PollService, func(next echo.HandlerFunc) echo.HandlerFunc {
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
	e.POST("/client/auth", controller.AuthClient)

	//route /server/client*
	e.GET("/server/client", controller.FindClient)
	e.GET("/server/client/:uuid", controller.GetClient)
	e.DELETE("/server/client/:uuid", controller.DeleteClient)
	//route /server/cluster*
	e.GET("/server/cluster", controller.FindCluster)
	e.GET("/server/cluster/:uuid", controller.GetCluster)
	e.POST("/server/cluster", controller.CreateCluster)
	e.PUT("/server/cluster/:uuid", controller.UpdateCluster)
	e.PUT("/server/cluster/:uuid/polling/raguler", controller.UpdateClusterPollingRaguler)
	e.PUT("/server/cluster/:uuid/polling/smart", controller.UpdateClusterPollingSmart)
	e.DELETE("/server/cluster/:uuid", controller.DeleteCluster)
	//route /server/template*
	e.GET("/server/template", controller.FindTemplate)
	e.GET("/server/template/:uuid", controller.GetTemplate)
	e.POST("/server/template", controller.CreateTemplate)
	e.PUT("/server/template/:uuid", controller.UpdateTemplate)
	e.DELETE("/server/template/:uuid", controller.DeleteTemplate)
	//route /server/template/:template_uuid/command*
	e.GET("/server/template/:template_uuid/command", controller.FindTemplateCommand)
	e.GET("/server/template/:template_uuid/command/:uuid", controller.GetTemplateCommand)
	e.POST("/server/template/:template_uuid/command", controller.CreateTemplateCommand)
	e.PUT("/server/template/:template_uuid/command/:uuid", controller.UpdateTemplateCommand)
	e.DELETE("/server/template/:template_uuid/command/:uuid", controller.DeleteTemplateCommand)
	//route /server/template_recipe*
	e.GET("/server/template_recipe", controller.FindTemplateRecipe)
	//route /server/service*
	e.GET("/server/service", controller.FindService)
	e.GET("/server/service/:uuid", controller.GetService)
	e.GET("/server/service/:uuid/result", controller.GetServiceResult)
	e.POST("/server/service", controller.CreateService)
	// router.e.PUT("/server/service/:uuid", controller.UpdateService)
	e.DELETE("/server/service/:uuid", controller.DeleteService)
	//route /server/service_step*
	e.GET("/server/service/:service_uuid/step", controller.FindServiceStep)
	e.GET("/server/service/:service_uuid/step/:uuid", controller.GetServiceStep)
	//route /server/environment*
	e.GET("/server/environment", controller.FindEnvironment)
	e.GET("/server/environment/:uuid", controller.GetEnvironment)
	e.PUT("/server/environment/:uuid", controller.UpdateEnvironmentValue)
	//route /server/session*
	e.GET("/server/session", controller.FindSession)
	e.GET("/server/session/:uuid", controller.GetSession)
	e.DELETE("/server/session/:uuid", controller.DeleteSession)
	//route /server/token*
	e.GET("/server/token", controller.FindToken)
	e.GET("/server/token/:uuid", controller.GetToken)
	e.PUT("/server/token/:uuid/label", controller.UpdateTokenLabel)
	e.DELETE("/server/token/:uuid", controller.DeleteToken)
	//route /server/token/cluster/*
	e.POST("/server/token/cluster", controller.CreateClusterToken)
	e.PUT("/server/token/cluster/:uuid/refresh", controller.RefreshClusterTokenTime)
	e.PUT("/server/token/cluster/:uuid/expire", controller.ExpireClusterToken)

	/*TODO: 라우트 연결 기능 구현
	e.POST("/server/cluster/:id/token", controller.CreateToken)

	//route /server/catalogue
	e.GET("/server/catalogue", controller.GetCatalogue)

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

// func echoErrorHandlerLogger(err error, ctx echo.Context) {
// 	code := http.StatusInternalServerError
// 	if he, ok := err.(*echo.HTTPError); ok {
// 		code = he.Code
// 		err = he.Internal
// 	}

// 	id := ctx.Response().Header().Get(echo.HeaderXRequestID)
// 	reqbody := echoutil.Body(ctx)

// 	logger.Error(fmt.Errorf("%w%v", err, logs.KVL(
// 		"id", id,
// 		"code", code,
// 		"reqbody", string(reqbody),
// 	)))
// }

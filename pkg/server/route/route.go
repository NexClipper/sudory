// @title SUDORY
// @version 0.0.1
// @description this is a sudory server.
// @contact.url https://nexclipper.io
// @contact.email jaehoon@nexclipper.io
package route

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/NexClipper/logger"
	"github.com/NexClipper/sudory/pkg/server/config"
	"github.com/NexClipper/sudory/pkg/server/control"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/version"
	"github.com/pkg/errors"

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
		e.Use(echoLogger(config.LoggerInfoOutput))
	}

	//echo error handler
	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		echoErrorHandlerResponse(err, ctx)
		echoErrorHandlerLogger(err, ctx)
	}
	//echo recover
	e.Use(echoRecover())

	//swago
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	//swago docs version
	docs.SwaggerInfo.Version = version.Version

	//"/client"
	{
		group := e.Group("/client")
		group.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) (err error) {
				switch c.Path() {
				case "/client/auth":
					//do noting
				default:
					if err := controller.VerifyClientSessionToken(c); err != nil {
						return err
					}
				}

				if err := next(c); err != nil {
					return err
				}

				return nil
			}
		})

		//route /client/service*
		group.GET("/service", controller.PollingService)
		group.PUT("/service", controller.UpdateService)
		//route /client/auth*
		group.POST("/auth", controller.AuthClient)
	}

	//"/server"
	{
		group := e.Group("/server")

		if cfg.Host.XAuthToken {
			group.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) (err error) {
					const key = "x_auth_token"
					header := c.Request().Header.Get(key)

					if len(header) == 0 {
						return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
							errors.Errorf("not found request header%s",
								logs.KVL(
									"key", key,
								)))
					}

					if strings.Compare(header, "SUDORY") != 0 {
						return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
							errors.Errorf("not found request header%s",
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

		//route /server/cluster*
		group.GET("/cluster", controller.FindCluster)
		group.GET("/cluster/:uuid", controller.GetCluster)
		group.POST("/cluster", controller.CreateCluster)
		group.PUT("/cluster/:uuid", controller.UpdateCluster)
		group.PUT("/cluster/:uuid/polling/raguler", controller.UpdateClusterPollingRaguler)
		group.PUT("/cluster/:uuid/polling/smart", controller.UpdateClusterPollingSmart)
		group.DELETE("/cluster/:uuid", controller.DeleteCluster)
		//route /server/template*
		group.GET("/template", controller.FindTemplate)
		group.GET("/template/:uuid", controller.GetTemplate)
		group.POST("/template", controller.CreateTemplate)
		group.PUT("/template/:uuid", controller.UpdateTemplate)
		group.DELETE("/template/:uuid", controller.DeleteTemplate)
		//route /server/template/:template_uuid/command*
		group.GET("/template/:template_uuid/command", controller.FindTemplateCommand)
		group.GET("/template/:template_uuid/command/:uuid", controller.GetTemplateCommand)
		group.POST("/template/:template_uuid/command", controller.CreateTemplateCommand)
		group.PUT("/template/:template_uuid/command/:uuid", controller.UpdateTemplateCommand)
		group.DELETE("/template/:template_uuid/command/:uuid", controller.DeleteTemplateCommand)
		//route /server/template_recipe*
		group.GET("/template_recipe", controller.FindTemplateRecipe)
		//route /server/service*
		group.GET("/service", controller.FindService)
		group.GET("/service/:uuid", controller.GetService)
		group.GET("/service/:uuid/result", controller.GetServiceResult)
		group.POST("/service", controller.CreateService)
		// router.e.PUT("/service/:uuid", controller.UpdateService)
		group.DELETE("/service/:uuid", controller.DeleteService)
		//route /server/service_step*
		group.GET("/service/:service_uuid/step", controller.FindServiceStep)
		group.GET("/service/:service_uuid/step/:uuid", controller.GetServiceStep)
		//route /server/global_variant*
		group.GET("/global_variant", controller.FindGlobalVariant)
		group.GET("/global_variant/:uuid", controller.GetGlobalVariant)
		group.PUT("/global_variant/:uuid", controller.UpdateGlobalVariantValue)
		//route /server/session*
		group.GET("/session", controller.FindSession)
		group.GET("/session/:uuid", controller.GetSession)
		group.DELETE("/session/:uuid", controller.DeleteSession)
		//route /server/cluster_token*
		group.GET("/cluster_token", controller.FindClusterToken)
		group.GET("/cluster_token/:uuid", controller.GetClusterToken)
		group.PUT("/cluster_token/:uuid/label", controller.UpdateClusterTokenLabel)
		group.DELETE("/cluster_token/:uuid", controller.DeleteClusterToken)
		group.POST("/cluster_token", controller.CreateClusterToken)
		group.PUT("/cluster_token/:uuid/refresh", controller.RefreshClusterTokenTime)
		group.PUT("/cluster_token/:uuid/expire", controller.ExpireClusterToken)
		//server/event*
		group.POST("/channel", controller.CreateChannel)
		group.GET("/channel", controller.FindChannel)
		group.GET("/channel/:uuid", controller.GetChannel)
		group.PUT("/channel/:uuid", controller.UpdateChannel)
		group.GET("/channel/:uuid/notifier_edges", controller.ListChannelNotifierEdges)
		group.PUT("/channel/:uuid/notifier_edges/add", controller.AddChannelNotifierEdge)
		group.PUT("/channel/:uuid/notifier_edges/sub", controller.SubChannelNotifierEdge)
		group.DELETE("/channel/:uuid", controller.DeleteChannel)
		//server/channel_notifier*
		group.POST("/channel_notifier/console", controller.CreateChannelNotifierConsole)
		group.POST("/channel_notifier/webhook", controller.CreateChannelNotifierWebhook)
		group.POST("/channel_notifier/rabbitmq", controller.CreateChannelNotifierRabbitMq)
		group.GET("/channel_notifier/console", controller.FindChannelNotifierConsole)
		group.GET("/channel_notifier/webhook", controller.FindChannelNotifierWebhook)
		group.GET("/channel_notifier/rabbitmq", controller.FindChannelNotifierRabbitmq)
		group.GET("/channel_notifier/console/:uuid", controller.GetChannelNotifierConsole)
		group.GET("/channel_notifier/webhook/:uuid", controller.GetChannelNotifierWebhook)
		group.GET("/channel_notifier/rabbitmq/:uuid", controller.GetChannelNotifierRabbitmq)
		group.PUT("/channel_notifier/console/:uuid", controller.UpdateChannelNotifierConsole)
		group.PUT("/channel_notifier/webhook/:uuid", controller.UpdateChannelNotifierWebhook)
		group.PUT("/channel_notifier/rabbitmq/:uuid", controller.UpdateChannelNotifierRabbitMq)
		group.DELETE("/channel_notifier/console/:uuid", controller.DeleteChannelNotifierConsole)
		group.DELETE("/channel_notifier/webhook/:uuid", controller.DeleteChannelNotifierWebhook)
		group.DELETE("/channel_notifier/rabbitmq/:uuid", controller.DeleteChannelNotifierRabbitmq)
		//server/channel_notifier_status*
		group.GET("/channel_notifier_status", controller.FindChannelNotifierStatus)
		group.DELETE("/channel_notifier_status/:uuid", controller.DeleteChannelNotifierStatus)
	}

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
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := r.e.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func echoCORSConfig(_config *config.Config) echo.MiddlewareFunc {
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

	fmt.Fprintf(os.Stdout, "ECHO CORS Config:\n")

	tabwrite := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	tabwrite.Write([]byte(strings.Join([]string{
		"", "allow-origins",
	}, "\t") + "\n"))
	tabwrite.Write([]byte(strings.Join([]string{
		"-", strings.Join(CORSConfig.AllowOrigins, ", "),
	}, "\t") + "\n"))
	tabwrite.Write([]byte(strings.Join([]string{
		"", "allow-methods",
	}, "\t") + "\n"))
	tabwrite.Write([]byte(strings.Join([]string{
		"-", strings.Join(CORSConfig.AllowMethods, ", "),
	}, "\t") + "\n"))

	tabwrite.Flush()

	fmt.Fprintln(os.Stdout, strings.Repeat("_", 40))

	// fmt.Fprintf(os.Stdout, "-   allow-origins: %v\n", strings.Join(CORSConfig.AllowOrigins, ", "))
	// fmt.Fprintf(os.Stdout, "-   allow-methods: %v\n", strings.Join(CORSConfig.AllowMethods, ", "))
	// fmt.Fprintf(os.Stdout, "%s\n", strings.Repeat("_", 40))

	return middleware.CORSWithConfig(CORSConfig)

}

func echoErrorHandlerResponse(err error, ctx echo.Context) {
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
	nullstring := func(p *string) (s string) {
		s = fmt.Sprintf("%v", p)
		if p != nil {
			s = *p
		}
		return
	}

	code := -1
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		err = he.Internal
	}

	var stack *string
	//stack for surface
	logs.StackIter(err, func(s string) {
		stack = &s
	})
	//stack for internal
	logs.CauseIter(err, func(err error) {
		logs.StackIter(err, func(s string) {
			stack = &s
		})
	})

	id := ctx.Response().Header().Get(echo.HeaderXRequestID)

	reqbody := echoutil.Body(ctx)

	logger.Error(fmt.Errorf("%w%v", err, logs.KVL(
		"id", id,
		"code", code,
		"reqbody", reqbody,
		"stack", nullstring(stack),
	)))
}

func echoLogger(w io.Writer) echo.MiddlewareFunc {
	//echo logger
	format := `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
		`"host":"${host}","method":"${method}","uri":"${uri}",` +
		`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}",` +
		`"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n"
	logconfig := middleware.DefaultLoggerConfig
	logconfig.Output = w
	logconfig.Format = format

	return middleware.LoggerWithConfig(logconfig)
}

func echoRecover(skipper ...middleware.Skipper) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			for _, skipper := range skipper {
				if skipper(c) {
					return next(c)
				}
			}

			defer func() {
				if r := recover(); r != nil {
					if r == http.ErrAbortHandler {
						panic(r)
					}

					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					err = errors.Wrapf(err, "echo recovered")

					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}

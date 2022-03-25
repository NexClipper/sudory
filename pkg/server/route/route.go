// @title SUDORY
// @version 0.0.1
// @description this is a sudory server.
// @contact.url https://nexclipper.io
// @contact.email jaehoon@nexclipper.io
package route

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
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
	echo_ := echo.New()
	router := &Route{e: echo_}
	controller := control.New(db)

	//echo cors config
	echoCORSConfig(echo_, cfg)

	//swago
	router.e.GET("/swagger/*", echoSwagger.WrapHandler)
	//swago docs version
	docs.SwaggerInfo.Version = version.Version

	//route /client/service*
	router.e.PUT("/client/service", controller.PollService())
	//route /client/auth*
	router.e.POST("/client/auth", controller.AuthClient())

	//route /server/client*
	router.e.GET("/server/client", controller.FindClient())
	router.e.GET("/server/client/:uuid", controller.GetClient())
	router.e.DELETE("/server/client/:uuid", controller.DeleteClient())
	//route /server/cluster*
	router.e.GET("/server/cluster", controller.FindCluster())
	router.e.GET("/server/cluster/:uuid", controller.GetCluster())
	router.e.POST("/server/cluster", controller.CreateCluster())
	router.e.PUT("/server/cluster/:uuid", controller.UpdateCluster())
	router.e.DELETE("/server/cluster/:uuid", controller.DeleteCluster())
	//route /server/template*
	router.e.GET("/server/template", controller.FindTemplate())
	router.e.GET("/server/template/:uuid", controller.GetTemplate())
	router.e.POST("/server/template", controller.CreateTemplate())
	router.e.PUT("/server/template/:uuid", controller.UpdateTemplate())
	router.e.DELETE("/server/template/:uuid", controller.DeleteTemplate())
	//route /server/template/:template_uuid/command*
	router.e.GET("/server/template/:template_uuid/command", controller.FindTemplateCommand())
	router.e.GET("/server/template/:template_uuid/command/:uuid", controller.GetTemplateCommand())
	router.e.POST("/server/template/:template_uuid/command", controller.CreateTemplateCommand())
	router.e.PUT("/server/template/:template_uuid/command/:uuid", controller.UpdateTemplateCommand())
	router.e.DELETE("/server/template/:template_uuid/command/:uuid", controller.DeleteTemplateCommand())
	//route /server/template_recipe*
	router.e.GET("/server/template_recipe", controller.FindTemplateRecipe())
	//route /server/service*
	router.e.GET("/server/service", controller.FindService())
	router.e.GET("/server/service/:uuid", controller.GetService())
	router.e.GET("/server/service/:uuid/result", controller.GetServiceResult())
	router.e.POST("/server/service", controller.CreateService())
	// router.e.PUT("/server/service/:uuid", controller.UpdateService())
	router.e.DELETE("/server/service/:uuid", controller.DeleteService())
	//route /server/service_step*
	router.e.GET("/server/service/:service_uuid/step", controller.FindServiceStep())
	router.e.GET("/server/service/:service_uuid/step/:uuid", controller.GetServiceStep())
	//route /server/environment*
	router.e.GET("/server/environment", controller.FindEnvironment())
	router.e.GET("/server/environment/:uuid", controller.GetEnvironment())
	router.e.PUT("/server/environment/:uuid", controller.UpdateEnvironmentValue())
	//route /server/session*
	router.e.GET("/server/session", controller.FindSession())
	router.e.GET("/server/session/:uuid", controller.GetSession())
	router.e.DELETE("/server/session/:uuid", controller.DeleteSession())
	//route /server/token*
	router.e.GET("/server/token", controller.FindToken())
	router.e.GET("/server/token/:uuid", controller.GetToken())
	router.e.PUT("/server/token/:uuid/label", controller.UpdateTokenLabel())
	router.e.DELETE("/server/token/:uuid", controller.DeleteToken())
	//route /server/token/cluster/*
	router.e.POST("/server/token/cluster", controller.CreateClusterToken())
	router.e.PUT("/server/token/cluster/:uuid/refresh", controller.RefreshClusterTokenTime())
	router.e.PUT("/server/token/cluster/:uuid/expire", controller.ExpireClusterToken())

	/*TODO: 라우트 연결 기능 구현
	router.e.POST("/server/cluster/:id/token", controller.CreateToken)

	//route /server/catalogue
	router.e.GET("/server/catalogue", controller.GetCatalogue)

	router.e.POST("/client/regist", controller.CreateClient)
	*/

	return router
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

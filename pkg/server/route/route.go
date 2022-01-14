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
	"time"

	"github.com/NexClipper/sudory/pkg/server/config"
	"github.com/NexClipper/sudory/pkg/server/control"
	"github.com/NexClipper/sudory/pkg/server/database"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/NexClipper/sudory/pkg/server/route/docs"
)

type Route struct {
	e  *echo.Echo
	db *database.DBManipulator
}

func New(cfg *config.Config, db *database.DBManipulator) *Route {
	router := &Route{e: echo.New(), db: db}
	controller := control.New(db)

	router.e.POST("/server/cluster", controller.CreateCluster)
	router.e.GET("/server/cluster/:id", controller.GetCluster)
	router.e.POST("/server/cluster/:id/token", controller.CreateToken)

	//route /server/template
	router.e.GET("/server/template", controller.FindTemplate)
	router.e.GET("/server/template/:uuid", controller.GetTemplate)
	router.e.POST("/server/template", controller.CreateTemplate)
	router.e.PUT("/server/template", controller.UpdateTemplate)
	router.e.DELETE("/server/template/:uuid", controller.DeleteTemplate)

	//route /server/template_command
	router.e.GET("/server/template_command/:uuid", controller.GetTemplateCommand)
	router.e.POST("/server/template_command", controller.CreateTemplateCommand)
	router.e.PUT("/server/template_command", controller.UpdateTemplateCommand)
	router.e.DELETE("/server/template_command/:uuid", controller.DeleteTemplateCommand)

	router.e.GET("/server/catalogue", controller.GetCatalogue)
	router.e.POST("/server/service", controller.CreateService)

	router.e.POST("/client/regist", controller.CreateClient)
	router.e.PUT("/client/service", controller.GetService)

	router.e.GET("/swagger/*", echoSwagger.WrapHandler)

	return router
}

func (r *Route) Start(port int32) error {
	go func() {
		address := fmt.Sprintf(":%d", port)
		if err := r.e.Start(address); err != nil {
			r.e.Logger.Info("shut down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	select {
	case <-quit:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := r.e.Shutdown(ctx); err != nil {
			r.db.Close()
			return err
		}
		r.db.Close()
	}

	return nil
}

package route

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/NexClipper/sudory-prototype-r1/pkg/config"
	"github.com/NexClipper/sudory-prototype-r1/pkg/control"
	"github.com/NexClipper/sudory-prototype-r1/pkg/database"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/NexClipper/sudory-prototype-r1/pkg/route/docs"
)

type Route struct {
	e *echo.Echo
}

func New(cfg *config.Config, db *database.DBManipulator) *Route {
	router := &Route{e: echo.New()}
	controller := control.New(db)

	router.e.POST("/cluster", controller.CreateCluster)

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
			return err
		}
	}

	return nil
}

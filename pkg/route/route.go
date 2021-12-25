package route

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/NexClipper/sudory-prototype-r1/pkg/control"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Route struct {
	e *echo.Echo
}

func New() *Route {
	router := &Route{e: echo.New()}

	controller := control.New()

	router.e.POST("/clusters", controller.CreateCluster)

	router.e.GET("/swagger/*", echoSwagger.WrapHandler)

	return router
}

func (r *Route) Start() error {
	go func() {
		port := fmt.Sprintf(":%d", 8099)
		if err := r.e.Start(port); err != nil {
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

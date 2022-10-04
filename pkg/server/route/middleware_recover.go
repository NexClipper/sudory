package route

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

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

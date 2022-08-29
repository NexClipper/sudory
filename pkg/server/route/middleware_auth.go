package route

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/config"
	"github.com/NexClipper/sudory/pkg/server/control"
	flavor "github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/resolvers/mysql"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func XAuthToken(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if cfg.Host.XAuthToken {
				// x_auth_token
				const (
					key   = "x_auth_token"
					value = "SUDORY"
				)
				header_value := c.Request().Header.Get(key)

				if len(header_value) == 0 || strings.Compare(header_value, value) != 0 {
					err = errors.Errorf("not found request header%s",
						logs.KVL(
							"key", key,
						))
				}
				if err != nil {
					return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(errors.Wrapf(err,
						http.StatusText(http.StatusUnauthorized)))
				}
				return
			}

			if err := next(c); err != nil {
				return err
			}

			return
		}
	}
}

func ClientSessionToken(db *sql.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if _, err := control.GetClientSessionClaims(c, db, flavor.Dialect()); err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(errors.Wrapf(err,
					http.StatusText(http.StatusUnauthorized)))
			}

			if err := next(c); err != nil {
				return err
			}
			return
		}
	}
}

func ServiceAuthorizationBearerToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			// service Authorization: Bearer
			if _, err := control.GetServiceAuthorizationClaims(c); err != nil {
				// set http header "WWW-Authenticate" "Bearer"
				echoutil.SeHttpHeader(c.Response().Header(), "WWW-Authenticate", "Bearer")

				return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(errors.Wrapf(err,
					http.StatusText(http.StatusUnauthorized)))
			}

			if err := next(c); err != nil {
				return err
			}

			return
		}
	}
}

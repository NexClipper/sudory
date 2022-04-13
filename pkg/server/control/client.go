package control

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Find Client
// @Description Find client
// @Accept      json
// @Produce     json
// @Tags        server/client
// @Router      /server/client [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success 200 {array} v1.Client
func (ctl Control) FindClient(ctx echo.Context) error {
	r, err := vault.NewClient(ctl.NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "find client"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// Get Client
// @Description Get a client
// @Accept      json
// @Produce     json
// @Tags        server/client
// @Router      /server/client/{uuid} [get]
// @Param       x_auth_token  header string false "client session token"
// @Param       uuid          path   string true  "Client 의 Uuid"
// @Success 200 {object} v1.Client
func (ctl Control) GetClient(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := vault.NewClient(ctl.NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(errors.Wrapf(err, "get client"))
	}

	return ctx.JSON(http.StatusOK, *r)
}

// Delete Client
// @Description Delete a client
// @Accept      json
// @Produce     json
// @Tags        server/client
// @Router      /server/client/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Client 의 Uuid"
// @Success 200
func (ctl Control) DeleteClient(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		if err := vault.NewClient(db).Delete(uuid); err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "delete client"))
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, OK())
}

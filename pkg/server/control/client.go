package control

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	clientv1 "github.com/NexClipper/sudory/pkg/server/model/client/v1"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Find Client
// @Description Find client
// @Accept      json
// @Produce     json
// @Tags        server/client
// @Router      /server/client [get]
// @Param       q query string false "query  pkg/server/database/prepared/README.md"
// @Param       o query string false "order  pkg/server/database/prepared/README.md"
// @Param       p query string false "paging pkg/server/database/prepared/README.md"
// @Success 200 {array} v1.HttpRspClient
func (ctl Control) FindClient(ctx echo.Context) error {
	// r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
	r, err := vault.NewClient(ctl.NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrapf(err, "find client%s",
			logs.KVL(
				"query", echoutil.QueryParamString(ctx),
			)))
	}

	return ctx.JSON(http.StatusOK, clientv1.TransToHttpRsp(r))
}

// Get Client
// @Description Get a client
// @Accept      json
// @Produce     json
// @Tags        server/client
// @Router      /server/client/{uuid} [get]
// @Param       uuid          path string true "Client 의 Uuid"
// @Success 200 {object} v1.HttpRspClient
func (ctl Control) GetClient(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := vault.NewClient(ctl.NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrapf(err, "get client%s",
			logs.KVL(
				"uuid", uuid,
			)))
	}

	return ctx.JSON(http.StatusOK, clientv1.HttpRspClient{DbSchema: *r})
}

// Delete Client
// @Description Delete a client
// @Accept      json
// @Produce     json
// @Tags        server/client
// @Router      /server/client/{uuid} [delete]
// @Param       uuid path string true "Client 의 Uuid"
// @Success 200
func (ctl Control) DeleteClient(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		err := vault.NewClient(db).Delete(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "delete client%s",
				logs.KVL(
					"uuid", uuid,
				))
		}

		return nil, nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, OK())
}

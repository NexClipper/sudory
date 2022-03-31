package control

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Find Session
// @Description Find Session
// @Accept      json
// @Produce     json
// @Tags        server/session
// @Router      /server/session [get]
// @Param       q query string false "query  pkg/server/database/prepared/README.md"
// @Param       o query string false "order  pkg/server/database/prepared/README.md"
// @Param       p query string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.HttpRspSession
func (ctl Control) FindSession(ctx echo.Context) error {
	r, err := vault.NewSession(ctl.NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			errors.Wrapf(err, "find session%s",
				logs.KVL(
					"query", echoutil.QueryParamString(ctx),
				)))
	}

	return ctx.JSON(http.StatusOK, sessionv1.TransToHttpRsp(r))
}

// Get Session
// @Description Get a Session
// @Accept      json
// @Produce     json
// @Tags        server/session
// @Router      /server/session/{uuid} [get]
// @Param       uuid          path string true "Session 의 Uuid"
// @Success     200 {object} v1.HttpRspSession
func (ctl Control) GetSession(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := vault.NewSession(ctl.NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			errors.Wrapf(err, "get session%s",
				logs.KVL(
					"uuid", uuid,
				)))
	}

	return ctx.JSON(http.StatusOK, sessionv1.HttpRspSession{DbSchema: *r})
}

// Delete Session
// @Description Delete a Session
// @Accept      json
// @Produce     json
// @Tags        server/session
// @Router      /server/session/{uuid} [delete]
// @Param       uuid path string true "Session 의 Uuid"
// @Success     200
func (ctl Control) DeleteSession(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		err := vault.NewSession(db).Delete(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "delete session%s",
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

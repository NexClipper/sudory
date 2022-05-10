package control

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

// Find Session
// @Description Find Session
// @Accept      json
// @Produce     json
// @Tags        server/session
// @Router      /server/session [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.Session
func (ctl Control) FindSession(ctx echo.Context) error {
	r, err := vault.NewSession(ctl.db.Engine().NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "find session"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// Get Session
// @Description Get a Session
// @Accept      json
// @Produce     json
// @Tags        server/session
// @Router      /server/session/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Session 의 Uuid"
// @Success     200 {object} v1.Session
func (ctl Control) GetSession(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := vault.NewSession(ctl.db.Engine().NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get session"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// Delete Session
// @Description Delete a Session
// @Accept      json
// @Produce     json
// @Tags        server/session
// @Router      /server/session/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Session 의 Uuid"
// @Success     200
func (ctl Control) DeleteSession(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		err := vault.NewSession(tx).Delete(uuid)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "delete session"))
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, OK())
}

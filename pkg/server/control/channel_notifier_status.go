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

// @deprecated
// @Description Find channel notifier status
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/channel_notifier_status
// @Router      /server/channel_notifier_status [get]
// @Param       q            query  string false "query  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       o            query  string false "order  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       p            query  string false "paging github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Success     200 {array} v1.NotifierStatus
func (ctl Control) FindChannelNotifierStatus(ctx echo.Context) error {
	//find event
	status, err := vault.NewNotifierStatus(ctl.db.Engine().NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "query event"))
	}

	return ctx.JSON(http.StatusOK, status)

}

// @deprecated
// @Description Delete a channel notifier status
// @Security    XAuthToken
// @Accept json
// @Produce json
// @Tags server/channel_notifier_status
// @Router /server/channel_notifier_status/{uuid} [delete]
// @Param       uuid         path   string true  "EventNotifierStatus 의 Uuid"
// @Success 200
func (ctl Control) DeleteChannelNotifierStatus(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		//delete event
		if err := vault.NewNotifierStatus(tx).Delete(uuid); err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "delete event notifier status"))
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, OK())
}

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

// @Description Find Event Nofitier Status
// @Accept      json
// @Produce     json
// @Tags        server/event_notifier_status
// @Router      /server/event_notifier_status [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.EventNotifierStatus
func (ctl Control) FindEventNofitierStatus(ctx echo.Context) error {
	//find event
	status, err := vault.NewEventNotifierStatus(ctl.db.Engine().NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "query event"))
	}

	return ctx.JSON(http.StatusOK, status)

}

// @Description Delete a Event Nofitier Status
// @Accept json
// @Produce json
// @Tags server/event_notifier_status
// @Router /server/event_notifier_status/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "EventNofitierStatus 의 Uuid"
// @Success 200
func (ctl Control) DeleteEventNofitierStatus(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		//delete event
		if err := vault.NewEventNotifierStatus(tx).Delete(uuid); err != nil {
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

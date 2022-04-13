package control

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// FindServiceStep
// @Description Find Service Steps
// @Accept      json
// @Produce     json
// @Tags        server/service_step
// @Router      /server/service/{service_uuid}/step [get]
// @Param       x_auth_token header string false "client session token"
// @Param       service_uuid path   string true  "ServiceStep 의 service_uuid"
// @Success     200 {array} v1.ServiceStep
func (ctl Control) FindServiceStep(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__SERVICE_UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __SERVICE_UUID__,
				)))
	}

	where := "service_uuid = ?"
	service_uuid := echoutil.Param(ctx)[__SERVICE_UUID__]

	r, err := vault.NewServiceStep(ctl.NewSession()).Find(where, service_uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "find service step"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// Get ServiceStep
// @Description Get a Service Step
// @Accept      json
// @Produce     json
// @Tags        server/service_step
// @Router      /server/service/{service_uuid}/step/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       service_uuid path   string true  "ServiceStep 의 service_uuid"
// @Param       uuid         path   string true  "ServiceStep 의 Uuid"
// @Success     200 {object} v1.ServiceStep
func (ctl Control) GetServiceStep(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__SERVICE_UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __SERVICE_UUID__,
				)))
	}

	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	_ = echoutil.Param(ctx)[__SERVICE_UUID__]
	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := vault.NewServiceStep(ctl.NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get service step"))
	}

	return ctx.JSON(http.StatusOK, r)
}

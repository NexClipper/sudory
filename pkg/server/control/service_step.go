package control

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/control/vanilla"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	servicev2 "github.com/NexClipper/sudory/pkg/server/model/service/v2"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// // FindServiceStep
// // @Description Find Service Steps
// // @Accept      json
// // @Produce     json
// // @Tags        server/service_step
// // @Router      /server/service/{service_uuid}/step [get]
// // @Param       x_auth_token header string false "client session token"
// // @Param       service_uuid path   string true  "ServiceStep 의 service_uuid"
// // @Success     200 {array} v2.HttpRsp_ServiceStep
// func (ctl Control) FindServiceStep(ctx echo.Context) error {
// 	if len(echoutil.Param(ctx)[__SERVICE_UUID__]) == 0 {
// 		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
// 			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
// 				logs.KVL(
// 					ParamLog(__SERVICE_UUID__, echoutil.Param(ctx)[__SERVICE_UUID__])...,
// 				)))
// 	}

// 	where := "service_uuid = ?"
// 	service_uuid := echoutil.Param(ctx)[__SERVICE_UUID__]

// 	r, err := vault.NewServiceStep(ctl.db.Engine().NewSession()).Find(where, service_uuid)
// 	if err != nil {
// 		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
// 			errors.Wrapf(err, "find service step"))
// 	}

// 	return ctx.JSON(http.StatusOK, r)
// }

// FindServiceStep
// @Description Find Service Steps
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/step [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v2.HttpRsp_ServiceStep
func (ctl ControlVanilla) FindServiceStep(ctx echo.Context) error {
	conditions, args, err := ParseDecoration(echoutil.QueryParam(ctx))
	err = errors.Wrapf(err, "ParseDecoration%v", logs.KVL(
		"query", echoutil.QueryParamString(ctx),
	))

	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	rsps := make([]servicev2.HttpRsp_ServiceStep, 0, __INIT_RECORD_CAPACITY__)

	Do(&err, func() (err error) {
		cond := vanilla.NewCond(
			strings.Join(conditions, "\n"),
			args...,
		)

		rsps, err = find_service_steps(ctl.DB(), *cond)
		err = errors.Wrapf(err, "find service steps")
		return
	})

	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, rsps)
}

// Get []ServiceStep
// @Description Get Service Steps
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid}/step [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid path   string true  "ServiceStep 의 uuid"
// @Success     200 {array} v2.HttpRsp_ServiceStep
func (ctl ControlVanilla) GetServiceSteps(ctx echo.Context) (err error) {
	Do(&err, func() (err error) {
		if len(echoutil.Param(ctx)[__UUID__]) == 0 {
			err = ErrorInvalidRequestParameter()
		}
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
			))
		return
	})

	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	steps, err := get_service_steps(ctl.DB(), uuid)
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, steps)
}

// Get ServiceStep
// @Description Get Service Step
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid}/step/{sequence} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid     path string true  "ServiceStep 의 uuid"
// @Param       sequence path string true  "ServiceStep 의 sequence"
// @Success     200 {object} v2.HttpRsp_ServiceStep
func (ctl ControlVanilla) GetServiceStep(ctx echo.Context) (err error) {
	Do(&err, func() (err error) {
		if len(echoutil.Param(ctx)[__UUID__]) == 0 {
			err = ErrorInvalidRequestParameter()
		}
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
			))
		return
	})

	Do(&err, func() (err error) {
		if len(echoutil.Param(ctx)[__SEQUENCE__]) == 0 {
			err = ErrorInvalidRequestParameter()
		}
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(__SEQUENCE__, echoutil.Param(ctx)[__SEQUENCE__])...,
			))
		return
	})

	uuid := echoutil.Param(ctx)[__UUID__]
	var sequence int
	Do(&err, func() (err error) {
		s := echoutil.Param(ctx)[__SEQUENCE__]
		sequence, err = strconv.Atoi(s)
		err = errors.Wrapf(err, "sequence atoi%s", logs.KVL(
			"a", s,
		))
		return
	})

	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	var step servicev2.ServiceStep_tangled
	Do(&err, func() (err error) {
		step, err = get_service_step(ctl.DB(), uuid, sequence)
		return
	})

	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, step)
}

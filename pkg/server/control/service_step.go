package control

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	service "github.com/NexClipper/sudory/pkg/server/model/service/v3"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @Description Find Service Steps
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/step [get]
// @Param       q            query  string false "query  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       o            query  string false "order  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       p            query  string false "paging github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Success     200 {array} service.HttpRsp_ServiceStep
func (ctl ControlVanilla) FindServiceStep(ctx echo.Context) error {
	q, err := stmt.ConditionLexer.Parse(echoutil.QueryParam(ctx)["q"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	o, err := stmt.OrderLexer.Parse(echoutil.QueryParam(ctx)["o"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	p, err := stmt.PaginationLexer.Parse(echoutil.QueryParam(ctx)["p"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	if p == nil {
		p = stmt.Limit(__DEFAULT_DECORATION_LIMIT__)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}
	rsp := make([]service.HttpRsp_ServiceStep, 0, state.ENV__INIT_SLICE_CAPACITY__())
	step_table := service.TableNameWithTenant_ServiceStep(claims.Hash)
	var step service.ServiceStep

	err = ctl.dialect.QueryRows(step_table, step.ColumnNames(), q, o, p)(ctx.Request().Context(), ctl)(
		func(scan excute.Scanner, _ int) error {
			err := step.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			rsp = append(rsp, step)

			return err
		})
	if err != nil {
		return err
	}

	// make response body
	return ctx.JSON(http.StatusOK, []service.HttpRsp_ServiceStep(rsp))
}

// @Description Get Service Steps
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid}/step [get]
// @Param       uuid         path   string true  "ServiceStep 의 uuid"
// @Success     200 {array} service.HttpRsp_ServiceStep
func (ctl ControlVanilla) GetServiceSteps(ctx echo.Context) (err error) {
	var uuid string
	if uuid = echoutil.Param(ctx)[__UUID__]; len(uuid) == 0 {
		err = ErrorInvalidRequestParameter
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
			))
		return HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	// get service step
	step_table := service.TableNameWithTenant_ServiceStep(claims.Hash)
	var step service.ServiceStep
	step_cond := stmt.And(
		stmt.Equal("uuid", uuid),
	)

	steps := make([]service.HttpRsp_ServiceStep, 0, state.ENV__INIT_SLICE_CAPACITY__())

	err = ctl.dialect.QueryRows(step_table, step.ColumnNames(), step_cond, nil, nil)(ctx.Request().Context(), ctl)(
		func(scan excute.Scanner, _ int) error {
			err := step.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			steps = append(steps, step)

			return err
		})
	if err != nil {
		return err
	}

	sort.Slice(steps, func(i, j int) bool {
		return steps[i].Sequence < steps[j].Sequence
	})

	// make response body
	return ctx.JSON(http.StatusOK, []service.HttpRsp_ServiceStep(steps))
}

// @Description Get Service Step
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid}/step/{sequence} [get]
// @Param       uuid     path string true "ServiceStep 의 uuid"
// @Param       sequence path string true "ServiceStep 의 sequence"
// @Success     200 {object} service.HttpRsp_ServiceStep
func (ctl ControlVanilla) GetServiceStep(ctx echo.Context) (err error) {
	var uuid string
	if uuid = echoutil.Param(ctx)[__UUID__]; len(uuid) == 0 {
		err = ErrorInvalidRequestParameter
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
			))
		return HttpError(err, http.StatusBadRequest)
	}
	var sequence int
	if sequence, err = strconv.Atoi(echoutil.Param(ctx)[__SEQUENCE__]); err != nil {
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(__SEQUENCE__, echoutil.Param(ctx)[__SEQUENCE__])...,
			))
		return HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	// get service step
	step_table := service.TableNameWithTenant_ServiceStep(claims.Hash)
	var step service.ServiceStep
	step_cond := stmt.And(
		stmt.Equal("uuid", uuid),
		stmt.Equal("seq", sequence),
	)

	err = ctl.dialect.QueryRow(step_table, step.ColumnNames(), step_cond, nil, nil)(ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := step.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return err
	}

	// make response body
	return ctx.JSON(http.StatusOK, service.HttpRsp_ServiceStep(step))
}

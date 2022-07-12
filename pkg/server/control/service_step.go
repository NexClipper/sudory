package control

import (
	"net/http"
	"strconv"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	servicev2 "github.com/NexClipper/sudory/pkg/server/model/service/v2"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

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
	q, o, p, err := ParseDecoration(echoutil.QueryParam(ctx))
	err = errors.Wrapf(err, "ParseDecoration%v", logs.KVL(
		"query", echoutil.QueryParamString(ctx),
	))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	rsps := make([]servicev2.HttpRsp_ServiceStep, 0, __INIT_SLICE_CAPACITY__())

	step := servicev2.ServiceStep_tangled{}
	stmt := vanilla.Stmt.Select(step.TableName(), step.ColumnNames(), q, o, p)
	err = stmt.QueryRows(ctl)(func(scan vanilla.Scanner, _ int) (err error) {
		err = step.Scan(scan)
		if err == nil {
			rsps = append(rsps, servicev2.HttpRsp_ServiceStep{
				ServiceStep_tangled: step,
			})
		}
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
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter()
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))

	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	steps := make([]servicev2.HttpRsp_ServiceStep, 0, __INIT_SLICE_CAPACITY__())

	eq_uuid := vanilla.Equal("uuid", uuid).Parse()

	step := servicev2.ServiceStep_tangled{}
	stmt := vanilla.Stmt.Select(step.TableName(), step.ColumnNames(), eq_uuid, nil, nil)
	err = stmt.QueryRows(ctl)(func(scan vanilla.Scanner, _ int) (err error) {
		err = step.Scan(scan)
		if err == nil {
			steps = append(steps, servicev2.HttpRsp_ServiceStep{
				ServiceStep_tangled: step,
			})
		}
		return
	})

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
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter()
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
		))

	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	if len(echoutil.Param(ctx)[__SEQUENCE__]) == 0 {
		err = ErrorInvalidRequestParameter()
	}
	err = errors.Wrapf(err, "valid param%s",
		logs.KVL(
			ParamLog(__SEQUENCE__, echoutil.Param(ctx)[__SEQUENCE__])...,
		))

	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	s := echoutil.Param(ctx)[__SEQUENCE__]
	sequence, err := strconv.Atoi(s)
	err = errors.Wrapf(err, "sequence atoi%s", logs.KVL(
		"a", s,
	))

	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	eq_uuid := vanilla.Equal("uuid", uuid)
	eq_sequence := vanilla.Equal("sequence", sequence)
	q := vanilla.And(eq_uuid, eq_sequence).Parse()

	step := servicev2.ServiceStep_tangled{}
	stmt := vanilla.Stmt.Select(step.TableName(), step.ColumnNames(), q, nil, nil)
	err = stmt.QueryRow(ctl)(func(s vanilla.Scanner) (err error) {
		err = step.Scan(s)
		return
	})

	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, servicev2.HttpRsp_ServiceStep{ServiceStep_tangled: step})
}

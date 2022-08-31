package control

import (
	"net/http"
	"strconv"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	servicev3 "github.com/NexClipper/sudory/pkg/server/model/service/v3"
	"github.com/NexClipper/sudory/pkg/server/status/state"
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
// @Success     200 {array} v3.HttpRsp_ServiceStep
func (ctl ControlVanilla) FindServiceStep(ctx echo.Context) error {
	q, o, p, err := ParseDecoration(echoutil.QueryParam(ctx))
	if err != nil {
		err = errors.Wrapf(err, "ParseDecoration%v", logs.KVL(
			"query", echoutil.QueryParamString(ctx),
		))
		return HttpError(err, http.StatusBadRequest)
	}

	// find service
	tablename := `(
SELECT A.cluster_uuid,A.uuid,A.seq,A.pdate,A.timestamp,A.name,A.summary,A.method,A.args,A.result_filter,A.status,A.started,A.ended,A.created
       FROM service_step A
 INNER JOIN ( SELECT C.cluster_uuid,C.uuid,C.seq,C.pdate,MAX(C.timestamp) AS timestamp FROM service_step AS C
                     GROUP BY C.cluster_uuid,C.uuid,C.seq,C.pdate
            ) B ON B.cluster_uuid = A.cluster_uuid AND B.uuid = A.uuid AND B.seq = A.seq AND B.pdate = A.pdate AND B.timestamp = A.timestamp 
) X`
	stepSet := make(map[string]map[int]servicev3.ServiceStep)
	step := servicev3.ServiceStep{}
	err = vanilla.Stmt.Select(tablename, step.ColumnNames(), q, o, p).
		QueryRowsContext(ctx.Request().Context(), ctl)(func(scan vanilla.Scanner, _ int) error {

		err := step.Scan(scan)
		if err != nil {
			return errors.Wrapf(err, "scan service")
		}

		if _, ok := stepSet[step.Uuid]; !ok {
			stepSet[step.Uuid] = make(map[int]servicev3.ServiceStep)
		}

		stepSet[step.Uuid][step.Sequence] = step
		return nil
	})
	if err != nil {
		return err
	}

	// make response body
	rsp := make([]servicev3.HttpRsp_ServiceStep, 0, state.ENV__INIT_SLICE_CAPACITY__())
	for _, seqSet := range stepSet {
		for _, step := range seqSet {
			rsp = append(rsp, servicev3.HttpRsp_ServiceStep{
				ServiceStep: step,
			})
		}
	}

	return ctx.JSON(http.StatusOK, rsp)
}

// Get []ServiceStep
// @Description Get Service Steps
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid}/step [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "ServiceStep 의 uuid"
// @Success     200 {array} v3.HttpRsp_ServiceStep
func (ctl ControlVanilla) GetServiceSteps(ctx echo.Context) (err error) {
	var uuid string
	if uuid = echoutil.Param(ctx)[__UUID__]; len(uuid) == 0 {
		err = ErrorInvalidRequestParameter()
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
			))
		return HttpError(err, http.StatusBadRequest)
	}

	// get service steps
	stepSet, err := vault.Servicev3.GetServiceSteps(ctx.Request().Context(), ctl, "", uuid)
	if err != nil {
		return err
	}

	// make response body
	rsp := make([]servicev3.HttpRsp_ServiceStep, 0, state.ENV__INIT_SLICE_CAPACITY__())
	for _, seqSet := range stepSet {
		for _, step := range seqSet {
			rsp = append(rsp, servicev3.HttpRsp_ServiceStep{
				ServiceStep: step,
			})
		}
	}

	return ctx.JSON(http.StatusOK, rsp)
}

// Get ServiceStep
// @Description Get Service Step
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid}/step/{sequence} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid     path string true "ServiceStep 의 uuid"
// @Param       sequence path string true "ServiceStep 의 sequence"
// @Success     200 {object} v3.HttpRsp_ServiceStep
func (ctl ControlVanilla) GetServiceStep(ctx echo.Context) (err error) {
	var uuid string
	if uuid = echoutil.Param(ctx)[__UUID__]; len(uuid) == 0 {
		err = ErrorInvalidRequestParameter()
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

	// get service step
	step, err := vault.Servicev3.GetServiceStep(ctx.Request().Context(), ctl, "", uuid, sequence)
	if err != nil {
		return err
	}

	// make response body
	rsp := servicev3.HttpRsp_ServiceStep{ServiceStep: *step}

	return ctx.JSON(http.StatusOK, rsp)
}

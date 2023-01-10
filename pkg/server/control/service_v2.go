package control

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/sqlex"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/qri-io/jsonschema"

	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"

	"github.com/NexClipper/sudory/pkg/server/model/service/v4"
	"github.com/NexClipper/sudory/pkg/server/model/template/v3"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @Description Create a Service (v2)
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /v2/server/service [post]
// @Param       service     body service.HttpReq_Service_create true  "HttpReq_Service_create"
// @Success     200 {array} service.HttpRsp_Service_create
func (ctl ControlVanilla) CreateService_v2(ctx echo.Context) error {
	var body = new(service.HttpReq_Service_create)
	if err := echoutil.Bind(ctx, body); err != nil {
		err = errors.Wrapf(err, "bind%s",
			logs.KVL(
				"type", TypeName(body),
			))
		return HttpError(err, http.StatusBadRequest)
	}
	// if len(body.Name) == 0 {
	// 	err := ErrorInvalidRequestParameter
	// 	err = errors.Wrapf(err, "valid param%s",
	// 		logs.KVL(
	// 			ParamLog(fmt.Sprintf("%s.Name", TypeName(body)), body.Name)...,
	// 		))
	// 	return HttpError(err, http.StatusBadRequest)
	// }
	if len(body.TemplateUuid) == 0 {
		err := ErrorInvalidRequestParameter
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(fmt.Sprintf("%s.TemplateUuid", TypeName(body)), body.TemplateUuid)...,
			))
		return HttpError(err, http.StatusBadRequest)
	}

	if len(body.ClusterUuid) == 0 {
		err := ErrorInvalidRequestParameter
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(fmt.Sprintf("%s.ClusterUuid", TypeName(body)), "empty")...,
			))
		return HttpError(err, http.StatusBadRequest)
	}

	for i := range body.ClusterUuid {
		if len(body.ClusterUuid[i]) == 0 {
			err := ErrorInvalidRequestParameter
			err = errors.Wrapf(err, "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.ClusterUuid", TypeName(body)), body.ClusterUuid[i])...,
				))
			return HttpError(err, http.StatusBadRequest)
		}
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	for i := range body.ClusterUuid {
		// check cluster
		err = vault.CheckCluster(ctx.Request().Context(), ctl, ctl.dialect, claims, body.ClusterUuid[i])
		if err != nil {
			return err
		}
	}

	// get template
	tmpl_cond := stmt.And(
		stmt.Equal("uuid", body.TemplateUuid),
		stmt.IsNull("deleted"),
	)

	tmpl := template.Template{}
	err = ctl.dialect.QueryRow(tmpl.TableName(), tmpl.ColumnNames(), tmpl_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := tmpl.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		err = errors.Wrapf(err, "failed to get template")
		return err
	}

	// set service data
	now_time := time.Now()
	// new service && validation
	err = ValidServiceInput(tmpl, now_time, *body)
	if err != nil {
		err = HttpError(err, http.StatusBadRequest)
		return err
	}

	new_servs, new_status := NewService_v2(tmpl, now_time, *body)

	var servs = make([]vault.Table, 0, len(new_servs))
	var status = make([]vault.Table, 0, len(new_status))
	for i := 0; i < len(new_servs); i++ {
		servs = append(servs, new_servs[i])
	}
	for i := 0; i < len(new_status); i++ {
		status = append(status, new_status[i])
	}

	// save
	err = sqlex.ScopeTx(ctx.Request().Context(), ctl, func(tx *sql.Tx) error {
		// save service
		if err := vault.SaveMultiTable(tx, ctl.dialect, servs); err != nil {
			err = errors.Wrapf(err, "failed to save service")
			return err
		}

		// save status
		if err := vault.SaveMultiTable(tx, ctl.dialect, status); err != nil {
			err = errors.Wrapf(err, "failed to save service")
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	// make response body
	switch body.IsMultiCluster {
	case true:
		var rsp = []service.HttpRsp_Service_create{}
		for i := range new_servs {
			rsp = append(rsp, service.HttpRsp_Service_create(new_servs[i]))
		}

		return ctx.JSON(http.StatusOK, []service.HttpRsp_Service_create(rsp))
	default:
		rsp := service.HttpRsp_Service_create(new_servs[0])
		return ctx.JSON(http.StatusOK, service.HttpRsp_Service_create(rsp))
	}
}

// @Description Find []Service (v2)
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /v2/server/service [get]
// @Param       q           query string false "query  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       o           query string false "order  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       p           query string false "paging github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Success     200 {array} service.HttpRsp_Service
func (ctl ControlVanilla) FindService_v2(ctx echo.Context) (err error) {
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
	// default pagination
	if p == nil {
		p = stmt.Limit(__DEFAULT_DECORATION_LIMIT__)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	service_table := service.TableNameWithTenant_Service(claims.Hash)

	servs := make([]service.Service, 0, state.ENV__INIT_SLICE_CAPACITY__())
	serv := service.Service{}
	err = ctl.dialect.QueryRows(service_table, serv.ColumnNames(), q, o, p)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner, _ int) error {
			err := serv.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			servs = append(servs, serv)

			return err
		})
	if err != nil {
		return err
	}

	// get service records
	var records map[string][]service.ServiceStatus = make(map[string][]service.ServiceStatus)

	if false {
		for _, serv := range servs {
			status_table := service.TableNameWithTenant_ServiceStatus(claims.Hash)
			var status service.ServiceStatus
			status_cond := stmt.And(
				stmt.Equal("uuid", serv.Uuid),
			)

			err = ctl.dialect.QueryRows(status_table, status.ColumnNames(), status_cond, nil, nil)(
				ctx.Request().Context(), ctl)(
				func(scan excute.Scanner, _ int) error {
					err := status.Scan(scan)
					if err != nil {
						err = errors.WithStack(err)
						return err
					}

					if records[serv.Uuid] == nil {
						records[serv.Uuid] = make([]service.ServiceStatus, 0, state.ENV__INIT_SLICE_CAPACITY__())
					}

					records[serv.Uuid] = append(records[serv.Uuid], status)

					return err
				})
			if err != nil {
				return err
			}
		}
	}

	// make response body
	rsp := make([]service.HttpRsp_Service, len(servs))
	for i, serv := range servs {
		rsp[i].Service = serv
		uuid := serv.Uuid

		sort.Slice(records[uuid], func(i, j int) bool {
			return records[uuid][i].Created.Before(records[uuid][j].Created)
		})

		rsp[i].StatusRecords = records[uuid]
	}

	return ctx.JSON(http.StatusOK, []service.HttpRsp_Service(rsp))
}

// @Description Get a Service (v2)
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /v2/server/service/{uuid} [get]
// @Param       uuid         path string true "service's UUID"
// @Success     200 {object} service.HttpRsp_Service
func (ctl ControlVanilla) GetService_v2(ctx echo.Context) (err error) {
	var uuid string
	if uuid = echoutil.Param(ctx)[__UUID__]; len(echoutil.Param(ctx)[__UUID__]) == 0 {
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

	service_cond := stmt.Equal("uuid", uuid)
	service_table := service.TableNameWithTenant_Service(claims.Hash)
	var serv service.Service

	err = ctl.dialect.QueryRow(service_table, serv.ColumnNames(), service_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := serv.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		err = errors.Wrapf(err, "failed to get a service")
		return
	}

	// get service recodes
	var recodes = make([]service.ServiceStatus, 0, state.ENV__INIT_SLICE_CAPACITY__())
	status_table := service.TableNameWithTenant_ServiceStatus(claims.Hash)
	var service_status service.ServiceStatus
	status_cond := stmt.Equal("uuid", uuid)

	err = ctl.dialect.QueryRows(status_table, service_status.ColumnNames(), status_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner, _ int) error {
			err = service_status.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			recodes = append(recodes, service_status)

			return err
		})
	if err != nil {
		return
	}

	sort.Slice(recodes, func(i, j int) bool {
		return recodes[i].Created.Before(recodes[j].Created)
	})

	// make response body
	rst := new(service.HttpRsp_Service)
	rst.Service = serv
	rst.StatusRecords = recodes

	return ctx.JSON(http.StatusOK, (*service.HttpRsp_Service)(rst))
}

// @Description Get a Service Result (v2)
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /v2/server/service/{uuid}/result [get]
// @Param       uuid         path string true "service's UUID"
// @Success     200 {object} service.HttpRsp_ServiceResult
func (ctl ControlVanilla) GetServiceResult_v2(ctx echo.Context) (err error) {
	var uuid string
	if uuid = echoutil.Param(ctx)[__UUID__]; len(echoutil.Param(ctx)[__UUID__]) == 0 {
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

	// get service result
	result_table := service.TableNameWithTenant_ServiceResult(claims.Hash)
	result := service.ServiceResult{}
	result_cond := stmt.Equal("uuid", uuid)

	err = ctl.dialect.QueryRows(result_table, result.ColumnNames(), result_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner, i int) error {
			err := result.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		err = errors.Wrapf(err, "failed to get a service result")
		return err
	}

	return ctx.JSON(http.StatusOK, service.HttpRsp_ServiceResult(result))
}

func ValidServiceInput(
	tmpl template.Template,
	now_time time.Time, body service.HttpReq_Service_create,

) (err error) {
	body_inputs := body.Inputs
	tmpl_inputs := tmpl.Inputs

	if body_inputs == nil {
		err = errors.New("step.Args must have value")
		return
	}

	json_schema_validator := &jsonschema.Schema{}
	// // var command_args_json []byte
	// command_args_json, err := json.Marshal(tmpl_inputs)
	// if err != nil {
	// 	err = errors.Wrapf(err, "command.Args convert to json")
	// 	return
	// }

	if err := json.Unmarshal([]byte(tmpl_inputs), json_schema_validator); err != nil {
		err = errors.Wrapf(err, "command.Args convert to json schema validator")
		return err
	}

	var step_args_json []byte
	if step_args_json, err = json.Marshal(body_inputs); err != nil {
		err = errors.Wrapf(err, "step.Args convert to json")
		return
	}

	timeout, cancel := context.WithTimeout(context.Background(), 333*time.Millisecond)
	defer cancel()

	verr, err := json_schema_validator.ValidateBytes(timeout, step_args_json)
	if err != nil {
		err = errors.Wrapf(err, "json schema validatebytes%s", logs.KVL(
			"step.args", string(step_args_json),
		))
		return
	}
	iter_verr := func() (err error) {
		for _, iter := range verr {
			if err == nil {
				err = iter
				continue
			}
			err = errors.Wrap(err, iter.Error())
		}
		return
	}
	if err = iter_verr(); err != nil {
		return
	}

	return nil
}

func NewService_v2(
	tmpl template.Template,
	now_time time.Time, body service.HttpReq_Service_create,

) ([]service.Service, []service.ServiceStatus) {

	getStringOrDefault := func(a, b string) string {
		if 0 < len(a) {
			return a
		}
		return b
	}

	getPriority := func(tmpl template.Template) service.Priority {
		if tmpl.Origin == template.OriginSystem.String() {
			return service.PriorityHigh // system
		}
		return service.PriorityLow
	}

	BuildService := func(body service.HttpReq_Service_create, cluster_uuid string) (new_service service.Service, new_status service.ServiceStatus) {

		// if uuid is empty then generate uuid
		uuid := genUuidString("")

		// compute flow
		var flow = []interface{}{}
		json.Unmarshal([]byte(tmpl.Flow), &flow)

		// property service
		new_service.PartitionDate = now_time
		new_service.ClusterUuid = cluster_uuid
		new_service.Uuid = uuid
		new_service.Name = getStringOrDefault(body.Name, tmpl.Name)
		new_service.Summary = *vanilla.NewNullString(getStringOrDefault(body.Summary, tmpl.Summary.String))
		new_service.TemplateUuid = body.TemplateUuid
		new_service.Flow = tmpl.Flow
		new_service.Inputs = body.Inputs
		new_service.StepMax = len(flow)
		new_service.SubscribedChannel = *vanilla.NewNullString(body.SubscribedChannel)
		new_service.Priority = getPriority(tmpl)
		new_service.Created = now_time

		// property status
		new_status.PartitionDate = new_service.PartitionDate
		new_status.Created = new_service.Created
		new_status.ClusterUuid = new_service.ClusterUuid
		new_status.Uuid = new_service.Uuid
		new_status.StepMax = new_service.StepMax

		return
	}

	var services = []service.Service{}
	var status = []service.ServiceStatus{}

	for i := range body.ClusterUuid {

		new_service, new_status := BuildService(body, body.ClusterUuid[i])

		services = append(services, new_service)
		status = append(status, new_status)
	}

	return services, status
}

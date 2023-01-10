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
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/sqlex"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/status/state"

	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/model/service/v3"
	templatev2 "github.com/NexClipper/sudory/pkg/server/model/template/v2"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/qri-io/jsonschema"
)

// @Description Create a Service
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service [post]
// @Param       service     body service.HttpReq_Service_create true "HttpReq_Service_create"
// @Success     200 {array} service.HttpRsp_Service_create
func (ctl ControlVanilla) CreateService(ctx echo.Context) error {
	var body = new(service.HttpReq_Service_create)
	if err := echoutil.Bind(ctx, body); err != nil {
		err = errors.Wrapf(err, "bind%s",
			logs.KVL(
				"type", TypeName(body),
			))
		return HttpError(err, http.StatusBadRequest)
	}
	if len(body.Name) == 0 {
		err := ErrorInvalidRequestParameter
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(fmt.Sprintf("%s.Name", TypeName(body)), body.Name)...,
			))
		return HttpError(err, http.StatusBadRequest)
	}
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

	// get template && commands
	template, commands, err := vault.GetTemplate(ctx.Request().Context(), ctl, ctl.dialect, body.TemplateUuid)
	if err != nil {
		return err
	}

	// check steps length
	if len(body.Steps) != len(commands) {
		err = errors.Errorf("diff length of steps and commands%s",
			logs.KVL(
				"expected", len(commands),
				"actual", len(body.Steps),
			))
		return HttpError(err, http.StatusBadRequest)
	}

	// set service data
	now_time := time.Now()
	// new service && validation
	new_services, new_steps, err := newCreateServiceWithValid(*template, commands, now_time, *body)
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// save
	err = sqlex.ScopeTx(ctx.Request().Context(), ctl, func(tx *sql.Tx) error {
		for i := range new_services {
			// save service
			if err := vault.SaveMultiTable(tx, ctl.dialect, []vault.Table{new_services[i]}); err != nil {
				return errors.Wrapf(err, "failed to save service")
			}

			// save service steps_
			steps_ := func(steps []service.ServiceStep_create) []vault.Table {
				tables := make([]vault.Table, len(steps))
				for i := range steps {
					tables[i] = steps[i]
				}
				return tables
			}(new_steps[new_services[i].Uuid])

			if err := vault.SaveMultiTable(tx, ctl.dialect, steps_); err != nil {
				return errors.Wrapf(err, "failed to save service steps")
			}
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
		for i := range new_services {
			rsp = append(rsp, service.HttpRsp_Service_create{
				Service_create: new_services[i],
				Steps:          new_steps[new_services[i].Uuid],
			})
		}

		return ctx.JSON(http.StatusOK, []service.HttpRsp_Service_create(rsp))
	default:
		rsp := service.HttpRsp_Service_create{
			Service_create: new_services[0],
			Steps:          new_steps[new_services[0].Uuid],
		}
		return ctx.JSON(http.StatusOK, service.HttpRsp_Service_create(rsp))
	}
}

// @Description Find []Service
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service [get]
// @Param       q           query string false "query  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       o           query string false "order  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       p           query string false "paging github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Success     200 {array} service.HttpRsp_Service
func (ctl ControlVanilla) FindService(ctx echo.Context) (err error) {
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

	serviceSet := make(map[string]service.Service)
	serv := service.Service{}
	err = ctl.dialect.QueryRows(service_table, serv.ColumnNames(), q, o, p)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner, _ int) error {
			err := serv.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			serviceSet[serv.Uuid] = serv

			return err
		})
	if err != nil {
		return err
	}

	// get service stepSet
	var stepSet map[string][]service.ServiceStep = make(map[string][]service.ServiceStep)

	for _, serv := range serviceSet {
		step_table := service.TableNameWithTenant_ServiceStep(claims.Hash)
		var step service.ServiceStep
		step_cond := stmt.And(
			stmt.Equal("uuid", serv.Uuid),
		)

		err = ctl.dialect.QueryRows(step_table, step.ColumnNames(), step_cond, nil, nil)(
			ctx.Request().Context(), ctl)(
			func(scan excute.Scanner, _ int) error {
				err := step.Scan(scan)
				if err != nil {
					err = errors.WithStack(err)
					return err
				}

				if stepSet[serv.Uuid] == nil {
					stepSet[serv.Uuid] = make([]service.ServiceStep, 0, state.ENV__INIT_SLICE_CAPACITY__())
				}

				stepSet[serv.Uuid] = append(stepSet[serv.Uuid], step)

				return err
			})
		if err != nil {
			return err
		}
	}
	// make response body
	rsp := make([]service.HttpRsp_Service, len(serviceSet))
	var i int
	for uuid, service := range serviceSet {
		rsp[i].Service = service

		sort.Slice(stepSet[uuid], func(i, j int) bool {
			return stepSet[uuid][i].Sequence < stepSet[uuid][j].Sequence
		})

		rsp[i].Steps = stepSet[uuid]
		i++
	}

	return ctx.JSON(http.StatusOK, []service.HttpRsp_Service(rsp))
}

// @Description Get a Service
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid} [get]
// @Param       uuid         path string true "service's UUID"
// @Success     200 {object} service.HttpRsp_Service
func (ctl ControlVanilla) GetService(ctx echo.Context) (err error) {
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
		err = errors.Wrapf(err, "failed to get service")
		return
	}

	// get service steps
	var steps = make([]service.ServiceStep, 0, state.ENV__INIT_SLICE_CAPACITY__())
	step_table := service.TableNameWithTenant_ServiceStep(claims.Hash)
	var service_step service.ServiceStep
	step_cond := stmt.Equal("uuid", uuid)

	err = ctl.dialect.QueryRows(step_table, service_step.ColumnNames(), step_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner, _ int) error {
			err = service_step.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			steps = append(steps, service_step)

			return err
		})
	if err != nil {
		err = errors.Wrapf(err, "failed to found service step")
		return
	}

	sort.Slice(steps, func(i, j int) bool {
		return steps[i].Sequence < steps[j].Sequence
	})

	// make response body
	rst := new(service.HttpRsp_Service)
	rst.Service = serv
	rst.Steps = steps

	return ctx.JSON(http.StatusOK, (*service.HttpRsp_Service)(rst))
}

// @Description Get a Service Result
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid}/result [get]
// @Param       uuid         path string true "service's UUID"
// @Success     200 {object} service.HttpRsp_ServiceResult
func (ctl ControlVanilla) GetServiceResult(ctx echo.Context) (err error) {
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

	err = ctl.dialect.QueryRow(result_table, result.ColumnNames(), result_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
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

const __DEFAULT_DECORATION_LIMIT__ = 20

// func ParseDecoration(m map[string]string) (q *prepare.Condition, o *prepare.Orders, p *prepare.Pagination, err error) {
// 	q, o, p, err = prepare.NewParser(m)
// 	if p == nil {
// 		p = vanilla.Limit(__DEFAULT_DECORATION_LIMIT__).Parse()
// 	}

// 	return
// }

func newCreateServiceWithValid(
	template templatev2.Template, commands []templatev2.TemplateCommand,
	now_time time.Time, body service.HttpReq_Service_create,

) ([]service.Service_create, map[string][]service.ServiceStep_create, error) {

	// build service step
	for i := range body.Steps {
		step_args := body.Steps[i].Args
		command_args := commands[i].Args

		if step_args == nil {
			err := errors.New("step.Args must have value")
			return nil, nil, err
		}

		json_schema_validator := &jsonschema.Schema{}
		// var command_args_json []byte
		command_args_json, err := json.Marshal(command_args)
		if err != nil {
			err = errors.Wrapf(err, "command.Args convert to json")
			return nil, nil, err
		}

		if err := json.Unmarshal([]byte(command_args_json), json_schema_validator); err != nil {
			err = errors.Wrapf(err, "command.Args convert to json schema validator")
			return nil, nil, err
		}
		var step_args_json []byte
		if step_args_json, err = json.Marshal(step_args); err != nil {
			err = errors.Wrapf(err, "step.Args convert to json")
			return nil, nil, err
		}

		timeout, cancel := context.WithTimeout(context.Background(), 333*time.Millisecond)
		defer cancel()

		verr, err := json_schema_validator.ValidateBytes(timeout, step_args_json)
		if err != nil {
			err = errors.Wrapf(err, "json schema validatebytes%s", logs.KVL(
				"step.args", string(step_args_json),
			))
			return nil, nil, err
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
		if err := iter_verr(); err != nil {
			return nil, nil, err
		}
	}

	getPriority := func(template templatev2.Template) service.Priority {
		if template.Origin == templatev2.OriginSystem.String() {
			return service.PriorityHigh // system
		}
		return service.PriorityLow
	}

	BuildService := func(body service.HttpReq_Service_create, cluster_uuid string) (new_service service.Service_create, new_steps []service.ServiceStep_create) {

		// if uuid is empty then generate uuid
		uuid := genUuidString("")

		// property service
		// new_service := service.Service_create{}
		new_service.PartitionDate = now_time
		new_service.ClusterUuid = cluster_uuid
		new_service.Uuid = uuid
		new_service.Timestamp = now_time
		new_service.Name = body.Name
		new_service.Summary = *vanilla.NewNullString(body.Summary)
		new_service.TemplateUuid = body.TemplateUuid
		new_service.StepCount = len(body.Steps)
		new_service.SubscribedChannel = *vanilla.NewNullString(body.SubscribedChannel)
		new_service.StepPosition = 0
		new_service.Status = service.StepStatusRegist
		new_service.Priority = getPriority(template)
		new_service.Created = now_time

		// create steps
		new_steps = make([]service.ServiceStep_create, 0, len(body.Steps))
		for i := range body.Steps {
			command := commands[i]
			body_step := body.Steps[i]

			//property step
			new_step := service.ServiceStep_create{}
			new_step.PartitionDate = now_time
			new_step.ClusterUuid = cluster_uuid
			new_step.Uuid = uuid
			new_step.Sequence = i
			new_step.Timestamp = now_time
			new_step.Name = command.Name
			new_step.Summary = command.Summary
			new_step.Method = command.Method
			new_step.Args = body_step.Args
			new_step.ResultFilter = command.ResultFilter
			new_step.Status = service.StepStatusRegist
			new_step.Created = now_time

			// append service step
			new_steps = append(new_steps, new_step)
		}

		return
	}

	var services = []service.Service_create{}
	var service_steps = map[string][]service.ServiceStep_create{}

	for i := range body.ClusterUuid {

		new_service, new_steps := BuildService(body, body.ClusterUuid[i])

		services = append(services, new_service)
		service_steps[new_service.Uuid] = new_steps
	}

	return services, service_steps, nil
}

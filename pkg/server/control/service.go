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
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmtex"
	"github.com/NexClipper/sudory/pkg/server/status/state"

	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	clusterv3 "github.com/NexClipper/sudory/pkg/server/model/cluster/v3"
	servicev3 "github.com/NexClipper/sudory/pkg/server/model/service/v3"
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
// @Param       service      body   v3.HttpReq_Service_create true  "HttpReq_Service_create"
// @Success     200 {object} v3.HttpRsp_Service_create
func (ctl ControlVanilla) CreateService(ctx echo.Context) error {
	var body = new(servicev3.HttpReq_Service_create)
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
				ParamLog(fmt.Sprintf("%s.ClusterUuid", TypeName(body)), body.ClusterUuid)...,
			))
		return HttpError(err, http.StatusBadRequest)
	}

	// check cluster
	// cluster_cond := vanilla.And(
	// 	vanilla.Equal("uuid", body.ClusterUuid),
	// 	vanilla.IsNull("deleted"),
	// )

	// cluster := clusterv2.Cluster{}
	// cluster_found, err := vanilla.Stmt.Exist(cluster.TableName(), cluster_cond.Parse())(ctx.Request().Context(), ctl)
	// if err != nil {
	// 	return errors.Wrapf(err, "faild to check cluster")
	// }
	// if !cluster_found {
	// 	return errors.Wrapf(err, "cluster does not exist ")
	// }

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}
	err = func() error {
		// check cluster
		cluster_table := clusterv3.TableNameWithTenant(claims.Hash)
		cluster_cond := stmt.And(
			stmt.Equal("uuid", body.ClusterUuid),
			stmt.IsNull("deleted"),
		)
		cluster_exist, err := stmtex.ExistContext(cluster_table, cluster_cond)(ctx.Request().Context(), ctl, ctl.Dialect())
		if err != nil {
			return errors.Wrapf(err, "check cluster")
		}
		if !cluster_exist {
			return errors.Wrapf(database.ErrorRecordWasNotFound, "check cluster")
		}
		return nil
	}()
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// get template
	template_cond := stmt.And(
		stmt.Equal("uuid", body.TemplateUuid),
		stmt.IsNull("deleted"),
	)

	template := templatev2.Template{}
	err = stmtex.Select(template.TableName(), template.ColumnNames(), template_cond, nil, nil).
		QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(
		func(scan stmtex.Scanner) error {
			return template.Scan(scan)
		})
	if err != nil {
		return errors.Wrapf(err, "faild to get template")
	}

	// get commands
	commandSet := make(map[int]templatev2.TemplateCommand)
	command_cond := stmt.And(
		stmt.Equal("template_uuid", body.TemplateUuid),
		stmt.IsNull("deleted"),
	)

	command := templatev2.TemplateCommand{}
	err = stmtex.Select(command.TableName(), command.ColumnNames(), command_cond, nil, nil).
		QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(
		func(scan stmtex.Scanner) error {
			err := command.Scan(scan)

			if err != nil {
				return err
			}
			commandSet[command.Sequence] = command

			return nil
		})
	if err != nil {
		return errors.Wrapf(err, "failed to get template commands")
	}

	if len(body.Steps) != len(commandSet) {
		err = errors.Errorf("diff length of steps and commands%s",
			logs.KVL(
				"expected", len(commandSet),
				"actual", len(body.Steps),
			))
		return HttpError(err, http.StatusBadRequest)
	}

	// convert commnad set to slice
	commands := make([]templatev2.TemplateCommand, 0, len(commandSet))
	for _, command := range commandSet {
		commands = append(commands, command)
	}

	// build service step
	for i := range body.Steps {
		step_args := body.Steps[i].Args
		command_args := commands[i].Args

		if step_args == nil {
			err = errors.New("step.Args must have value")
			return HttpError(err, http.StatusBadRequest) // bad request
		}

		json_schema_validator := &jsonschema.Schema{}
		var command_args_json []byte
		if command_args_json, err = json.Marshal(command_args); err != nil {
			err = errors.Wrapf(err, "command.Args convert to json")
			return HttpError(err, http.StatusBadRequest)
		}

		if err := json.Unmarshal([]byte(command_args_json), json_schema_validator); err != nil {
			err = errors.Wrapf(err, "command.Args convert to json schema validator")
			return HttpError(err, http.StatusBadRequest)
		}
		var step_args_json []byte
		if step_args_json, err = json.Marshal(step_args); err != nil {
			err = errors.Wrapf(err, "step.Args convert to json")
			return HttpError(err, http.StatusBadRequest)
		}

		timeout, cancel := context.WithTimeout(context.Background(), 333*time.Millisecond)
		defer cancel()

		verr, err := json_schema_validator.ValidateBytes(timeout, step_args_json)
		if err != nil {
			err = errors.Wrapf(err, "json schema validatebytes%s", logs.KVL(
				"step.args", string(step_args_json),
			))
			return HttpError(err, http.StatusBadRequest)
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
			return HttpError(err, http.StatusBadRequest)
		}
	}

	// set service data
	now_time := time.Now()
	// gen uuid
	body.Uuid = genUuidString(body.Uuid)

	getPriority := func(template templatev2.Template) servicev3.Priority {
		if template.Origin == templatev2.OriginSystem.String() {
			return servicev3.PriorityHigh // system
		}
		return servicev3.PriorityLow
	}

	// property service
	service := servicev3.Service_create{}
	service.PartitionDate = now_time
	service.ClusterUuid = body.ClusterUuid
	service.Uuid = body.Uuid
	service.Timestamp = now_time
	service.Name = body.Name
	service.Summary = *vanilla.NewNullString(body.Summary)
	service.TemplateUuid = body.TemplateUuid
	service.StepCount = len(body.Steps)
	service.SubscribedChannel = *vanilla.NewNullString(body.SubscribedChannel)
	service.StepPosition = 0
	service.Status = servicev3.StepStatusRegist
	service.Priority = getPriority(template)
	service.Created = now_time

	// create steps
	steps := make([]servicev3.ServiceStep_create, 0, len(body.Steps))
	for i := range body.Steps {
		command := commands[i]
		body_step := body.Steps[i]

		//property step
		step := servicev3.ServiceStep_create{}
		step.PartitionDate = now_time
		step.ClusterUuid = body.ClusterUuid
		step.Uuid = body.Uuid
		step.Sequence = i
		step.Timestamp = now_time
		step.Name = command.Name
		step.Summary = command.Summary
		step.Method = command.Method
		step.Args = body_step.Args
		step.ResultFilter = command.ResultFilter
		step.Status = servicev3.StepStatusRegist
		step.Created = now_time

		// append service step
		steps = append(steps, step)
	}

	// save
	err = stmtex.ScopeTx(ctx.Request().Context(), ctl, func(tx *sql.Tx) error {

		// save service
		if err := vault.SaveMultiTable(tx, ctl.Dialect(), []vault.Table{service}); err != nil {
			return errors.Wrapf(err, "faild to save service")
		}

		// save service steps_
		steps_ := func(steps []servicev3.ServiceStep_create) []vault.Table {
			tables := make([]vault.Table, len(steps))
			for i := range steps {
				tables[i] = steps[i]
			}
			return tables
		}(steps)

		if err := vault.SaveMultiTable(tx, ctl.Dialect(), steps_); err != nil {
			return errors.Wrapf(err, "faild to save service steps")
		}

		return nil
	})
	if err != nil {
		return err
	}

	// make response body
	rsp := servicev3.HttpRsp_Service_create{}
	rsp.Service_create = service
	rsp.Steps = steps

	return ctx.JSON(http.StatusOK, servicev3.HttpRsp_Service_create(rsp))
}

// @Description Find []Service
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service [get]
// @Param       q            query  string false "query  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       o            query  string false "order  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       p            query  string false "paging github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Success     200 {array} v3.HttpRsp_Service
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

	service_table := servicev3.TableNameWithTenant_Service(claims.Hash)

	serviceSet := make(map[string]servicev3.Service)
	service := servicev3.Service{}
	err = stmtex.Select(service_table, service.ColumnNames(), q, o, p).
		QueryRowsContext(ctx.Request().Context(), ctl, ctl.Dialect())(
		func(scan stmtex.Scanner, _ int) error {

			err := service.Scan(scan)
			if err != nil {
				return errors.Wrapf(err, "scan service")
			}

			serviceSet[service.Uuid] = service
			return nil
		})
	if err != nil {
		return err
	}

	// get service stepSet
	var stepSet map[string][]servicev3.ServiceStep = make(map[string][]servicev3.ServiceStep)

	for _, service := range serviceSet {
		step_table := servicev3.TableNameWithTenant_ServiceStep(claims.Hash)
		var step servicev3.ServiceStep
		step_cond := stmt.And(
			stmt.Equal("uuid", service.Uuid),
		)

		err = stmtex.Select(step_table, step.ColumnNames(), step_cond, nil, nil).
			QueryRowsContext(ctx.Request().Context(), ctl, ctl.Dialect())(
			func(scan stmtex.Scanner, _ int) error {
				err := step.Scan(scan)
				if err != nil {
					return errors.Wrapf(err, "scan service step")
				}
				if stepSet[service.Uuid] == nil {
					stepSet[service.Uuid] = make([]servicev3.ServiceStep, 0, state.ENV__INIT_SLICE_CAPACITY__())
				}

				stepSet[service.Uuid] = append(stepSet[service.Uuid], step)
				return nil
			})
		if err != nil {
			return err
		}
	}
	// make response body
	rsp := make([]servicev3.HttpRsp_Service, len(serviceSet))
	var i int
	for uuid, service := range serviceSet {
		rsp[i].Service = service

		sort.Slice(stepSet[uuid], func(i, j int) bool {
			return stepSet[uuid][i].Sequence < stepSet[uuid][j].Sequence
		})

		rsp[i].Steps = stepSet[uuid]
		i++
	}

	return ctx.JSON(http.StatusOK, []servicev3.HttpRsp_Service(rsp))
}

// @Description Get a Service
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid} [get]
// @Param       uuid         path   string true  "Service 의 Uuid"
// @Success     200 {object} v3.HttpRsp_Service
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
	service_table := servicev3.TableNameWithTenant_Service(claims.Hash)
	var service servicev3.Service

	err = stmtex.Select(service_table, service.ColumnNames(), service_cond, nil, nil).
		QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(
		func(scan stmtex.Scanner) error {
			return service.Scan(scan)
		})
	if err != nil {
		err = errors.Wrapf(err, "failed to get service")
		return
	}

	// get service steps
	var steps = make([]servicev3.ServiceStep, 0, state.ENV__INIT_SLICE_CAPACITY__())
	step_table := servicev3.TableNameWithTenant_ServiceStep(claims.Hash)
	var service_step servicev3.ServiceStep
	step_cond := stmt.Equal("uuid", uuid)

	err = stmtex.Select(step_table, service_step.ColumnNames(), step_cond, nil, nil).
		QueryRowsContext(ctx.Request().Context(), ctl, ctl.Dialect())(
		func(scan stmtex.Scanner, _ int) error {

			err = service_step.Scan(scan)
			if err != nil {
				return errors.Wrapf(err, "scan service step")
			}

			steps = append(steps, service_step)
			return nil
		})
	if err != nil {
		err = errors.Wrapf(err, "failed to found service step")
		return
	}

	sort.Slice(steps, func(i, j int) bool {
		return steps[i].Sequence < steps[j].Sequence
	})

	// make response body
	rst := new(servicev3.HttpRsp_Service)
	rst.Service = service
	rst.Steps = steps

	return ctx.JSON(http.StatusOK, (*servicev3.HttpRsp_Service)(rst))
}

// @Description Get a Service Result
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid}/result [get]
// @Param       uuid         path   string true  "Service 의 Uuid"
// @Success     200 {object} v3.HttpRsp_ServiceResult
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
	result_table := servicev3.TableNameWithTenant_ServiceResult(claims.Hash)
	result := servicev3.ServiceResult{}
	result_cond := stmt.Equal("uuid", uuid)

	err = stmtex.Select(result_table, result.ColumnNames(), result_cond, nil, nil).
		QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(
		func(scan stmtex.Scanner) error {
			return result.Scan(scan)
		})
	if err != nil {
		err = errors.Wrapf(err, "failed to get a service result")
		return err
	}

	return ctx.JSON(http.StatusOK, servicev3.HttpRsp_ServiceResult(result))
}

const __DEFAULT_DECORATION_LIMIT__ = 20

// func ParseDecoration(m map[string]string) (q *prepare.Condition, o *prepare.Orders, p *prepare.Pagination, err error) {
// 	q, o, p, err = prepare.NewParser(m)
// 	if p == nil {
// 		p = vanilla.Limit(__DEFAULT_DECORATION_LIMIT__).Parse()
// 	}

// 	return
// }

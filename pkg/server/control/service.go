package control

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro"

	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	clusterv2 "github.com/NexClipper/sudory/pkg/server/model/cluster/v2"
	servicev3 "github.com/NexClipper/sudory/pkg/server/model/service/v3"
	v3 "github.com/NexClipper/sudory/pkg/server/model/service/v3"
	templatev2 "github.com/NexClipper/sudory/pkg/server/model/template/v2"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/qri-io/jsonschema"
)

// Create Service
// @Description Create a Service
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service [post]
// @Param       x_auth_token header string                    false "client session token"
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
		err := ErrorInvalidRequestParameter()
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(fmt.Sprintf("%s.Name", TypeName(body)), body.Name)...,
			))
		return HttpError(err, http.StatusBadRequest)
	}
	if len(body.TemplateUuid) == 0 {
		err := ErrorInvalidRequestParameter()
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(fmt.Sprintf("%s.TemplateUuid", TypeName(body)), body.TemplateUuid)...,
			))
		return HttpError(err, http.StatusBadRequest)
	}
	if len(body.ClusterUuid) == 0 {
		err := ErrorInvalidRequestParameter()
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(fmt.Sprintf("%s.ClusterUuid", TypeName(body)), body.ClusterUuid)...,
			))
		return HttpError(err, http.StatusBadRequest)
	}

	// check cluster
	cluster_cond := vanilla.And(
		vanilla.Equal("uuid", body.ClusterUuid),
		vanilla.IsNull("deleted"),
	)

	cluster := clusterv2.Cluster{}
	cluster_found, err := vanilla.Stmt.Exist(cluster.TableName(), cluster_cond.Parse())(ctx.Request().Context(), ctl)
	if err != nil {
		return errors.Wrapf(err, "faild to check cluster")
	}
	if !cluster_found {
		return errors.Wrapf(err, "cluster does not exist ")
	}

	// get template
	template_cond := vanilla.And(
		vanilla.Equal("uuid", body.TemplateUuid),
		vanilla.IsNull("deleted"),
	)

	template := templatev2.Template{}
	template_scan := func(scan vanilla.Scanner) (err error) {
		return template.Scan(scan)
	}
	stmt := vanilla.Stmt.Select(template.TableName(), template.ColumnNames(), template_cond.Parse(), nil, nil)
	if err = stmt.QueryRow(ctl)(template_scan); err != nil {
		return errors.Wrapf(err, "faild to get template")
	}

	// get commands
	commandSet := make(map[int]templatev2.TemplateCommand)
	command_cond := vanilla.And(
		vanilla.Equal("template_uuid", body.TemplateUuid),
		vanilla.IsNull("deleted"),
	)

	command := templatev2.TemplateCommand{}
	command_scan := func(scan vanilla.Scanner, _ int) (err error) {
		err = command.Scan(scan)
		if err != nil {
			return err
		}

		commandSet[command.Sequence] = command
		return
	}
	stmt = vanilla.Stmt.Select(command.TableName(), command.ColumnNames(), command_cond.Parse(), nil, nil)
	if err := stmt.QueryRows(ctl)(command_scan); err != nil {
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
		if iter_verr(); err != nil {
			return HttpError(err, http.StatusBadRequest)
		}
	}

	// set service data
	now_time := time.Now()
	uuid := body.Uuid
	if len(uuid) == 0 {
		uuid = macro.NewUuidString()
	}

	// property service
	service := servicev3.Service_create{}
	service.PartitionDate = now_time
	service.ClusterUuid = body.ClusterUuid
	service.Uuid = uuid
	service.Timestamp = now_time
	service.Name = body.Name
	service.Summary = *vanilla.NewNullString(body.Summary)
	service.TemplateUuid = body.TemplateUuid
	service.StepCount = len(body.Steps)
	service.SubscribedChannel = *vanilla.NewNullString(body.SubscribedChannel)
	service.StepPosition = 0
	service.Status = servicev3.StepStatusRegist
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
		step.Uuid = uuid
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
	err = ctl.ScopeTx(ctx.Request().Context(), func(tx *sql.Tx) error {

		// save service
		if err := vault.SaveMultiTable(tx, []vault.Table{service}); err != nil {
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

		if err := vault.SaveMultiTable(tx, steps_); err != nil {
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

	return ctx.JSON(http.StatusOK, rsp)
}

// Find []Service
// @Description Find []Service
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v3.HttpRsp_Service
func (ctl ControlVanilla) FindService(ctx echo.Context) (err error) {
	q, o, p, err := ParseDecoration(echoutil.QueryParam(ctx))
	if err != nil {
		err = errors.Wrapf(err, "ParseDecoration%v", logs.KVL(
			"query", echoutil.QueryParamString(ctx),
		))
		return HttpError(err, http.StatusBadRequest)
	}

	// find service
	tablename := `(
SELECT A.pdate,A.cluster_uuid,A.uuid,A.created,A.name,A.summary,A.template_uuid,A.step_count,A.subscribed_channel,A.assigned_client_uuid,A.step_position,A.status,A.message,A.service_created
       FROM service A
 INNER JOIN service B
         ON B.pdate = A.pdate AND B.cluster_uuid = A.cluster_uuid AND B.uuid = A.uuid
        AND B.created = ( SELECT MAX(C.created)
                            FROM service C
                           WHERE C.pdate = B.pdate AND C.cluster_uuid = B.cluster_uuid AND B.uuid = C.uuid )
) X`
	serviceSet := make(map[string]servicev3.Service)
	service := servicev3.Service{}
	err = vanilla.Stmt.Select(tablename, service.ColumnNames(), q, o, p).
		QueryRowsContext(ctx.Request().Context(), ctl)(func(scan vanilla.Scanner, _ int) error {

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

	// get service steps
	var steps map[string]map[int]servicev3.ServiceStep = make(map[string]map[int]servicev3.ServiceStep)
	for _, service := range serviceSet {

		stepSet, err := vault.Servicev3.GetServiceSteps(ctx.Request().Context(), ctl, service.ClusterUuid, service.Uuid)
		if err != nil {
			return errors.Wrapf(err, "failed to found service steps%v", logs.KVL(
				"cluster_uuid", service.ClusterUuid,
				"uuid", service.Uuid,
			))
		}

		for uuid, seqSet := range stepSet {
			// init sub record set
			if steps[uuid] == nil {
				steps[uuid] = make(map[int]v3.ServiceStep)
			}
			// move to buffer set
			for seq, step := range seqSet {
				steps[uuid][seq] = step
			}
		}
	}

	// make response body
	rsp := make([]servicev3.HttpRsp_Service, len(serviceSet))
	var i int
	for uuid, service := range serviceSet {
		steps_ := make([]servicev3.ServiceStep, 0, len(steps[uuid]))
		for _, step := range steps[uuid] {
			steps_ = append(steps_, step)
		}

		rsp[i].Service = service
		rsp[i].Steps = steps_
		i++
	}

	return ctx.JSON(http.StatusOK, rsp)
}

// Get Service
// @Description Get a Service
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Service 의 Uuid"
// @Success     200 {object} v3.HttpRsp_Service
func (ctl ControlVanilla) GetService(ctx echo.Context) (err error) {
	var uuid string
	if uuid = echoutil.Param(ctx)[__UUID__]; len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter()
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
			))
		return HttpError(err, http.StatusBadRequest)
	}

	// get service
	service, err := vault.Servicev3.GetService(ctx.Request().Context(), ctl, "", uuid)
	if err != nil {
		return err
	}

	// get service steps
	stepSet, err := vault.Servicev3.GetServiceSteps(ctx.Request().Context(), ctl, "", uuid)
	if err != nil {
		return err
	}

	// make response body
	rst := new(servicev3.HttpRsp_Service)
	rst.Service = *service
	for _, seqSet := range stepSet {
		rst.Steps = make([]servicev3.ServiceStep, 0, len(stepSet))
		for _, step := range seqSet {
			rst.Steps = append(rst.Steps, step)
		}
	}

	return ctx.JSON(http.StatusOK, rst)
}

// @Description Get a Service Result
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid}/result [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Service 의 Uuid"
// @Success     200 {object} v3.HttpRsp_ServiceResult
func (ctl ControlVanilla) GetServiceResult(ctx echo.Context) (err error) {
	var uuid string
	if uuid = echoutil.Param(ctx)[__UUID__]; len(echoutil.Param(ctx)[__UUID__]) == 0 {
		err = ErrorInvalidRequestParameter()
		err = errors.Wrapf(err, "valid param%s",
			logs.KVL(
				ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
			))
		return HttpError(err, http.StatusBadRequest)
	}

	// get service result
	cond_service := vanilla.Equal("uuid", uuid)

	service_result, err := vault.Servicev3.GetServiceResult(ctx.Request().Context(), ctl, cond_service.Parse(), nil, nil)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, service_result)
}

const __DEFAULT_DECORATION_LIMIT__ = 20

func ParseDecoration(m map[string]string) (q *prepare.Condition, o *prepare.Orders, p *prepare.Pagination, err error) {
	q, o, p, err = prepare.NewParser(m)
	if p == nil {
		p = vanilla.Limit(__DEFAULT_DECORATION_LIMIT__).Parse()
	}

	return
}

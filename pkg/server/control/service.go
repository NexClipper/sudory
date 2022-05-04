package control

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
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
// @Param       x_auth_token header string                   false "client session token"
// @Param       service      body   v1.HttpReqService_Create true  "HttpReqService_Create"
// @Success     200 {object} v1.HttpRspService
func (ctl Control) CreateService(ctx echo.Context) error {
	map_step_create := func(elems []stepv1.HttpReqServiceStep_Create_ByService, mapper func(int, stepv1.HttpReqServiceStep_Create_ByService) stepv1.ServiceStep) []stepv1.ServiceStep {
		rst := make([]stepv1.ServiceStep, len(elems))
		for i := range elems {
			rst[i] = mapper(i, elems[i])
		}
		return rst
	}

	body := new(servicev1.HttpReqService_Create)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}
	if len(body.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.Name", TypeName(body)), body.Name)...,
				)))
	}
	if len(body.TemplateUuid) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.TemplateUuid", TypeName(body)), body.TemplateUuid)...,
				)))
	}
	if len(body.ClusterUuid) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.ClusterUuid", TypeName(body)), body.ClusterUuid)...,
				)))
	}
	//valid cluster
	if _, err := vault.NewCluster(ctl.NewSession()).Get(body.ClusterUuid); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(errors.Wrapf(err, "found cluster"))
	}
	//valid template
	if _, err := vault.NewTemplate(ctl.NewSession()).Get(body.TemplateUuid); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(errors.Wrapf(err, "found template"))
	}

	//valid commands
	where := "template_uuid = ?"
	commands, err := vault.NewTemplateCommand(ctl.NewSession()).Find(where, body.TemplateUuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(err, "NewTemplateCommand Find"))
	}
	//valid steps
	if len(body.Steps) != len(commands) {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Errorf("diff length of steps and commands expect=%d actual=%d", len(commands), len(body.Steps)))
	}
	for i := range body.Steps {
		step_args := body.Steps[i].Args
		if step_args == nil {
			return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
				errors.Wrapf(ErrorBindRequestObject(), "valid%s",
					logs.KVL(
						ParamLog(fmt.Sprintf("%s.Args", TypeName(body.Steps[i])), body.Steps[i].Args)...,
					)))
		}
		//JSON SCHEMA 유효 검사
		args_data, err := json.Marshal(step_args)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "json marshal template_command args"))
		}

		command_args := commands[i].Args
		args_schema, err := json.Marshal(command_args)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "json marshal template_command args"))
		}
		//스키마 검사 객체 생성
		validator := &jsonschema.Schema{}
		if err := json.Unmarshal([]byte(args_schema), validator); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "convert template_command args to jsonschema schema"))
		}

		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		//스키마로 유효성 검사
		valid_errors, err := validator.ValidateBytes(timeout, args_data)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "validate json schema service_step args"))
		}
		for _, valid_error := range valid_errors {
			return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(valid_error, "validate json schema service_step args"))
		}
	}

	//property
	//property service
	service := servicev1.Service{}
	service.UuidMeta = NewUuidMeta()                          //new service uuid
	service.LabelMeta = NewLabelMeta(body.Name, body.Summary) //label
	service.ClusterUuid = body.ClusterUuid
	service.TemplateUuid = body.TemplateUuid
	service.Status = newist.Int32(int32(servicev1.StatusRegist)) //init service Status(Regist)
	service.StepCount = newist.Int32(int32(len(body.Steps)))
	service.SubscribeChannel = nullable.String(body.SubscribeChannel).Value()

	//property step
	steps := map_step_create(body.Steps, func(i int, sse stepv1.HttpReqServiceStep_Create_ByService) stepv1.ServiceStep {
		command := commands[i]

		step := stepv1.ServiceStep{}
		step.UuidMeta = NewUuidMeta()                                //new step uuid
		step.LabelMeta = NewLabelMeta(command.Name, command.Summary) //label
		step.ServiceUuid = service.Uuid                              //new service uuid
		step.Sequence = newist.Int32(int32(i))                       //sequence 0 to len(steps)
		step.Status = newist.Int32(int32(servicev1.StatusRegist))    //init step Status(Regist)
		step.Method = command.Method                                 //command method
		step.Args = sse.Args                                         //step args
		step.ResultFilter = command.ResultFilter                     //command result filter
		// step.Result = newist.String("")
		return step
	})

	//property service; chaining step
	service.ServiceProperty = service.ChaniningStep(steps)

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		service, err := vault.NewService(db).Create(service)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "database create"))
		}

		if err := foreach_step(steps, func(i int, step stepv1.ServiceStep) error {
			step_, err := vault.NewServiceStep(db).Create(step)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
					errors.Wrapf(err, "database create"))
			}

			steps[i] = *step_

			return nil
		}); err != nil {
			return nil, err
		}

		return &servicev1.HttpRspService{Service: *service, Steps: steps}, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
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
// @Success     200 {array} v1.HttpRspService
func (ctl Control) FindService(ctx echo.Context) error {
	tx := ctl.NewSession()
	defer tx.Close()

	services, err := vault.NewService(tx).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "find service"))
	}
	//필터링 적용
	// 서비스 Result 필드 값 제거
	services = map_service(services, service_exclude_result)

	rsp := make([]servicev1.HttpRspService, len(services))
	if err := foreach_service(services, func(i int, s servicev1.Service) error {
		where := "service_uuid = ?"
		steps, err := vault.NewServiceStep(tx).Find(where, s.Uuid)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "find service step"))
		}

		sort.Slice(steps, func(i, j int) bool {
			var a, b int32 = 0, 0
			if steps[i].Sequence != nil {
				a = *steps[i].Sequence
			}
			if steps[j].Sequence != nil {
				b = *steps[j].Sequence
			}
			return a < b
		})

		rsp[i] = servicev1.HttpRspService{Service: s, Steps: steps}

		return nil
	}); err != nil {
		return err
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
// @Success     200 {object} v1.HttpRspService
func (ctl Control) GetService(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	service, err := vault.NewService(ctl.NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get service"))
	}
	//서비스 조회에 결과 필드는 제거
	*service = service_exclude_result(*service)

	where := "service_uuid = ?"
	steps, err := vault.NewServiceStep(ctl.NewSession()).Find(where, uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "find service step"))
	}

	sort.Slice(steps, func(i, j int) bool {
		var a, b int32 = 0, 0
		if steps[i].Sequence != nil {
			a = *steps[i].Sequence
		}
		if steps[j].Sequence != nil {
			b = *steps[j].Sequence
		}
		return a < b
	})

	return ctx.JSON(http.StatusOK, servicev1.HttpRspService{Service: *service, Steps: steps})
}

// Get Service Result
// @Description Get a Service with Result
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid}/result [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Service 의 Uuid"
// @Success     200 {object} v1.HttpRspService
func (ctl Control) GetServiceResult(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	service, err := vault.NewService(ctl.NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(errors.Wrapf(err, "get service%s",
			logs.KVL(
				"uuid", uuid,
			)))
	}

	where := "service_uuid = ?"
	steps, err := vault.NewServiceStep(ctl.NewSession()).Find(where, uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "find service step"))
	}

	return ctx.JSON(http.StatusOK, servicev1.HttpRspService{Service: *service, Steps: steps})
}

// Delete Service
// @Description Delete a Service
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Service 의 Uuid"
// @Success     200
func (ctl Control) DeleteService(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		if err := vault.NewServiceStep(db).Delete_ByService(uuid); err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "delete service step"))
		}
		if err := vault.NewService(db).Delete(uuid); err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "delete service"))
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, OK())
}

//service_exclude_result
//  서비스 조회에 결과 필드 제거
func service_exclude_result(service servicev1.Service) servicev1.Service {
	service.Result = nil //서비스 조회에 결과 필드 제거
	return service
}

func foreach_step(elems []stepv1.ServiceStep, mapper func(int, stepv1.ServiceStep) error) error {
	for n := range elems {
		if err := mapper(n, elems[n]); err != nil {
			return err
		}
	}
	return nil
}

func foreach_service(elems []servicev1.Service, mapper func(int, servicev1.Service) error) error {
	for n := range elems {
		if err := mapper(n, elems[n]); err != nil {
			return err
		}
	}
	return nil
}

func map_service(elems []servicev1.Service, mapper func(servicev1.Service) servicev1.Service) []servicev1.Service {
	rst := make([]servicev1.Service, len(elems))
	for n := range elems {
		rst[n] = mapper(elems[n])
	}
	return rst
}

func map_step(elems []stepv1.ServiceStep, mapper func(stepv1.ServiceStep) stepv1.ServiceStep) []stepv1.ServiceStep {
	rst := make([]stepv1.ServiceStep, len(elems))
	for n := range elems {
		rst[n] = mapper(elems[n])
	}
	return rst
}

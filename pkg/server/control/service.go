package control

import (
	"net/http"
	"sort"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Create Service
// @Description Create a Service
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service [post]
// @Param       service body v1.HttpReqService_Create true "HttpReqService_Create"
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
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", TypeName(body.Name),
				)))
	}
	if len(body.TemplateUuid) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", TypeName(body.TemplateUuid),
				)))
	}
	if len(body.ClusterUuid) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", TypeName(body.ClusterUuid),
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
	for _, step := range body.Steps {
		if step.Args == nil {
			return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
				errors.Wrapf(ErrorBindRequestObject(), "valid%s",
					logs.KVL(
						"pram", TypeName(step.Args),
					)))
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
	service.SubscribeEvent = body.SubscribeEvent

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
// @Param       q query string false "query  pkg/server/database/prepared/README.md"
// @Param       o query string false "order  pkg/server/database/prepared/README.md"
// @Param       p query string false "paging pkg/server/database/prepared/README.md"
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
// @Param       uuid path string true "Service 의 Uuid"
// @Success     200 {object} v1.HttpRspService
func (ctl Control) GetService(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
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
// @Param       uuid path string true "Service 의 Uuid"
// @Success     200 {object} v1.HttpRspService
func (ctl Control) GetServiceResult(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
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
// @Param       uuid path string true "Service 의 Uuid"
// @Success     200
func (ctl Control) DeleteService(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
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

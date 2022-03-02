package control

import (
	"fmt"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	. "github.com/NexClipper/sudory/pkg/server/macro"
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
// @Param       service body v1.HttpReqServiceCreate true "HttpReqServiceCreate"
// @Success     200 {object} v1.HttpRspServiceWithSteps
func (c *Control) CreateService() func(ctx echo.Context) error {
	binder := func(ctx Contexter) error {
		body := new(servicev1.HttpReqServiceCreate)
		err := ctx.Bind(body)
		if err != nil {
			return ErrorBindRequestObject(err)
		}
		if body.Name == nil {
			return ErrorInvaliedRequestParameterName("Name")
		}
		if len(body.OriginKind) == 0 {
			return ErrorInvaliedRequestParameterName("OriginKind")
		}
		if len(body.OriginUuid) == 0 {
			return ErrorInvaliedRequestParameterName("OriginUuid")
		}
		if len(body.ClusterUuid) == 0 {
			return ErrorInvaliedRequestParameterName("ClusterUuid")
		}

		err = foreach_step_essential(body.Steps, func(ss stepv1.ServiceStepEssential) error {
			//Name
			if ss.Name == nil {
				return ErrorInvaliedRequestParameterName("Name")
			}
			//Method
			if ss.Args == nil {
				return ErrorInvaliedRequestParameterName("Args")
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		body, ok := ctx.Object().(*servicev1.HttpReqServiceCreate)
		if !ok {
			return nil, ErrorFailedCast()
		}

		service_essential := body.ServiceEssential
		steps_essential := body.Steps

		//property service
		service := servicev1.Service{}
		service.Name = service_essential.Name
		service.Summary = service_essential.Summary
		service.OriginKind = &service_essential.OriginKind
		service.OriginUuid = &service_essential.OriginUuid
		service.ClusterUuid = &service_essential.ClusterUuid
		//origin template uuid
		trace, err := TraceServiceOrigin(ctx, service)
		if err != nil {
			return nil, err
		}
		if len(trace) == 0 {
			return nil, fmt.Errorf("not found origin template")
		}

		template_uuid := strings.Split(trace[len(trace)-1], ":")[1] //last
		service.TemplateUuid = newist.String(template_uuid)
		//meta
		service.UuidMeta = NewUuidMeta()
		service.LabelMeta = NewLabelMeta(service.Name, service.Summary)
		//Status
		service.Status = newist.Int32(int32(servicev1.StatusRegist))

		where := "template_uuid = ?"
		commands, err := operator.NewTemplateCommand(ctx.Database()).Find(where, template_uuid)
		if err != nil {
			return nil, errors.WithMessage(err, "not found origin template commands")
		}
		if len(steps_essential) != len(commands) {
			return nil, fmt.Errorf("diff length of steps and commands expect=%d actual=%d", len(commands), len(steps_essential))
		}

		//property step
		seq := 0
		steps := map_step_essential_to_step(steps_essential, func(ss stepv1.ServiceStepEssential) stepv1.ServiceStep {

			command := commands[seq]

			step := stepv1.ServiceStep{}
			//LabelMeta
			step.UuidMeta = NewUuidMeta()
			step.LabelMeta = NewLabelMeta(ss.Name, ss.Summary)
			//ServiceUuid
			step.ServiceUuid = service.Uuid
			//Sequence
			step.Sequence = newist.Int32(int32(seq))
			//Status = Regist
			step.Status = newist.Int32(int32(servicev1.StatusRegist))
			//Method Args
			step.Method = command.Method
			step.Args = ss.Args
			// step.Result = newist.String("")

			seq++
			return step
		})

		//create service
		err = operator.NewService(ctx.Database()).
			Create(service)
		if err != nil {
			return nil, err
		}
		//create steps
		err = foreach_step(steps, func(ss stepv1.ServiceStep) error {
			if err := operator.NewServiceStep(ctx.Database()).
				Create(ss); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		//Service Chaining
		if err := operator.NewService(ctx.Database()).Chaining(service.Uuid); err != nil {
			return nil, err
		}

		return servicev1.HttpRspServiceWithSteps{ServiceAndSteps: servicev1.ServiceAndSteps{Service: service, Steps: steps}}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Lock(c.db.Engine()),
	})
}

// Find []Service
// @Description Find []Service
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service [get]
// @Param       cluster_uuid query string false "Service 의 ClusterUuid"
// @Param       uuid         query string false "Service 의 Uuid"
// @Param       status       query string false "Service 의 Status"
// @Success     200 {array} v1.HttpRspServiceWithSteps
func (c *Control) FindService() func(ctx echo.Context) error {
	binder := func(ctx Contexter) error {
		// if len(ctx.Query()) == 0 {
		// 	return ErrorInvaliedRequestParameter()
		// }
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {

		//make condition
		args := make([]interface{}, 0)
		add, build := StringBuilder()

		for key, val := range ctx.Querys() {
			switch key {
			case "status":
				args = append(args, val)
				add(fmt.Sprintf("%s in (?)", key))
			default:
				args = append(args, val)
				add(fmt.Sprintf("%s = ?", key))
			}
		}
		//find service
		where := build(" AND ")
		services, err := operator.NewService(ctx.Database()).
			Find(where, args...)
		if err != nil {
			return nil, err
		}

		//서비스 조회에 결과 필드는 제거
		services = map_service(services, service_exclude_result)

		//make respose
		rspadd, rspbuild := servicev1.HttpRspBuilder(len(services))
		err = foreach_service(services, func(service servicev1.Service) error {
			service_uuid := service.Uuid
			where := "service_uuid = ?"
			//find steps
			steps, err := operator.NewServiceStep(ctx.Database()).
				Find(where, service_uuid)
			if err != nil {
				return err
			}
			rspadd(service, steps) //push
			return nil
		})
		if err != nil {
			return nil, err
		}

		return rspbuild(), nil //pop
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
	})
}

// Get Service
// @Description Get a Service
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid} [get]
// @Param       uuid path string true "Service 의 Uuid"
// @Success     200 {object} v1.HttpRspServiceWithSteps
func (c *Control) GetService() func(ctx echo.Context) error {
	binder := func(ctx Contexter) error {
		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}

		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		uuid := ctx.Params()[__UUID__]

		//get service
		service, err := operator.NewService(ctx.Database()).
			Get(uuid)
		if err != nil {
			return nil, err
		}

		//find step
		where := "service_uuid = ?"
		service_uuid := ctx.Params()[__UUID__]
		steps, err := operator.NewServiceStep(ctx.Database()).
			Find(where, service_uuid)
		if err != nil {
			return nil, err
		}

		//서비스 조회에 결과 필드는 제거
		*service = service_exclude_result(*service)

		rsp := &servicev1.HttpRspServiceWithSteps{ServiceAndSteps: servicev1.ServiceAndSteps{Service: *service, Steps: steps}}

		return rsp, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
	})
}

// Get Service Result
// @Description Get a Service with Result
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid}/result [get]
// @Param       uuid path string true "Service 의 Uuid"
// @Success     200 {object} v1.HttpRspServiceWithSteps
func (c *Control) GetServiceResult() func(ctx echo.Context) error {
	binder := func(ctx Contexter) error {
		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}

		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		//get service
		uuid := ctx.Params()[__UUID__]
		service, err := operator.NewService(ctx.Database()).
			Get(uuid)
		if err != nil {
			return nil, err
		}

		//find step
		where := "service_uuid = ?"
		service_uuid := ctx.Params()[__UUID__]
		steps, err := operator.NewServiceStep(ctx.Database()).
			Find(where, service_uuid)
		if err != nil {
			return nil, err
		}

		rsp := &servicev1.HttpRspServiceWithSteps{ServiceAndSteps: servicev1.ServiceAndSteps{Service: *service, Steps: steps}}

		return rsp, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
	})
}

// Update Service
// @Description Update a Service
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid} [put]
// @Param       uuid    path string true "Service 의 Uuid"
// @Param       service body v1.HttpReqService true "HttpReqService"
// @Success     200 {object} v1.HttpRspService
func (c *Control) UpdateService() func(ctx echo.Context) error {
	binder := func(ctx Contexter) error {
		body := new(servicev1.HttpReqService)
		if err := ctx.Bind(body); err != nil {
			return ErrorBindRequestObject(err)
		}

		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}

		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		body, ok := ctx.Object().(*servicev1.HttpReqService)
		if !ok {
			return nil, ErrorFailedCast()
		}

		service := body.Service

		uuid := ctx.Params()[__UUID__]

		//set uuid from path
		service.Uuid = uuid

		//update service
		err := operator.NewService(ctx.Database()).
			Update(service)
		if err != nil {
			return nil, err
		}

		return servicev1.HttpRspService{Service: service}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Lock(c.db.Engine()),
	})
}

// Delete Service
// @Description Delete a Service
// @Accept      json
// @Produce     json
// @Tags        server/service
// @Router      /server/service/{uuid} [delete]
// @Param       uuid path string true "Service 의 Uuid"
// @Success     200
func (c *Control) DeleteService() func(ctx echo.Context) error {
	binder := func(ctx Contexter) error {
		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		uuid := ctx.Params()[__UUID__]

		//steps 테이블에 데이터 있는 경우 삭제 방지
		where := "service_uuid = ?"
		steps, err := operator.NewServiceStep(ctx.Database()).Find(where, uuid)
		if err != nil {
			return nil, err
		}
		if len(steps) == 0 {
			return nil, fmt.Errorf("steps not empty")
		}

		//service 삭제
		if err := operator.NewService(ctx.Database()).Delete(uuid); err != nil {
			return nil, err
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Lock(c.db.Engine()),
	})
}

// service_exclude_result
//  서비스 조회에 결과 필드는 제거
func service_exclude_result(service servicev1.Service) servicev1.Service {
	service.Result = nil //서비스 조회에 결과 필드는 제거
	return service
}

func foreach_step(elems []stepv1.ServiceStep, fn func(stepv1.ServiceStep) error) error {
	for _, it := range elems {
		if err := fn(it); err != nil {
			return err
		}
	}
	return nil
}

func map_step(elems []stepv1.ServiceStep, mapper func(stepv1.ServiceStep) stepv1.ServiceStep) []stepv1.ServiceStep {
	rst := make([]stepv1.ServiceStep, len(elems))
	for n := range elems {
		rst[n] = mapper(elems[n])
	}
	return rst
}

func foreach_service(elems []servicev1.Service, fn func(servicev1.Service) error) error {
	for _, it := range elems {
		if err := fn(it); err != nil {
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

func foreach_client_service_and_steps(elems []servicev1.ServiceAndSteps, fn func(servicev1.Service, []stepv1.ServiceStep) error) error {
	for _, it := range elems {
		if err := fn(it.Service, it.Steps); err != nil {
			return err
		}
	}
	return nil
}

func map_client_service_and_steps(elems []servicev1.ServiceAndSteps, mapper func(servicev1.Service, []stepv1.ServiceStep) servicev1.ServiceAndSteps) []servicev1.ServiceAndSteps {
	rst := make([]servicev1.ServiceAndSteps, len(elems))
	for n := range elems {
		rst[n] = mapper(elems[n].Service, elems[n].Steps)
	}
	return rst
}

func foreach_step_essential(elems []stepv1.ServiceStepEssential, fn func(stepv1.ServiceStepEssential) error) error {
	for _, it := range elems {
		if err := fn(it); err != nil {
			return err
		}
	}
	return nil
}

func map_step_essential_to_step(elems []stepv1.ServiceStepEssential, mapper func(stepv1.ServiceStepEssential) stepv1.ServiceStep) []stepv1.ServiceStep {
	rst := make([]stepv1.ServiceStep, len(elems))
	for n := range elems {
		rst[n] = mapper(elems[n])
	}
	return rst
}

func TraceServiceOrigin(ctx Contexter, service servicev1.Service) ([]string, error) {
	trace := make([]string, 0)

	apnd := func(kind, uuid string) {
		trace = append(trace, fmt.Sprintf("%s:%s", kind, uuid))
	}

LOOP:
	for {
		switch *service.OriginKind {
		case "template":
			template, err := operator.NewTemplate(ctx.Database()).Get(*service.OriginUuid)
			if err != nil {
				return nil, err
			}
			apnd("template", template.Uuid)
			break LOOP
		case "service":
			service_, err := operator.NewService(ctx.Database()).Get(*service.OriginUuid)
			if err != nil {
				return nil, err
			}
			apnd("template", service_.Uuid)
			service = *service_
		}
	}

	return trace, nil
}

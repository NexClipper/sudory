package control

import (
	"fmt"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/database/prepared"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
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

		//origin template uuid
		trace, err := TraceServiceOrigin(ctx, service_essential.OriginKind, service_essential.OriginUuid)
		if err != nil {
			return nil, err
		}
		if len(trace) == 0 {
			return nil, fmt.Errorf("not found origin template")
		}
		template_uuid := strings.Split(trace[len(trace)-1], ":")[1] //last

		//property service
		service := servicev1.Service{}
		service.Name = nullable.String(service_essential.Name).Ptr()
		service.Summary = nullable.String(service_essential.Summary).Ptr()
		service.OriginKind = newist.String(service_essential.OriginKind)
		service.OriginUuid = newist.String(service_essential.OriginUuid)
		service.ClusterUuid = newist.String(service_essential.ClusterUuid)
		service.TemplateUuid = newist.String(template_uuid)
		service.UuidMeta = NewUuidMeta()                                //meta uuid
		service.LabelMeta = NewLabelMeta(service.Name, service.Summary) //meta label
		service.Status = newist.Int32(int32(servicev1.StatusRegist))    //Status
		service.StepCount = newist.Int32(int32(len(steps_essential)))   //step count

		//get commands
		where := "template_uuid = ?"
		commands, err := operator.NewTemplateCommand(ctx.Database()).Find(where, template_uuid)
		if err != nil {
			return nil, errors.WithMessage(err, "not found origin template commands")
		}
		if len(steps_essential) != len(commands) {
			return nil, fmt.Errorf("diff length of steps and commands expect=%d actual=%d", len(commands), len(steps_essential))
		}

		//property step
		steps := map_step_essential_to_step(steps_essential, func(sse stepv1.ServiceStepEssential, i int) stepv1.ServiceStep {

			command := commands[i]

			step := stepv1.ServiceStep{}
			step.UuidMeta = NewUuidMeta()                                //meta uuid
			step.LabelMeta = NewLabelMeta(command.Name, command.Summary) //meta label
			step.ServiceUuid = newist.String(service.Uuid)               //ServiceUuid
			step.Sequence = newist.Int32(int32(i))                       //Sequence
			step.Status = newist.Int32(int32(servicev1.StatusRegist))    //Status(Regist)
			step.Method = command.Method                                 //Method
			step.Args = sse.Args                                         //Args
			// step.Result = newist.String("")
			return step
		})

		//create service
		if err := operator.NewService(ctx.Database()).Create(service); err != nil {
			return nil, err
		}
		//create steps
		if err := foreach_step(steps, operator.NewServiceStep(ctx.Database()).Create); err != nil {
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
// @Param       q query string false "query  pkg/server/database/prepared/README.md"
// @Param       o query string false "order  pkg/server/database/prepared/README.md"
// @Param       p query string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.HttpRspServiceWithSteps
func (c *Control) FindService() func(ctx echo.Context) error {
	binder := func(ctx Contexter) error {
		// if len(ctx.Query()) == 0 {
		// 	return ErrorInvaliedRequestParameter()
		// }
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		preparer, err := prepared.NewParser(ctx.Queries())
		if err != nil {
			return nil, errors.Wrapf(err, "NewParser queries=%+v", ctx.Queries())
		}

		records := make([]servicev1.DbSchemaService, 0)
		if err := ctx.Database().Prepared(preparer).Find(&records); err != nil {
			return nil, errors.Wrapf(err, "Database Find")
		}
		services := servicev1.TransFormDbSchema(records)

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
				return errors.Wrapf(err, "NewServiceStep Find where=%s service_uuid=%s", where, service_uuid)
			}
			rspadd(service, steps) //push
			return nil
		})
		if err != nil {
			return nil, errors.Wrapf(err, "foreach_service")
		}

		return rspbuild(), nil //pop
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Contexter) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "FindService binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "FindService operator")
			}
			return v, nil
		},
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

func map_step_essential_to_step(elems []stepv1.ServiceStepEssential, mapper func(stepv1.ServiceStepEssential, int) stepv1.ServiceStep) []stepv1.ServiceStep {
	rst := make([]stepv1.ServiceStep, len(elems))
	for n := range elems {
		rst[n] = mapper(elems[n], n)
	}
	return rst
}

func TraceServiceOrigin(ctx Contexter, originkind, originuuid string) ([]string, error) {
	trace := make([]string, 0)

	apnd := func(kind, uuid string) {
		trace = append(trace, fmt.Sprintf("%s:%s", kind, uuid))
	}

LOOP:
	for {
		switch originkind {
		case "template":
			template, err := operator.NewTemplate(ctx.Database()).Get(originuuid)
			if err != nil {
				return nil, err
			}
			apnd("template", template.Uuid)
			break LOOP
		case "service":
			service, err := operator.NewService(ctx.Database()).Get(originuuid)
			if err != nil {
				return nil, err
			}
			apnd("template", service.Uuid)
			originkind, originuuid = *service.OriginKind, *service.OriginUuid
		}
	}

	return trace, nil
}

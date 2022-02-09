package control

import (
	"fmt"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/database"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	"github.com/labstack/echo/v4"
)

// Create Service
// @Description Create a Service
// @Accept json
// @Produce json
// @Tags server/service
// @Router /server/service [post]
// @Param service body v1.HttpReqServiceCreate true "HttpReqServiceCreate"
// @Success 200
func (c *Control) CreateService() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := new(servicev1.HttpReqServiceCreate)
		err := ctx.Bind(req)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}
		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(*servicev1.HttpReqServiceCreate)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//create service
		err := operator.NewService(ctx).
			Create(req.Service)
		if err != nil {
			return nil, err
		}
		//create step
		err = foreach_step(req.Steps, func(step stepv1.ServiceStep) error {
			if err := operator.NewServiceStep(ctx).
				Create(step); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		//Service Chaining
		operator.NewService(ctx).
			Chaining(req.Service.Uuid)

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    Lock,
	})
}

// Find []Service
// @Description Find []Service
// @Accept json
// @Produce json
// @Tags server/service
// @Router /server/service [get]
// @Param cluster_uuid query string false "Service 의 ClusterUuid"
// @Param uuid         query string false "Service 의 Uuid"
// @Param status       query string false "Service 의 Status"
// @Success 200 {array} v1.HttpRspService
func (c *Control) FindService() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for key := range ctx.QueryParams() {
			req[key] = ctx.QueryParam(key)
		}
		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		args := make([]interface{}, 0)
		add, build := StringBuilder()

		for key, val := range req {
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
		services, err := operator.NewService(ctx).
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
			steps, err := operator.NewServiceStep(ctx).
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
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    NoLock,
	})
}

// Get Service
// @Description Get a Service
// @Accept json
// @Produce json
// @Tags server/service
// @Router /server/service/{uuid} [get]
// @Param uuid path string true "Service 의 Uuid"
// @Success 200 {object} v1.HttpRspService
func (c *Control) GetService() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req[__UUID__]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		//get service
		uuid := req[__UUID__]
		service, err := operator.NewService(ctx).
			Get(uuid)
		if err != nil {
			return nil, err
		}

		//find step
		where := "service_uuid = ?"
		service_uuid := req[__UUID__]
		steps, err := operator.NewServiceStep(ctx).Find(where, service_uuid)
		if err != nil {
			return nil, err
		}

		//서비스 조회에 결과 필드는 제거
		*service = service_exclude_result(*service)

		rsp := &servicev1.HttpRspService{Service: *service, Steps: steps}

		return rsp, nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    NoLock,
	})
}

// Get Service Result
// @Description Get a Service with Result
// @Accept json
// @Produce json
// @Tags server/service
// @Router /server/service/{uuid}/result [get]
// @Param uuid path string true "Service 의 Uuid"
// @Success 200 {object} v1.HttpRspService
func (c *Control) GetServiceResult() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req[__UUID__]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		//get service
		uuid := req[__UUID__]
		service, err := operator.NewService(ctx).
			Get(uuid)
		if err != nil {
			return nil, err
		}

		//find step
		where := "service_uuid = ?"
		service_uuid := req[__UUID__]
		steps, err := operator.NewServiceStep(ctx).Find(where, service_uuid)
		if err != nil {
			return nil, err
		}

		rsp := &servicev1.HttpRspService{Service: *service, Steps: steps}

		return rsp, nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    NoLock,
	})
}

// Update Service
// @Description Update a Service
// @Accept json
// @Produce json
// @Tags server/service
// @Router /server/service/{uuid} [put]
// @Param uuid    path string true "Service 의 Uuid"
// @Param service body v1.HttpReqService true "HttpReqService"
// @Success 200
func (c *Control) UpdateService() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req[__UUID__].(string)) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		body := new(servicev1.HttpReqService)
		err := ctx.Bind(body)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}
		req[__BODY__] = body

		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}
		uuid, ok := req[__UUID__].(string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		body, ok := req[__BODY__].(*servicev1.HttpReqService)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//set uuid from path
		body.Service.Uuid = uuid

		//update service
		err := operator.NewService(ctx).
			Update(body.Service)
		if err != nil {
			return nil, err
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    Lock,
	})
}

// Delete Service
// @Description Delete a Service
// @Accept json
// @Produce json
// @Tags server/service
// @Router /server/service/{uuid} [delete]
// @Param uuid path string true "Service 의 Uuid"
// @Success 200
func (c *Control) DeleteService() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req[__UUID__]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//delete service
		uuid := req[__UUID__]
		err := operator.NewService(ctx).
			Delete(uuid)
		if err != nil {
			return nil, err
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    Lock,
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

func foreach_service(elems []servicev1.Service, fn func(servicev1.Service) error) error {
	for _, it := range elems {
		if err := fn(it); err != nil {
			return err
		}
	}
	return nil
}

func foreach_client_service(elems []servicev1.HttpReqClientSideService, fn func(servicev1.Service, []stepv1.ServiceStep) error) error {
	for _, it := range elems {
		if err := fn(it.Service, it.Steps); err != nil {
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

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

		return nil, nil
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
		for key, _ := range ctx.QueryParams() {
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
		join, build := StringJoin()

		for key, val := range req {
			switch key {
			case "status":
				args = append(args, val)
				join(fmt.Sprintf("%s in (?)", key))
			default:
				args = append(args, val)
				join(fmt.Sprintf("%s = ?", key))
			}
		}
		//find service
		where := build(" AND ")
		services, err := operator.NewService(ctx).
			Find(where, args...)
		if err != nil {
			return nil, err
		}
		//make respose
		push, pop := servicev1.HttpRspBuilder(len(services))
		err = foreach_service(services, func(service servicev1.Service) error {
			service_uuid := service.Uuid
			where := "service_uuid = ?"
			//find steps
			steps, err := operator.NewServiceStep(ctx).
				Find(where, service_uuid)
			if err != nil {
				return err
			}
			push(service, steps) //push
			return nil
		})
		if err != nil {
			return nil, err
		}

		return pop(), nil //pop
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
		if len(req["uuid"]) == 0 {
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
		uuid := req["uuid"]
		service, err := operator.NewService(ctx).
			Get(uuid)
		if err != nil {
			return nil, err
		}
		//find step
		where := "service_uuid = ?"
		service_uuid := req["uuid"]
		steps, err := operator.NewServiceStep(ctx).Find(where, service_uuid)
		if err != nil {
			return nil, err
		}

		return &servicev1.HttpRspService{Service: *service, Steps: steps}, nil
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
		if len(req["uuid"].(string)) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		body := new(servicev1.HttpReqService)
		err := ctx.Bind(body)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}
		req["_"] = body

		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}
		uuid, ok := req["uuid"].(string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		body, ok := req["_"].(*servicev1.HttpReqService)
		if !ok {
			return nil, ErrorFailedCast()
		}

		body.Service.Uuid = uuid
		//update service
		err := operator.NewService(ctx).
			Update(body.Service)
		if err != nil {
			return nil, err
		}
		return nil, nil
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
		if len(req["uuid"]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//find step
		where := "service_uuid = ?"
		service_uuid := req["uuid"]
		steps, err := operator.NewServiceStep(ctx).
			Find(where, service_uuid)
		if err != nil {
			return nil, err
		}
		//delete step
		err = foreach_step(steps, func(step stepv1.ServiceStep) error {
			err := operator.NewServiceStep(ctx).
				Delete(step.Uuid)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		//delete service
		uuid := req["uuid"]
		err = operator.NewService(ctx).
			Delete(uuid)
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    Lock,
	})
}

// Pull []Service (client)
// @Description Get a Service
// @Accept json
// @Produce json
// @Tags client/service
// @Router /client/service [put]
// @Param cluster_uuid query string false "Client 의 ClusterUuid"
// @Param service      body []v1.HttpReqClientSideService true "HttpReqClientSideService"
// @Success 200 {array} v1.HttpRspClientSideService
func (c *Control) PullClientServices() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for key, _ := range ctx.QueryParams() {
			req[key] = ctx.QueryParam(key)
		}
		if len(req["cluster_uuid"].(string)) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		body := make([]servicev1.HttpReqClientSideService, 0)
		err := ctx.Bind(&body)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}
		req["_"] = body
		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}
		body, ok := req["_"].([]servicev1.HttpReqClientSideService)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//update service
		err := foreach_client_service(body, func(service servicev1.Service, steps []stepv1.ServiceStep) error {

			//update service
			err := operator.NewService(ctx).
				Update(service)
			if err != nil {
				return err
			}

			//update step
			err = foreach_step(steps, func(step stepv1.ServiceStep) error {
				err := operator.NewServiceStep(ctx).
					Update(step)
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			return nil, err
		}

		//find service
		where := "cluster_uuid = ? AND status < ?"
		cluster_uuid := req["cluster_uuid"].(string)
		status := servicev1.StatusSuccess //상태 값이 완료 상태보다 작은것

		services, err := operator.NewService(ctx).
			Find(where, cluster_uuid, status)
		if err != nil {
			return nil, err
		}
		//make response
		push, pop := servicev1.HttpRspBuilder(len(services))
		err = foreach_service(services, func(service servicev1.Service) error {
			service_uuid := service.Uuid
			where := "service_uuid = ?"
			//find steps
			steps, err := operator.NewServiceStep(ctx).
				Find(where, service_uuid)
			if err != nil {
				return err
			}
			push(service, steps) //push
			return nil
		})
		if err != nil {
			return nil, err
		}

		return pop(), nil //pop
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    Lock,
	})
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

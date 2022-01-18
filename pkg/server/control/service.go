package control

import (
	"fmt"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	"github.com/labstack/echo/v4"
)

// Create Service
// @Description Create a Service
// @Accept json
// @Produce json
// @Tags server
// @Router /server/service [post]
// @Param service body v1.HttpReqService true "HttpReqService"
// @Success 200
func (c *Control) CreateService() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := new(servicev1.HttpReqService)
		err := ctx.Bind(req)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}
		return req, nil
	}
	operator := func(v interface{}) (interface{}, error) {
		req, ok := v.(*servicev1.HttpReqService)
		if !ok {
			return nil, ErrorFailedCast()
		}

		err := operator.NewService(c.db).
			Create(req.Service)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	return MakeMiddlewareFunc(binder, operator, HttpResponse)
}

// Get []Service
// @Description Get a Servicies
// @Accept json
// @Produce json
// @Tags server
// @Router /server/service [get]
// @Param cluster_uuid query string false "Service 의 ClusterUuid"
// @Param uuid         query string false "Service 의 Uuid"
// @Param status       query string false "Service 의 Status"
// @Success 200 {array} v1.HttpRspService
func (c *Control) GetServicies() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for key, _ := range ctx.QueryParams() {
			req[key] = ctx.QueryParam(key)
		}
		// if len(req["cluster_uuid"]) == 0 {
		// 	return nil, ErrorInvaliedRequestParameter()
		// }
		return req, nil
	}
	operator := func(v interface{}) (interface{}, error) {
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

		where := build(" AND ")

		records, err := operator.NewService(c.db).
			Find(where, args...)
		if err != nil {
			return nil, err
		}

		get_step := func(service_uuid string) ([]stepv1.ServiceStep, error) {

			where := "service_uuid = ?"
			records, err := operator.NewServiceStep(c.db).
				Find(where, service_uuid)
			if err != nil {
				return nil, err
			}
			return records, nil
		}

		rsp := make([]servicev1.HttpRspService, len(records))
		for n, it := range records {
			steps, err := get_step(it.Uuid)
			if err != nil {
				return nil, err
			}
			rsp[n].Service = it
			rsp[n].Steps = steps
		}

		return rsp, nil
	}
	return MakeMiddlewareFunc(binder, operator, HttpResponse)
}

// Get Service
// @Description Get a Service
// @Accept json
// @Produce json
// @Tags server
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
	operator := func(v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		uuid := req["uuid"]

		record, err := operator.NewService(c.db).
			Get(uuid)
		if err != nil {
			return nil, err
		}

		where := "service_uuid = ?"
		service_uuid := req["uuid"]

		steps, err := operator.NewServiceStep(c.db).Find(where, service_uuid)
		if err != nil {
			return nil, err
		}
		for _, it := range steps {
			err := operator.NewServiceStep(c.db).Delete(it.Uuid)
			if err != nil {
				return nil, err
			}
		}

		return &servicev1.HttpRspService{Service: *record, Steps: steps}, nil
	}
	return MakeMiddlewareFunc(binder, operator, HttpResponse)
}

// Update Service
// @Description Update a Service
// @Accept json
// @Produce json
// @Tags server
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
	operator := func(v interface{}) (interface{}, error) {
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
		err := operator.NewService(c.db).
			Update(body.Service)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}

	return MakeMiddlewareFunc(binder, operator, HttpResponse)
}

// Delete Service
// @Description Delete a Service
// @Accept json
// @Produce json
// @Tags server
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
	operator := func(v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		uuid := req["uuid"]

		err := operator.NewService(c.db).
			Delete(uuid)
		if err != nil {
			return nil, err
		}

		where := "service_uuid = ?"
		service_uuid := req["uuid"]

		steps, err := operator.NewServiceStep(c.db).Find(where, service_uuid)
		if err != nil {
			return nil, err
		}
		for _, it := range steps {
			err := operator.NewServiceStep(c.db).Delete(it.Uuid)
			if err != nil {
				return nil, err
			}
		}

		return nil, nil
	}
	return MakeMiddlewareFunc(binder, operator, HttpResponse)
}

// Get []Service (client)
// @Description Get a Servicies
// @Accept json
// @Produce json
// @Tags server
// @Router /client/service [get]
// @Param cluster_uuid query string false "Service 의 ClusterUuid"
// @Success 200 {array} v1.HttpRspService
func (c *Control) GetClientServicies() func(ctx echo.Context) error {
	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for key, _ := range ctx.QueryParams() {
			req[key] = ctx.QueryParam(key)
		}
		if len(req["cluster_uuid"]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//서비스 상태 값이 완료 상태보다 작은것
		where := "cluster_uuid = ? AND status < ?"
		cluster_uuid := req["cluster_uuid"]
		status := servicev1.StatusDone

		records, err := operator.NewService(c.db).
			Find(where, cluster_uuid, status)
		if err != nil {
			return nil, err
		}

		get_step := func(service_uuid string) ([]stepv1.ServiceStep, error) {

			where := "service_uuid = ?"
			records, err := operator.NewServiceStep(c.db).
				Find(where, service_uuid)
			if err != nil {
				return nil, err
			}
			return records, nil
		}

		rsp := make([]servicev1.HttpRspService, len(records))
		for n, it := range records {
			steps, err := get_step(it.Uuid)
			if err != nil {
				return nil, err
			}
			rsp[n].Service = it
			rsp[n].Steps = steps
		}

		return rsp, nil
	}
	return MakeMiddlewareFunc(binder, operator, HttpResponse)
}

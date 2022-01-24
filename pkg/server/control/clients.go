package control

import (
	"errors"
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/database"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	"github.com/labstack/echo/v4"
)

// Poll []Service (client)
// @Description Poll a Service
// @Accept json
// @Produce json
// @Tags client/service
// @Router /client/service [put]
// @Param cluster_uuid query string false "Client 의 ClusterUuid"
// @Param service      body []v1.HttpReqClientSideService true "HttpReqClientSideService"
// @Success 200 {array} v1.HttpRspClientSideService
func (c *Control) PollService() func(ctx echo.Context) error {
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

// Auth Client
// @Description Auth a client
// @Accept json
// @Produce json
// @Tags client/auth
// @Router /client/auth [post]
// @Param client_uuid  query string true "Client 의 Uuid"
// @Param cluster_uuid query string true "Cluster 의 Uuid"
// @Success 200 {object} v1.HttpRspClient
func (c *Control) AuthClient() func(ctx echo.Context) error {

	set_cookie := func(ctx echo.Context, status int, v interface{}) error {
		const (
			exp = 60 //minute
		)
		cookie := new(http.Cookie)
		cookie.Name = "client-token"
		cookie.Value = "jon"
		cookie.Expires = time.Now().Add(exp * time.Minute)

		err := HttpResponse(ctx, status, v)
		if err != nil {
			return err
		}
		return nil
	}

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for key, _ := range ctx.QueryParams() {
			req[key] = ctx.QueryParam(key)
		}
		if len(req["client_uuid"]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		if len(req["cluster_uuid"]) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		client_uuid := req["client_uuid"]
		client, err := operator.NewClient(ctx).
			Get(client_uuid)
		if err != nil {
			return nil, err
		}

		cluster_uuid := req["cluster_uuid"]
		cluster, err := operator.NewCluster(ctx).
			Get(cluster_uuid)
		if err != nil {
			return nil, err
		}

		if client == nil || cluster == nil {
			return nil, errors.New("")
		}

		return nil, nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: set_cookie,
		BlockMaker:    Lock,
	})
}

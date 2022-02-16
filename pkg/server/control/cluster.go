package control

import (
	"fmt"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/macro"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	clusterv1 "github.com/NexClipper/sudory/pkg/server/model/cluster/v1"
	"github.com/labstack/echo/v4"
)

// Create Cluster
// @Description Create a cluster
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster [post]
// @Param       client body v1.HttpReqCluster true "HttpReqCluster"
// @Success     200 {object} v1.HttpRspCluster
func (c *Control) CreateCluster() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {

		body := new(clusterv1.HttpReqCluster)
		err := ctx.Bind(body)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}

		return body, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {

		body, ok := ctx.Req.(*clusterv1.HttpReqCluster)
		if !ok {
			return nil, ErrorFailedCast()
		}

		cluster := body.Cluster

		//property
		cluster.LabelMeta = NewLabelMeta(body.Name, body.Summary)

		//create
		err := operator.NewCluster(ctx.Database).
			Create(cluster)
		if err != nil {
			return nil, err
		}

		return cluster, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Lock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// Find Cluster
// @Description Find cluster
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster [get]
// @Param       uuid query string false "Cluster 의 Uuid"
// @Param       name query string false "Cluster 의 Name"
// @Success     200 {array} v1.HttpRspCluster
func (c *Control) FindCluster() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]string)
		for key := range ctx.QueryParams() {
			req[key] = ctx.QueryParam(key)
		}
		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//make condition
		args := make([]interface{}, 0)
		add, build := StringBuilder()

		for key, val := range req {
			switch key {
			case "uuid":
				args = append(args, fmt.Sprintf("%s%%", val)) //앞 부분 부터 일치 해야함
			default:
				args = append(args, fmt.Sprintf("%%%s%%", val))
			}
			add(fmt.Sprintf("%s LIKE ?", key)) //조건문 만들기
		}
		where := build(" AND ")

		//find client
		clusters, err := operator.NewCluster(ctx.Database).
			Find(where, args...)
		if err != nil {
			return nil, err
		}
		return clusterv1.TransToHttpRsp(clusters), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Nolock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// Get Cluster
// @Description Get a cluster
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid} [get]
// @Param       uuid path string true "Cluster 의 Uuid"
// @Success     200 {object} v1.HttpRspCluster
func (c *Control) GetCluster() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if _, ok := macro.MapString(req, __UUID__); !ok {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}

		uuid, _ := macro.MapString(req, __UUID__)
		cluster, err := operator.NewCluster(ctx.Database).
			Get(uuid)
		if err != nil {
			return nil, err
		}
		return clusterv1.HttpRspCluster{Cluster: *cluster}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Nolock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// Update Cluster
// @Description Update a cluster
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid} [put]
// @Param       uuid   path string true "Cluster 의 Uuid"
// @Param       client body v1.HttpReqCluster true "HttpReqCluster"
// @Success     200 {object} v1.HttpRspCluster
func (c *Control) UpdateCluster() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if _, ok := macro.MapString(req, __UUID__); !ok {
			return nil, ErrorInvaliedRequestParameter()
		}

		body := new(clusterv1.HttpReqCluster)
		err := ctx.Bind(body)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}
		req[__BODY__] = body

		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}
		uuid, _ := macro.MapString(req, __UUID__)

		body, ok := req[__BODY__].(*clusterv1.HttpReqCluster)
		if !ok {
			return nil, ErrorFailedCast()
		}

		cluster := body.Cluster

		//set uuid from path
		cluster.Uuid = uuid

		err := operator.NewCluster(ctx.Database).
			Update(cluster)
		if err != nil {
			return nil, err
		}

		return cluster, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Lock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// Delete Cluster
// @Description Delete a cluster
// @Accept json
// @Produce json
// @Tags server/cluster
// @Router /server/cluster/{uuid} [delete]
// @Param uuid path string true "Cluster 의 Uuid"
// @Success 200
func (c *Control) DeleteCluster() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if _, ok := macro.MapString(req, __UUID__); !ok {
			return nil, ErrorInvaliedRequestParameter()
		}
		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}

		uuid, _ := macro.MapString(req, __UUID__)

		err := operator.NewCluster(ctx.Database).
			Delete(uuid)
		if err != nil {
			return nil, err
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Lock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

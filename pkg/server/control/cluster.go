package control

import (
	"fmt"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/database"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	clusterv1 "github.com/NexClipper/sudory/pkg/server/model/cluster/v1"
	"github.com/labstack/echo/v4"
)

// Create Cluster
// @Description Create a cluster
// @Accept json
// @Produce json
// @Tags server/cluster
// @Router /server/cluster [post]
// @Param client body v1.HttpReqCluster true "HttpReqCluster"
// @Success 200
func (c *Control) CreateCluster() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}

		body := new(clusterv1.HttpReqCluster)
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
		body, ok := req[__BODY__].(*clusterv1.HttpReqCluster)
		if !ok {
			return nil, ErrorFailedCast()
		}

		err := operator.NewCluster(ctx).
			Create(body.Cluster)
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

// Find Cluster
// @Description Find cluster
// @Accept json
// @Produce json
// @Tags server/cluster
// @Router /server/cluster [get]
// @Param uuid query string false "Cluster 의 Uuid"
// @Param name query string false "Cluster 의 Name"
// @Success 200 {array} v1.HttpRspCluster
func (c *Control) FindCluster() func(ctx echo.Context) error {

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

		//make condition
		args := make([]interface{}, 0)
		join, build := StringJoin()

		for key, val := range req {
			switch key {
			case "uuid":
				args = append(args, fmt.Sprintf("%s%%", val)) //앞 부분 부터 일치 해야함
			default:
				args = append(args, fmt.Sprintf("%%%s%%", val))
			}
			join(fmt.Sprintf("%s LIKE ?", key)) //조건문 만들기
		}
		where := build(" AND ")

		//find client
		rst, err := operator.NewCluster(ctx).
			Find(where, args...)

		return clusterv1.TransToHttpRsp(rst), err
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    NoLock,
	})
}

// Get Cluster
// @Description Get a cluster
// @Accept json
// @Produce json
// @Tags server/cluster
// @Router /server/cluster/{uuid} [get]
// @Param uuid          path string true "Cluster 의 Uuid"
// @Success 200 {object} v1.HttpRspCluster
func (c *Control) GetCluster() func(ctx echo.Context) error {

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

		uuid := req[__UUID__]
		rst, err := operator.NewCluster(ctx).
			Get(uuid)
		if err != nil {
			return nil, err
		}
		return clusterv1.HttpRspCluster{Cluster: *rst}, nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    NoLock,
	})
}

// Update Cluster
// @Description Update a cluster
// @Accept json
// @Produce json
// @Tags server/cluster
// @Router /server/cluster/{uuid} [put]
// @Param uuid   path string true "Cluster 의 Uuid"
// @Param client body v1.HttpReqCluster true "HttpReqCluster"
// @Success 200
func (c *Control) UpdateCluster() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req[__UUID__].(string)) == 0 {
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
	operator := func(ctx database.Context, v interface{}) (interface{}, error) {
		req, ok := v.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}
		uuid, ok := req[__UUID__].(string)
		if !ok {
			return nil, ErrorFailedCast()
		}
		body, ok := req[__BODY__].(*clusterv1.HttpReqCluster)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//set uuid from path
		body.Cluster.Uuid = uuid

		err := operator.NewCluster(ctx).
			Update(body.Cluster)
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

		uuid := req[__UUID__]
		err := operator.NewCluster(ctx).
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

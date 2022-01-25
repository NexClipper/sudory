package control

import (
	"fmt"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/database"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	clinetv1 "github.com/NexClipper/sudory/pkg/server/model/client/v1"
	"github.com/labstack/echo/v4"
)

// Create Client
// @Description Create a client
// @Accept json
// @Produce json
// @Tags server/client
// @Router /server/client [post]
// @Param client body v1.HttpReqClient true "HttpReqClient"
// @Success 200
func (c *Control) CreateClient() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}

		body := new(clinetv1.HttpReqClient)
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
		body, ok := req["_"].(*clinetv1.HttpReqClient)
		if !ok {
			return nil, ErrorFailedCast()
		}

		err := operator.NewClient(ctx).
			Create(body.Client)
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

// Find Client
// @Description Find client
// @Accept json
// @Produce json
// @Tags server/client
// @Router /server/client [get]
// @Param uuid         query string false "Client 의 Uuid"
// @Param name         query string false "Client 의 Name"
// @Param cluster_uuid query string false "Client 의 ClusterUuid"
// @Success 200 {array} v1.HttpRspClient
func (c *Control) FindClient() func(ctx echo.Context) error {

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
		rst, err := operator.NewClient(ctx).
			Find(where, args...)

		return clinetv1.TransToHttpRsp(rst), err
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    NoLock,
	})
}

// Get Client
// @Description Get a client
// @Accept json
// @Produce json
// @Tags server/client
// @Router /server/client/{uuid} [get]
// @Param uuid          path string true "Client 의 Uuid"
// @Success 200 {object} v1.HttpRspClient
func (c *Control) GetClient() func(ctx echo.Context) error {

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

		uuid := req["uuid"]
		rst, err := operator.NewClient(ctx).
			Get(uuid)
		if err != nil {
			return nil, err
		}
		return clinetv1.HttpRspClient{Client: *rst}, nil
	}

	return MakeMiddlewareFunc(Option{
		Engine:        c.db.Engine(),
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		BlockMaker:    NoLock,
	})
}

// Update Client
// @Description Update a client
// @Accept json
// @Produce json
// @Tags server/client
// @Router /server/client/{uuid} [put]
// @Param uuid   path string true "Client 의 Uuid"
// @Param client body v1.HttpReqClient true "HttpReqClient"
// @Success 200
func (c *Control) UpdateClient() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if len(req["uuid"].(string)) == 0 {
			return nil, ErrorInvaliedRequestParameter()
		}

		body := new(clinetv1.HttpReqClient)
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
		body, ok := req["_"].(*clinetv1.HttpReqClient)
		if !ok {
			return nil, ErrorFailedCast()
		}

		body.Client.Uuid = uuid

		err := operator.NewClient(ctx).
			Update(body.Client)
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

// Delete Client
// @Description Delete a client
// @Accept json
// @Produce json
// @Tags server/client
// @Router /server/client/{uuid} [delete]
// @Param uuid path string true "Client 의 Uuid"
// @Success 200
func (c *Control) DeleteClient() func(ctx echo.Context) error {

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

		uuid := req["uuid"]
		err := operator.NewClient(ctx).
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

package control

import (
	"fmt"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/macro"
	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	envv1 "github.com/NexClipper/sudory/pkg/server/model/environment/v1"
	"github.com/labstack/echo/v4"
)

// Find Environment
// @Description Find Environment
// @Accept      json
// @Produce     json
// @Tags        server/environment
// @Router      /server/environment [get]
// @Param       uuid    query string false "Environment 의 Uuid"
// @Param       summary query string false "Environment 의 Summary"
// @Param       name    query string false "Environment 의 Name"
// @Param       value   query string false "Environment 의 Value"
// @Success 200 {array} v1.HttpRspEnvironment
func (c *Control) FindEnvironment() func(ctx echo.Context) error {

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

		//find Environment
		rst, err := operator.NewEnvironment(ctx.Database).
			Find(where, args...)
		if err != nil {
			return nil, err
		}
		return envv1.TransToHttpRsp(rst), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Nolock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// Get Environment
// @Description Get a Environment
// @Accept      json
// @Produce     json
// @Tags        server/environment
// @Router      /server/environment/{uuid} [get]
// @Param       uuid path string true "Environment 의 Uuid"
// @Success 200 {object} v1.HttpRspEnvironment
func (c *Control) GetEnvironment() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if _, ok := macro.MapString(req, __UUID__); !ok {
			return nil, ErrorInvaliedRequestParameterName(__UUID__)
		}
		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}

		uuid, _ := macro.MapString(req, __UUID__)
		rst, err := operator.NewEnvironment(ctx.Database).
			Get(uuid)
		if err != nil {
			return nil, err
		}
		return envv1.HttpRspEnvironment{Environment: *rst}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Nolock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// UpdateEnvironmentValue
// @Description Update Environment Value
// @Accept      x-www-form-urlencoded
// @Produce     json
// @Tags        server/environment
// @Router      /server/environment/{uuid}/value [put]
// @Param       uuid   path     string true "Environment 의 Uuid"
// @Param       value  formData string true "Environment 의 Value"
// @Success 200 {object} v1.HttpRspEnvironment
func (c *Control) UpdateEnvironmentValue() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		for key := range ctx.QueryParams() {
			req[key] = ctx.QueryParam(key)
		}
		formdatas, err := ctx.FormParams()
		if err != nil {
			return nil, err
		}
		for key := range formdatas {
			req[key] = ctx.FormValue(key)
		}
		if _, ok := macro.MapString(req, __UUID__); !ok {
			return nil, ErrorInvaliedRequestParameterName(__UUID__)
		}
		if _, ok := macro.MapString(req, __VALUE__); !ok {
			return nil, ErrorInvaliedRequestParameterName(__VALUE__)
		}
		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}

		uuid, _ := macro.MapString(req, __UUID__)
		value, _ := macro.MapString(req, __VALUE__)

		//get record
		env, err := operator.NewEnvironment(ctx.Database).Get(uuid)
		if err != nil {
			return nil, err
		}

		//udate value
		env.Value = newist.String(value)

		//update record
		err = operator.NewEnvironment(ctx.Database).
			Update(*env)
		if err != nil {
			return nil, err
		}

		return env, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Lock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

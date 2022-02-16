package control

import (
	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/database/query_parser"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
	"github.com/labstack/echo/v4"
)

// Find Session
// @Description Find Session
// @Accept json
// @Produce json
// @Tags server/session
// @Router /server/session [get]
// @Param uuid         query string false "Session 의 Uuid"
// @Param name         query string false "Session 의 Name"
// @Param user_kind    query string false "Session 의 user_kind"
// @Param user_uuid    query string false "Session 의 user_uuid"
// #Param limit        query int    false "Pagination 의 limit"
// #Param page         query int    false "Pagination 의 page"
// #Param order        query string false "Pagination 의 order"
// @Success 200 {array} v1.HttpRspSession
func (c *Control) FindSession() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for key := range ctx.QueryParams() {
			req[key] = ctx.QueryParam(key)
		}
		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}

		cond := query_parser.NewCondition(req, func(key string) (string, string, bool) {
			switch key {
			case "uuid", "user_uuid", "user_kind":
				return "=", "%s", true
			case "name":
				return "LIKE", "%%%s%%", true
			default:
				return "", "", false
			}
		})

		//find Session
		records, err := operator.NewSession(ctx.Database).
			Find(cond.Where(), cond.Args()...)
		if err != nil {
			return nil, err
		}

		return sessionv1.TransToHttpRsp(records), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Nolock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// Get Session
// @Description Get a Session
// @Accept json
// @Produce json
// @Tags server/session
// @Router /server/session/{uuid} [get]
// @Param uuid          path string true "Session 의 Uuid"
// @Success 200 {object} v1.HttpRspSession
func (c *Control) GetSession() func(ctx echo.Context) error {

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
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		uuid := req[__UUID__]
		rst, err := operator.NewSession(ctx.Database).
			Get(uuid)
		if err != nil {
			return nil, err
		}
		return sessionv1.HttpRspSession{Session: *rst}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Nolock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// Delete Session
// @Description Delete a Session
// @Accept json
// @Produce json
// @Tags server/session
// @Router /server/session/{uuid} [delete]
// @Param uuid path string true "Session 의 Uuid"
// @Success 200
func (c *Control) DeleteSession() func(ctx echo.Context) error {

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
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]string)
		if !ok {
			return nil, ErrorFailedCast()
		}

		uuid := req[__UUID__]
		err := operator.NewSession(ctx.Database).
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

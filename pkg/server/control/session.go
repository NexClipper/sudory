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

	binder := func(ctx Contexter) error {
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {

		cond := query_parser.NewCondition(ctx.Querys(), func(key string) (string, string, bool) {
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
		records, err := operator.NewSession(ctx.Database()).
			Find(cond.Where(), cond.Args()...)
		if err != nil {
			return nil, err
		}

		return sessionv1.TransToHttpRsp(records), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
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
		rst, err := operator.NewSession(ctx.Database()).
			Get(uuid)
		if err != nil {
			return nil, err
		}
		return sessionv1.HttpRspSession{Session: *rst}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
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

		err := operator.NewSession(ctx.Database()).
			Delete(uuid)
		if err != nil {
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

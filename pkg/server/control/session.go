package control

import (
	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/database/prepared"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Find Session
// @Description Find Session
// @Accept      json
// @Produce     json
// @Tags        server/session
// @Router      /server/session [get]
// @Param       q query string false "query  pkg/server/database/prepared/README.md"
// @Param       o query string false "order  pkg/server/database/prepared/README.md"
// @Param       p query string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.HttpRspSession
func (c *Control) FindSession() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		preparer, err := prepared.NewParser(ctx.Queries())
		if err != nil {
			return nil, errors.Wrapf(err, "NewParser queries=%+v", ctx.Queries())
		}

		records := make([]sessionv1.DbSchemaSession, 0)
		if err := ctx.Database().Prepared(preparer).Find(&records); err != nil {
			return nil, errors.Wrapf(err, "Database Find")
		}
		return sessionv1.TransToHttpRsp(sessionv1.TransFormDbSchema(records)), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Contexter) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "FindSession binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "FindSession operator")
			}
			return v, nil
		},
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
	})
}

// Get Session
// @Description Get a Session
// @Accept      json
// @Produce     json
// @Tags        server/session
// @Router      /server/session/{uuid} [get]
// @Param       uuid          path string true "Session 의 Uuid"
// @Success     200 {object} v1.HttpRspSession
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
// @Accept      json
// @Produce     json
// @Tags        server/session
// @Router      /server/session/{uuid} [delete]
// @Param       uuid path string true "Session 의 Uuid"
// @Success     200
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

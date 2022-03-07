package control

import (
	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/database/prepared"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
	envv1 "github.com/NexClipper/sudory/pkg/server/model/environment/v1"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Find Environment
// @Description Find Environment
// @Accept      json
// @Produce     json
// @Tags        server/environment
// @Router      /server/environment [get]
// @Param       q query string false "query  pkg/server/database/prepared/README.md"
// @Param       o query string false "order  pkg/server/database/prepared/README.md"
// @Param       p query string false "paging pkg/server/database/prepared/README.md"
// @Success 200 {array} v1.HttpRspEnvironment
func (c *Control) FindEnvironment() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		preparer, err := prepared.NewParser(ctx.Queries())
		if err != nil {
			return nil, errors.Wrapf(err, "NewParser queries=%+v", ctx.Queries())
		}

		records := make([]envv1.DbSchemaEnvironment, 0)
		if err := ctx.Database().Prepared(preparer).Find(&records); err != nil {
			return nil, errors.Wrapf(err, "Database Find")
		}
		return envv1.TransToHttpRsp(envv1.TransFormDbSchema(records)), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Contexter) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "FindEnvironment binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "FindEnvironment operator")
			}
			return v, nil
		},
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
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

		rst, err := operator.NewEnvironment(ctx.Database()).
			Get(uuid)
		if err != nil {
			return nil, err
		}
		return envv1.HttpRspEnvironment{Environment: *rst}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
	})
}

// UpdateEnvironment
// @Description Update Environment Value
// @Accept      json
// @Produce     json
// @Tags        server/environment
// @Router      /server/environment/{uuid} [put]
// @Param       uuid       path string                      true  "Environment 의 Uuid"
// @Param       enviroment body v1.HttpReqUpdateEnvironment false "HttpReqUpdateEnvironment"
// @Success 200 {object} v1.HttpRspEnvironment
func (c *Control) UpdateEnvironmentValue() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {
		body := new(envv1.HttpReqUpdateEnvironment)
		if err := ctx.Bind(body); err != nil {
			return ErrorBindRequestObject(err)
		}

		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}

		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}

		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		body, ok := ctx.Object().(*envv1.HttpReqUpdateEnvironment)
		if !ok {
			return nil, ErrorFailedCast()
		}
		update_env := body.UpdateEnvironment

		uuid := ctx.Params()[__UUID__]

		//get record
		env, err := operator.NewEnvironment(ctx.Database()).Get(uuid)
		if err != nil {
			return nil, err
		}

		//update property
		//value
		env.Value = nullable.String(*update_env.Value).Ptr()

		//update record
		err = operator.NewEnvironment(ctx.Database()).
			Update(*env)
		if err != nil {
			return nil, err
		}

		return envv1.HttpRspEnvironment{Environment: *env}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
	})
}

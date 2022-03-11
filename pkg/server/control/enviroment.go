package control

import (
	"github.com/NexClipper/sudory/pkg/server/control/vault"
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
		records, err := vault.NewEnvironment(ctx.Database()).Query(ctx.Queries())
		if err != nil {
			return nil, errors.Wrapf(err, "NewEnvironment Query")
		}

		return envv1.TransToHttpRsp(records), nil
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
		HttpResponsor: HttpJsonResponsor,
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

		rst, err := vault.NewEnvironment(ctx.Database()).Get(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "NewEnvironment Get")
		}
		return envv1.HttpRspEnvironment{DbSchema: *rst}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Contexter) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "GetEnvironment binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "GetEnvironment operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
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
// @Param       enviroment body v1.HttpReqEnvironmentUpdate false "HttpReqEnvironmentUpdate"
// @Success 200 {object} v1.HttpRspEnvironment
func (c *Control) UpdateEnvironmentValue() func(ctx echo.Context) error {
	binder := func(ctx Contexter) error {
		body := new(envv1.HttpReqEnvironmentUpdate)
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
		body, ok := ctx.Object().(*envv1.HttpReqEnvironmentUpdate)
		if !ok {
			return nil, ErrorFailedCast()
		}
		update_env := body.EnvironmentUpdate

		uuid := ctx.Params()[__UUID__]

		//get record
		env, err := vault.NewEnvironment(ctx.Database()).Get(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "NewEnvironment Get")
		}

		//update property
		//value
		env.Value = nullable.String(update_env.Value).Ptr()

		//update record
		env_, err := vault.NewEnvironment(ctx.Database()).Update(env.Environment)
		if err != nil {
			return nil, errors.Wrapf(err, "NewEnvironment Update")
		}

		return envv1.HttpRspEnvironment{DbSchema: *env_}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Contexter) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "UpdateEnvironmentValue binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "UpdateEnvironmentValue operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Nolock(c.db.Engine()),
	})
}

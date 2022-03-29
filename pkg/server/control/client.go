package control

import (
	"github.com/NexClipper/sudory/pkg/server/control/vault"
	clientv1 "github.com/NexClipper/sudory/pkg/server/model/client/v1"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Find Client
// @Description Find client
// @Accept      json
// @Produce     json
// @Tags        server/client
// @Router      /server/client [get]
// @Param       q query string false "query  pkg/server/database/prepared/README.md"
// @Param       o query string false "order  pkg/server/database/prepared/README.md"
// @Param       p query string false "paging pkg/server/database/prepared/README.md"
// @Success 200 {array} v1.HttpRspClient
func (c *Control) FindClient() func(ctx echo.Context) error {

	binder := func(ctx Context) error {
		return nil
	}
	operator := func(ctx Context) (interface{}, error) {
		records, err := vault.NewClient(ctx.Database()).Query(ctx.Queries())
		if err != nil {
			return nil, errors.Wrapf(err, "NewClient Query")
		}
		return clientv1.TransToHttpRsp(records), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Context) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "FindClient binder")
			}
			return nil
		},
		Operator: func(ctx Context) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "FindClient operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Nolock(c.db.Engine()),
	})
}

// Get Client
// @Description Get a client
// @Accept      json
// @Produce     json
// @Tags        server/client
// @Router      /server/client/{uuid} [get]
// @Param       uuid          path string true "Client 의 Uuid"
// @Success 200 {object} v1.HttpRspClient
func (c *Control) GetClient() func(ctx echo.Context) error {

	binder := func(ctx Context) error {
		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}

		return nil
	}
	operator := func(ctx Context) (interface{}, error) {
		uuid := ctx.Params()[__UUID__]

		rst, err := vault.NewClient(ctx.Database()).Get(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "NewClient Get")
		}
		return clientv1.HttpRspClient{DbSchema: *rst}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Context) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "GetClient binder")
			}
			return nil
		},
		Operator: func(ctx Context) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "GetClient operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Nolock(c.db.Engine()),
	})
}

// Delete Client
// @Description Delete a client
// @Accept      json
// @Produce     json
// @Tags        server/client
// @Router      /server/client/{uuid} [delete]
// @Param       uuid path string true "Client 의 Uuid"
// @Success 200
func (c *Control) DeleteClient() func(ctx echo.Context) error {

	binder := func(ctx Context) error {
		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}

		return nil
	}
	operator := func(ctx Context) (interface{}, error) {
		uuid := ctx.Params()[__UUID__]

		if err := vault.NewClient(ctx.Database()).Delete(uuid); err != nil {
			return nil, errors.Wrapf(err, "NewClient Delete")
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Context) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "GetClient binder")
			}
			return nil
		},
		Operator: func(ctx Context) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "GetClient operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Lock(c.db.Engine()),
	})
}

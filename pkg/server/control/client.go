package control

import (
	"fmt"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
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
// @Success 200 {object} v1.HttpRspClient
func (c *Control) CreateClient() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {
		body := new(clinetv1.HttpReqClient)
		if err := ctx.Bind(body); err != nil {
			return ErrorBindRequestObject(err)
		}
		if body.Name == nil {
			return ErrorInvaliedRequestParameterName("Name")
		}
		if len(body.ClusterUuid) == 0 {
			return ErrorInvaliedRequestParameterName("ClusterUuid")
		}

		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		body, ok := ctx.Object().(*clinetv1.HttpReqClient)
		if !ok {
			return nil, ErrorFailedCast()
		}

		client := body.Client

		//property
		client.UuidMeta = NewUuidMeta()
		client.LabelMeta = NewLabelMeta(client.Name, client.Summary)

		err := operator.NewClient(ctx.Database()).
			Create(client)
		if err != nil {
			return nil, err
		}

		return clinetv1.HttpRspClient{Client: client}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Lock(c.db.Engine()),
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

	binder := func(ctx Contexter) error {
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		//make condition
		args := make([]interface{}, 0)
		add, build := StringBuilder()

		for key, val := range ctx.Querys() {
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
		rst, err := operator.NewClient(ctx.Database()).
			Find(where, args...)
		if err != nil {
			return nil, err
		}
		return clinetv1.TransToHttpRsp(rst), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
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

		rst, err := operator.NewClient(ctx.Database()).
			Get(uuid)
		if err != nil {
			return nil, err
		}
		return clinetv1.HttpRspClient{Client: *rst}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
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
// @Success 200 {object} v1.HttpRspClient
func (c *Control) UpdateClient() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {
		body := new(clinetv1.HttpReqClient)
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
		body, ok := ctx.Object().(*clinetv1.HttpReqClient)
		if !ok {
			return nil, ErrorFailedCast()
		}

		uuid := ctx.Params()[__UUID__]

		client := body.Client

		//set uuid from path
		client.Uuid = uuid

		err := operator.NewClient(ctx.Database()).
			Update(client)
		if err != nil {
			return nil, err
		}

		return clinetv1.HttpRspClient{Client: client}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Lock(c.db.Engine()),
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

		err := operator.NewClient(ctx.Database()).
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

package control

import (
	"fmt"
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepared"
	"github.com/NexClipper/sudory/pkg/server/macro"
	labelv1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	tokenv1 "github.com/NexClipper/sudory/pkg/server/model/token/v1"
	"github.com/NexClipper/sudory/pkg/server/status/env"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// token user kind
const (
	// cluster
	token_user_kind_cluster = "cluster"
)

// CreateClusterToken
// @Description Create a Cluster Token
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token/cluster [post]
// @Param       object body v1.HttpReqToken true "HttpReqToken"
// @Success     200 {object} v1.HttpRspToken
func (c *Control) CreateClusterToken() func(ctx echo.Context) error {

	const user_kind = token_user_kind_cluster

	binder := func(ctx Contexter) error {

		body := new(tokenv1.HttpReqToken)
		if err := ctx.Bind(body); err != nil {
			return ErrorBindRequestObject(err)
		}

		if body.Name == nil {
			return ErrorInvaliedRequestParameterName("Name")
		}
		if len(body.UserUuid) == 0 {
			return ErrorInvaliedRequestParameterName("UserUuid")
		}
		// if len(body.Token) == 0 {
		// 	return nil, ErrorInvaliedRequestParameterName("Token")
		// }
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		body, ok := ctx.Object().(*tokenv1.HttpReqToken)
		if !ok {
			return nil, ErrorFailedCast()
		}

		token := body.Token

		//vaild token user
		err := vaildTokenUser(ctx.Database(), user_kind, token.UserUuid)
		if err != nil {
			return nil, err
		}

		//property
		token.UuidMeta = NewUuidMeta()
		token.LabelMeta = NewLabelMeta(token.Name, token.Summary)
		token.UserKind = user_kind
		token.IssuedAtTime, token.ExpirationTime = bearerTokenTimeIssueNow()
		token.Token = macro.NewUuidString()

		//create
		err = operator.NewToken(ctx.Database()).
			Create(token)
		if err != nil {
			return nil, err
		}

		return tokenv1.HttpRspToken{Token: token}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Lock(c.db.Engine()),
	})
}

// FindToken
// @Description Find Token
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token [get]
// @Param       q query string false "query  pkg/server/database/prepared/README.md"
// @Param       o query string false "order  pkg/server/database/prepared/README.md"
// @Param       p query string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.HttpRspToken
func (c *Control) FindToken() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		preparer, err := prepared.NewParser(ctx.Queries())
		if err != nil {
			return nil, errors.Wrapf(err, "NewParser queries=%+v", ctx.Queries())
		}

		records := make([]tokenv1.DbSchemaToken, 0)
		if err := ctx.Database().Prepared(preparer).Find(&records); err != nil {
			return nil, errors.Wrapf(err, "Database Find")
		}

		return tokenv1.TransToHttpRsp(tokenv1.TransFormDbSchema(records)), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Contexter) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "FindToken binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "FindToken operator")
			}
			return v, nil
		},
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
	})
}

// GetToken
// @Description Get a Token
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token/{uuid} [get]
// @Param       uuid path string true "Token 의 Uuid"
// @Success     200 {object} v1.HttpRspToken
func (c *Control) GetToken() func(ctx echo.Context) error {

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

		rst, err := operator.NewToken(ctx.Database()).
			Get(uuid)
		if err != nil {
			return nil, err
		}
		return tokenv1.HttpRspToken{Token: *rst}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
	})
}

// UpdateTokenLabel
// @Description Update Token Label
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router		/server/token/{uuid}/label [put]
// @Param       uuid   path string true "Token 의 Uuid"
// @Param       object body v1.LabelMeta true "Token 의 LabelMeta"
// @Success 	200 {object} v1.HttpRspToken
func (c *Control) UpdateTokenLabel() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {

		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}

		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}

		body := new(labelv1.LabelMeta)
		if err := ctx.Bind(body); err != nil {
			return ErrorBindRequestObject(err)
		}

		// if body.Name == nil {
		// 	return ErrorInvaliedRequestParameterName("Name")
		// }

		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {

		body, ok := ctx.Object().(*labelv1.LabelMeta)
		if !ok {
			return nil, ErrorFailedCast()
		}
		uuid := ctx.Params()[__UUID__]

		label := body

		//get token
		token, err := operator.NewToken(ctx.Database()).
			Get(uuid)
		if err != nil {
			return nil, err
		}

		if token == nil {
			return nil, fmt.Errorf("record was not found: token")
		}

		//property
		token.Name = label.Name
		token.Summary = label.Summary

		//update record
		err = operator.NewToken(ctx.Database()).
			Update(*token)
		if err != nil {
			return nil, err
		}

		return tokenv1.HttpRspToken{Token: *token}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
	})
}

// RefreshClusterTokenTime
// @Description Refresh Cluster Token Time
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token/cluster/{uuid}/refresh [put]
// @Param       uuid    path string true "Token 의 Uuid"
// @Success     200 {object} v1.HttpRspToken
func (c *Control) RefreshClusterTokenTime() func(ctx echo.Context) error {

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

		token, err := operator.NewToken(ctx.Database()).
			Get(uuid)
		if err != nil {
			return nil, err
		}

		//property
		token.IssuedAtTime, token.ExpirationTime = bearerTokenTimeIssueNow() //만료시간 연장

		err = operator.NewToken(ctx.Database()).
			Update(*token)
		if err != nil {
			return nil, err
		}

		return tokenv1.HttpRspToken{Token: *token}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
	})
}

// ExpireClusterToken
// @Description Expire Cluster Token
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token/cluster/{uuid}/expire [put]
// @Param       uuid    path string true "Token 의 Uuid"
// @Success     200 {object} v1.HttpRspToken
func (c *Control) ExpireClusterToken() func(ctx echo.Context) error {

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

		token, err := operator.NewToken(ctx.Database()).
			Get(uuid)
		if err != nil {
			return nil, err
		}

		//property
		token.IssuedAtTime, token.ExpirationTime = time.Now(), time.Now() //현재 시간으로 만료시간 설정

		err = operator.NewToken(ctx.Database()).
			Update(*token)
		if err != nil {
			return nil, err
		}

		return tokenv1.HttpRspToken{Token: *token}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      operator,
		HttpResponser: HttpResponse,
		Behavior:      Nolock(c.db.Engine()),
	})
}

// DeleteToken
// @Description Delete a Token
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token/{uuid} [delete]
// @Param       uuid path string true "Token 의 Uuid"
// @Success     200
func (c *Control) DeleteToken() func(ctx echo.Context) error {

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

		err := operator.NewToken(ctx.Database()).
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
		Behavior:      Nolock(c.db.Engine()),
	})
}

// vaildTokenUser
func vaildTokenUser(ctx database.Context, user_kind, user_uuid string) error {
	switch user_kind {
	case token_user_kind_cluster:
		//get cluster
		_, err := operator.NewCluster(ctx).
			Get(user_uuid)
		if err != nil {
			return fmt.Errorf("record was not found cluster: %w", err) //can't exist
		}
	default:
		return fmt.Errorf("invalid user_kind")
	}

	return nil
}

func bearerTokenTimeIssueNow() (time.Time, time.Time) {
	iat := time.Now()
	exp := env.BearerTokenExpirationTime(iat)
	return iat, exp
}

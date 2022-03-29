package control

import (
	"fmt"
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro"
	labelv1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	tokenv1 "github.com/NexClipper/sudory/pkg/server/model/token/v1"
	"github.com/NexClipper/sudory/pkg/server/status/env"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

const token_user_kind_cluster = "cluster"

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
	binder := func(ctx Context) error {

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
	operator := func(ctx Context) (interface{}, error) {
		body, ok := ctx.Object().(*tokenv1.HttpReqToken)
		if !ok {
			return nil, ErrorFailedCast()
		}

		token := body.Token

		//vaild token user
		if err := vaildTokenUser(ctx.Database(), user_kind, token.UserUuid); err != nil {
			return nil, errors.Wrapf(err, "vaildTokenUser CreateClusterToken user_kind=%s user_uuid=%s", user_kind, token.UserUuid)
		}

		//property
		token.UuidMeta = NewUuidMeta()
		token.LabelMeta = NewLabelMeta(token.Name, token.Summary)
		token.UserKind = user_kind
		token.IssuedAtTime, token.ExpirationTime = bearerTokenTimeIssueNow()
		token.Token = macro.NewUuidString()

		new_token, err := vault.NewToken(ctx.Database()).CreateClusterToken(token)
		if err != nil {
			return nil, errors.Wrapf(err, "NewToken CreateClusterToken")
		}

		return tokenv1.HttpRspToken{DbSchema: *new_token}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Context) error {
			if err := binder(ctx); err != nil {
				return errors.Wrapf(err, "CreateClusterToken binder")
			}
			return nil
		},
		Operator: func(ctx Context) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "CreateClusterToken operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
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

	binder := func(ctx Context) error {
		return nil
	}
	operator := func(ctx Context) (interface{}, error) {
		records, err := vault.NewToken(ctx.Database()).Query(ctx.Queries())
		if err != nil {
			return nil, errors.Wrapf(err, "NewToken Query")
		}

		return tokenv1.TransToHttpRsp(records), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Context) error {
			if err := binder(ctx); err != nil {
				return errors.Wrapf(err, "FindToken binder")
			}
			return nil
		},
		Operator: func(ctx Context) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "FindToken operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
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

		rst, err := vault.NewToken(ctx.Database()).Get(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "NewToken Get")
		}
		return tokenv1.HttpRspToken{DbSchema: *rst}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Context) error {
			if err := binder(ctx); err != nil {
				return errors.Wrapf(err, "GetToken binder")
			}
			return nil
		},
		Operator: func(ctx Context) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "GetToken operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
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
	binder := func(ctx Context) error {

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
	operator := func(ctx Context) (interface{}, error) {

		body, ok := ctx.Object().(*labelv1.LabelMeta)
		if !ok {
			return nil, ErrorFailedCast()
		}
		uuid := ctx.Params()[__UUID__]

		label := body

		//get token
		token, err := vault.NewToken(ctx.Database()).Get(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "NewToken Get")
		}

		//property
		token.Name = label.Name
		token.Summary = label.Summary

		//update record
		token_, err := vault.NewToken(ctx.Database()).Update(token.Token)
		if err != nil {
			return nil, errors.Wrapf(err, "NewToken Update")
		}

		return tokenv1.HttpRspToken{DbSchema: *token_}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Context) error {
			if err := binder(ctx); err != nil {
				return errors.Wrapf(err, "UpdateTokenLabel binder")
			}
			return nil
		},
		Operator: func(ctx Context) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "UpdateTokenLabel operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
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

		token, err := vault.NewToken(ctx.Database()).Get(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "NewToken Get")
		}

		//property
		token.IssuedAtTime, token.ExpirationTime = bearerTokenTimeIssueNow() //만료시간 연장
		token_, err := vault.NewToken(ctx.Database()).Update(token.Token)
		if err != nil {
			return nil, errors.Wrapf(err, "NewToken Update")
		}

		return tokenv1.HttpRspToken{DbSchema: *token_}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Context) error {
			if err := binder(ctx); err != nil {
				return errors.Wrapf(err, "RefreshClusterTokenTime binder")
			}
			return nil
		},
		Operator: func(ctx Context) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "RefreshClusterTokenTime operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
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

		token, err := vault.NewToken(ctx.Database()).Get(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "NewToken Get")
		}

		//property
		token.IssuedAtTime, token.ExpirationTime = time.Now(), time.Now() //현재 시간으로 만료시간 설정
		token_, err := vault.NewToken(ctx.Database()).Update(token.Token)
		if err != nil {
			return nil, errors.Wrapf(err, "NewToken Update")
		}

		return tokenv1.HttpRspToken{DbSchema: *token_}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Context) error {
			if err := binder(ctx); err != nil {
				return errors.Wrapf(err, "ExpireClusterToken binder")
			}
			return nil
		},
		Operator: func(ctx Context) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "ExpireClusterToken operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
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

		if err := vault.NewToken(ctx.Database()).Delete(uuid); err != nil {
			return nil, errors.Wrapf(err, "NewToken Delete")
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Context) error {
			if err := binder(ctx); err != nil {
				return errors.Wrapf(err, "DeleteToken binder")
			}
			return nil
		},
		Operator: func(ctx Context) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "DeleteToken operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Nolock(c.db.Engine()),
	})
}

// vaildTokenUser
func vaildTokenUser(ctx database.Context, user_kind, user_uuid string) error {
	switch user_kind {
	case token_user_kind_cluster:
		//get cluster
		if _, err := vault.NewCluster(ctx).Get(user_uuid); err != nil {
			return errors.Wrapf(err, "NewCluster Get") //can't exist
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

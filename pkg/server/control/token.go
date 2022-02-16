package control

import (
	"fmt"
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/query_parser"
	"github.com/NexClipper/sudory/pkg/server/macro"
	labelv1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	tokenv1 "github.com/NexClipper/sudory/pkg/server/model/token/v1"
	"github.com/NexClipper/sudory/pkg/server/status/env"
	"github.com/labstack/echo/v4"
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

	binder := func(ctx echo.Context) (interface{}, error) {
		body := new(tokenv1.HttpReqToken)
		err := ctx.Bind(body)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}
		if len(body.Token.Name) == 0 {
			return nil, ErrorInvaliedRequestParameterName("Token.Name")
		}
		if len(body.Token.UserUuid) == 0 {
			return nil, ErrorInvaliedRequestParameterName("Token.UserUuid")
		}
		if len(body.Token.Token) == 0 {
			return nil, ErrorInvaliedRequestParameterName("Token.Token")
		}
		return body, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		body, ok := ctx.Req.(*tokenv1.HttpReqToken)
		if !ok {
			return nil, ErrorFailedCast()
		}

		token := body.Token

		//vaild
		err := vaildTokenUser(ctx.Database, token_user_kind_cluster, token.UserUuid)
		if err != nil {
			return nil, err
		}

		//property
		token.LabelMeta = NewLabelMeta(token.Name, token.Summary)
		token.UserKind = token_user_kind_cluster
		token.IssuedAtTime, token.ExpirationTime = BearerTokenTimeIssueNow()

		//create
		err = operator.NewToken(ctx.Database).
			Create(token)
		if err != nil {
			return nil, err
		}

		return tokenv1.HttpRspToken{Token: token}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Lock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// FindToken
// @Description Find Token
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token [get]
// @Param       uuid         query string false "Token 의 Uuid"
// @Param       name         query string false "Token 의 Name"
// @Param       user_kind    query string false "Token 의 user_kind"
// @Param       user_uuid    query string false "Token 의 user_uuid"
// @Param       token        query string false "Token 의 token"
// #Param       limit        query int    false "Pagination 의 limit"
// #Param       page         query int    false "Pagination 의 page"
// #Param       order        query string false "Pagination 의 order"
// @Success     200 {array} v1.HttpRspToken
func (c *Control) FindToken() func(ctx echo.Context) error {

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

		//make condition
		cond := query_parser.NewCondition(req, func(key string) (string, string, bool) {
			switch key {
			case "uuid", "user_uuid", "user_kind", "token":
				return "=", "%s", true
			case "name":
				return "LIKE", "%%%s%%", true
			default:
				return "", "", false
			}
		})

		//find Token
		rst, err := operator.NewToken(ctx.Database).
			Find(cond.Where(), cond.Args()...)
		if err != nil {
			return nil, err
		}
		return tokenv1.TransToHttpRsp(rst), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Nolock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
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

		rst, err := operator.NewToken(ctx.Database).
			Get(uuid)
		if err != nil {
			return nil, err
		}
		return tokenv1.HttpRspToken{Token: *rst}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Nolock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
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

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}

		body := new(labelv1.LabelMeta)
		err := ctx.Bind(body)
		if err != nil {
			return nil, ErrorBindRequestObject(err)
		}
		req[__BODY__] = body

		if _, ok := macro.MapString(req, __UUID__); !ok {
			return nil, ErrorInvaliedRequestParameterName(__UUID__)
		}

		if len(body.Name) == 0 {
			return nil, ErrorInvaliedRequestParameterName("LabelMeta.Name")
		}

		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}

		uuid, _ := macro.MapString(req, __UUID__)

		label, _ := req[__BODY__].(*labelv1.LabelMeta)
		if label == nil {
			return nil, ErrorFailedCast()
		}

		//get token
		token, err := operator.NewToken(ctx.Database).
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
		err = operator.NewToken(ctx.Database).
			Update(*token)
		if err != nil {
			return nil, err
		}

		return tokenv1.HttpRspToken{Token: *token}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Lock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
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

		token, err := operator.NewToken(ctx.Database).
			Get(uuid)
		if err != nil {
			return nil, err
		}

		//property
		//만료시간 연장
		token.IssuedAtTime, token.ExpirationTime = BearerTokenTimeIssueNow()

		err = operator.NewToken(ctx.Database).
			Update(*token)
		if err != nil {
			return nil, err
		}

		return tokenv1.HttpRspToken{Token: *token}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Lock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
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

		token, err := operator.NewToken(ctx.Database).
			Get(uuid)
		if err != nil {
			return nil, err
		}

		//property
		//현재 시간으로 만료시간 설정
		token.IssuedAtTime, token.ExpirationTime = time.Now(), time.Now()

		err = operator.NewToken(ctx.Database).
			Update(*token)
		if err != nil {
			return nil, err
		}

		return tokenv1.HttpRspToken{Token: *token}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder:        binder,
		Operator:      Lock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
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

		err := operator.NewToken(ctx.Database).
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

func BearerTokenTimeIssueNow() (time.Time, time.Time) {
	iat := time.Now()
	exp := env.BearerTokenExpirationTime(iat)
	return iat, exp
}

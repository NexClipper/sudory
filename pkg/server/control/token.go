package control

import (
	"fmt"
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/operator"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/query_parser"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
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
// @Accept x-www-form-urlencoded
// @Produce json
// @Tags server/token
// @Router /server/token/cluster [post]
// @Param name         formData string true  "Token 의 Name"
// @Param summary      formData string false "Token 의 Summary"
// @Param token        formData string true  "Token 의 Token"
// @Param user_uuid    formData string true  "Token 의 UserUuid"
// @Success 200 {object} v1.HttpRspToken
func (c *Control) CreateClusterToken() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		// for _, it := range ctx.ParamNames() {
		// 	req[it] = ctx.Param(it)
		// }
		// for key := range ctx.QueryParams() {
		// 	req[key] = ctx.QueryParam(key)
		// }
		formdatas, err := ctx.FormParams()
		if err != nil {
			return nil, err
		}
		for key := range formdatas {
			req[key] = ctx.FormValue(key)
		}
		if !macro.MapFound(req, __NAME__) {
			return nil, ErrorInvaliedRequestParameter()
		}

		//lint:ignore SA9003 auto-generated
		if !macro.MapFound(req, __SUMMARY__) {
			//(optional)
		}

		if !macro.MapFound(req, __TOKEN__) {
			return nil, ErrorInvaliedRequestParameter()
		}

		if !macro.MapFound(req, __USER_UUID__) {
			return nil, ErrorInvaliedRequestParameter()
		}

		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}

		//new uuid
		user_kind := token_user_kind_cluster

		name, ok := macro.MapString(req, __NAME__)
		if !ok {
			return nil, ErrorFailedCast()
		}
		summary, _ := macro.MapString(req, __SUMMARY__)

		token, ok := macro.MapString(req, __TOKEN__)
		if !ok {
			return nil, ErrorFailedCast()
		}
		user_uuid, ok := macro.MapString(req, __USER_UUID__)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//vaild token_user
		err := vaildTokenUser(ctx.Database, user_kind, user_uuid)
		if err != nil {
			return nil, err
		}

		new_token := tokenv1.Token{}
		new_token.LabelMeta = NewLabelMeta(name, summary)
		new_token.UserKind = user_kind
		new_token.UserUuid = user_uuid
		new_token.Token = token

		//property
		new_token = tokenExpirationTime(new_token)

		//create
		err = operator.NewToken(ctx.Database).
			Create(new_token)
		if err != nil {
			return nil, err
		}

		return tokenv1.HttpRspToken{Token: new_token}, nil
	}

	return MakeMiddlewareFunc_experimental(Option_experimental{
		Binder:        binder,
		Operator:      MakeBlockWithLock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// FindToken
// @Description Find Token
// @Accept json
// @Produce json
// @Tags server/token
// @Router /server/token [get]
// @Param uuid         query string false "Token 의 Uuid"
// @Param name         query string false "Token 의 Name"
// @Param user_kind    query string false "Token 의 user_kind"
// @Param user_uuid    query string false "Token 의 user_uuid"
// @Param token        query string false "Token 의 token"
// @Param limit        query int    false "Pagination 의 limit"
// @Param page         query int    false "Pagination 의 page"
// @Param order        query string false "Pagination 의 order"
// @Success 200 {array} v1.HttpRspToken
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
		query := query_parser.NewQueryParser(req, func(key string) (string, string) {
			switch key {
			case "uuid", "user_uuid", "user_kind", "token":
				return "=", "%s"
			default:
				return "LIKE", "%%%s%%"
			}
		})

		//find Token
		rst, err := operator.NewToken(ctx.Database).
			Query(query)
		if err != nil {
			return nil, err
		}
		return tokenv1.TransToHttpRsp(rst), nil
	}

	return MakeMiddlewareFunc_experimental(Option_experimental{
		Binder:        binder,
		Operator:      MakeBlockNoLock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// GetToken
// @Description Get a Token
// @Accept json
// @Produce json
// @Tags server/token
// @Router /server/token/{uuid} [get]
// @Param uuid          path string true "Token 의 Uuid"
// @Success 200 {object} v1.HttpRspToken
func (c *Control) GetToken() func(ctx echo.Context) error {

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
		rst, err := operator.NewToken(ctx.Database).
			Get(uuid)
		if err != nil {
			return nil, err
		}
		return tokenv1.HttpRspToken{Token: *rst}, nil
	}

	return MakeMiddlewareFunc_experimental(Option_experimental{
		Binder:        binder,
		Operator:      MakeBlockNoLock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// UpdateClusterToken
// @Description Update Cluster Token
// @Accept x-www-form-urlencoded
// @Produce json
// @Tags server/token
// @Router /server/token/cluster/{uuid} [put]
// @Param uuid         path     string true  "Token 의 Uuid"
// @Param name         formData string false "Token 의 Name"
// @Param summary      formData string false "Token 의 Summary"
// @Success 200 {object} v1.HttpRspToken
func (c *Control) UpdateClusterToken() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		for key := range ctx.QueryParams() {
			req[key] = ctx.QueryParam(key)
		}
		formdatas, err := ctx.FormParams()
		if err != nil {
			return nil, err
		}
		for key := range formdatas {
			req[key] = ctx.FormValue(key)
		}

		if !macro.MapFound(req, __UUID__) {
			return nil, ErrorInvaliedRequestParameter()
		}

		//lint:ignore SA9003 auto-generated
		if !macro.MapFound(req, __NAME__) {
			//(optional)
		}

		//lint:ignore SA9003 auto-generated
		if !macro.MapFound(req, __SUMMARY__) {
			//(optional)
		}

		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}
		user_kind := token_user_kind_cluster

		uuid, ok := macro.MapString(req, __UUID__)
		if !ok {
			return nil, ErrorFailedCast()
		}

		//get token
		where := `uuid = ? AND user_kind = ?`
		tokens, err := operator.NewToken(ctx.Database).
			Find(where, uuid, user_kind)
		if err != nil {
			return nil, err
		}

		first := func() *tokenv1.Token {
			for _, it := range tokens {
				return &it
			}
			return nil
		}
		token := first()

		if token == nil {
			return nil, fmt.Errorf("record was not found: token")
		}

		//vaild token_user
		err = vaildTokenUser(ctx.Database, token.UserKind, token.UserUuid)
		if err != nil {
			return nil, err
		}

		//set uuid from path
		token.Uuid = uuid

		//update value

		if name, ok := macro.MapString(req, __NAME__); ok {
			token.Name = name
		}

		if summary, ok := macro.MapString(req, __SUMMARY__); ok {
			token.Summary = newist.String(summary)
		}

		//update record
		err = operator.NewToken(ctx.Database).
			Update(*token)
		if err != nil {
			return nil, err
		}

		return tokenv1.HttpRspToken{Token: *token}, nil
	}

	return MakeMiddlewareFunc_experimental(Option_experimental{
		Binder:        binder,
		Operator:      MakeBlockWithLock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// RefreshClusterTokenExpirationTime
// @Description Refresh Cluster Token Expiration Time
// @Accept json
// @Produce json
// @Tags server/token
// @Router /server/token/cluster/{uuid}/exp [put]
// @Param uuid    path string true "Token 의 Uuid"
// @Success 200 {object} v1.HttpRspToken
func (c *Control) RefreshClusterTokenExpirationTime() func(ctx echo.Context) error {

	binder := func(ctx echo.Context) (interface{}, error) {
		req := make(map[string]interface{})
		for _, it := range ctx.ParamNames() {
			req[it] = ctx.Param(it)
		}
		if _, ok := macro.MapString(req, __UUID__); !ok {
			return nil, ErrorInvaliedRequestParameter()
		}

		return req, nil
	}
	operator := func(ctx OperateContext) (interface{}, error) {
		req, ok := ctx.Req.(map[string]interface{})
		if !ok {
			return nil, ErrorFailedCast()
		}
		uuid, ok := macro.MapString(req, __UUID__)
		if !ok {
			return nil, ErrorFailedCast()
		}

		token, err := operator.NewToken(ctx.Database).
			Get(uuid)
		if err != nil {
			return nil, err
		}

		//vaild token_user
		err = vaildTokenUser(ctx.Database, token.UserKind, token.UserUuid)
		if err != nil {
			return nil, err
		}

		//property
		//만료시간 연장
		*token = tokenExpirationTime(*token)

		err = operator.NewToken(ctx.Database).
			Update(*token)
		if err != nil {
			return nil, err
		}

		return tokenv1.HttpRspToken{Token: *token}, nil
	}

	return MakeMiddlewareFunc_experimental(Option_experimental{
		Binder:        binder,
		Operator:      MakeBlockWithLock(c.db.Engine(), operator),
		HttpResponser: HttpResponse,
	})
}

// DeleteToken
// @Description Delete a Token
// @Accept json
// @Produce json
// @Tags server/token
// @Router /server/token/{uuid} [delete]
// @Param uuid path string true "Token 의 Uuid"
// @Success 200
func (c *Control) DeleteToken() func(ctx echo.Context) error {

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
		err := operator.NewToken(ctx.Database).
			Delete(uuid)
		if err != nil {
			return nil, err
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc_experimental(Option_experimental{
		Binder:        binder,
		Operator:      MakeBlockWithLock(c.db.Engine(), operator),
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

func tokenExpirationTime(token tokenv1.Token) tokenv1.Token {
	iat := time.Now()
	exp := env.BearerTokenExpirationTime(iat)

	token.IssuedAtTime = iat
	token.ExpirationTime = exp

	return token
}

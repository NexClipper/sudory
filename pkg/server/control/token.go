//go:generate go-enum --file=token.go --names --nocase=true
package control

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
	labelv1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	tokenv1 "github.com/NexClipper/sudory/pkg/server/model/token/v1"
	"github.com/NexClipper/sudory/pkg/server/status/env"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

/* ENUM(
cluster
)
*/
type TokenUserKind int32

// CreateClusterToken
// @Description Create a Cluster Token
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token/cluster [post]
// @Param       object body v1.HttpReqToken true "HttpReqToken"
// @Success     200 {object} v1.HttpRspToken
func (ctl Control) CreateClusterToken(ctx echo.Context) error {
	const user_kind = TokenUserKindCluster

	body := new(tokenv1.HttpReqToken)
	if err := ctx.Bind(body); err != nil {
		return HttpError(http.StatusBadRequest,
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if len(nullable.String(body.Name).Value()) == 0 {
		return HttpError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					"param", TypeName(body.Name),
				)))
	}

	if len(nullable.String(body.UserUuid).Value()) == 0 {
		return HttpError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					"param", TypeName(body.UserUuid),
				)))
	}

	token := body.Token

	//valvalidied token user
	if err := validTokenUser(ctl.NewSession(), user_kind, *token.UserUuid); err != nil {
		return HttpError(http.StatusInternalServerError,
			errors.Wrapf(err, "valid token user%s",
				logs.KVL(
					"user_kind", user_kind,
					"user_uuid", token.UserUuid,
				)))
	}

	//property
	token.UuidMeta = NewUuidMeta()
	token.LabelMeta = NewLabelMeta(token.Name, token.Summary)
	token.UserKind = newist.String(user_kind.String())
	token.IssuedAtTime, token.ExpirationTime = bearerTokenTimeIssueNow()
	token.Token = newist.String(macro.NewUuidString())

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		token_, err := vault.NewToken(db).CreateToken(token)
		if err != nil {
			return nil, errors.Wrapf(err, "create cluster token")
		}

		return tokenv1.HttpRspToken{DbSchema: *token_}, nil
	})
	if err != nil {
		return HttpError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, r)
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
func (ctl Control) FindToken(ctx echo.Context) error {
	r, err := vault.NewToken(ctl.NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return HttpError(http.StatusInternalServerError,
			errors.Wrapf(err, "find token%s",
				logs.KVL(
					"query", echoutil.QueryParamString(ctx),
				)))
	}

	return ctx.JSON(http.StatusOK, tokenv1.TransToHttpRsp(r))
}

// GetToken
// @Description Get a Token
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token/{uuid} [get]
// @Param       uuid path string true "Token 의 Uuid"
// @Success     200 {object} v1.HttpRspToken
func (ctl Control) GetToken(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return HttpError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := vault.NewToken(ctl.NewSession()).Get(uuid)
	if err != nil {
		return HttpError(http.StatusInternalServerError,
			errors.Wrapf(err, "get token%s",
				logs.KVL(
					"uuid", uuid,
				)))
	}

	return ctx.JSON(http.StatusOK, tokenv1.HttpRspToken{DbSchema: *r})
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
func (ctl Control) UpdateTokenLabel(ctx echo.Context) error {
	body := new(labelv1.LabelMeta)
	if err := ctx.Bind(body); err != nil {
		return HttpError(http.StatusBadRequest,
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return HttpError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	label := body
	uuid := echoutil.Param(ctx)[__UUID__]

	// //get token
	// token, err := vault.NewToken(ctl.NewSession()).Get(uuid)
	// if err != nil {
	// 	return HttpError(http.StatusInternalServerError,
	// 		errors.Wrapf(err, "get token"))
	// }

	//property
	token := tokenv1.DbSchema{}
	token.Uuid = uuid
	token.Name = label.Name
	token.Summary = label.Summary

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		//update record
		token_, err := vault.NewToken(db).Update(token.Token)
		if err != nil {
			return nil, errors.Wrapf(err, "update token")
		}

		return tokenv1.HttpRspToken{DbSchema: *token_}, nil
	})
	if err != nil {
		return HttpError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, r)
}

// RefreshClusterTokenTime
// @Description Refresh Cluster Token Time
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token/cluster/{uuid}/refresh [put]
// @Param       uuid    path string true "Token 의 Uuid"
// @Success     200 {object} v1.HttpRspToken
func (ctl Control) RefreshClusterTokenTime(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return HttpError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	// token, err := vault.NewToken(ctl.NewSession()).Get(uuid)
	// if err != nil {
	// 	return HttpError(http.StatusInternalServerError,
	// 		errors.Wrapf(err, "get token"))
	// }

	//property
	token := tokenv1.DbSchema{}
	token.Uuid = uuid
	token.IssuedAtTime, token.ExpirationTime = bearerTokenTimeIssueNow() //만료시간 연장

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		token_, err := vault.NewToken(db).Update(token.Token)
		if err != nil {
			return nil, errors.Wrapf(err, "update token")
		}

		return tokenv1.HttpRspToken{DbSchema: *token_}, nil
	})
	if err != nil {
		return HttpError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, r)
}

// ExpireClusterToken
// @Description Expire Cluster Token
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token/cluster/{uuid}/expire [put]
// @Param       uuid    path string true "Token 의 Uuid"
// @Success     200 {object} v1.HttpRspToken
func (ctl Control) ExpireClusterToken(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return HttpError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	// token, err := vault.NewToken(ctl.NewSession()).Get(uuid)
	// if err != nil {
	// 	return HttpError(http.StatusInternalServerError,
	// 		errors.Wrapf(err, "get token"))
	// }

	//property
	token := tokenv1.DbSchema{}
	token.Uuid = uuid
	token.IssuedAtTime, token.ExpirationTime = newist.Time(time.Now()), newist.Time(time.Now()) //현재 시간으로 만료시간 설정

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		token_, err := vault.NewToken(db).Update(token.Token)
		if err != nil {
			return nil, errors.Wrapf(err, "update token")
		}

		return tokenv1.HttpRspToken{DbSchema: *token_}, nil
	})
	if err != nil {
		return HttpError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, r)
}

// DeleteToken
// @Description Delete a Token
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token/{uuid} [delete]
// @Param       uuid path string true "Token 의 Uuid"
// @Success     200
func (ctl Control) DeleteToken(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return HttpError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		err := vault.NewToken(db).Delete(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "delete token%s",
				logs.KVL(
					"uuid", uuid,
				))
		}

		return nil, nil
	})
	if err != nil {
		return HttpError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, OK())
}

// validTokenUser
func validTokenUser(ctx database.Context, user_kind TokenUserKind, user_uuid string) error {
	switch user_kind {
	case TokenUserKindCluster:
		//get cluster
		if _, err := vault.NewCluster(ctx).Get(user_uuid); err != nil {
			return errors.Wrapf(err, "found cluster token user")
		}
	default:
		return fmt.Errorf("invalid token user kind")
	}

	return nil
}

func bearerTokenTimeIssueNow() (*time.Time, *time.Time) {
	iat := time.Now()
	exp := env.BearerTokenExpirationTime(iat)
	return newist.Time(iat), newist.Time(exp)
}

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
	cryptov1 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v1"
	tokenv1 "github.com/NexClipper/sudory/pkg/server/model/token/v1"
	"github.com/NexClipper/sudory/pkg/server/status/env"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// CreateClusterToken
// @Description Create a Cluster Token
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token/cluster [post]
// @Param       x_auth_token header string                             false "client session token"
// @Param       object       body   v1.HttpReqToken_CreateClusterToken true  "HttpReqToken_CreateClusterToken"
// @Success     200 {object} v1.Token
func (ctl Control) CreateClusterToken(ctx echo.Context) error {
	const user_kind = tokenv1.TokenUserKindCluster

	body := new(tokenv1.HttpReqToken_CreateClusterToken)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if len(body.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.Name", TypeName(body)), body.Name)...,
				)))
	}

	if len(body.UserUuid) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.UserUuid", TypeName(body)), body.UserUuid)...,
				)))
	}

	//valvalidied token user
	if err := validTokenUser(ctl.NewSession(), user_kind, body.UserUuid); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "valid token user%s",
				logs.KVL(
					"user_kind", user_kind,
					"user_uuid", body.UserUuid,
				)))
	}

	//property
	token := tokenv1.Token{}
	token.UuidMeta = NewUuidMeta()
	token.LabelMeta = NewLabelMeta(body.Name, body.Summary)
	token.UserKind = user_kind.String()
	token.UserUuid = body.UserUuid
	token.IssuedAtTime, token.ExpirationTime = bearerTokenTimeIssueNow()
	token.Token = cryptov1.String(macro.NewUuidString())

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		token_, err := vault.NewToken(db).CreateToken(token)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "create cluster token"))
		}

		return token_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// FindToken
// @Description Find Token
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.Token
func (ctl Control) FindToken(ctx echo.Context) error {
	r, err := vault.NewToken(ctl.NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "find token"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// GetToken
// @Description Get a Token
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Token 의 Uuid"
// @Success     200 {object} v1.Token
func (ctl Control) GetToken(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := vault.NewToken(ctl.NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get token"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// UpdateTokenLabel
// @Description Update Token Label
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router		/server/token/{uuid}/label [put]
// @Param       x_auth_token header string                      false "client session token"
// @Param       uuid         path   string                      true  "Token 의 Uuid"
// @Param       object       body   v1.HttpReqToken_UpdateLabel true  "Token 의 HttpReqToken_UpdateLabel"
// @Success 	200 {object} v1.Token
func (ctl Control) UpdateTokenLabel(ctx echo.Context) error {
	body := new(tokenv1.HttpReqToken_UpdateLabel)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	//property
	token := tokenv1.Token{}
	token.Uuid = uuid
	token.LabelMeta = NewLabelMeta(body.Name, body.Summary)

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		//update record
		token_, err := vault.NewToken(db).Update(token)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "update token"))
		}

		return token_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// RefreshClusterTokenTime
// @Description Refresh Cluster Token Time
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token/cluster/{uuid}/refresh [put]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Token 의 Uuid"
// @Success     200 {object} v1.Token
func (ctl Control) RefreshClusterTokenTime(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	//property
	token := tokenv1.Token{}
	token.Uuid = uuid
	token.IssuedAtTime, token.ExpirationTime = bearerTokenTimeIssueNow() //만료시간 연장

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		token_, err := vault.NewToken(db).Update(token)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "update token"))
		}

		return token_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// ExpireClusterToken
// @Description Expire Cluster Token
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token/cluster/{uuid}/expire [put]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Token 의 Uuid"
// @Success     200 {object} v1.Token
func (ctl Control) ExpireClusterToken(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	//property
	token := tokenv1.Token{}
	token.Uuid = uuid
	token.IssuedAtTime, token.ExpirationTime = time.Now(), time.Now() //현재 시간으로 만료시간 설정

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		token_, err := vault.NewToken(db).Update(token)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "update token"))
		}

		return token_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// DeleteToken
// @Description Delete a Token
// @Accept      json
// @Produce     json
// @Tags        server/token
// @Router      /server/token/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Token 의 Uuid"
// @Success     200
func (ctl Control) DeleteToken(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		err := vault.NewToken(db).Delete(uuid)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "delete token"))
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, OK())
}

// validTokenUser
func validTokenUser(ctx database.Context, user_kind tokenv1.TokenUserKind, user_uuid string) error {
	switch user_kind {
	case tokenv1.TokenUserKindCluster:
		//get cluster
		if _, err := vault.NewCluster(ctx).Get(user_uuid); err != nil {
			return errors.Wrapf(err, "found cluster token user")
		}
	default:
		return errors.Errorf("invalid token user kind")
	}

	return nil
}

func bearerTokenTimeIssueNow() (time.Time, time.Time) {
	iat := time.Now()
	exp := env.BearerTokenExpirationTime(iat)
	return iat, exp
}

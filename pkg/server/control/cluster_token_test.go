package control_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	clustertokenv1 "github.com/NexClipper/sudory/pkg/server/model/cluster_token/v1"
	cryptov1 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v1"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

// CreateClusterToken
// @Description Create a Cluster Token
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router      /server/cluster_token [post]
// @Param       x_auth_token header string                             false "client session token"
// @Param       object       body   v1.HttpReqClusterToken_Create true  "HttpReqClusterToken_Create"
// @Success     200 {object} v1.ClusterToken
func (ctl Control) CreateClusterToken(ctx echo.Context) error {
	body := new(clustertokenv1.HttpReqClusterToken_Create)
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

	if len(body.ClusterUuid) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.ClusterUuid", TypeName(body)), body.ClusterUuid)...,
				)))
	}

	//valvalidied token user
	if _, err := vault.NewCluster(ctl.db.Engine().NewSession()).Get(body.ClusterUuid); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "valid cluster user%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.ClusterUuid", TypeName(body)), body.ClusterUuid)...,
				)))
	}

	//property
	token := clustertokenv1.ClusterToken{}
	token.UuidMeta = metav1.NewUuidMeta()
	token.LabelMeta = metav1.NewLabelMeta(body.Name, body.Summary)
	token.ClusterUuid = body.ClusterUuid
	token.IssuedAtTime, token.ExpirationTime = bearerTokenTimeIssueNow()
	token.Token = cryptov1.String(macro.NewUuidString())

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		token_, err := vault.NewClusterToken(tx).CreateToken(token)
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

// FindClusterToken
// @Description Find Token
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router      /server/cluster_token [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.ClusterToken
func (ctl Control) FindClusterToken(ctx echo.Context) error {
	r, err := vault.NewClusterToken(ctl.db.Engine().NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "find token"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// GetClusterToken
// @Description Get a Token
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router      /server/cluster_token/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Token 의 Uuid"
// @Success     200 {object} v1.ClusterToken
func (ctl Control) GetClusterToken(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := vault.NewClusterToken(ctl.db.Engine().NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get token"))
	}

	return ctx.JSON(http.StatusOK, r)
}

// UpdateClusterTokenLabel
// @Description Update Token Label
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router		/server/cluster_token/{uuid}/label [put]
// @Param       x_auth_token header string                      false "client session token"
// @Param       uuid         path   string                      true  "Token 의 Uuid"
// @Param       object       body   v1.HttpReqClusterToken_UpdateLabel true  "Token 의 HttpReqClusterToken_UpdateLabel"
// @Success 	200 {object} v1.ClusterToken
func (ctl Control) UpdateClusterTokenLabel(ctx echo.Context) error {
	body := new(clustertokenv1.HttpReqClusterToken_UpdateLabel)
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
	token := clustertokenv1.ClusterToken{}
	token.Uuid = uuid
	token.LabelMeta = metav1.NewLabelMeta(body.Name, body.Summary)

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		//update record
		token_, err := vault.NewClusterToken(tx).Update(token)
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
// @Tags        server/cluster_token
// @Router      /server/cluster_token/{uuid}/refresh [put]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Token 의 Uuid"
// @Success     200 {object} v1.ClusterToken
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
	token := clustertokenv1.ClusterToken{}
	token.Uuid = uuid
	token.IssuedAtTime, token.ExpirationTime = bearerTokenTimeIssueNow() //만료시간 연장

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		token_, err := vault.NewClusterToken(tx).Update(token)
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
// @Tags        server/cluster_token
// @Router      /server/cluster_token/{uuid}/expire [put]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Token 의 Uuid"
// @Success     200 {object} v1.ClusterToken
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
	token := clustertokenv1.ClusterToken{}
	token.Uuid = uuid
	token.IssuedAtTime, token.ExpirationTime = time.Now(), time.Now() //현재 시간으로 만료시간 설정

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		token_, err := vault.NewClusterToken(tx).Update(token)
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

// DeleteClusterToken
// @Description Delete a Token
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router      /server/cluster_token/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Token 의 Uuid"
// @Success     200
func (ctl Control) DeleteClusterToken(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		err := vault.NewClusterToken(tx).Delete(uuid)
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

func bearerTokenTimeIssueNow() (time.Time, time.Time) {
	iat := time.Now()
	exp := globvar.BearerTokenExpirationTime(iat)
	return iat, exp
}

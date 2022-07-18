package control

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	clusterv2 "github.com/NexClipper/sudory/pkg/server/model/cluster/v2"
	clustertokenv2 "github.com/NexClipper/sudory/pkg/server/model/cluster_token/v2"
	cryptov2 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v2"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @Description Create a Cluster Token
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router      /server/cluster_token [post]
// @Param       x_auth_token header string                         false "client session token"
// @Param       object       body   v2.HttpReq_ClusterToken_create true  "ClusterToken HttpReq_ClusterToken_create"
// @Success     200 {object} v2.HttpRsp_ClusterToken
func (ctl ControlVanilla) CreateClusterToken(ctx echo.Context) (err error) {
	body := new(clustertokenv2.HttpReq_ClusterToken_create)
	err = func() (err error) {
		if err := echoutil.Bind(ctx, body); err != nil {
			return errors.Wrapf(err, "bind%s",
				logs.KVL(
					"type", TypeName(body),
				))
		}
		if len(body.Name) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.Name", TypeName(body)), body.Name)...,
				))
		}
		if len(body.ClusterUuid) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.ClusterUuid", TypeName(body)), body.ClusterUuid)...,
				))
		}

		return
	}()
	if err != nil {
		HttpError(err, http.StatusBadRequest)
	}

	//valvalidied token user
	cluster := clusterv2.Cluster{}
	cluster_uuid := vanilla.Equal("uuid", body.ClusterUuid)
	err = vanilla.Stmt.Select(cluster.TableName(), cluster.ColumnNames(), cluster_uuid.Parse(), nil, nil).
		QueryRow(ctl)(func(scan vanilla.Scanner) (err error) {
		err = cluster.Scan(scan)
		return
	})
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	//property
	token := clustertokenv2.ClusterToken{}
	token.Uuid = func() string {
		if 0 < len(body.Uuid) {
			return body.Uuid
		}
		return macro.NewUuidString()
	}()
	token.Name = body.Name
	token.Summary = func() vanilla.NullString {
		if body.Summary != nil {
			return *vanilla.NewNullString(*body.Summary)
		}
		return token.Summary
	}()
	token.ClusterUuid = body.ClusterUuid
	token.IssuedAtTime = time.Now()
	token.ExpirationTime = globvar.BearerTokenExpirationTime(token.IssuedAtTime)
	token.Token = cryptov2.CryptoString(macro.NewUuidString())
	token.Created = time.Now()

	err = func() (err error) {
		stmt, err := vanilla.Stmt.Insert(token.TableName(), token.ColumnNames(), token.Values())
		if err != nil {
			return
		}

		affected, err := stmt.Exec(ctl)
		if err != nil {
			return
		}
		if affected == 0 {
			err = errors.Errorf("no affected")
			return
		}
		return
	}()
	err = errors.Wrapf(err, "failed to create cluster token")
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	token, err = vault.GetClusterToken(ctl, token.Uuid)
	err = errors.Wrapf(err, "not found cluster token")
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, clustertokenv2.HttpRsp_ClusterToken{ClusterToken: token})
}

// @Description Find Cluster Token
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router      /server/cluster_token [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v2.HttpRsp_ClusterToken
func (ctl ControlVanilla) FindClusterToken(ctx echo.Context) error {
	q, o, p, err := ParseDecoration(echoutil.QueryParam(ctx))
	err = errors.Wrapf(err, "ParseDecoration%v", logs.KVL(
		"query", echoutil.QueryParamString(ctx),
	))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	rsp := make([]clustertokenv2.HttpRsp_ClusterToken, 0, __INIT_SLICE_CAPACITY__())

	var token clustertokenv2.ClusterToken
	err = vanilla.Stmt.Select(token.TableName(), token.ColumnNames(), q, o, p).
		QueryRows(ctl)(func(scan vanilla.Scanner, _ int) (err error) {
		err = token.Scan(scan)
		if err != nil {
			return
		}

		rsp = append(rsp, clustertokenv2.HttpRsp_ClusterToken{ClusterToken: token})
		return
	})
	err = errors.Wrapf(err, "not found cluster token")
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, rsp)

}

// @Description Get a Cluster Token
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router      /server/cluster_token/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "ClusterToken Uuid"
// @Success     200 {object} v2.HttpRsp_ClusterToken
func (ctl ControlVanilla) GetClusterToken(ctx echo.Context) (err error) {
	err = func() (err error) {
		if len(echoutil.Param(ctx)[__UUID__]) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				))
		}
		return
	}()
	if err != nil {
		HttpError(err, http.StatusBadRequest)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	token, err := vault.GetClusterToken(ctl, uuid)
	err = errors.Wrapf(err, "not found cluster token")
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, clustertokenv2.HttpRsp_ClusterToken{ClusterToken: token})
}

// UpdateClusterTokenLabel
// @Description Update Label of Cluster Token
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router		/server/cluster_token/{uuid}/label [put]
// @Param       x_auth_token header string                         false "client session token"
// @Param       uuid         path   string                         true  "ClusterToken Uuid"
// @Param       object       body   v2.HttpReq_ClusterToken_update true  "ClusterToken HttpReq_ClusterToken_update"
// @Success 	200 {object} v2.HttpRsp_ClusterToken
func (ctl ControlVanilla) UpdateClusterTokenLabel(ctx echo.Context) (err error) {
	body := new(clustertokenv2.HttpReq_ClusterToken_update)
	err = func() (err error) {
		if err := echoutil.Bind(ctx, body); err != nil {
			return errors.Wrapf(err, "bind%s",
				logs.KVL(
					"type", TypeName(body),
				))
		}
		if len(echoutil.Param(ctx)[__UUID__]) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				))
		}
		return
	}()
	if err != nil {
		HttpError(err, http.StatusBadRequest)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	// get token
	token, err := vault.GetClusterToken(ctl, uuid)
	err = errors.Wrapf(err, "not found cluster token")
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}
	if token.Deleted.Valid {
		err = errors.Errorf("this cluster token is already deleted%v", logs.KVL(
			"uuid", uuid,
		))
		return HttpError(err, http.StatusBadRequest)
	}

	//property
	token.Name = body.Name
	token.Summary = func() vanilla.NullString {
		if body.Summary != nil {
			return *vanilla.NewNullString(*body.Summary)
		}
		return token.Summary
	}()
	token.Updated = *vanilla.NewNullTime(time.Now())

	eq_uuid := vanilla.Equal("uuid", token.Uuid)

	updateSet := map[string]interface{}{}
	var activated bool = false
	active := func() {
		updateSet["updated"] = token.Updated
		activated = true
	}
	if 0 < len(token.Name) {
		updateSet["name"] = token.Name
		active()
	}
	if token.Summary.Valid {
		updateSet["summary"] = token.Summary
		active()
	}

	err = func() (err error) {
		if !activated {
			return
		}

		affected, err := vanilla.Stmt.Update(token.TableName(), updateSet, eq_uuid.Parse()).
			Exec(ctl)
		if err != nil {
			return
		}
		if affected == 0 {
			err = errors.New("no affectd")
		}
		return
	}()
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, clustertokenv2.HttpRsp_ClusterToken{ClusterToken: token})
}

// @Description Refresh Time of Cluster Token
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router      /server/cluster_token/{uuid}/refresh [put]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "ClusterToken Uuid"
// @Success     200 {object} v2.HttpRsp_ClusterToken
func (ctl ControlVanilla) RefreshClusterTokenTime(ctx echo.Context) (err error) {
	err = func() (err error) {
		if len(echoutil.Param(ctx)[__UUID__]) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				))
		}
		return
	}()
	if err != nil {
		HttpError(err, http.StatusBadRequest)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	token, err := vault.GetClusterToken(ctl, uuid)
	err = errors.Wrapf(err, "not found cluster token")
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}
	if token.Deleted.Valid {
		err = errors.Errorf("this cluster token is already deleted%v", logs.KVL(
			"uuid", uuid,
		))
		return HttpError(err, http.StatusBadRequest)
	}

	//property
	token.IssuedAtTime = time.Now()
	token.ExpirationTime = globvar.BearerTokenExpirationTime(time.Now()) //만료시간 연장
	token.Updated = *vanilla.NewNullTime(time.Now())

	eq_uuid := vanilla.Equal("uuid", token.Uuid)
	updateSet := map[string]interface{}{}
	if true {
		updateSet["issued_at_time"] = token.IssuedAtTime
	}
	updateSet["expiration_time"] = token.ExpirationTime
	updateSet["updated"] = token.Updated

	err = func() (err error) {
		affected, err := vanilla.Stmt.Update(token.TableName(), updateSet, eq_uuid.Parse()).
			Exec(ctl)
		if err != nil {
			return
		}
		if affected == 0 {
			err = errors.New("no affectd")
		}
		return
	}()
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, clustertokenv2.HttpRsp_ClusterToken{ClusterToken: token})
}

// @Description Expire Cluster Token
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router      /server/cluster_token/{uuid}/expire [put]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "ClusterToken Uuid"
// @Success     200 {object} v2.HttpRsp_ClusterToken
func (ctl ControlVanilla) ExpireClusterToken(ctx echo.Context) (err error) {
	err = func() (err error) {
		if len(echoutil.Param(ctx)[__UUID__]) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				))
		}
		return
	}()
	if err != nil {
		HttpError(err, http.StatusBadRequest)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	// get a cluster token
	token, err := vault.GetClusterToken(ctl, uuid)
	err = errors.Wrapf(err, "not found cluster token")
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}
	if token.Deleted.Valid {
		err = errors.Errorf("this cluster token is already deleted%v", logs.KVL(
			"uuid", uuid,
		))
		return HttpError(err, http.StatusBadRequest)
	}

	//property
	token.IssuedAtTime = time.Now()
	token.ExpirationTime = time.Now() // 현 시간으로 만료시간 설정
	token.Updated = *vanilla.NewNullTime(time.Now())

	eq_uuid := vanilla.Equal("uuid", token.Uuid)
	updateSet := map[string]interface{}{}
	if false {
		updateSet["issued_at_time"] = token.IssuedAtTime
	}
	updateSet["expiration_time"] = token.ExpirationTime
	updateSet["updated"] = token.Updated

	err = func() (err error) {
		affected, err := vanilla.Stmt.Update(token.TableName(), updateSet, eq_uuid.Parse()).
			Exec(ctl)
		if err != nil {
			return
		}
		if affected == 0 {
			err = errors.New("no affectd")
		}
		return
	}()
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, clustertokenv2.HttpRsp_ClusterToken{ClusterToken: token})
}

// @Description Delete a Cluster Token
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router      /server/cluster_token/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "ClusterToken Uuid"
// @Success     200
func (ctl ControlVanilla) DeleteClusterToken(ctx echo.Context) (err error) {
	err = func() (err error) {
		if len(echoutil.Param(ctx)[__UUID__]) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				))
		}
		return
	}()
	if err != nil {
		HttpError(err, http.StatusBadRequest)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	// property
	token := clustertokenv2.ClusterToken{}
	token.Uuid = uuid
	token.Deleted = *vanilla.NewNullTime(time.Now())

	eq_uuid := vanilla.Equal("uuid", token.Uuid)
	updateSet := map[string]interface{}{}
	updateSet["deleted"] = token.Deleted

	err = func() (err error) {
		affected, err := vanilla.Stmt.Update(token.TableName(), updateSet, eq_uuid.Parse()).
			Exec(ctl)
		if err != nil {
			return
		}
		if affected == 0 {
			err = errors.New("no affectd")
		}
		return
	}()
	if err != nil {
		return HttpError(err, http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK, OK())
}

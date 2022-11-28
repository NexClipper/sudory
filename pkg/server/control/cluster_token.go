package control

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	clusterv3 "github.com/NexClipper/sudory/pkg/server/model/cluster/v3"
	"github.com/NexClipper/sudory/pkg/server/model/cluster_token/v3"
	cryptov2 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v2"
	"github.com/NexClipper/sudory/pkg/server/status/globvar"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @Description Create a Cluster Token
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router      /server/cluster_token [post]
// @Param       object       body   cluster_token.HttpReq_ClusterToken_create true  "ClusterToken HttpReq_ClusterToken_create"
// @Success     200 {object} cluster_token.HttpRsp_ClusterToken
func (ctl ControlVanilla) CreateClusterToken(ctx echo.Context) (err error) {
	body := new(cluster_token.HttpReq_ClusterToken_create)
	err = func() (err error) {
		if err := echoutil.Bind(ctx, body); err != nil {
			return errors.Wrapf(err, "bind%s",
				logs.KVL(
					"type", TypeName(body),
				))
		}
		if len(body.Name) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.Name", TypeName(body)), body.Name)...,
				))
		}
		if len(body.ClusterUuid) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
				logs.KVL(
					ParamLog(fmt.Sprintf("%s.ClusterUuid", TypeName(body)), body.ClusterUuid)...,
				))
		}

		return
	}()
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}
	err = func() error {
		// check cluster
		cluster_table := clusterv3.TableNameWithTenant(claims.Hash)
		cluster_cond := stmt.And(
			stmt.Equal("uuid", body.ClusterUuid),
			stmt.IsNull("deleted"),
		)
		cluster_exist, err := ctl.dialect.Exist(cluster_table, cluster_cond)(ctx.Request().Context(), ctl)
		if err != nil {
			return errors.Wrapf(err, "check cluster")
		}
		if !cluster_exist {
			return errors.Wrapf(database.ErrorRecordWasNotFound, "check cluster")
		}
		return nil
	}()
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	// gen uuid
	body.Uuid = genUuidString(body.Uuid)
	time_now := time.Now()
	tokenstr := macro.NewUuidString()

	//property
	new_token := cluster_token.ClusterToken{}
	new_token.Uuid = body.Uuid
	new_token.Name = body.Name
	new_token.Summary = body.Summary
	new_token.ClusterUuid = body.ClusterUuid
	new_token.IssuedAtTime = time_now
	new_token.ExpirationTime = globvar.BearerToken.ExpirationTime(time_now)
	new_token.Token = cryptov2.CryptoString(tokenstr)
	new_token.Created = time_now
	new_token.Updated = *vanilla.NewNullTime(time_now)

	err = func() error {
		var affected int64
		// save cluster token
		affected, new_token.ID, err = ctl.dialect.Insert(new_token.TableName(), new_token.ColumnNames(), new_token.Values())(
			ctx.Request().Context(), ctl)
		if err != nil {
			return errors.Wrapf(err, "save cluster token")
		}
		if new_token.ID == 0 {
			return errors.Wrapf(database.ErrorNoLastInsertId, "save cluster token")
		}
		if affected == 0 {
			return errors.Wrapf(database.ErrorNoAffected, "save cluster token")
		}

		return nil
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to create cluster token")
	}

	return ctx.JSON(http.StatusOK, cluster_token.HttpRsp_ClusterToken(new_token))
}

// @Description Find Cluster Token
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router      /server/cluster_token [get]
// @Param       q            query  string false "query  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       o            query  string false "order  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       p            query  string false "paging github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Success     200 {array} cluster_token.HttpRsp_ClusterToken
func (ctl ControlVanilla) FindClusterToken(ctx echo.Context) error {
	q, err := stmt.ConditionLexer.Parse(echoutil.QueryParam(ctx)["q"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	o, err := stmt.OrderLexer.Parse(echoutil.QueryParam(ctx)["o"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	p, err := stmt.PaginationLexer.Parse(echoutil.QueryParam(ctx)["p"])
	if err != nil && !logs.DeepCompare(err, stmt.ErrorInvalidArgumentEmptyString) {
		return HttpError(err, http.StatusBadRequest)
	}
	// additional conditon
	q = stmt.And(q,
		stmt.IsNull("deleted"),
	)
	// default pagination
	if p == nil {
		p = stmt.Limit(__DEFAULT_DECORATION_LIMIT__)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	rsp := make([]cluster_token.HttpRsp_ClusterToken, 0, state.ENV__INIT_SLICE_CAPACITY__())
	token := cluster_token.ClusterToken{}
	token_table := cluster_token.TableNameWithTenant(claims.Hash)

	err = ctl.dialect.QueryRows(token_table, token.ColumnNames(), q, o, p)(ctx.Request().Context(), ctl)(
		func(scan excute.Scanner, _ int) error {
			err := token.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			rsp = append(rsp, token)

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "failed to find cluster token")
	}

	return ctx.JSON(http.StatusOK, []cluster_token.HttpRsp_ClusterToken(rsp))

}

// @Description Get a Cluster Token
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router      /server/cluster_token/{uuid} [get]
// @Param       uuid         path   string true  "ClusterToken Uuid"
// @Success     200 {object} cluster_token.HttpRsp_ClusterToken
func (ctl ControlVanilla) GetClusterToken(ctx echo.Context) (err error) {
	err = func() (err error) {
		if len(echoutil.Param(ctx)[__UUID__]) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				))
		}
		return
	}()
	if err != nil {
		HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	token_table := cluster_token.TableNameWithTenant(claims.Hash)
	token := cluster_token.ClusterToken{}
	token.Uuid = uuid

	token_cond := stmt.And(
		stmt.Equal("uuid", token.Uuid),
		stmt.IsNull("deleted"),
	)
	err = ctl.dialect.QueryRow(token_table, token.ColumnNames(), token_cond, nil, nil)(ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := token.Scan(scan)
			err = errors.WithStack(err)

			return err
		})

	if err != nil {
		return errors.Wrapf(err, "failed to get cluster token")
	}

	return ctx.JSON(http.StatusOK, cluster_token.HttpRsp_ClusterToken(token))
}

// @Description Update Label of Cluster Token
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router		/server/cluster_token/{uuid}/label [put]
// @Param       uuid         path   string                         true  "ClusterToken Uuid"
// @Param       object       body   cluster_token.HttpReq_ClusterToken_update true  "ClusterToken HttpReq_ClusterToken_update"
// @Success 	200 {object} cluster_token.HttpRsp_ClusterToken
func (ctl ControlVanilla) UpdateClusterTokenLabel(ctx echo.Context) (err error) {
	body := new(cluster_token.HttpReq_ClusterToken_update)
	err = func() (err error) {
		if err := echoutil.Bind(ctx, body); err != nil {
			return errors.Wrapf(err, "bind%s",
				logs.KVL(
					"type", TypeName(body),
				))
		}
		if len(echoutil.Param(ctx)[__UUID__]) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				))
		}
		return
	}()
	if err != nil {
		HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	// get cluster token
	token_table := cluster_token.TableNameWithTenant(claims.Hash)
	var token cluster_token.ClusterToken
	token.Uuid = uuid
	token_cond := stmt.And(
		stmt.Equal("uuid", token.Uuid),
		stmt.IsNull("deleted"),
	)

	err = ctl.dialect.QueryRow(token_table, token.ColumnNames(), token_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := token.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "get cluster token")
	}

	//property
	time_now := time.Now()

	updateSet := map[string]interface{}{}

	if 0 < len(body.Name) {
		token.Name = body.Name
		updateSet["name"] = token.Name
	}
	if body.Summary.Valid {
		token.Summary = body.Summary
		updateSet["summary"] = token.Summary
	}

	// valid update column counts
	if len(updateSet) == 0 {
		return HttpError(errors.New("noting to update"), http.StatusBadRequest)
	}

	token.Updated = *vanilla.NewNullTime(time_now)
	updateSet["updated"] = token.Updated

	err = func() error {
		// update cluster token
		_, err := ctl.dialect.Update(token.TableName(), updateSet, token_cond)(
			ctx.Request().Context(), ctl)
		if err != nil {
			return errors.Wrapf(err, "update cluster token")
		}
		// if affected == 0 {
		// 	return errors.Wrapf(database.ErrorNoAffected, "update cluster token")
		// }
		return nil
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to update cluster token")
	}

	return ctx.JSON(http.StatusOK, cluster_token.HttpRsp_ClusterToken(token))
}

// @Description Refresh Time of Cluster Token
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router      /server/cluster_token/{uuid}/refresh [put]
// @Param       uuid         path   string true  "ClusterToken Uuid"
// @Success     200 {object} cluster_token.HttpRsp_ClusterToken
func (ctl ControlVanilla) RefreshClusterTokenTime(ctx echo.Context) (err error) {
	err = func() (err error) {
		if len(echoutil.Param(ctx)[__UUID__]) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				))
		}
		return
	}()
	if err != nil {
		HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	// get cluster token
	token_table := cluster_token.TableNameWithTenant(claims.Hash)
	var token cluster_token.ClusterToken
	token.Uuid = uuid
	token_cond := stmt.And(
		stmt.Equal("uuid", token.Uuid),
		stmt.IsNull("deleted"),
	)

	err = ctl.dialect.QueryRow(token_table, token.ColumnNames(), token_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := token.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "get cluster token")
	}

	//property
	time_now := time.Now()

	updateSet := map[string]interface{}{}
	if true {
		token.IssuedAtTime = time_now
		updateSet["issued_at_time"] = token.IssuedAtTime
	}

	token.ExpirationTime = globvar.BearerToken.ExpirationTime(time_now) //만료시간 연장
	updateSet["expiration_time"] = token.ExpirationTime

	token.Updated = *vanilla.NewNullTime(time_now)
	updateSet["updated"] = token.Updated

	err = func() error {
		// update cluster token
		_, err := ctl.dialect.Update(token.TableName(), updateSet, token_cond)(
			ctx.Request().Context(), ctl)
		if err != nil {
			return errors.Wrapf(err, "update cluster token")
		}
		// if affected == 0 {
		// 	return errors.Wrapf(database.ErrorNoAffected, "update cluster token")
		// }
		return nil
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to update cluster token")
	}

	return ctx.JSON(http.StatusOK, cluster_token.HttpRsp_ClusterToken(token))
}

// @Description Expire Cluster Token
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router      /server/cluster_token/{uuid}/expire [put]
// @Param       uuid         path   string true  "ClusterToken Uuid"
// @Success     200 {object} cluster_token.HttpRsp_ClusterToken
func (ctl ControlVanilla) ExpireClusterToken(ctx echo.Context) (err error) {
	err = func() (err error) {
		if len(echoutil.Param(ctx)[__UUID__]) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				))
		}
		return
	}()
	if err != nil {
		HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	// get cluster token
	token_table := cluster_token.TableNameWithTenant(claims.Hash)
	var token cluster_token.ClusterToken
	token.Uuid = uuid
	token_cond := stmt.And(
		stmt.Equal("uuid", token.Uuid),
		stmt.IsNull("deleted"),
	)

	err = ctl.dialect.QueryRow(token_table, token.ColumnNames(), token_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := token.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "get cluster token")
	}

	//property
	time_now := time.Now()

	updateSet := map[string]interface{}{}
	if false {
		token.IssuedAtTime = time_now
		updateSet["issued_at_time"] = token.IssuedAtTime
	}

	token.ExpirationTime = time_now // 현 시간으로 만료시간 설정
	updateSet["expiration_time"] = token.ExpirationTime

	token.Updated = *vanilla.NewNullTime(time_now)
	updateSet["updated"] = token.Updated

	err = func() error {
		// update cluster token
		_, err := ctl.dialect.Update(token.TableName(), updateSet, token_cond)(
			ctx.Request().Context(), ctl)
		if err != nil {
			return errors.Wrapf(err, "update cluster token")
		}
		// if affected == 0 {
		// 	return errors.Wrapf(database.ErrorNoAffected, "update cluster token")
		// }
		return nil
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to update cluster token")
	}

	return ctx.JSON(http.StatusOK, cluster_token.HttpRsp_ClusterToken(token))
}

// @Description Delete a Cluster Token
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/cluster_token
// @Router      /server/cluster_token/{uuid} [delete]
// @Param       uuid         path   string true  "ClusterToken Uuid"
// @Success     200
func (ctl ControlVanilla) DeleteClusterToken(ctx echo.Context) (err error) {
	err = func() (err error) {
		if len(echoutil.Param(ctx)[__UUID__]) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				))
		}
		return
	}()
	if err != nil {
		HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	// get cluster token
	token_table := cluster_token.TableNameWithTenant(claims.Hash)
	var token cluster_token.ClusterToken
	token.Uuid = uuid
	token_cond := stmt.And(
		stmt.Equal("uuid", token.Uuid),
	)

	err = ctl.dialect.QueryRow(token_table, token.ColumnNames(), token_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := token.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return errors.Wrapf(err, "get cluster token")
	}

	// property
	time_now := time.Now()

	token.Deleted = *vanilla.NewNullTime(time_now)
	updateSet := map[string]interface{}{}
	updateSet["deleted"] = token.Deleted

	err = func() error {
		// update cluster token
		_, err := ctl.dialect.Update(token.TableName(), updateSet, token_cond)(ctx.Request().Context(), ctl)
		if err != nil {
			return errors.Wrapf(err, "update cluster token")
		}
		// if affected == 0 {
		// 	return errors.Wrapf(database.ErrorNoAffected, "update cluster token")
		// }
		return nil
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to delete cluster token")
	}

	return ctx.JSON(http.StatusOK, OK())
}

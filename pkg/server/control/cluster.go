package control

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmtex"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	clusterv3 "github.com/NexClipper/sudory/pkg/server/model/cluster/v3"
	"github.com/NexClipper/sudory/pkg/server/model/tenants/v3"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @Description Create a cluster
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster [post]
// @Param       cluster      body   v3.HttpReq_Cluster_create true  "HttpReq_Cluster_create"
// @Success     200 {object} v3.HttpRsp_Cluster
func (ctl ControlVanilla) CreateCluster(ctx echo.Context) (err error) {
	body := new(clusterv3.HttpReq_Cluster_create)
	err = echoutil.WrapHttpError(http.StatusBadRequest,
		func() error {
			err := echoutil.Bind(ctx, body)
			return errors.Wrapf(err, "bind%v", logs.KVL(
				"type", TypeName(body),
			))
		},
		func() (err error) {
			if len(body.Name) == 0 {
				err = ErrorInvalidRequestParameter
			}
			return errors.Wrapf(err, "valid param%v", logs.KVL(
				ParamLog(fmt.Sprintf("%s.Name", TypeName(body)), body.Name)...,
			))
		},
	)
	if err != nil {
		return errors.WithStack(err)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	time_now := time.Now()
	// gen uuid
	body.Uuid = genUuidString(body.Uuid)

	// property
	// new new_cluster
	new_cluster := clusterv3.Cluster{}
	new_cluster.Uuid = body.Uuid
	new_cluster.Name = body.Name
	new_cluster.Summary = body.Summary
	new_cluster.PollingOption = *vanilla.NewNullObject(clusterv3.ConvPollingOption(body.PollingOption).ToMap())
	new_cluster.PoliingLimit = body.PoliingLimit.Int()
	new_cluster.Created = time_now

	err = stmtex.ScopeTx(ctx.Request().Context(), ctl, func(tx *sql.Tx) error {
		var affected int64
		// save cluster
		affected, new_cluster.ID, err = stmtex.Insert(new_cluster.TableName(), new_cluster.ColumnNames(), new_cluster.Values()).
			ExecContext(ctx.Request().Context(), tx, ctl.Dialect())
		if err != nil {
			return errors.Wrapf(err, "save cluster")
		}
		if new_cluster.ID == 0 {
			return errors.Wrapf(database.ErrorNoLastInsertId, "save cluster")
		}
		if affected == 0 {
			return errors.Wrapf(database.ErrorNoAffected, "save cluster")
		}

		// save tenant_clusters
		tenant_clusters := new(tenants.TenantClusters)
		tenant_clusters.TenantId = claims.ID
		tenant_clusters.ClusterId = new_cluster.ID
		affected, _, err = stmtex.Insert(tenant_clusters.TableName(), tenant_clusters.ColumnNames(), tenant_clusters.Values()).
			ExecContext(ctx.Request().Context(), tx, ctl.Dialect())
		if err != nil {
			return errors.Wrapf(err, "save tenant clusters")
		}
		if affected == 0 {
			return errors.Wrapf(database.ErrorNoAffected, "save tenant clusters")
		}

		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "failed to create cluster")
	}

	return ctx.JSON(http.StatusOK, clusterv3.HttpRsp_Cluster(new_cluster))
}

// @Description Find cluster
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster [get]
// @Param       q            query  string false "query  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       o            query  string false "order  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       p            query  string false "paging github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Success     200 {array} v3.HttpRsp_Cluster
func (ctl ControlVanilla) FindCluster(ctx echo.Context) (err error) {
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

	rsp := make([]clusterv3.HttpRsp_Cluster, 0, state.ENV__INIT_SLICE_CAPACITY__())

	cluster_table := clusterv3.TableNameWithTenant(claims.Hash)
	cluster := clusterv3.Cluster{}

	err = stmtex.Select(cluster_table, cluster.ColumnNames(), q, o, p).
		QueryRowsContext(ctx.Request().Context(), ctl, ctl.Dialect())(
		func(scan stmtex.Scanner, _ int) error {
			err = cluster.Scan(scan)
			if err != nil {
				return errors.Wrapf(err, "failed to scan")
			}
			rsp = append(rsp, cluster)
			return nil
		})
	if err != nil {
		return errors.Wrapf(err, "failed to find clusters")
	}

	return ctx.JSON(http.StatusOK, []clusterv3.HttpRsp_Cluster(rsp))
}

// @Description Get a cluster
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid} [get]
// @Param       uuid         path   string true  "Cluster Uuid"
// @Success     200 {object} v3.HttpRsp_Cluster
func (ctl ControlVanilla) GetCluster(ctx echo.Context) (err error) {
	err = echoutil.WrapHttpError(http.StatusBadRequest,
		func() (err error) {
			if len(echoutil.Param(ctx)[__UUID__]) == 0 {
				err = ErrorInvalidRequestParameter
			}
			return errors.Wrapf(err, "valid param%v", logs.KVL(
				ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
			))
		})
	if err != nil {
		return
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	cluster_table := clusterv3.TableNameWithTenant(claims.Hash)
	cluster := clusterv3.Cluster{}
	cluster.Uuid = uuid

	cond := stmt.And(
		stmt.Equal("uuid", cluster.Uuid),
		stmt.IsNull("deleted"),
	)

	err = stmtex.Select(cluster_table, cluster.ColumnNames(), cond, nil, nil).
		QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(
		func(scan stmtex.Scanner) (err error) {
			err = cluster.Scan(scan)
			if err != nil {
				return errors.Wrapf(err, "failed to scan")
			}
			return
		})
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusOK, clusterv3.HttpRsp_Cluster(cluster))
}

// @Description Update a cluster
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid} [put]
// @Param       uuid         path   string                    true  "Cluster Uuid"
// @Param       cluster      body   v3.HttpReq_Cluster_update true  "HttpReq_Cluster_update"
// @Success     200 {object} v3.HttpRsp_Cluster
func (ctl ControlVanilla) UpdateCluster(ctx echo.Context) (err error) {
	body := new(clusterv3.HttpReq_Cluster_update)

	err = echoutil.WrapHttpError(http.StatusBadRequest,
		func() error {
			err := echoutil.Bind(ctx, body)
			return errors.Wrapf(err, "bind%s",
				logs.KVL(
					"type", TypeName(body),
				))
		},
		func() (err error) {
			if len(echoutil.Param(ctx)[__UUID__]) == 0 {
				err = ErrorInvalidRequestParameter
			}
			return errors.Wrapf(err, "valid param%v", logs.KVL(
				ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
			))
		},
	)
	if err != nil {
		return
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]
	//get cluster
	var cluster clusterv3.Cluster
	cluster.Uuid = uuid //set uuid from path
	cluster_cond := stmt.And(
		stmt.Equal("uuid", cluster.Uuid),
		stmt.IsNull("deleted"),
	)

	cluster_table := clusterv3.TableNameWithTenant(claims.Hash)
	err = stmtex.Select(cluster_table, cluster.ColumnNames(), cluster_cond, nil, nil).
		QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(
		func(scan stmtex.Scanner) error {
			return cluster.Scan(scan)
		})
	if err != nil {
		return errors.Wrapf(err, "get cluster")
	}

	time_now := time.Now()

	//property
	updateSet := map[string]interface{}{}

	if 0 < len(body.Name) {
		cluster.Name = body.Name
		updateSet["name"] = cluster.Name
	}
	if body.Summary.Valid {
		cluster.Summary = body.Summary
		updateSet["summary"] = cluster.Summary
	}
	if body.PollingOption.Valid {
		cluster.PollingOption = *vanilla.NewNullObject(clusterv3.ConvPollingOption(body.PollingOption).ToMap())
		updateSet["polling_option"] = cluster.PollingOption
	}
	if body.PoliingLimit.Valid {
		cluster.PoliingLimit = body.PoliingLimit.Int()
		updateSet["polling_limit"] = cluster.PoliingLimit
	}

	// valied update column counts
	if len(updateSet) == 0 {
		return HttpError(errors.New("noting to update"), http.StatusBadRequest)
	}

	cluster.Updated = *vanilla.NewNullTime(time_now)
	updateSet["updated"] = cluster.Updated

	// update
	err = func() error {
		// // check cluster
		// cluster_table := clusterv3.TableNameWithTenant(claims.Hash)
		// tc_exist, err := stmtex.Exist(cluster_table, cluster_cond)(ctx.Request().Context(), ctl, ctl.Dialect())
		// if err != nil {
		// 	return errors.Wrapf(err, "check cluster")
		// }
		// if !tc_exist {
		// 	return errors.Wrapf(database.ErrorRecordWasNotFound, "check cluster")
		// }
		// update cluster
		_, err := stmtex.Update(cluster.TableName(), updateSet, cluster_cond).
			ExecContext(ctx.Request().Context(), ctl, ctl.Dialect())
		if err != nil {
			return errors.Wrapf(err, "update cluster")
		}
		// if affected == 0 {
		// 	return errors.Wrapf(database.ErrorNoAffected, "update cluster")
		// }
		return nil
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to update cluster")
	}

	return ctx.JSON(http.StatusOK, clusterv3.HttpRsp_Cluster(cluster))
}

// @Description Update a cluster Polling Reguar
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid}/polling/regular [put]
// @Param       uuid           path   string                  true  "Cluster Uuid"
// @Param       polling_option body   v3.RegularPollingOption true  "RegularPollingOption"
// @Success     200 {object} v3.HttpRsp_Cluster
func (ctl ControlVanilla) UpdateClusterPollingRegular(ctx echo.Context) (err error) {
	body := new(clusterv3.RegularPollingOption)
	err = echoutil.WrapHttpError(http.StatusBadRequest,
		func() error {
			err := echoutil.Bind(ctx, body)
			return errors.Wrapf(err, "bind%s",
				logs.KVL(
					"type", TypeName(body),
				))
		},
		func() (err error) {
			if len(echoutil.Param(ctx)[__UUID__]) == 0 {
				err = ErrorInvalidRequestParameter
			}
			return errors.Wrapf(err, "valid param%v", logs.KVL(
				ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
			))
		},
	)
	if err != nil {
		return
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	// get cluster
	var cluster clusterv3.Cluster
	cluster.Uuid = uuid //set uuid from path
	cluster_cond := stmt.And(
		stmt.Equal("uuid", cluster.Uuid),
		stmt.IsNull("deleted"),
	)

	cluster_table := clusterv3.TableNameWithTenant(claims.Hash)
	err = stmtex.Select(cluster_table, cluster.ColumnNames(), cluster_cond, nil, nil).
		QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(
		func(scan stmtex.Scanner) error {
			return cluster.Scan(scan)
		})
	if err != nil {
		return errors.Wrapf(err, "get cluster")
	}

	//property
	updateSet := map[string]interface{}{}

	cluster.PollingOption = *vanilla.NewNullObject(body.ToMap())
	updateSet["polling_option"] = cluster.PollingOption

	cluster.Updated = *vanilla.NewNullTime(time.Now())
	updateSet["updated"] = cluster.Updated

	// update
	err = func() error {

		// update cluster
		_, err := stmtex.Update(cluster.TableName(), updateSet, cluster_cond).
			ExecContext(ctx.Request().Context(), ctl, ctl.Dialect())
		if err != nil {
			return errors.Wrapf(err, "update cluster")
		}
		// if affected == 0 {
		// 	return errors.Wrapf(database.ErrorNoAffected, "update cluster")
		// }
		return nil
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to update cluster")
	}

	return ctx.JSON(http.StatusOK, clusterv3.HttpRsp_Cluster(cluster))
}

// @Description Update a cluster Polling Smart
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid}/polling/smart [put]
// @Param       uuid           path   string                true  "Cluster Uuid"
// @Param       polling_option body   v3.SmartPollingOption true  "SmartPollingOption"
// @Success     200 {object} v3.HttpRsp_Cluster
func (ctl ControlVanilla) UpdateClusterPollingSmart(ctx echo.Context) (err error) {
	body := new(clusterv3.SmartPollingOption)
	err = echoutil.WrapHttpError(http.StatusBadRequest,
		func() error {
			err := echoutil.Bind(ctx, body)
			return errors.Wrapf(err, "bind%s",
				logs.KVL(
					"type", TypeName(body),
				))
		},
		func() (err error) {
			if len(echoutil.Param(ctx)[__UUID__]) == 0 {
				err = ErrorInvalidRequestParameter
			}
			return errors.Wrapf(err, "valid param%v", logs.KVL(
				ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
			))
		},
	)
	if err != nil {
		return
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	//get cluster
	var cluster clusterv3.Cluster
	cluster.Uuid = uuid //set uuid from path
	cluster_cond := stmt.And(
		stmt.Equal("uuid", cluster.Uuid),
		stmt.IsNull("deleted"),
	)

	cluster_table := clusterv3.TableNameWithTenant(claims.Hash)
	err = stmtex.Select(cluster_table, cluster.ColumnNames(), cluster_cond, nil, nil).
		QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(
		func(scan stmtex.Scanner) error {
			return cluster.Scan(scan)
		})
	if err != nil {
		return errors.Wrapf(err, "get cluster")
	}

	//property

	updateSet := map[string]interface{}{}

	cluster.PollingOption = *vanilla.NewNullObject(body.ToMap())
	updateSet["polling_option"] = cluster.PollingOption

	cluster.Updated = *vanilla.NewNullTime(time.Now())
	updateSet["updated"] = cluster.Updated

	// update
	err = func() error {
		// update cluster
		_, err := stmtex.Update(cluster.TableName(), updateSet, cluster_cond).
			ExecContext(ctx.Request().Context(), ctl, ctl.Dialect())
		if err != nil {
			return errors.Wrapf(err, "update cluster")
		}
		// if affected == 0 {
		// 	return errors.Wrapf(database.ErrorNoAffected, "update cluster")
		// }
		return nil
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to update cluster")
	}

	return ctx.JSON(http.StatusOK, clusterv3.HttpRsp_Cluster(cluster))
}

// @Description Delete a cluster
// @Security    ServiceAuth
// @Accept json
// @Produce json
// @Tags server/cluster
// @Router /server/cluster/{uuid} [delete]
// @Param       uuid         path   string true  "Cluster Uuid"
// @Success 200
func (ctl ControlVanilla) DeleteCluster(ctx echo.Context) (err error) {
	err = echoutil.WrapHttpError(http.StatusBadRequest,
		func() (err error) {
			if len(echoutil.Param(ctx)[__UUID__]) == 0 {
				err = ErrorInvalidRequestParameter
			}
			return errors.Wrapf(err, "valid param%v", logs.KVL(
				ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
			))
		})
	if err != nil {
		return
	}
	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	var cluster clusterv3.Cluster
	cluster.Uuid = uuid //set uuid from path
	cluster_cond := stmt.And(
		stmt.Equal("uuid", cluster.Uuid),
		// stmt.IsNull("deleted"),
	)

	cluster_table := clusterv3.TableNameWithTenant(claims.Hash)
	err = stmtex.Select(cluster_table, cluster.ColumnNames(), cluster_cond, nil, nil).
		QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(
		func(scan stmtex.Scanner) error {
			return cluster.Scan(scan)
		})
	if err != nil {
		return errors.Wrapf(err, "get cluster")
	}

	if cluster.Deleted.Valid {
		return ctx.JSON(http.StatusOK, OK())
	}

	//property
	time_now := time.Now()

	cluster.Deleted = *vanilla.NewNullTime(time_now)
	updateSet := map[string]interface{}{
		"deleted": cluster.Deleted,
	}

	err = func() error {
		// update cluster
		_, err := stmtex.Update(cluster.TableName(), updateSet, cluster_cond).
			ExecContext(ctx.Request().Context(), ctl, ctl.Dialect())
		if err != nil {
			return errors.Wrapf(err, "update cluster")
		}
		// if affected == 0 {
		// 	return errors.Wrapf(database.ErrorNoAffected, "update cluster")
		// }
		return nil
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to delete cluster")
	}

	return ctx.JSON(http.StatusOK, OK())
}

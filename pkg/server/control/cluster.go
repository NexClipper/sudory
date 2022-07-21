package control

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	clusterv2 "github.com/NexClipper/sudory/pkg/server/model/cluster/v2"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @Description Create a cluster
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster [post]
// @Param       x_auth_token header string                    false "client session token"
// @Param       cluster      body   v2.HttpReq_Cluster_create true  "HttpReq_Cluster_create"
// @Success     200 {object} v2.Cluster
func (ctl ControlVanilla) CreateCluster(ctx echo.Context) (err error) {
	body := new(clusterv2.HttpReq_Cluster_create)
	err = echoutil.WrapHttpError(http.StatusBadRequest,
		func() error {
			err := echoutil.Bind(ctx, body)
			return errors.Wrapf(err, "bind%v", logs.KVL(
				"type", TypeName(body),
			))
		},
		func() (err error) {
			if len(body.Name) == 0 {
				err = ErrorInvalidRequestParameter()
			}
			return errors.Wrapf(err, "valid param%v", logs.KVL(
				ParamLog(fmt.Sprintf("%s.Name", TypeName(body)), body.Name)...,
			))
		},
	)
	if err != nil {
		return err
	}

	//property
	cluster := clusterv2.Cluster{}
	cluster.Uuid = func() string {
		if 0 < len(body.Uuid) {
			return body.Uuid
		}
		return macro.NewUuidString()
	}()
	cluster.Name = body.Name
	cluster.Summary = body.Summary

	if body.PollingOption.Valid {
		cluster.PollingOption = *vanilla.NewNullObject(body.GetPollingOption().ToMap())
	}
	cluster.PoliingLimit = body.PoliingLimit
	cluster.Created = time.Now()

	eq_uuid := vanilla.And(
		vanilla.Equal("uuid", cluster.Uuid),
		vanilla.IsNull("deleted"),
	)

	// insert
	stmt, err := vanilla.Stmt.Insert(cluster.TableName(), cluster.ColumnNames(), cluster.Values())
	if err != nil {
		return err
	}
	affected, err := stmt.Exec(ctl)
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("no affected")
	}

	// get
	err = vanilla.Stmt.Select(cluster.TableName(), cluster.ColumnNames(), eq_uuid.Parse(), nil, nil).
		QueryRow(ctl)(func(scan vanilla.Scanner) (err error) {
		err = cluster.Scan(scan)
		if err != nil {
			return errors.Wrapf(err, "failed to scan")
		}
		return
	})
	if err != nil {
		return errors.Wrapf(err, "not found record%v", logs.KVL(
			"table", cluster.TableName(),
		))
	}

	return ctx.JSON(http.StatusOK, cluster)
}

// @Description Find cluster
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v2.HttpRsp_Cluster
func (ctl ControlVanilla) FindCluster(ctx echo.Context) (err error) {
	q, o, p, err := ParseDecoration(echoutil.QueryParam(ctx))
	err = errors.Wrapf(err, "ParseDecoration%v", logs.KVL(
		"query", echoutil.QueryParamString(ctx),
	))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	rsp := make([]clusterv2.HttpRsp_Cluster, 0, state.ENV__INIT_SLICE_CAPACITY__())

	cluster := clusterv2.Cluster{}
	err = vanilla.Stmt.Select(cluster.TableName(), cluster.ColumnNames(), q, o, p).
		QueryRows(ctl)(func(scan vanilla.Scanner, _ int) (err error) {
		err = cluster.Scan(scan)
		if err != nil {
			return errors.Wrapf(err, "failed to scan")
		}
		rsp = append(rsp, clusterv2.HttpRsp_Cluster{Cluster: cluster})
		return
	})
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusOK, rsp)
}

// Get Cluster
// @Description Get a cluster
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Cluster Uuid"
// @Success     200 {object} v2.HttpRsp_Cluster
func (ctl ControlVanilla) GetCluster(ctx echo.Context) (err error) {
	err = echoutil.WrapHttpError(http.StatusBadRequest,
		func() (err error) {
			if len(echoutil.Param(ctx)[__UUID__]) == 0 {
				err = ErrorInvalidRequestParameter()
			}
			return errors.Wrapf(err, "valid param%v", logs.KVL(
				ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
			))
		})
	if err != nil {
		return
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	cluster := clusterv2.Cluster{}
	cluster.Uuid = uuid
	eq_uuid := vanilla.And(
		vanilla.Equal("uuid", cluster.Uuid),
		// vanilla.IsNull("deleted"),
	)

	err = vanilla.Stmt.Select(cluster.TableName(), cluster.ColumnNames(), eq_uuid.Parse(), nil, nil).
		QueryRow(ctl)(func(scan vanilla.Scanner) (err error) {
		err = cluster.Scan(scan)
		if err != nil {
			return errors.Wrapf(err, "failed to scan")
		}
		return
	})
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusOK, clusterv2.HttpRsp_Cluster{Cluster: cluster})
}

// @Description Update a cluster
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid} [put]
// @Param       x_auth_token header string                    false "client session token"
// @Param       uuid         path   string                    true  "Cluster Uuid"
// @Param       cluster      body   v2.HttpReq_Cluster_update true  "HttpReq_Cluster_update"
// @Success     200 {object} v2.HttpRsp_Cluster
func (ctl ControlVanilla) UpdateCluster(ctx echo.Context) (err error) {
	body := new(clusterv2.HttpReq_Cluster_update)

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
				err = ErrorInvalidRequestParameter()
			}
			return errors.Wrapf(err, "valid param%v", logs.KVL(
				ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
			))
		},
	)
	if err != nil {
		return
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	//property
	cluster := clusterv2.Cluster{}
	cluster.Uuid = uuid //set uuid from path
	cluster.Name = body.Name
	cluster.Summary = body.Summary
	if body.PollingOption.Valid {
		cluster.PollingOption = *vanilla.NewNullObject(body.GetPollingOption().ToMap())
	}
	cluster.PoliingLimit = body.PoliingLimit
	cluster.Updated = *vanilla.NewNullTime(time.Now())

	eq_uuid := vanilla.And(
		vanilla.Equal("uuid", cluster.Uuid),
		vanilla.IsNull("deleted"),
	)

	updateSet := map[string]interface{}{}
	if 0 < len(body.Name) {
		updateSet["name"] = cluster.Name
	}
	if cluster.Summary.Valid {
		updateSet["summary"] = cluster.Summary
	}
	if cluster.PollingOption.Valid {
		updateSet["polling_option"] = cluster.PollingOption
	}
	if 0 <= cluster.PoliingLimit {
		updateSet["polling_limit"] = cluster.PoliingLimit
	}
	updateSet["updated"] = cluster.Updated

	// update
	affected, err := vanilla.Stmt.Update(cluster.TableName(), updateSet, eq_uuid.Parse()).
		Exec(ctl)
	if err != nil {
		return
	}
	if affected == 0 {
		return errors.New("no affected")
	}

	// get
	err = vanilla.Stmt.Select(cluster.TableName(), cluster.ColumnNames(), eq_uuid.Parse(), nil, nil).
		QueryRow(ctl)(func(scan vanilla.Scanner) (err error) {
		err = cluster.Scan(scan)
		if err != nil {
			return errors.Wrapf(err, "failed to scan")
		}
		return
	})
	if err != nil {
		return errors.Wrapf(err, "not found record%v", logs.KVL(
			"table", cluster.TableName(),
		))
	}

	return ctx.JSON(http.StatusOK, clusterv2.HttpRsp_Cluster{Cluster: cluster})
}

// @Description Update a cluster Polling Reguar
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid}/polling/regular [put]
// @Param       x_auth_token   header string                  false "client session token"
// @Param       uuid           path   string                  true  "Cluster Uuid"
// @Param       polling_option body   v2.RegularPollingOption true  "RegularPollingOption"
// @Success     200 {object} v2.HttpRsp_Cluster
func (ctl ControlVanilla) UpdateClusterPollingRegular(ctx echo.Context) (err error) {
	body := new(clusterv2.RegularPollingOption)
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
				err = ErrorInvalidRequestParameter()
			}
			return errors.Wrapf(err, "valid param%v", logs.KVL(
				ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
			))
		},
	)
	if err != nil {
		return
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	cluster, err := ctl.updateClusterPollingOptions(uuid, body)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, clusterv2.HttpRsp_Cluster{Cluster: cluster})
}

// @Description Update a cluster Polling Smart
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid}/polling/smart [put]
// @Param       x_auth_token   header string                false "client session token"
// @Param       uuid           path   string                true  "Cluster Uuid"
// @Param       polling_option body   v2.SmartPollingOption true  "SmartPollingOption"
// @Success     200 {object} v2.HttpRsp_Cluster
func (ctl ControlVanilla) UpdateClusterPollingSmart(ctx echo.Context) (err error) {
	body := new(clusterv2.SmartPollingOption)
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
				err = ErrorInvalidRequestParameter()
			}
			return errors.Wrapf(err, "valid param%v", logs.KVL(
				ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
			))
		},
	)
	if err != nil {
		return
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	cluster, err := ctl.updateClusterPollingOptions(uuid, body)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, clusterv2.HttpRsp_Cluster{Cluster: cluster})
}

// @Description Delete a cluster
// @Accept json
// @Produce json
// @Tags server/cluster
// @Router /server/cluster/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Cluster Uuid"
// @Success 200
func (ctl ControlVanilla) DeleteCluster(ctx echo.Context) (err error) {
	err = echoutil.WrapHttpError(http.StatusBadRequest,
		func() (err error) {
			if len(echoutil.Param(ctx)[__UUID__]) == 0 {
				err = ErrorInvalidRequestParameter()
			}
			return errors.Wrapf(err, "valid param%v", logs.KVL(
				ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
			))
		})
	if err != nil {
		return
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	channel := clusterv2.Cluster{}
	channel.Uuid = uuid
	channel.Deleted = *vanilla.NewNullTime(time.Now())
	updateSet := map[string]interface{}{
		"deleted": channel.Deleted,
	}
	eq_uuid := vanilla.And(
		vanilla.Equal("uuid", uuid),
	)

	affected, err := vanilla.Stmt.Update(channel.TableName(), updateSet, eq_uuid.Parse()).
		Exec(ctl)
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.Errorf("no affected")
	}

	return ctx.JSON(http.StatusOK, OK())
}

func (ctl ControlVanilla) updateClusterPollingOptions(uuid string, polling_option clusterv2.PollingHandler) (cluster clusterv2.Cluster, err error) {

	//property
	cluster.Uuid = uuid                      //set uuid from path
	cluster.SetPollingOption(polling_option) //update polling option

	cluster.Updated = *vanilla.NewNullTime(time.Now())

	eq_uuid := vanilla.And(
		vanilla.Equal("uuid", cluster.Uuid),
		vanilla.IsNull("deleted"),
	)

	updateSet := map[string]interface{}{}

	if cluster.PollingOption.Valid {
		updateSet["polling_option"] = cluster.PollingOption
	}

	updateSet["updated"] = cluster.Updated

	// update
	affected, err := vanilla.Stmt.Update(cluster.TableName(), updateSet, eq_uuid.Parse()).
		Exec(ctl)
	if err != nil {
		return
	}
	if affected == 0 {
		return cluster, errors.New("no affected")
	}

	// get
	err = vanilla.Stmt.Select(cluster.TableName(), cluster.ColumnNames(), eq_uuid.Parse(), nil, nil).
		QueryRow(ctl)(func(scan vanilla.Scanner) (err error) {
		err = cluster.Scan(scan)
		if err != nil {
			return errors.Wrapf(err, "failed to scan")
		}
		return
	})
	if err != nil {
		return cluster, errors.Wrapf(err, "not found record%v", logs.KVL(
			"table", cluster.TableName(),
		))
	}

	return
}

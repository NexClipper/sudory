package control

import (
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmtex"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	sessionv3 "github.com/NexClipper/sudory/pkg/server/model/session/v3"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @Description Find Session
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/session
// @Router      /server/session [get]
// @Param       q            query  string false "query  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       o            query  string false "order  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       p            query  string false "paging github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Success     200 {array} v3.Session
func (ctl ControlVanilla) FindSession(ctx echo.Context) (err error) {
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

	rsp := make([]sessionv3.Session, 0, state.ENV__INIT_SLICE_CAPACITY__())
	session := sessionv3.Session{}
	session_table := sessionv3.TableNameWithTenant(claims.Hash)

	err = stmtex.Select(session_table, session.ColumnNames(), q, o, p).
		QueryRowsContext(ctx.Request().Context(), ctl, ctl.Dialect())(
		func(scan stmtex.Scanner, _ int) error {
			err = session.Scan(scan)
			if err != nil {
				return errors.Wrapf(err, "failed to scan")
			}
			rsp = append(rsp, session)
			return nil
		})
	if err != nil {
		return errors.Wrapf(err, "failed to find sessions")
	}

	return ctx.JSON(http.StatusOK, []sessionv3.Session(rsp))
}

// @Description Get a Session
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/session
// @Router      /server/session/{uuid} [get]
// @Param       uuid         path   string true  "Session Uuid"
// @Success     200 {object} v3.Session
func (ctl ControlVanilla) GetSession(ctx echo.Context) (err error) {
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
		return HttpError(err, http.StatusBadRequest)
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	session_table := sessionv3.TableNameWithTenant(claims.Hash)
	session := sessionv3.Session{}
	session.Uuid = uuid

	cond := stmt.And(
		stmt.Equal("uuid", session.Uuid),
		stmt.IsNull("deleted"),
	)

	err = stmtex.Select(session_table, session.ColumnNames(), cond, nil, nil).
		QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(
		func(scan stmtex.Scanner) (err error) {
			err = session.Scan(scan)
			return
		})
	if err != nil {
		return errors.Wrapf(err, "failed to get session")
	}

	return ctx.JSON(http.StatusOK, sessionv3.Session(session))
}

// @Description Alive Cluster Session
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/session
// @Router      /server/session/cluster/{cluster_uuid}/alive [get]
// @Param       cluster_uuid path   string true  "Cluster Uuid"
// @Success     200 {object} boolean
func (ctl ControlVanilla) AliveClusterSession(ctx echo.Context) (err error) {
	const __CLUSTER_UUID__ = "cluster_uuid"

	err = func() (err error) {
		if len(echoutil.Param(ctx)[__CLUSTER_UUID__]) == 0 {
			return errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
				logs.KVL(
					ParamLog(__CLUSTER_UUID__, echoutil.Param(ctx)[__CLUSTER_UUID__])...,
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

	cluster_uuid := echoutil.Param(ctx)[__CLUSTER_UUID__]

	session_table := sessionv3.TableNameWithTenant(claims.Hash)
	session := sessionv3.Session{}

	session.ClusterUuid = cluster_uuid
	cond := stmt.And(
		stmt.Equal("cluster_uuid", session.ClusterUuid),
		stmt.IsNull("deleted"),
	)
	order := stmt.Desc("expiration_time")
	page := stmt.Limit(1)
	err = stmtex.Select(session_table, session.ColumnNames(), cond, order, page).
		QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(func(scan stmtex.Scanner) error {
		return session.Scan(scan)
	})
	if err != nil {
		return err
	}

	var expt bool = false
	if session.ExpirationTime.Valid {
		expt = time.Now().Before(session.ExpirationTime.Time)
	}

	return ctx.JSON(http.StatusOK, expt)
}

// @Description Delete a Session
// @Security    ServiceAuth
// @Accept      json
// @Produce     json
// @Tags        server/session
// @Router      /server/session/{uuid} [delete]
// @Param       uuid         path   string true  "Session Uuid"
// @Success     200
func (ctl ControlVanilla) DeleteSession(ctx echo.Context) (err error) {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter, "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	// get tenant claims
	claims, err := GetServiceAuthorizationClaims(ctx)
	if err != nil {
		return HttpError(err, http.StatusForbidden)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	// get session
	session_table := sessionv3.TableNameWithTenant(claims.Hash)
	var session sessionv3.Session
	session.Uuid = uuid

	session_cond := stmt.Equal("uuid", session.Uuid)
	err = stmtex.Select(session_table, session.ColumnNames(), session_cond, nil, nil).
		QueryRowContext(ctx.Request().Context(), ctl, ctl.Dialect())(
		func(scan stmtex.Scanner) error {
			return session.Scan(scan)
		})
	if err != nil {
		return errors.Wrapf(err, "check session")
	}

	if session.Deleted.Valid {
		return ctx.JSON(http.StatusOK, OK())
	}

	//property
	time_now := time.Now()

	session.Deleted = *vanilla.NewNullTime(time_now)
	updateSet := map[string]interface{}{
		"deleted": session.Deleted,
	}

	err = func() error {
		// delete session
		affected, err := stmtex.Update(session.TableName(), updateSet, session_cond).
			ExecContext(ctx.Request().Context(), ctl, ctl.Dialect())
		if err != nil {
			return errors.Wrapf(err, "delete session")
		}
		if affected == 0 {
			return errors.Wrapf(database.ErrorNoAffected, "delete session")
		}

		return nil
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to delete session")
	}

	return ctx.JSON(http.StatusOK, OK())
}

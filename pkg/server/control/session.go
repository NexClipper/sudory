package control

import (
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	sessionv2 "github.com/NexClipper/sudory/pkg/server/model/session/v2"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Find Session
// @Description Find Session
// @Accept      json
// @Produce     json
// @Tags        server/session
// @Router      /server/session [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v2.Session
func (ctl ControlVanilla) FindSession(ctx echo.Context) (err error) {
	q, o, p, err := ParseDecoration(echoutil.QueryParam(ctx))
	err = errors.Wrapf(err, "ParseDecoration%v", logs.KVL(
		"query", echoutil.QueryParamString(ctx),
	))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	rsp := make([]sessionv2.Session, 0, state.ENV__INIT_SLICE_CAPACITY__())
	session := sessionv2.Session{}
	err = vanilla.Stmt.Select(session.TableName(), session.ColumnNames(), q, o, p).QueryRows(ctl)(func(scan vanilla.Scanner, _ int) (err error) {
		err = session.Scan(scan)
		if err == nil {
			rsp = append(rsp, session)
		}
		return
	})
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, rsp)
}

// Get Session
// @Description Get a Session
// @Accept      json
// @Produce     json
// @Tags        server/session
// @Router      /server/session/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Session Uuid"
// @Success     200 {object} v2.Session
func (ctl ControlVanilla) GetSession(ctx echo.Context) (err error) {
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
		return HttpError(err, http.StatusBadRequest)
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	session := sessionv2.Session{}
	session.Uuid = uuid
	eq_uuid := vanilla.Equal("uuid", session.Uuid)
	err = vanilla.Stmt.Select(session.TableName(), session.ColumnNames(), eq_uuid.Parse(), nil, nil).
		QueryRow(ctl)(func(scan vanilla.Scanner) (err error) {
		err = session.Scan(scan)

		return
	})
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, session)
}

// Delete Session
// @Description Delete a Session
// @Accept      json
// @Produce     json
// @Tags        server/session
// @Router      /server/session/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Session Uuid"
// @Success     200
func (ctl ControlVanilla) DeleteSession(ctx echo.Context) (err error) {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	session := sessionv2.Session{}
	session.Uuid = uuid
	session.Deleted = *vanilla.NewNullTime(time.Now())
	eq_uuid := vanilla.Equal("uuid", session.Uuid)
	updateSet := map[string]interface{}{}
	updateSet["deleted"] = session.Deleted

	err = func() (err error) {
		affected, err := vanilla.Stmt.Update(session.TableName(), updateSet, eq_uuid.Parse()).
			Exec(ctl)
		if err != nil {
			return errors.Wrapf(err, "could not found session%v", logs.KVL(
				"uuid", session.Uuid,
			))
		}

		if affected == 0 {
			return errors.Errorf("no affected")
		}

		return
	}()
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	return ctx.JSON(http.StatusOK, OK())
}

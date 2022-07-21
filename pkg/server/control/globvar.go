package control

import (
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	globvarv2 "github.com/NexClipper/sudory/pkg/server/model/global_variables/v2"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @Description Find GlobalVariables
// @Accept      json
// @Produce     json
// @Tags        server/global_variables
// @Router      /server/global_variables [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success 200 {array} v2.GlobalVariables
func (ctl ControlVanilla) FindGlobalVariables(ctx echo.Context) (err error) {
	q, o, p, err := ParseDecoration(echoutil.QueryParam(ctx))
	err = errors.Wrapf(err, "ParseDecoration%v", logs.KVL(
		"query", echoutil.QueryParamString(ctx),
	))
	if err != nil {
		return HttpError(err, http.StatusBadRequest)
	}

	rsp := make([]globvarv2.GlobalVariables, 0, state.ENV__INIT_SLICE_CAPACITY__())

	globvar := globvarv2.GlobalVariables{}
	err = vanilla.Stmt.Select(globvar.TableName(), globvar.ColumnNames(), q, o, p).
		QueryRows(ctl)(func(scan vanilla.Scanner, _ int) (err error) {
		err = globvar.Scan(scan)
		if err != nil {
			return errors.Wrapf(err, "failed to scan")
		}
		rsp = append(rsp, globvar)
		return
	})
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusOK, rsp)
}

// @Description Get a GlobalVariables
// @Accept      json
// @Produce     json
// @Tags        server/global_variables
// @Router      /server/global_variables/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "GlobalVariables Uuid"
// @Success 200 {object} v2.GlobalVariables
func (ctl ControlVanilla) GetGlobalVariables(ctx echo.Context) (err error) {
	err = echoutil.WrapHttpError(http.StatusBadRequest,
		func() (err error) {
			if len(echoutil.Param(ctx)[__UUID__]) == 0 {
				return errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%v", logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				))
			}
			return
		})
	if err != nil {
		return
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	globvar := globvarv2.GlobalVariables{}
	globvar.Uuid = uuid
	eq_uuid := vanilla.And(
		vanilla.Equal("uuid", globvar.Uuid),
		// vanilla.IsNull("deleted"),
	)
	err = vanilla.Stmt.Select(globvar.TableName(), globvar.ColumnNames(), eq_uuid.Parse(), nil, nil).
		QueryRow(ctl)(func(scan vanilla.Scanner) (err error) {
		err = globvar.Scan(scan)
		if err != nil {
			return errors.Wrapf(err, "failed to scan")
		}
		return
	})
	if err != nil {
		return
	}

	return ctx.JSON(http.StatusOK, globvar)
}

// @Description Update GlobalVariables Value
// @Accept      json
// @Produce     json
// @Tags        server/global_variables
// @Router      /server/global_variables/{uuid} [put]
// @Param       x_auth_token header string                       false "client session token"
// @Param       uuid         path   string                       true  "GlobalVariables Uuid"
// @Param       enviroment   body   v2.HttpReq_GlobalVariables_update false "HttpReq_GlobalVariables_update"
// @Success 200 {object} v2.GlobalVariables
func (ctl ControlVanilla) UpdateGlobalVariablesValue(ctx echo.Context) (err error) {
	body := new(globvarv2.HttpReq_GlobalVariables_update)
	err = echoutil.WrapHttpError(http.StatusBadRequest,
		func() error {
			err := echoutil.Bind(ctx, body)
			err = errors.Wrapf(err, "bind%s",
				logs.KVL(
					"type", TypeName(body),
				))
			return err
		},
		func() error {
			if len(echoutil.Param(ctx)[__UUID__]) == 0 {
				return errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%v", logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				))
			}
			return nil
		},
	)
	if err != nil {
		return
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	//property
	globvar := globvarv2.GlobalVariables{}
	globvar.Uuid = uuid
	eq_uuid := vanilla.And(
		vanilla.Equal("uuid", globvar.Uuid),
		vanilla.IsNull("deleted"),
	)
	globvar.Value = *vanilla.NewNullString(body.Value)
	globvar.Updated = *vanilla.NewNullTime(time.Now())

	updateSet := map[string]interface{}{}
	updateSet["value"] = globvar.Value
	updateSet["updated"] = globvar.Updated

	// update
	affected, err := vanilla.Stmt.Update(globvar.TableName(), updateSet, eq_uuid.Parse()).
		Exec(ctl)
	if err != nil {
		return
	}
	if affected == 0 {
		return errors.New("no affected")
	}

	// get
	err = vanilla.Stmt.Select(globvar.TableName(), globvar.ColumnNames(), eq_uuid.Parse(), nil, nil).
		QueryRow(ctl)(func(scan vanilla.Scanner) (err error) {
		err = globvar.Scan(scan)
		if err != nil {
			return errors.Wrapf(err, "failed to scan")
		}
		return
	})
	if err != nil {
		return errors.Wrapf(err, "not found record%v", logs.KVL(
			"table", globvar.TableName(),
		))
	}

	return ctx.JSON(http.StatusOK, globvar)
}

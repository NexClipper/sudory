package control

import (
	"net/http"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	globvarv2 "github.com/NexClipper/sudory/pkg/server/model/global_variables/v2"
	"github.com/NexClipper/sudory/pkg/server/status/state"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// @Description Find GlobalVariables
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/global_variables
// @Router      /server/global_variables [get]
// @Param       q            query  string false "query  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       o            query  string false "order  github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Param       p            query  string false "paging github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/README.md"
// @Success 200 {array} v2.GlobalVariables
func (ctl ControlVanilla) FindGlobalVariables(ctx echo.Context) error {
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
	// default pagination
	if p == nil {
		p = stmt.Limit(__DEFAULT_DECORATION_LIMIT__)
	}

	rsp := make([]globvarv2.GlobalVariables, 0, state.ENV__INIT_SLICE_CAPACITY__())

	globvar := globvarv2.GlobalVariables{}
	err = ctl.dialect.QueryRows(globvar.TableName(), globvar.ColumnNames(), q, o, p)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner, _ int) error {
			err := globvar.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			rsp = append(rsp, globvar)

			return err
		})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, rsp)
}

// @Description Get a GlobalVariables
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/global_variables
// @Router      /server/global_variables/{uuid} [get]
// @Param       uuid         path   string true  "GlobalVariables Uuid"
// @Success 200 {object} v2.GlobalVariables
func (ctl ControlVanilla) GetGlobalVariables(ctx echo.Context) error {
	err := echoutil.WrapHttpError(http.StatusBadRequest,
		func() (err error) {
			if len(echoutil.Param(ctx)[__UUID__]) == 0 {
				return errors.Wrapf(ErrorInvalidRequestParameter, "valid param%v", logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				))
			}
			return
		})
	if err != nil {
		return err
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	globvar := globvarv2.GlobalVariables{}
	globvar.Uuid = uuid
	eq_uuid := stmt.And(
		stmt.Equal("uuid", globvar.Uuid),
		stmt.IsNull("deleted"),
	)
	err = ctl.dialect.QueryRow(globvar.TableName(), globvar.ColumnNames(), eq_uuid, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := globvar.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, globvar)
}

// @Description Update GlobalVariables Value
// @Security    XAuthToken
// @Accept      json
// @Produce     json
// @Tags        server/global_variables
// @Router      /server/global_variables/{uuid} [put]
// @Param       uuid         path   string                       true  "GlobalVariables Uuid"
// @Param       enviroment   body   v2.HttpReq_GlobalVariables_update false "HttpReq_GlobalVariables_update"
// @Success 200 {object} v2.GlobalVariables
func (ctl ControlVanilla) UpdateGlobalVariablesValue(ctx echo.Context) error {
	body := new(globvarv2.HttpReq_GlobalVariables_update)
	err := echoutil.WrapHttpError(http.StatusBadRequest,
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
				return errors.Wrapf(ErrorInvalidRequestParameter, "valid param%v", logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				))
			}
			return nil
		},
	)
	if err != nil {
		return err
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	// get globvar
	var globvar globvarv2.GlobalVariables
	globvar.Uuid = uuid
	globvar_cond := stmt.And(
		stmt.Equal("uuid", globvar.Uuid),
		stmt.IsNull("deleted"),
	)
	err = ctl.dialect.QueryRow(globvar.TableName(), globvar.ColumnNames(), globvar_cond, nil, nil)(
		ctx.Request().Context(), ctl)(
		func(scan excute.Scanner) error {
			err := globvar.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return errors.WithStack(err)
	}

	//property
	updateSet := map[string]interface{}{}

	globvar.Value = *vanilla.NewNullString(body.Value)
	updateSet["value"] = globvar.Value

	globvar.Updated = *vanilla.NewNullTime(time.Now())
	updateSet["updated"] = globvar.Updated

	// update
	affected, err := ctl.dialect.Update(globvar.TableName(), updateSet, globvar_cond)(
		ctx.Request().Context(), ctl)
	if err != nil {
		err = errors.Wrapf(err, "failed to update globvar %v=%+v %v=%+v",
			"cond", globvar_cond,
			"values", updateSet,
		)
		return err
	}
	if affected == 0 {
		return errors.New("no affected")
	}

	return ctx.JSON(http.StatusOK, globvar)
}

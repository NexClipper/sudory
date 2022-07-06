package control

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	globvarv1 "github.com/NexClipper/sudory/pkg/server/model/global_variables/v1"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

// Find GlobalVariables
// @Description Find GlobalVariables
// @Accept      json
// @Produce     json
// @Tags        server/global_variables
// @Router      /server/global_variables [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success 200 {array} v1.GlobalVariables
func (ctl Control) FindGlobalVariables(ctx echo.Context) error {
	env, err := vault.NewGlobalVariables(ctl.db.Engine().NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "find environment"))
	}

	return ctx.JSON(http.StatusOK, env)
}

// Get GlobalVariables
// @Description Get a GlobalVariables
// @Accept      json
// @Produce     json
// @Tags        server/global_variables
// @Router      /server/global_variables/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "GlobalVariables 의 Uuid"
// @Success 200 {object} v1.GlobalVariables
func (ctl Control) GetGlobalVariables(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	env, err := vault.NewGlobalVariables(ctl.db.Engine().NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get environment"))
	}

	return ctx.JSON(http.StatusOK, env)
}

// Update GlobalVariables
// @Description Update GlobalVariables Value
// @Accept      json
// @Produce     json
// @Tags        server/global_variables
// @Router      /server/global_variables/{uuid} [put]
// @Param       x_auth_token header string                       false "client session token"
// @Param       uuid         path   string                       true  "GlobalVariables 의 Uuid"
// @Param       enviroment   body   v1.HttpReqGlobalVariables_update false "HttpReqGlobalVariables_update"
// @Success 200 {object} v1.GlobalVariables
func (ctl Control) UpdateGlobalVariablesValue(ctx echo.Context) error {
	update_env := new(globvarv1.HttpReqGlobalVariables_update)
	if err := echoutil.Bind(ctx, update_env); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(update_env),
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
	env := globvarv1.GlobalVariables{}
	env.Uuid = uuid
	env.Value = update_env.Value

	r, err := ctl.ScopeSession(func(tx *xorm.Session) (interface{}, error) {
		env, err := vault.NewGlobalVariables(tx).Update(env)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "update environment"))
		}

		return env, err
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

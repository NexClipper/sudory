package control

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	globvarv1 "github.com/NexClipper/sudory/pkg/server/model/global_variant/v1"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Find global_variant
// @Description Find global_variant
// @Accept      json
// @Produce     json
// @Tags        server/global_variant
// @Router      /server/global_variant [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success 200 {array} v1.GlobalVariant
func (ctl Control) FindGlobalVariant(ctx echo.Context) error {
	env, err := vault.NewEnvironment(ctl.NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "find environment"))
	}

	return ctx.JSON(http.StatusOK, env)
}

// Get global_variant
// @Description Get a global_variant
// @Accept      json
// @Produce     json
// @Tags        server/global_variant
// @Router      /server/global_variant/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "GlobalVariant 의 Uuid"
// @Success 200 {object} v1.GlobalVariant
func (ctl Control) GetGlobalVariant(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					ParamLog(__UUID__, echoutil.Param(ctx)[__UUID__])...,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	env, err := vault.NewEnvironment(ctl.NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get environment"))
	}

	return ctx.JSON(http.StatusOK, env)
}

// Update global_variant
// @Description Update global_variant Value
// @Accept      json
// @Produce     json
// @Tags        server/global_variant
// @Router      /server/global_variant/{uuid} [put]
// @Param       x_auth_token header string                       false "client session token"
// @Param       uuid         path   string                       true  "GlobalVariant 의 Uuid"
// @Param       enviroment   body   v1.HttpReqGlobalVariant_Update false "HttpReqGlobalVariant_Update"
// @Success 200 {object} v1.GlobalVariant
func (ctl Control) UpdateGlobalVariantValue(ctx echo.Context) error {
	update_env := new(globvarv1.HttpReqGlobalVariant_Update)
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
	env := globvarv1.GlobalVariant{}
	env.Uuid = uuid
	env.Value = update_env.Value

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		env, err := vault.NewEnvironment(db).Update(env)
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

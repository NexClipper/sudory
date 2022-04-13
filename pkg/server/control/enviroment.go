package control

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	envv1 "github.com/NexClipper/sudory/pkg/server/model/environment/v1"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Find Environment
// @Description Find Environment
// @Accept      json
// @Produce     json
// @Tags        server/environment
// @Router      /server/environment [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success 200 {array} v1.Environment
func (ctl Control) FindEnvironment(ctx echo.Context) error {
	env, err := vault.NewEnvironment(ctl.NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "find environment"))
	}

	return ctx.JSON(http.StatusOK, env)
}

// Get Environment
// @Description Get a Environment
// @Accept      json
// @Produce     json
// @Tags        server/environment
// @Router      /server/environment/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Environment 의 Uuid"
// @Success 200 {object} v1.Environment
func (ctl Control) GetEnvironment(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
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

// UpdateEnvironment
// @Description Update Environment Value
// @Accept      json
// @Produce     json
// @Tags        server/environment
// @Router      /server/environment/{uuid} [put]
// @Param       x_auth_token header string                       false "client session token"
// @Param       uuid         path   string                       true  "Environment 의 Uuid"
// @Param       enviroment   body   v1.HttpReqEnvironment_Update false "HttpReqEnvironment_Update"
// @Success 200 {object} v1.Environment
func (ctl Control) UpdateEnvironmentValue(ctx echo.Context) error {
	update_env := new(envv1.HttpReqEnvironment_Update)
	if err := echoutil.Bind(ctx, update_env); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(update_env),
				)))
	}

	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	//property
	env := envv1.Environment{}
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

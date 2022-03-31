package control

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
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
// @Param       q query string false "query  pkg/server/database/prepared/README.md"
// @Param       o query string false "order  pkg/server/database/prepared/README.md"
// @Param       p query string false "paging pkg/server/database/prepared/README.md"
// @Success 200 {array} v1.HttpRspEnvironment
func (ctl Control) FindEnvironment(ctx echo.Context) error {
	r, err := vault.NewEnvironment(ctl.NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			errors.Wrapf(err, "find environment%s",
				logs.KVL(
					"query", echoutil.QueryParamString(ctx),
				)))
	}

	return ctx.JSON(http.StatusOK, envv1.TransToHttpRsp(r))
}

// Get Environment
// @Description Get a Environment
// @Accept      json
// @Produce     json
// @Tags        server/environment
// @Router      /server/environment/{uuid} [get]
// @Param       uuid path string true "Environment 의 Uuid"
// @Success 200 {object} v1.HttpRspEnvironment
func (ctl Control) GetEnvironment(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := vault.NewEnvironment(ctl.NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			errors.Wrapf(err, "get environment%s",
				logs.KVL(
					"uuid", uuid,
				)))
	}

	return ctx.JSON(http.StatusOK, envv1.HttpRspEnvironment{DbSchema: *r})
}

// UpdateEnvironment
// @Description Update Environment Value
// @Accept      json
// @Produce     json
// @Tags        server/environment
// @Router      /server/environment/{uuid} [put]
// @Param       uuid       path string                      true  "Environment 의 Uuid"
// @Param       enviroment body v1.HttpReqEnvironmentUpdate false "HttpReqEnvironmentUpdate"
// @Success 200 {object} v1.HttpRspEnvironment
func (ctl Control) UpdateEnvironmentValue(ctx echo.Context) error {
	body := new(envv1.HttpReqEnvironmentUpdate)
	if err := ctx.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}
	update_env := body.EnvironmentUpdate

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		//get record
		env, err := vault.NewEnvironment(db).Get(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "get environment%s",
				logs.KVL(
					"uuid", uuid,
				))
		}

		//property
		env.Value = nullable.String(update_env.Value).Ptr() //value

		env_, err := vault.NewEnvironment(db).Update(env.Environment)
		if err != nil {
			return nil, errors.Wrapf(err, "update environment%s",
				logs.KVL(
					"environment", env,
				))
		}

		return envv1.HttpRspEnvironment{DbSchema: *env_}, err
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, r)
}

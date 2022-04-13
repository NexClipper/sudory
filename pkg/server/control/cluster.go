package control

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	clusterv1 "github.com/NexClipper/sudory/pkg/server/model/cluster/v1"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// Create Cluster
// @Description Create a cluster
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster [post]
// @Param       x_auth_token header string                   false "client session token"
// @Param       client       body   v1.HttpReqCluster_Create true  "HttpReqCluster_Create"
// @Success     200 {object} v1.Cluster
func (ctl Control) CreateCluster(ctx echo.Context) error {
	body := new(clusterv1.HttpReqCluster_Create)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if len(body.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					"param", TypeName(body.Name),
				)))
	}

	//property
	cluster := clusterv1.Cluster{}
	cluster.UuidMeta = NewUuidMeta()
	cluster.LabelMeta = NewLabelMeta(body.Name, body.Summary)
	cluster.ClusterProperty = body.ClusterProperty

	if cluster.PollingOption == nil {
		cluster.PollingOption = new(clusterv1.RagulerPollingOption).ToMap()
	}

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		cluster_, err := vault.NewCluster(db).Create(cluster)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "create cluster"))
		}
		return cluster_, err
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// Find Cluster
// @Description Find cluster
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster [get]
// @Param       x_auth_token header string false "client session token"
// @Param       q            query  string false "query  pkg/server/database/prepared/README.md"
// @Param       o            query  string false "order  pkg/server/database/prepared/README.md"
// @Param       p            query  string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.Cluster
func (ctl Control) FindCluster(ctx echo.Context) error {
	clusters, err := vault.NewCluster(ctl.NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "find cluster"))
	}

	return ctx.JSON(http.StatusOK, clusters)
}

// Get Cluster
// @Description Get a cluster
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid} [get]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Cluster 의 Uuid"
// @Success     200 {object} v1.Cluster
func (ctl Control) GetCluster(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	cluster, err := vault.NewCluster(ctl.NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
			errors.Wrapf(err, "get cluster"))
	}

	return ctx.JSON(http.StatusOK, cluster)
}

// Update Cluster
// @Description Update a cluster
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid} [put]
// @Param       x_auth_token header string                   false "client session token"
// @Param       uuid         path   string                   true  "Cluster 의 Uuid"
// @Param       client       body   v1.HttpReqCluster_Update true  "HttpReqCluster_Update"
// @Success     200 {object} v1.Cluster
func (ctl Control) UpdateCluster(ctx echo.Context) error {
	body := new(clusterv1.HttpReqCluster_Update)
	if err := echoutil.Bind(ctx, body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
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
	cluster := clusterv1.Cluster{}
	cluster.Uuid = uuid //set uuid from path
	cluster.LabelMeta = NewLabelMeta(body.Name, body.Summary)
	cluster.ClusterProperty = body.ClusterProperty

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		cluster_, err := vault.NewCluster(db).Update(cluster)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "update cluster"))
		}

		return cluster_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// UpdateClusterPollingRaguler
// @Description Update a cluster Polling Reguar
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid}/polling/raguler [put]
// @Param       x_auth_token   header string                  false "client session token"
// @Param       uuid           path   string                  true  "Cluster 의 Uuid"
// @Param       polling_option body   v1.RagulerPollingOption true  "RagulerPollingOption"
// @Success     200 {object} v1.Cluster
func (ctl Control) UpdateClusterPollingRaguler(ctx echo.Context) error {
	polling_option := new(clusterv1.RagulerPollingOption)
	if err := echoutil.Bind(ctx, polling_option); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(polling_option),
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
	cluster := clusterv1.Cluster{}
	cluster.Uuid = uuid                      //set uuid from path
	cluster.SetPollingOption(polling_option) //update polling option

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		cluster_, err := vault.NewCluster(db).Update(cluster)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "update cluster"))
		}

		return cluster_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// UpdateClusterPollingSmart
// @Description Update a cluster Polling Smart
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid}/polling/smart [put]
// @Param       x_auth_token   header string                false "client session token"
// @Param       uuid           path   string                true  "Cluster 의 Uuid"
// @Param       polling_option body   v1.SmartPollingOption true  "SmartPollingOption"
// @Success     200 {object} v1.Cluster
func (ctl Control) UpdateClusterPollingSmart(ctx echo.Context) error {
	polling_option := new(clusterv1.SmartPollingOption)
	if err := echoutil.Bind(ctx, polling_option); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(polling_option),
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
	cluster := clusterv1.Cluster{}
	cluster.Uuid = uuid                      //set uuid from path
	cluster.SetPollingOption(polling_option) //update polling option

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		cluster_, err := vault.NewCluster(db).Update(cluster)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "update cluster"))
		}

		return cluster_, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, r)
}

// Delete Cluster
// @Description Delete a cluster
// @Accept json
// @Produce json
// @Tags server/cluster
// @Router /server/cluster/{uuid} [delete]
// @Param       x_auth_token header string false "client session token"
// @Param       uuid         path   string true  "Cluster 의 Uuid"
// @Success 200
func (ctl Control) DeleteCluster(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		if err := vault.NewCluster(db).Delete(uuid); err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError).SetInternal(
				errors.Wrapf(err, "delete cluster"))
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, OK())
}

package control

import (
	"net/http"

	"github.com/NexClipper/sudory/pkg/server/control/vault"
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/echoutil"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
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
// @Param       client body v1.HttpReqCluster true "HttpReqCluster"
// @Success     200 {object} v1.HttpRspCluster
func (ctl Control) CreateCluster(ctx echo.Context) error {
	body := new(clusterv1.HttpReqCluster)
	if err := ctx.Bind(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorBindRequestObject(), "bind%s",
				logs.KVL(
					"type", TypeName(body),
				)))
	}

	if len(nullable.String(body.Name).Value()) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid param%s",
				logs.KVL(
					"param", TypeName(body.Name),
				)))
	}

	cluster := body.Cluster

	//property
	cluster.UuidMeta = NewUuidMeta()
	cluster.LabelMeta = NewLabelMeta(cluster.Name, cluster.Summary)

	if cluster.PollingOption == nil {
		cluster.PollingOption = new(clusterv1.RagulerPollingOption).ToMap()
	}

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		r, err := vault.NewCluster(db).Create(cluster)
		if err != nil {
			return nil, errors.Wrapf(err, "create cluster%s",
				logs.KVL(
					"query", echoutil.QueryParamString(ctx),
				))
		}
		return clusterv1.HttpRspCluster{DbSchema: *r}, err
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, r)
}

// Find Cluster
// @Description Find cluster
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster [get]
// @Param       q query string false "query  pkg/server/database/prepared/README.md"
// @Param       o query string false "order  pkg/server/database/prepared/README.md"
// @Param       p query string false "paging pkg/server/database/prepared/README.md"
// @Success     200 {array} v1.HttpRspCluster
func (ctl Control) FindCluster(ctx echo.Context) error {
	r, err := vault.NewCluster(ctl.NewSession()).Query(echoutil.QueryParam(ctx))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrapf(err, "find cluster%s",
			logs.KVL(
				"query", echoutil.QueryParamString(ctx),
			)))
	}

	return ctx.JSON(http.StatusOK, clusterv1.TransToHttpRsp(r))
}

// Get Cluster
// @Description Get a cluster
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid} [get]
// @Param       uuid path string true "Cluster 의 Uuid"
// @Success     200 {object} v1.HttpRspCluster
func (ctl Control) GetCluster(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := vault.NewCluster(ctl.NewSession()).Get(uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrapf(err, "get cluster%s",
			logs.KVL(
				"uuid", uuid,
			)))
	}

	return ctx.JSON(http.StatusOK, clusterv1.HttpRspCluster{DbSchema: *r})
}

// Update Cluster
// @Description Update a cluster
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid} [put]
// @Param       uuid   path string true "Cluster 의 Uuid"
// @Param       client body v1.HttpReqCluster true "HttpReqCluster"
// @Success     200 {object} v1.HttpRspCluster
func (ctl Control) UpdateCluster(ctx echo.Context) error {
	body := new(clusterv1.HttpReqCluster)
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

	cluster := body.Cluster

	uuid := echoutil.Param(ctx)[__UUID__]

	//property
	cluster.Uuid = uuid //set uuid from path

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		cluster_, err := vault.NewCluster(db).Update(cluster)
		if err != nil {
			return nil, errors.Wrapf(err, "update cluster%s",
				logs.KVL(
					"cluster", cluster,
				))
		}

		return clusterv1.HttpRspCluster{DbSchema: *cluster_}, err
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, r)
}

// UpdateClusterPollingRaguler
// @Description Update a cluster Polling Reguar
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid}/polling/raguler [put]
// @Param       uuid   path string true "Cluster 의 Uuid"
// @Param       polling_option body v1.RagulerPollingOption true "RagulerPollingOption"
// @Success     200 {object} v1.HttpRspCluster
func (ctl Control) UpdateClusterPollingRaguler(ctx echo.Context) error {
	body := new(clusterv1.RagulerPollingOption)
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

	polling_option := body

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {

		cluster, err := vault.NewCluster(db).Get(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "get cluster")
		}

		//property
		cluster.SetPollingOption(polling_option) //update polling option
		cluster.Uuid = uuid                      //set uuid from path

		cluster_, err := vault.NewCluster(db).Update(cluster.Cluster)
		if err != nil {
			return nil, errors.Wrapf(err, "update cluster%s",
				logs.KVL(
					"cluster", cluster,
				))
		}

		return clusterv1.HttpRspCluster{DbSchema: *cluster_}, err
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, r)
}

// UpdateClusterPollingSmart
// @Description Update a cluster Polling Smart
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid}/polling/smart [put]
// @Param       uuid   path string true "Cluster 의 Uuid"
// @Param       polling_option body v1.SmartPollingOption true "SmartPollingOption"
// @Success     200 {object} v1.HttpRspCluster
func (ctl Control) UpdateClusterPollingSmart(ctx echo.Context) error {
	body := new(clusterv1.SmartPollingOption)
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

	polling_option := body

	uuid := echoutil.Param(ctx)[__UUID__]

	r, err := ctl.Scope(func(db database.Context) (interface{}, error) {

		cluster, err := vault.NewCluster(db).Get(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "get cluster")
		}

		//property
		cluster.SetPollingOption(polling_option) //update polling option
		cluster.Uuid = uuid                      //set uuid from path

		cluster_, err := vault.NewCluster(db).Update(cluster.Cluster)
		if err != nil {
			return nil, errors.Wrapf(err, "update cluster%s",
				logs.KVL(
					"cluster", cluster,
				))
		}

		return clusterv1.HttpRspCluster{DbSchema: *cluster_}, err
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, r)
}

// Delete Cluster
// @Description Delete a cluster
// @Accept json
// @Produce json
// @Tags server/cluster
// @Router /server/cluster/{uuid} [delete]
// @Param uuid path string true "Cluster 의 Uuid"
// @Success 200
func (ctl Control) DeleteCluster(ctx echo.Context) error {
	if len(echoutil.Param(ctx)[__UUID__]) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest,
			errors.Wrapf(ErrorInvalidRequestParameter(), "valid%s",
				logs.KVL(
					"param", __UUID__,
				)))
	}

	uuid := echoutil.Param(ctx)[__UUID__]

	_, err := ctl.Scope(func(db database.Context) (interface{}, error) {
		err := vault.NewCluster(db).Delete(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "delete cluster%s",
				logs.KVL(
					"uuid", uuid,
				))
		}

		return nil, nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, OK())
}

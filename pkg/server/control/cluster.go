package control

import (
	"github.com/NexClipper/sudory/pkg/server/control/vault"
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
func (c *Control) CreateCluster() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {
		body := new(clusterv1.HttpReqCluster)
		if err := ctx.Bind(body); err != nil {
			return ErrorBindRequestObject(err)
		}
		if body.Name == nil {
			return ErrorInvaliedRequestParameterName("Name")
		}
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {

		body, ok := ctx.Object().(*clusterv1.HttpReqCluster)
		if !ok {
			return nil, ErrorFailedCast()
		}

		cluster := body.Cluster

		//property
		cluster.UuidMeta = NewUuidMeta()
		cluster.LabelMeta = NewLabelMeta(body.Name, body.Summary)

		//create
		cluster_, err := vault.NewCluster(ctx.Database()).Create(cluster)
		if err != nil {
			return nil, errors.Wrapf(err, "NewCluster Create")
		}

		return clusterv1.HttpRspCluster{DbSchema: *cluster_}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Contexter) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "CreateCluster binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "CreateCluster operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Lock(c.db.Engine()),
	})
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
func (c *Control) FindCluster() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		records, err := vault.NewCluster(ctx.Database()).Query(ctx.Queries())
		if err != nil {
			return nil, errors.Wrapf(err, "NewCluster Query")
		}

		return clusterv1.TransToHttpRsp(records), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Contexter) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "FindCluster binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "FindCluster operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Nolock(c.db.Engine()),
	})
}

// Get Cluster
// @Description Get a cluster
// @Accept      json
// @Produce     json
// @Tags        server/cluster
// @Router      /server/cluster/{uuid} [get]
// @Param       uuid path string true "Cluster 의 Uuid"
// @Success     200 {object} v1.HttpRspCluster
func (c *Control) GetCluster() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {
		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}

		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		uuid := ctx.Params()[__UUID__]

		cluster, err := vault.NewCluster(ctx.Database()).Get(uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "NewCluster Get")
		}
		return clusterv1.HttpRspCluster{DbSchema: *cluster}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Contexter) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "GetCluster binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "GetCluster operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Nolock(c.db.Engine()),
	})
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
func (c *Control) UpdateCluster() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {
		body := new(clusterv1.HttpReqCluster)
		if err := ctx.Bind(body); err != nil {
			return ErrorBindRequestObject(err)
		}

		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}
		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}

		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		body, ok := ctx.Object().(*clusterv1.HttpReqCluster)
		if !ok {
			return nil, ErrorFailedCast()
		}
		uuid := ctx.Params()[__UUID__]

		cluster := body.Cluster

		//set uuid from path
		cluster.Uuid = uuid
		cluster_, err := vault.NewCluster(ctx.Database()).Update(cluster)
		if err != nil {
			return nil, errors.Wrapf(err, "NewCluster Update")
		}

		return clusterv1.HttpRspCluster{DbSchema: *cluster_}, nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Contexter) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "UpdateCluster binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "UpdateCluster operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Lock(c.db.Engine()),
	})
}

// Delete Cluster
// @Description Delete a cluster
// @Accept json
// @Produce json
// @Tags server/cluster
// @Router /server/cluster/{uuid} [delete]
// @Param uuid path string true "Cluster 의 Uuid"
// @Success 200
func (c *Control) DeleteCluster() func(ctx echo.Context) error {

	binder := func(ctx Contexter) error {
		if len(ctx.Params()) == 0 {
			return ErrorInvaliedRequestParameter()
		}

		if len(ctx.Params()[__UUID__]) == 0 {
			return ErrorInvaliedRequestParameterName(__UUID__)
		}
		return nil
	}
	operator := func(ctx Contexter) (interface{}, error) {
		uuid := ctx.Params()[__UUID__]

		if err := vault.NewCluster(ctx.Database()).Delete(uuid); err != nil {
			return nil, errors.Wrapf(err, "NewCluster Delete")
		}

		return OK(), nil
	}

	return MakeMiddlewareFunc(Option{
		Binder: func(ctx Contexter) error {
			err := binder(ctx)
			if err != nil {
				return errors.Wrapf(err, "GetCluster binder")
			}
			return nil
		},
		Operator: func(ctx Contexter) (interface{}, error) {
			v, err := operator(ctx)
			if err != nil {
				return nil, errors.Wrapf(err, "GetCluster operator")
			}
			return v, nil
		},
		HttpResponsor: HttpJsonResponsor,
		Behavior:      Lock(c.db.Engine()),
	})
}

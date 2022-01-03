package operator

import (
	"github.com/NexClipper/sudory-prototype-r1/pkg/database"
	"github.com/NexClipper/sudory-prototype-r1/pkg/model"
	"github.com/labstack/echo/v4"
)

type Client struct {
	db *database.DBManipulator

	ID        uint64
	AgentID   string
	ClusterID uint64
	IP        string
	Port      int

	Response ResponseFn
}

func NewClient(d *database.DBManipulator) Operator {
	return &Client{db: d}
}

func (o *Client) toModel() *model.Client {
	m := &model.Client{
		ID:        o.ID,
		AgentID:   o.AgentID,
		ClusterID: o.ClusterID,
		IP:        o.IP,
		Port:      o.Port,
	}

	return m
}

func (o *Client) Create(ctx echo.Context) error {
	client := o.toModel()

	_, err := o.db.CreateClient(client)
	if err != nil {
		return err
	}

	if o.Response != nil {
		o.Response(ctx, nil)
	}

	return nil
}

func (o *Client) Get(ctx echo.Context) error {
	return nil
}

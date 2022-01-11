package operator

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/model"
	"github.com/labstack/echo/v4"
)

type Client struct {
	db *database.DBManipulator

	ID        uint64
	MachineID string
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
		MachineID: o.MachineID,
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

package operator

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	clientv1 "github.com/NexClipper/sudory/pkg/server/model/client/v1"
)

type Client struct {
	ctx database.Context
}

func NewClient(ctx database.Context) *Client {
	return &Client{ctx: ctx}
}

func (o *Client) Create(model clientv1.Client) error {
	err := o.ctx.CreateClient(clientv1.DbSchemaClient{Client: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *Client) Get(uuid string) (*clientv1.Client, error) {

	record, err := o.ctx.GetClient(uuid)
	if err != nil {
		return nil, err
	}

	return &record.Client, nil
}

func (o *Client) Find(where string, args ...interface{}) ([]clientv1.Client, error) {
	r, err := o.ctx.FindClient(where, args...)
	if err != nil {
		return nil, err
	}

	records := clientv1.TransFormDbSchema(r)

	return records, nil
}

func (o *Client) Update(model clientv1.Client) error {

	err := o.ctx.UpdateClient(clientv1.DbSchemaClient{Client: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *Client) Delete(uuid string) error {

	err := o.ctx.DeleteClient(uuid)
	if err != nil {
		return err
	}

	return nil
}

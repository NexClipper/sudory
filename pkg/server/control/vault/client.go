package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	clientv1 "github.com/NexClipper/sudory/pkg/server/model/client/v1"
	"github.com/pkg/errors"
)

type Client struct {
	ctx database.Context
}

func NewClient(ctx database.Context) *Client {
	return &Client{ctx: ctx}
}

func (vault Client) Create(model clientv1.Client) (*clientv1.DbSchema, error) {
	record := &clientv1.DbSchema{Client: model}
	if err := vault.ctx.Create(record); err != nil {
		return nil, errors.Wrapf(err, "database create")
	}

	return record, nil
}

func (vault Client) Get(uuid string) (*clientv1.DbSchema, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	record := &clientv1.DbSchema{}
	if err := vault.ctx.Where(where, args...).Get(record); err != nil {
		return nil, errors.Wrapf(err, "database get where=%s args=%+v", where, args)
	}

	return record, nil
}

func (vault Client) Find(where string, args ...interface{}) ([]clientv1.DbSchema, error) {
	records := make([]clientv1.DbSchema, 0)
	if err := vault.ctx.Where(where, args...).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find where=%s args=%+v", where, args)
	}

	return records, nil
}

func (vault Client) Query(query map[string]string) ([]clientv1.DbSchema, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser query=%+v", query)
	}

	//find service
	records := make([]clientv1.DbSchema, 0)
	if err := vault.ctx.Prepared(preparer).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find query=%+v", query)
	}

	return records, nil
}

func (vault Client) Update(model clientv1.Client) (*clientv1.DbSchema, error) {
	where := "uuid = ?"
	args := []interface{}{
		model.Uuid,
	}
	record := &clientv1.DbSchema{Client: model}
	if err := vault.ctx.Where(where, args...).Update(record); err != nil {
		return nil, errors.Wrapf(err, "database update where=%s args=%+v", where, args)
	}

	return record, nil
}

func (vault Client) Delete(uuid string) error {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	record := &clientv1.DbSchema{}
	if err := vault.ctx.Where(where, args...).Delete(record); err != nil {
		return errors.Wrapf(err, "database delete where=%s args=%+v", where, args)
	}

	return nil
}

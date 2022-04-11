package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	clientv1 "github.com/NexClipper/sudory/pkg/server/model/client/v1"
	"github.com/pkg/errors"
)

type Client struct {
	ctx database.Context
}

func NewClient(ctx database.Context) *Client {
	return &Client{ctx: ctx}
}

func (vault Client) Create(model clientv1.Client) (*clientv1.Client, error) {
	if err := vault.ctx.Create(&model); err != nil {
		return nil, errors.Wrapf(err, "database create")
	}

	return &model, nil
}

func (vault Client) Get(uuid string) (*clientv1.Client, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	record := &clientv1.Client{}
	if err := vault.ctx.Where(where, args...).Get(record); err != nil {
		return nil, errors.Wrapf(err, "database get%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return record, nil
}

func (vault Client) Find(where string, args ...interface{}) ([]clientv1.Client, error) {
	records := make([]clientv1.Client, 0)
	if err := vault.ctx.Where(where, args...).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return records, nil
}

func (vault Client) Query(query map[string]string) ([]clientv1.Client, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser%v",
			logs.KVL(
				"query", query,
			))
	}

	//find service
	records := make([]clientv1.Client, 0)
	if err := vault.ctx.Prepared(preparer).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"query", query,
			))
	}

	return records, nil
}

func (vault Client) Update(model clientv1.Client) (*clientv1.Client, error) {
	where := "uuid = ?"
	args := []interface{}{
		model.Uuid,
	}

	if err := vault.ctx.Where(where, args...).Update(&model); err != nil {
		return nil, errors.Wrapf(err, "database update%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return &model, nil
}

func (vault Client) Delete(uuid string) error {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	record := &clientv1.Client{}
	if err := vault.ctx.Where(where, args...).Delete(record); err != nil {
		return errors.Wrapf(err, "database delete%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return nil
}

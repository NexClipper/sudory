package vault

import (
	"sort"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
	"github.com/pkg/errors"
)

//ServiceStep
type ServiceStep struct {
	ctx database.Context
}

func NewServiceStep(ctx database.Context) *ServiceStep {
	return &ServiceStep{ctx: ctx}
}

func (vault ServiceStep) Create(model stepv1.ServiceStep) (*stepv1.DbSchema, error) {
	record := &stepv1.DbSchema{ServiceStep: model}
	if err := vault.ctx.Create(record); err != nil {
		return nil, errors.Wrapf(err, "database create")
	}

	return record, nil
}

func (vault ServiceStep) Get(uuid string) (*stepv1.DbSchema, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	record := &stepv1.DbSchema{}
	if err := vault.ctx.Where(where, args...).Get(record); err != nil {
		return nil, errors.Wrapf(err, "database get where=%s args=%+v", where, args)
	}

	return record, nil
}

func (vault ServiceStep) Find(where string, args ...interface{}) ([]stepv1.DbSchema, error) {
	records := make([]stepv1.DbSchema, 0)
	if err := vault.ctx.Where(where, args...).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find where=%s args=%+v", where, args)
	}

	return records, nil
}

func (vault ServiceStep) Update(model stepv1.ServiceStep) (*stepv1.DbSchema, error) {
	where := "uuid = ?"
	args := []interface{}{
		model.Uuid,
	}
	record := &stepv1.DbSchema{ServiceStep: model}
	if err := vault.ctx.Where(where, args...).Update(record); err != nil {
		return nil, errors.Wrapf(err, "database update where=%s args=%+v", where, args)
	}

	//make result
	record_, err := vault.Get(record.Uuid)
	if err != nil {
		return nil, errors.Wrapf(err, "make update result")
	}

	return record_, nil
}

func (vault ServiceStep) Delete(uuid string) error {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	record := &stepv1.DbSchema{}
	if err := vault.ctx.Where(where, args...).Delete(record); err != nil {
		return errors.Wrapf(err, "database delete where=%s args=%+v", where, args)
	}

	return nil
}

// ChainingSequence
//  uuid: 해당 객체는 대상에서 제외
//  대상 객체 외는 순서에 맞추어 Sequence 지정
func (vault ServiceStep) ChainingSequence(service_uuid, uuid string) error {
	where := "service_uuid = ?"
	steps, err := vault.Find(where, service_uuid)
	if err != nil {
		return errors.Wrapf(err, "Database Find")
	}

	//sort -> Sequence
	sort.Slice(steps, func(i, j int) bool {
		return nullable.Int32(steps[i].Sequence).Value() < nullable.Int32(steps[j].Sequence).Value()
	})

	seq := int32(0)
	for i := range steps {
		if steps[i].Uuid != uuid {
			steps[i].Sequence = newist.Int32(int32(seq))
		}
		seq++
	}
	for i := range steps {
		if _, err := vault.Update(steps[i].ServiceStep); err != nil {
			return errors.Wrapf(err, "Database Update")
		}
	}

	return nil
}

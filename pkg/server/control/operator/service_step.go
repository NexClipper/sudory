package operator

import (
	"sort"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
	stepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
)

//ServiceStep
type ServiceStep struct {
	ctx database.Context
}

func NewServiceStep(ctx database.Context) *ServiceStep {
	return &ServiceStep{ctx: ctx}
}

func (o *ServiceStep) Create(model stepv1.ServiceStep) error {
	err := o.ctx.CreateServiceStep(stepv1.DbSchemaServiceStep{ServiceStep: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *ServiceStep) Get(uuid string) (*stepv1.ServiceStep, error) {
	record, err := o.ctx.GetServiceStep(uuid)
	if err != nil {
		return nil, err
	}

	return &record.ServiceStep, nil
}

func (o *ServiceStep) Find(where string, args ...interface{}) ([]stepv1.ServiceStep, error) {
	r, err := o.ctx.FindServiceStep(where, args...)
	if err != nil {
		return nil, err
	}

	records := stepv1.TransFormDbSchema(r)

	return records, nil
}

func (o *ServiceStep) Update(model stepv1.ServiceStep) error {
	err := o.ctx.UpdateServiceStep(stepv1.DbSchemaServiceStep{ServiceStep: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *ServiceStep) Delete(uuid string) error {

	err := o.ctx.DeleteServiceStep(uuid)
	if err != nil {
		return err
	}

	return nil
}

// ChainingSequence
//  uuid: 해당 객체는 대상에서 제외
//  대상 객체 외는 순서에 맞추어 Sequence 지정
func (o *ServiceStep) ChainingSequence(service_uuid, uuid string) error {
	where := "service_uuid = ?"
	steps, err := o.ctx.FindServiceStep(where, service_uuid)
	if err != nil {
		return err
	}

	//sort -> Sequence
	sort.Slice(steps, func(i, j int) bool {
		return nullable.Int32(steps[i].Sequence).V() < nullable.Int32(steps[j].Sequence).V()
	})

	seq := int32(0)
	steps = map_step(steps, func(ss stepv1.DbSchemaServiceStep) stepv1.DbSchemaServiceStep {
		if ss.Uuid == uuid {
			seq++
			return ss
		}
		//Sequence
		ss.Sequence = newist.Int32(int32(seq))
		seq++
		return ss
	})

	for n := range steps {
		if err := o.ctx.UpdateServiceStep(steps[n]); err != nil {
			return err
		}
	}

	return nil
}

func map_step(elems []stepv1.DbSchemaServiceStep, mapper func(stepv1.DbSchemaServiceStep) stepv1.DbSchemaServiceStep) []stepv1.DbSchemaServiceStep {
	rst := make([]stepv1.DbSchemaServiceStep, len(elems))
	for n := range elems {
		rst[n] = mapper(elems[n])
	}
	return rst
}

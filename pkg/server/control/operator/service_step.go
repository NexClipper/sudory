package operator

import (
	"github.com/NexClipper/sudory/pkg/server/database"
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

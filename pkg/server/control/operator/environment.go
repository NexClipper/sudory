package operator

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	envv1 "github.com/NexClipper/sudory/pkg/server/model/environment/v1"
)

type Environment struct {
	ctx database.Context
}

func NewEnvironment(ctx database.Context) *Environment {
	return &Environment{ctx: ctx}
}

func (o *Environment) Create(model envv1.Environment) error {
	err := o.ctx.CreateEnvironment(envv1.DbSchemaEnvironment{Environment: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *Environment) Get(uuid string) (*envv1.Environment, error) {

	record, err := o.ctx.GetEnvironment(uuid)
	if err != nil {
		return nil, err
	}

	return &record.Environment, nil
}

func (o *Environment) Find(where string, args ...interface{}) ([]envv1.Environment, error) {
	r, err := o.ctx.FindEnvironment(where, args...)
	if err != nil {
		return nil, err
	}

	records := envv1.TransFormDbSchema(r)

	return records, nil
}

func (o *Environment) Update(model envv1.Environment) error {

	err := o.ctx.UpdateEnvironment(envv1.DbSchemaEnvironment{Environment: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *Environment) Delete(uuid string) error {

	err := o.ctx.DeleteEnvironment(uuid)
	if err != nil {
		return err
	}

	return nil
}

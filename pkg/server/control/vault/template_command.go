package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	commandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
	"github.com/pkg/errors"
)

//TemplateCommand
type TemplateCommand struct {
	// ctx *database.DBManipulator
	ctx database.Context
}

// func NewTemplateCommand(d *database.DBManipulator) *TemplateCommand {
// 	return &TemplateCommand{db: d}
// }
func NewTemplateCommand(ctx database.Context) *TemplateCommand {
	return &TemplateCommand{ctx: ctx}
}

func (vault TemplateCommand) Create(model commandv1.TemplateCommand) (*commandv1.DbSchema, error) {
	record := &commandv1.DbSchema{TemplateCommand: model}
	if err := vault.ctx.Create(record); err != nil {
		return nil, errors.Wrapf(err, "database create")
	}

	return record, nil
}

func (vault TemplateCommand) Get(uuid string) (*commandv1.DbSchema, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	record := &commandv1.DbSchema{}
	if err := vault.ctx.Where(where, args...).Get(record); err != nil {
		return nil, errors.Wrapf(err, "database get where=%s args=%+v", where, args)
	}

	return record, nil
}

func (vault TemplateCommand) Find(where string, args ...interface{}) ([]commandv1.DbSchema, error) {
	commands := make([]commandv1.DbSchema, 0)
	if err := vault.ctx.Where(where, args...).Find(&commands); err != nil {
		return nil, errors.Wrapf(err, "database find where=%s args=%+v", where, args)
	}

	return commands, nil
}

func (vault TemplateCommand) Query(query map[string]string) ([]commandv1.DbSchema, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser query=%+v", query)
	}

	//find service
	records := make([]commandv1.DbSchema, 0)
	if err := vault.ctx.Prepared(preparer).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find query=%+v", query)
	}

	return records, nil
}

func (vault TemplateCommand) Update(model commandv1.TemplateCommand) (*commandv1.DbSchema, error) {
	where := "uuid = ?"
	args := []interface{}{
		model.Uuid,
	}
	record := &commandv1.DbSchema{TemplateCommand: model}
	if err := vault.ctx.Where(where, args...).Update(record); err != nil {
		return nil, errors.Wrapf(err, "database update where=%s args=%+v", where, args)
	}

	return record, nil
}

func (vault TemplateCommand) Delete(uuid string) error {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	record := &commandv1.DbSchema{}
	if err := vault.ctx.Where(where, args...).Delete(record); err != nil {
		return errors.Wrapf(err, "database delete where=%s args=%+v", where, args)
	}

	return nil
}

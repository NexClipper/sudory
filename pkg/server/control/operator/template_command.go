package operator

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	tcommandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
)

//TemplateCommand
type TemplateCommand struct {
	db *database.DBManipulator
}

func NewTemplateCommand(d *database.DBManipulator) *TemplateCommand {
	return &TemplateCommand{db: d}
}

func (o *TemplateCommand) Create(model tcommandv1.TemplateCommand) error {

	err := o.db.CreateTemplateCommand(tcommandv1.DbSchemaTemplateCommand{TemplateCommand: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *TemplateCommand) Get(uuid string) (*tcommandv1.TemplateCommand, error) {

	record, err := o.db.GetTemplateCommand(uuid)
	if err != nil {
		return nil, err
	}

	return &record.TemplateCommand, nil
}

func (o *TemplateCommand) Find(where string, args ...interface{}) ([]tcommandv1.TemplateCommand, error) {

	r, err := o.db.FindTemplateCommand(where, args...)
	if err != nil {
		return nil, err
	}

	records := tcommandv1.TransFromDbSchema(r)

	return records, nil
}

func (o *TemplateCommand) Update(model tcommandv1.TemplateCommand) error {

	err := o.db.UpdateTemplateCommand(tcommandv1.DbSchemaTemplateCommand{TemplateCommand: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *TemplateCommand) Delete(uuid string) error {

	err := o.db.DeleteTemplateCommand(uuid)
	if err != nil {
		return err
	}

	return nil
}

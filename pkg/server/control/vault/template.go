package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/pkg/errors"

	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
	commandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
)

//Template
type Template struct {
	// ctx *database.DBManipulator
	ctx database.Context
}

func NewTemplate(ctx database.Context) *Template {
	return &Template{ctx: ctx}
}

func (vault Template) Create(model templatev1.TemplateWithCommands) (*templatev1.DbSchemaTemplateWithCommands, error) {
	//create service
	record_template := &templatev1.DbSchema{Template: model.Template}
	if err := vault.ctx.Create(record_template); err != nil {
		return nil, errors.Wrapf(err, "database create")
	}
	//create steps
	record_commands := make([]commandv1.DbSchema, len(model.Commands))
	for i := range model.Commands {
		record, err := NewTemplateCommand(vault.ctx).Create(model.Commands[i])
		if err != nil {
			return nil, errors.Wrapf(err, "NewServiceStep Create")
		}
		record_commands[i] = *record
	}

	return &templatev1.DbSchemaTemplateWithCommands{DbSchema: *record_template, Commands: record_commands}, nil
}

func (vault Template) Get(uuid string) (*templatev1.DbSchemaTemplateWithCommands, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	template := new(templatev1.DbSchema)
	if err := vault.ctx.Where(where, args...).Get(template); err != nil {
		return nil, errors.Wrapf(err, "database get where=%s args=%+v", where, args)
	}

	where = "template_uuid = ?"
	args = []interface{}{
		uuid,
	}
	commands := make([]commandv1.DbSchema, 0)
	if err := vault.ctx.Where(where, args...).Find(&commands); err != nil {
		return nil, errors.Wrapf(err, "database get where=%s args=%+v", where, args)
	}

	return &templatev1.DbSchemaTemplateWithCommands{DbSchema: *template, Commands: commands}, nil
}

func (vault Template) Find(where string, args ...interface{}) ([]templatev1.DbSchemaTemplateWithCommands, error) {
	//find template
	records := make([]templatev1.DbSchema, 0)
	if err := vault.ctx.Where(where, args...).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find where=%s args=%+v", where, args)
	}
	//make result
	var templates = make([]templatev1.DbSchemaTemplateWithCommands, len(records))
	for i := range records {
		template := records[i]
		//set template
		templates[i].DbSchema = template
		//find command
		commands, err := NewTemplateCommand(vault.ctx).Find("template_uuid = ?", template.Uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "NewTemplateCommand Find")
		}
		//set commands
		templates[i].Commands = commands
	}

	return templates, nil
}

func (vault Template) Query(query map[string]string) ([]templatev1.DbSchemaTemplateWithCommands, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser query=%+v", query)
	}

	//find service
	records := make([]templatev1.DbSchema, 0)
	if err := vault.ctx.Prepared(preparer).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find query=%+v", query)
	}

	//make result
	var templates = make([]templatev1.DbSchemaTemplateWithCommands, len(records))
	for i := range records {
		template := records[i]
		//set template
		templates[i].DbSchema = template
		//find command
		commands, err := NewTemplateCommand(vault.ctx).Find("template_uuid = ?", template.Uuid)
		if err != nil {
			return nil, errors.Wrapf(err, "NewTemplateCommand Find")
		}
		//set commands
		templates[i].Commands = commands
	}

	return templates, nil
}

func (vault Template) Update(model templatev1.Template) (*templatev1.DbSchema, error) {
	where := "uuid = ?"
	args := []interface{}{
		model.Uuid,
	}
	record := &templatev1.DbSchema{Template: model}
	if err := vault.ctx.Where(where, args...).Update(record); err != nil {
		return nil, errors.Wrapf(err, "database update where=%s args=%+v", where, args)
	}

	//make result
	record_, err := vault.Get(record.Uuid)
	if err != nil {
		return nil, errors.Wrapf(err, "make update result")
	}

	return &record_.DbSchema, nil
}

func (vault Template) Delete(uuid string) error {
	//delete template
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	template := &templatev1.DbSchema{}
	if err := vault.ctx.Where(where, args...).Delete(template); err != nil {
		return errors.Wrapf(err, "database delete where=%s args=%+v", where, args)
	}
	//delete command
	where = "template_uuid = ?"
	args = []interface{}{
		uuid,
	}
	command := &commandv1.DbSchema{}
	if err := vault.ctx.Where(where, args...).Delete(command); err != nil {
		return errors.Wrapf(err, "database delete where=%s args=%+v", where, args)
	}

	return nil
}

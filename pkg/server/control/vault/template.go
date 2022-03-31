package vault

import (
	"sort"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
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
	template := &templatev1.DbSchema{Template: model.Template}
	if err := vault.ctx.Create(template); err != nil {
		return nil, errors.Wrapf(err, "database create")
	}
	//create steps
	commands := make([]commandv1.DbSchema, len(model.Commands))
	for i := range model.Commands {
		command := &commandv1.DbSchema{TemplateCommand: model.Commands[i]}
		if err := vault.ctx.Create(command); err != nil {
			return nil, errors.Wrapf(err, "database create")
		}
		commands[i] = *command
	}

	return &templatev1.DbSchemaTemplateWithCommands{DbSchema: *template, Commands: commands}, nil
}

func (vault Template) Get(uuid string) (*templatev1.DbSchemaTemplateWithCommands, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	template := new(templatev1.DbSchema)
	if err := vault.ctx.Where(where, args...).Get(template); err != nil {
		return nil, errors.Wrapf(err, "database get%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	where = "template_uuid = ?"
	args = []interface{}{
		uuid,
	}
	commands := make([]commandv1.DbSchema, 0)
	if err := vault.ctx.Where(where, args...).Find(&commands); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	//sort -> Sequence ASC
	sort.Slice(commands, func(i, j int) bool {
		return nullable.Int32(commands[i].Sequence).Value() < nullable.Int32(commands[j].Sequence).Value()
	})

	return &templatev1.DbSchemaTemplateWithCommands{DbSchema: *template, Commands: commands}, nil
}

func (vault Template) Find(where string, args ...interface{}) ([]templatev1.DbSchemaTemplateWithCommands, error) {
	//find template
	records := make([]templatev1.DbSchema, 0)
	if err := vault.ctx.Where(where, args...).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}
	//make result
	var templates = make([]templatev1.DbSchemaTemplateWithCommands, len(records))
	for i := range records {
		template := records[i]
		//set template
		templates[i].DbSchema = template
		//find command
		where := "template_uuid = ?"
		args := []interface{}{
			template.Uuid,
		}
		commands := make([]commandv1.DbSchema, 0)
		if err := vault.ctx.Where(where, args...).Find(&commands); err != nil {
			return nil, errors.Wrapf(err, "database find%v",
				logs.KVL(
					"where", where,
					"args", args,
				))
		}
		//sort -> Sequence ASC
		sort.Slice(commands, func(i, j int) bool {
			return nullable.Int32(commands[i].Sequence).Value() < nullable.Int32(commands[j].Sequence).Value()
		})

		//set commands
		templates[i].Commands = commands
	}

	return templates, nil
}

func (vault Template) Query(query map[string]string) ([]templatev1.DbSchemaTemplateWithCommands, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser%v",
			logs.KVL(
				"query", query,
			))
	}

	//find service
	records := make([]templatev1.DbSchema, 0)
	if err := vault.ctx.Prepared(preparer).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"query", query,
			))
	}

	//make result
	var templates = make([]templatev1.DbSchemaTemplateWithCommands, len(records))
	for i := range records {
		template := records[i]
		//set template
		templates[i].DbSchema = template
		//find command
		where := "template_uuid = ?"
		args := []interface{}{
			template.Uuid,
		}
		commands := make([]commandv1.DbSchema, 0)
		if err := vault.ctx.Where(where, args...).Find(&commands); err != nil {
			return nil, errors.Wrapf(err, "database find%v",
				logs.KVL(
					"where", where,
					"args", args,
				))
		}
		//sort -> Sequence ASC
		sort.Slice(commands, func(i, j int) bool {
			return nullable.Int32(commands[i].Sequence).Value() < nullable.Int32(commands[j].Sequence).Value()
		})

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
		return nil, errors.Wrapf(err, "database update%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return record, nil
}

func (vault Template) Delete(uuid string) error {
	//delete template
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	template := &templatev1.DbSchema{}
	if err := vault.ctx.Where(where, args...).Delete(template); err != nil {
		return errors.Wrapf(err, "database delete%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}
	//delete command
	where = "template_uuid = ?"
	args = []interface{}{
		uuid,
	}
	command := &commandv1.DbSchema{}
	if err := vault.ctx.Where(where, args...).Delete(command); err != nil {
		return errors.Wrapf(err, "database delete%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return nil
}

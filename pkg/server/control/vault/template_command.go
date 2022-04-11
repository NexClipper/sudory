package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
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

func (vault TemplateCommand) Create(model commandv1.TemplateCommand) (*commandv1.TemplateCommand, error) {
	if err := vault.ctx.Create(&model); err != nil {
		return nil, errors.Wrapf(err, "database create")
	}

	return &model, nil
}

func (vault TemplateCommand) Get(uuid string) (*commandv1.TemplateCommand, error) {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &commandv1.TemplateCommand{}
	if err := vault.ctx.Where(where, args...).Get(model); err != nil {
		return nil, errors.Wrapf(err, "database get%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return model, nil
}

func (vault TemplateCommand) Find(where string, args ...interface{}) ([]commandv1.TemplateCommand, error) {
	models := make([]commandv1.TemplateCommand, 0)
	if err := vault.ctx.Where(where, args...).Find(&models); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return models, nil
}

func (vault TemplateCommand) Query(query map[string]string) ([]commandv1.TemplateCommand, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser%v",
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]commandv1.TemplateCommand, 0)
	if err := vault.ctx.Prepared(preparer).Find(&models); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"query", query,
			))
	}

	return models, nil
}

func (vault TemplateCommand) Update(model commandv1.TemplateCommand) (*commandv1.TemplateCommand, error) {
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

func (vault TemplateCommand) Delete(uuid string) error {
	where := "uuid = ?"
	args := []interface{}{
		uuid,
	}
	model := &commandv1.TemplateCommand{}
	if err := vault.ctx.Where(where, args...).Delete(model); err != nil {
		return errors.Wrapf(err, "database delete%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return nil
}

func (vault TemplateCommand) Delete_ByTemplate(template_uuid string) error {
	where := "template_uuid = ?"
	args := []interface{}{
		template_uuid,
	}
	model := &commandv1.TemplateCommand{}
	if err := vault.ctx.Where(where, args...).Delete(model); err != nil {
		return errors.Wrapf(err, "database delete%v",
			logs.KVL(
				"where", where,
				"args", args,
			))
	}

	return nil
}

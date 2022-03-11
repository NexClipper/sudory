package vault

import (
	"sort"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
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
	records := make([]commandv1.DbSchema, 0)
	if err := vault.ctx.Where(where, args...).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find where=%s args=%+v", where, args)
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

	//make result
	record_, err := vault.Get(record.Uuid)
	if err != nil {
		return nil, errors.Wrapf(err, "make update result")
	}

	return record_, nil
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

// ChainingSequence
//  uuid: 해당 객체는 대상에서 제외
//  대상 객체 외는 순서에 맞추어 Sequence 지정
func (vault TemplateCommand) ChainingSequence(template_uuid, uuid string) error {
	where := "template_uuid = ?"
	args := []interface{}{
		template_uuid,
	}
	commands, err := vault.Find(where, args...)
	if err != nil {
		return errors.Wrapf(err, "Database Find")
	}

	//sort -> Sequence
	sort.Slice(commands, func(i, j int) bool {
		return nullable.Int32(commands[i].Sequence).Value() < nullable.Int32(commands[j].Sequence).Value()
	})

	seq := int32(0)
	for i := range commands {
		if commands[i].Uuid != uuid {
			commands[i].Sequence = newist.Int32(int32(seq))
		}
		seq++
	}
	for i := range commands {
		if _, err := vault.Update(commands[i].TemplateCommand); err != nil {
			return errors.Wrapf(err, "Database Update")
		}
	}

	return nil
}

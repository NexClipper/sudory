package operator

import (
	"sort"

	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	"github.com/NexClipper/sudory/pkg/server/macro/nullable"
	tcommandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
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

func (o *TemplateCommand) Create(model tcommandv1.TemplateCommand) error {

	err := o.ctx.CreateTemplateCommand(tcommandv1.DbSchemaTemplateCommand{TemplateCommand: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *TemplateCommand) Get(uuid string) (*tcommandv1.TemplateCommand, error) {

	record, err := o.ctx.GetTemplateCommand(uuid)
	if err != nil {
		return nil, err
	}

	return &record.TemplateCommand, nil
}

func (o *TemplateCommand) Find(where string, args ...interface{}) ([]tcommandv1.TemplateCommand, error) {

	r, err := o.ctx.FindTemplateCommand(where, args...)
	if err != nil {
		return nil, err
	}

	records := tcommandv1.TransFromDbSchema(r)

	return records, nil
}

func (o *TemplateCommand) Update(model tcommandv1.TemplateCommand) error {

	err := o.ctx.UpdateTemplateCommand(tcommandv1.DbSchemaTemplateCommand{TemplateCommand: model})
	if err != nil {
		return err
	}

	return nil
}

func (o *TemplateCommand) Delete(uuid string) error {

	err := o.ctx.DeleteTemplateCommand(uuid)
	if err != nil {
		return err
	}

	return nil
}

// ChainingSequence
//  uuid: 해당 객체는 대상에서 제외
//  대상 객체 외는 순서에 맞추어 Sequence 지정
func (o *TemplateCommand) ChainingSequence(template_uuid, uuid string) error {
	where := "template_uuid = ?"
	commands, err := o.ctx.FindTemplateCommand(where, template_uuid)
	if err != nil {
		return err
	}

	//sort -> Sequence
	sort.Slice(commands, func(i, j int) bool {
		return nullable.Int32(commands[i].Sequence).V() < nullable.Int32(commands[j].Sequence).V()
	})

	seq := int32(0)
	commands = map_command(commands, func(ss tcommandv1.DbSchemaTemplateCommand) tcommandv1.DbSchemaTemplateCommand {
		if ss.Uuid == uuid {
			seq++
			return ss
		}
		//Sequence
		ss.Sequence = newist.Int32(int32(seq))
		seq++
		return ss
	})

	for n := range commands {
		if err := o.ctx.UpdateTemplateCommand(commands[n]); err != nil {
			return err
		}
	}

	return nil
}

func map_command(elems []tcommandv1.DbSchemaTemplateCommand, mapper func(tcommandv1.DbSchemaTemplateCommand) tcommandv1.DbSchemaTemplateCommand) []tcommandv1.DbSchemaTemplateCommand {
	rst := make([]tcommandv1.DbSchemaTemplateCommand, len(elems))
	for n := range elems {
		rst[n] = mapper(elems[n])
	}
	return rst
}

package vault

import (
	"context"
	"sort"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/model/template/v2"
	"github.com/pkg/errors"
)

// import (
// 	"github.com/NexClipper/sudory/pkg/server/database"
// 	"github.com/NexClipper/sudory/pkg/server/database/prepare"
// 	"github.com/NexClipper/sudory/pkg/server/macro/logs"
// 	"github.com/pkg/errors"

// 	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
// )

// //Template
// type Template struct {
// 	// ctx *database.DBManipulator
// 	ctx database.Context
// }

// func NewTemplate(ctx database.Context) *Template {
// 	return &Template{ctx: ctx}
// }

// func (vault Template) Create(model templatev1.Template) (*templatev1.Template, error) {
// 	//create service
// 	if err := vault.ctx.Create(&model); err != nil {
// 		return nil, errors.Wrapf(err, "database create")
// 	}

// 	return &model, nil
// }

// func (vault Template) Get(uuid string) (*templatev1.Template, error) {
// 	where := "uuid = ?"
// 	args := []interface{}{
// 		uuid,
// 	}
// 	model := new(templatev1.Template)
// 	if err := vault.ctx.Where(where, args...).Get(model); err != nil {
// 		return nil, errors.Wrapf(err, "database get%v",
// 			logs.KVL(
// 				"where", where,
// 				"args", args,
// 			))
// 	}

// 	return model, nil
// }

// func (vault Template) Find(where string, args ...interface{}) ([]templatev1.Template, error) {
// 	//find template
// 	models := make([]templatev1.Template, 0)
// 	if err := vault.ctx.Where(where, args...).Find(&models); err != nil {
// 		return nil, errors.Wrapf(err, "database find%v",
// 			logs.KVL(
// 				"where", where,
// 				"args", args,
// 			))
// 	}

// 	return models, nil
// }

// func (vault Template) Query(query map[string]string) ([]templatev1.Template, error) {
// 	//parse query
// 	preparer, err := prepare.NewParser(query)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "prepare newParser%v",
// 			logs.KVL(
// 				"query", query,
// 			))
// 	}

// 	//find service
// 	models := make([]templatev1.Template, 0)
// 	if err := vault.ctx.Prepared(preparer).Find(&models); err != nil {
// 		return nil, errors.Wrapf(err, "database find%v",
// 			logs.KVL(
// 				"query", query,
// 			))
// 	}

// 	return models, nil
// }

// func (vault Template) Update(model templatev1.Template) (*templatev1.Template, error) {
// 	where := "uuid = ?"
// 	args := []interface{}{
// 		model.Uuid,
// 	}
// 	if err := vault.ctx.Where(where, args...).Update(&model); err != nil {
// 		return nil, errors.Wrapf(err, "database update%v",
// 			logs.KVL(
// 				"where", where,
// 				"args", args,
// 			))
// 	}

// 	return &model, nil
// }

// func (vault Template) Delete(uuid string) error {
// 	//delete template
// 	where := "uuid = ?"
// 	args := []interface{}{
// 		uuid,
// 	}

// 	model := &templatev1.Template{}
// 	if err := vault.ctx.Where(where, args...).Delete(model); err != nil {
// 		return errors.Wrapf(err, "database delete%v",
// 			logs.KVL(
// 				"where", where,
// 				"args", args,
// 			))
// 	}

// 	return nil
// }

func GetTemplate(ctx context.Context, tx excute.Preparer, dialect excute.SqlExcutor,
	template_uuid string,
) (*template.Template, []template.TemplateCommand, error) {

	// get template
	tmpl_cond := stmt.And(
		stmt.Equal("uuid", template_uuid),
		stmt.IsNull("deleted"),
	)

	tmpl := template.Template{}
	err := dialect.QueryRow(tmpl.TableName(), tmpl.ColumnNames(), tmpl_cond, nil, nil)(ctx, tx)(
		func(scan excute.Scanner) error {
			err := tmpl.Scan(scan)
			err = errors.WithStack(err)

			return err
		})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get template")
	}

	// get commands
	commands := make([]template.TemplateCommand, 0)
	command_cond := stmt.And(
		stmt.Equal("template_uuid", template_uuid),
		stmt.IsNull("deleted"),
	)

	var command template.TemplateCommand
	err = dialect.QueryRow(command.TableName(), command.ColumnNames(), command_cond, nil, nil)(ctx, tx)(
		func(scan excute.Scanner) error {
			err := command.Scan(scan)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}

			commands = append(commands, command)

			return err
		})
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to get template commands")
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Sequence < commands[j].Sequence
	})

	return &tmpl, commands, nil
}

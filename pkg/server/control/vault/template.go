package vault

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

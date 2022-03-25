package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/pkg/errors"

	recipev1 "github.com/NexClipper/sudory/pkg/server/model/template_recipe/v1"
)

//TemplateRecipe
type TemplateRecipe struct {
	// ctx *database.DBManipulator
	ctx database.Context
}

func NewTemplateRecipe(ctx database.Context) *TemplateRecipe {
	return &TemplateRecipe{ctx: ctx}
}

func (vault TemplateRecipe) Query(query map[string]string) ([]recipev1.DbSchema, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser query=%+v", query)
	}

	//find service
	records := make([]recipev1.DbSchema, 0)
	if err := vault.ctx.Prepared(preparer).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find query=%+v", query)
	}

	return records, nil
}

func (vault TemplateRecipe) Prepare(condition map[string]interface{}) ([]recipev1.DbSchema, error) {
	//parse cond
	preparer, err := prepare.NewConditionMap(condition)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser condition=%+v", condition)
	}

	//find service
	records := make([]recipev1.DbSchema, 0)
	if err := vault.ctx.Prepared(preparer).Find(&records); err != nil {
		return nil, errors.Wrapf(err, "database find condition=%+v", preparer)
	}

	return records, nil
}

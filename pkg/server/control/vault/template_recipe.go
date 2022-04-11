package vault

import (
	"github.com/NexClipper/sudory/pkg/server/database"
	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
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

func (vault TemplateRecipe) Query(query map[string]string) ([]recipev1.TemplateRecipe, error) {
	//parse query
	preparer, err := prepare.NewParser(query)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser%v",
			logs.KVL(
				"query", query,
			))
	}

	//find service
	models := make([]recipev1.TemplateRecipe, 0)
	if err := vault.ctx.Prepared(preparer).Find(&models); err != nil {
		return nil, errors.Wrapf(err, "database find%v",
			logs.KVL(
				"query", query,
			))
	}

	return models, nil
}

func (vault TemplateRecipe) Prepare(condition map[string]interface{}) ([]recipev1.TemplateRecipe, error) {
	//parse cond
	preparer, err := prepare.NewConditionMap(condition)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare newParser condition=%+v", condition)
	}

	//find service
	models := make([]recipev1.TemplateRecipe, 0)
	if err := vault.ctx.Prepared(preparer).Find(&models); err != nil {
		return nil, errors.Wrapf(err, "database find condition=%+v", preparer)
	}

	return models, nil
}

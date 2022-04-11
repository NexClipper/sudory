package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

// TemplateRecipeProperty
type TemplateRecipeProperty struct {
	Method *string `json:"method,omitempty" xorm:"'method' varchar(255) notnull unique(recipe_value) comment('method')"`
	Args   *string `json:"args,omitempty"   xorm:"'args'   varchar(255) notnull unique(recipe_value) comment('args')"`
}

// DATABASE SCHEMA: TEMPLATE_RECIPE
type TemplateRecipe struct {
	metav1.DbMeta          `json:",inline" xorm:"extends"`
	metav1.LabelMeta       `json:",inline" xorm:"extends"` //inline labelmeta
	TemplateRecipeProperty `json:",inline" xorm:"extends"` //inline property
}

func (TemplateRecipe) TableName() string {
	return "template_recipe"
}

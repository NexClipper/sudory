package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

// TemplateRecipeProperty
type TemplateRecipeProperty struct {
	//method
	Method *string `json:"method,omitempty" xorm:"varchar(255) null 'method' comment('method')"`
	//arguments
	Args *string `json:"args,omitempty" xorm:"varchar(255) null 'args' comment('args')"`
}

// MODEL: TEMPLATE_RECIPE
type TemplateRecipe struct {
	metav1.LabelMeta       `json:",inline" xorm:"extends"` //inline labelmeta
	TemplateRecipeProperty `json:",inline" xorm:"extends"` //inline property
}

// DATABASE SCHEMA: TEMPLATE_RECIPE
type DbSchema struct {
	metav1.DbMeta  `xorm:"extends"`
	TemplateRecipe `xorm:"extends"`
}

func (DbSchema) TableName() string {
	return "template_recipe"
}

// HTTP RESPONSE BODY: TEMPLATE_RECIPE
type HttpRspTemplateRecipe struct {
	DbSchema `json:",inline"`
}

//변환 Token -> HttpRsp
func TransToHttpRsp(s []DbSchema) []HttpRspTemplateRecipe {
	var out = make([]HttpRspTemplateRecipe, len(s))
	for n, it := range s {
		out[n].DbSchema = it
	}
	return out
}

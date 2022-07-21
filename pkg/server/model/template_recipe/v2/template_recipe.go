package v2

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type TemplateRecipe struct {
	ID      int64              `column:"id"      json:"id,omitempty"` // pk
	Name    string             `column:"name"    json:"name,omitempty"`
	Summary vanilla.NullString `column:"summary" json:"summary,omitempty" swaggertype:"string"`
	Method  string             `column:"method"  json:"method,omitempty"`
	Args    string             `column:"args"    json:"args,omitempty"`
	Created time.Time          `column:"created" json:"created,omitempty"`
	Updated vanilla.NullTime   `column:"updated" json:"updated,omitempty" swaggertype:"string"`
	Deleted vanilla.NullTime   `column:"deleted" json:"deleted,omitempty" swaggertype:"string"`
}

func (TemplateRecipe) TableName() string {
	return "template_recipe"
}

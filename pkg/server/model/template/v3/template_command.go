package template

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type TemplateCommand struct {
	Name          string             `column:"name"           json:"name,omitempty"` // pk
	Summary       vanilla.NullString `column:"summary"        json:"summary,omitempty"        swaggertype:"string"`
	Inputs        string             `column:"inputs"         json:"inputs,omitempty"`
	Outputs       string             `column:"outputs"        json:"outputs,omitempty"`
	ClientVersion int                `column:"client_version" json:"client_version,omitempty"`
	Category      string             `column:"category"       json:"category,omitempty"`
	Created       time.Time          `column:"created"        json:"created,omitempty"`
	Updated       vanilla.NullTime   `column:"updated"        json:"updated,omitempty"        swaggertype:"string"`
	Deleted       vanilla.NullTime   `column:"deleted"        json:"deleted,omitempty"        swaggertype:"string"`
}

func (TemplateCommand) TableName() string {
	return "template_command_v2"
}

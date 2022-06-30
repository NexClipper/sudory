package v2

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type TemplateCommand_essential struct {
	Name         string             `column:"name"          json:"name,omitempty"`
	Summary      vanilla.NullString `column:"summary"       json:"summary,omitempty"`
	TemplateUuid string             `column:"template_uuid" json:"template_uuid"`
	Sequence     vanilla.NullInt    `column:"sequence"      json:"sequence,omitempty"`
	Method       vanilla.NullString `column:"method"        json:"method,omitempty"`
	Args         vanilla.NullObject `column:"args"          json:"args,omitempty"`
	ResultFilter vanilla.NullString `column:"result_filter" json:"result_filter,omitempty"`
}

func (TemplateCommand_essential) TableName() string {
	return "template_command"
}

type TemplateCommand struct {
	// Id   uint64 `column:"id"   json:"id,omitempty"`   //pk
	Uuid string `column:"uuid" json:"uuid,omitempty"` //pk

	TemplateCommand_essential `json:",inline"`

	Created time.Time        `column:"created" json:"created,omitempty"`
	Updated vanilla.NullTime `column:"updated" json:"updated,omitempty"`
	Deleted vanilla.NullTime `column:"deleted" json:"deleted,omitempty"`
}

package v2

import (
	"time"

	noxorm "github.com/NexClipper/sudory/pkg/server/model/noxorm/v2"
)

type TemplateCommand_essential struct {
	Name         string            `column:"name"          json:"name,omitempty"`
	Summary      noxorm.NullString `column:"summary"       json:"summary,omitempty"`
	TemplateUuid string            `column:"template_uuid" json:"template_uuid"`
	Sequence     noxorm.NullInt    `column:"sequence"      json:"sequence,omitempty"`
	Method       noxorm.NullString `column:"method"        json:"method,omitempty"`
	Args         noxorm.NullJson   `column:"args"          json:"args,omitempty"`
	ResultFilter noxorm.NullString `column:"result_filter" json:"result_filter,omitempty"`
}

func (TemplateCommand_essential) TableName() string {
	return "template_command"
}

type TemplateCommand struct {
	// Id   uint64 `column:"id"   json:"id,omitempty"`   //pk
	Uuid string `column:"uuid" json:"uuid,omitempty"` //pk

	TemplateCommand_essential `json:",inline"`

	Created time.Time       `column:"created" json:"created,omitempty"`
	Updated noxorm.NullTime `column:"updated" json:"updated,omitempty"`
	Deleted noxorm.NullTime `column:"deleted" json:"deleted,omitempty"`
}

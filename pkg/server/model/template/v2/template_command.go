package v2

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type TemplateCommand struct {
	Id           uint64             `column:"id"            json:"id,omitempty"`   // pk
	Uuid         string             `column:"uuid"          json:"uuid,omitempty"` // uuid
	Name         string             `column:"name"          json:"name,omitempty"`
	Summary      vanilla.NullString `column:"summary"       json:"summary,omitempty"       swaggertype:"string"`
	TemplateUuid string             `column:"template_uuid" json:"template_uuid"`
	Sequence     int                `column:"sequence"      json:"sequence,omitempty"`
	Method       string             `column:"method"        json:"method,omitempty"`
	Args         vanilla.NullObject `column:"args"          json:"args,omitempty"          swaggertype:"object"`
	ResultFilter vanilla.NullString `column:"result_filter" json:"result_filter,omitempty" swaggertype:"string"`
	Created      time.Time          `column:"created"       json:"created,omitempty"`
	Updated      vanilla.NullTime   `column:"updated"       json:"updated,omitempty"       swaggertype:"string"`
	Deleted      vanilla.NullTime   `column:"deleted"       json:"deleted,omitempty"       swaggertype:"string"`
}

func (TemplateCommand) TableName() string {
	return "template_command"
}

package template

import (
	"encoding/json"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type Template struct {
	Uuid    string             `column:"uuid"    json:"uuid,omitempty"` // pk
	Name    string             `column:"name"    json:"name,omitempty"`
	Summary vanilla.NullString `column:"summary" json:"summary,omitempty" swaggertype:"string"`
	Flow    string             `column:"flow"    json:"flow,omitempty"`
	Inputs  string             `column:"inputs"  json:"inputs,omitempty"`
	Origin  string             `column:"origin"  json:"origin,omitempty"`
	Created time.Time          `column:"created" json:"created,omitempty"`
	Updated vanilla.NullTime   `column:"updated" json:"updated,omitempty" swaggertype:"string"`
	Deleted vanilla.NullTime   `column:"deleted" json:"deleted,omitempty" swaggertype:"string"`
}

func (Template) TableName() string {
	return "template_v2"
}

func (obj Template) MarshalJSON() ([]byte, error) {

	type T struct {
		Uuid      string             `column:"uuid"    json:"uuid,omitempty"` // pk
		Name      string             `column:"name"    json:"name,omitempty"`
		Summary   vanilla.NullString `column:"summary" json:"summary,omitempty" swaggertype:"string"`
		Flow      string             `column:"flow"    json:"flow,omitempty"`
		Inputs    string             `column:"inputs"  json:"inputs,omitempty"`
		Origin    string             `column:"origin"  json:"origin,omitempty"`
		Created   time.Time          `column:"created" json:"created,omitempty"`
		Updated   vanilla.NullTime   `column:"updated" json:"updated,omitempty" swaggertype:"string"`
		Deleted   vanilla.NullTime   `column:"deleted" json:"deleted,omitempty" swaggertype:"string"`
		StepCount int                `json:"step_count,omitempty"`
	}

	var flow = []interface{}{}
	err := json.Unmarshal([]byte(obj.Flow), &flow)
	if err != nil {
		return nil, err
	}
	var t T

	t.Uuid = obj.Uuid
	t.Name = obj.Name
	t.Summary = obj.Summary
	t.Flow = obj.Flow
	t.Inputs = obj.Inputs
	t.Origin = obj.Origin
	t.Created = obj.Created
	t.Updated = obj.Updated
	t.Deleted = obj.Deleted
	t.StepCount = len(flow)

	return json.Marshal(t)
}

package v1

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/NexClipper/sudory/pkg/server/macro/newist"
)

func TestTemplateCommandJson(t *testing.T) {

	cmd := NewServiceCommand().TemplateCommand

	data, err := json.MarshalIndent(cmd, "", " ")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
	kubcmd_ := new(DbSchemaTemplateCommand)

	err = json.Unmarshal(data, kubcmd_)
	if err != nil {
		t.Error(err)
	}
}

func TestDbSchemaTemplateCommandJson(t *testing.T) {

	cmd := NewServiceCommand()

	data, err := json.MarshalIndent(cmd, "", " ")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))

	kubcmd_ := new(DbSchemaTemplateCommand)

	err = json.Unmarshal(data, kubcmd_)
	if err != nil {
		t.Error(err)
	}
}

func NewServiceCommand() DbSchemaTemplateCommand {

	out := DbSchemaTemplateCommand{}

	out.Id = 11112222333344445555
	out.Created = newist.Time(time.Now())
	out.Updated = newist.Time(time.Now())
	out.Deleted = nil
	out.Uuid = "00001111222233334444555566667777"
	out.Name = "test-name"
	out.Summary = newist.String("test: ...")
	out.ApiVersion = "v1"
	out.TemplateUuid = "00001111222233334444555566667777"
	out.Sequence = newist.Int32(0)
	out.Method = newist.String("test.method.get.v1")
	out.Args = map[string]string{
		"name":  "test-name",
		"arg-1": "test-arg-1",
	}

	return out
}

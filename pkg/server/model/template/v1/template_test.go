package v1

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/NexClipper/sudory/pkg/server/macro/newist"
)

func TestTemplateJson(t *testing.T) {

	m := NewTemplate()

	data, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		t.Error(err)
	}
	println(string(data))

	cmd_ := new(Template)

	err = json.Unmarshal(data, cmd_)
	if err != nil {
		t.Error(err)
	}
}

func TestDbSchemaTemplateJson(t *testing.T) {

	m := NewTemplate()

	data, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		t.Error(err)
	}
	println(string(data))

	cmd_ := new(Template)

	err = json.Unmarshal(data, cmd_)
	if err != nil {
		t.Error(err)
	}
}

func TestHttpReqTemplateJson(t *testing.T) {

	m := NewTemplate()

	data, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		t.Error(err)
	}
	println(string(data))

	cmd_ := new(Template)

	err = json.Unmarshal(data, cmd_)
	if err != nil {
		t.Error(err)
	}
}

// const templateUuid = "cda6498a235d4f7eae19661d41bc154c"

func NewTemplate() Template {

	out := Template{}
	out.Id = 11112222333344445555
	out.Created = newist.Time(time.Now())
	out.Updated = newist.Time(time.Now())
	out.Deleted = nil
	out.Uuid = "00001111222233334444555566667777"
	out.Name = "test-name"
	out.Summary = newist.String("test: ...")
	// out.ApiVersion = newist.String("v1")
	out.Origin = "origin"

	return out
}

func EmptyServiceTemplate() Template {
	return Template{}
}

// func LogPrint(out io.Writer, template DbSchemaTemplate) {
// 	w := tabwriter.NewWriter(out, 4, 4, 2, ' ', tabwriter.TabIndent|tabwriter.Debug)

// 	combine := stringCombiner("\t")
// 	fmt.Fprintln(w, combine("Id", template.Id))
// 	fmt.Fprintln(w, combine("CreatedAt", template.CreatedAt.String()))
// 	fmt.Fprintln(w, combine("UpdatedAt", template.UpdatedAt.String()))
// 	fmt.Fprintln(w, combine("Removed", template.DeletedAt))
// 	fmt.Fprintln(w, combine("Name", template.Name))
// 	fmt.Fprintln(w, combine("Summary", template.Summary))
// 	fmt.Fprintln(w, combine("ApiVersion", template.ApiVersion))
// 	fmt.Fprintln(w, combine("Origin", template.Origin))

// 	w.Flush()
// }

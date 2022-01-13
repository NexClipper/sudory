package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	_ "github.com/go-sql-driver/mysql" //justifying
)

func TestTemplateJson(t *testing.T) {

	m := NewTemplate().Template

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

const templateUuid = "cda6498a235d4f7eae19661d41bc154c"

func NewTemplate() DbSchemaTemplate {
	return DbSchemaTemplate{
		DbMeta: metav1.DbMeta{},
		Template: Template{
			LabelMeta: metav1.LabelMeta{
				Uuid:       templateUuid,
				Name:       "template_kube_get_pods",
				Summary:    "template_kube_get_pods: ...",
				ApiVersion: "v1",
			},
			TemplateProperty: TemplateProperty{
				Origin: "test_defined",
			},
		},
	}
}

var testdatetime = timeParse("2009-11-10 23:00:00 UTC")

func timeParse(s string) time.Time {
	const layout = "2006-01-02 15:04:05 MST"
	t, err := time.Parse(layout, s)
	if err != nil {
		panic(err)
	}

	return t
}

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func createKeyValueJson(m map[string]string) []byte {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return b
}

func stringCombiner(div string) func(s ...interface{}) string {

	tostrng := func(a interface{}) string {
		return fmt.Sprintf("%v", a)
	}

	return func(s ...interface{}) string {
		var tmp string
		if len(s) == 0 {
			return tmp
		}

		tmp = tostrng(s[0])
		for _, it := range s[1:] {
			tmp += div + tostrng(it)
		}
		return tmp
	}
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

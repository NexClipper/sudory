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

const templateUuid = "cda6498a235d4f7eae19661d41bc154c"
const commandUuid = "d1b8af587470407f9b4299b501979e00"

func NewServiceCommand() DbSchemaTemplateCommand {
	return DbSchemaTemplateCommand{
		DbMeta: metav1.DbMeta{},
		TemplateCommand: TemplateCommand{
			LabelMeta: metav1.LabelMeta{
				Uuid:       commandUuid,
				Name:       "template_kube_get_pods",
				Summary:    "template_kube_get_pods: ...",
				ApiVersion: "v1",
			},
			TemplateCommandProperty: TemplateCommandProperty{
				TemplateUuid: templateUuid,
				Method:       "kubernetes.deployment.get.v1",
				Args: map[string]string{
					"--output": "yaml",
				},
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

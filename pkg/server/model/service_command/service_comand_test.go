package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
	"text/tabwriter"
	"time"

	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	_ "github.com/go-sql-driver/mysql" //justifying
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

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

const templateUuid = "cda6498a235d4f7eae19661d41bc154c"
const commandUuid = "d1b8af587470407f9b4299b501979e00"

func NewServiceCommand() ServiceCommand {
	return ServiceCommand{
		DbMeta: metav1.DbMeta{
			Id:        0,
			Uuid:      commandUuid,
			CreatedBy: "hansel",
			// CreatedAt: serializev1.JSONTime{Time: testdatetime},
			UpdatedBy: "hansel",
			// UpdatedAt: serializev1.JSONTime{Time: testdatetime},
			// DeletedAt: serializev1.JSONTime{Time: testdatetime},
		},
		LabelMeta: metav1.LabelMeta{
			Name:       "template_kube_get_pods",
			Summary:    "template_kube_get_pods: ...",
			ApiVersion: "v1",
		},
		TemplateUuid: templateUuid,
		Method:       "kubernetes.deployment.get.v1",
		Args: map[string]string{
			"--output": "yaml",
		},
	}
}

func TestJsonSerialize(t *testing.T) {

	cmd := NewServiceCommand()

	cmd.Method = "kubernetes.deployment.get.v1"
	cmd.Args = map[string]string{
		"--output": "yaml",
	}

	data, err := json.MarshalIndent(cmd, "", " ")
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))

	kubcmd_ := new(ServiceCommand)

	err = json.Unmarshal(data, kubcmd_)
	if err != nil {
		t.Error(err)
	}

	LogPrint(os.Stdout, cmd)
}

func LogPrint(out io.Writer, command ServiceCommand) {
	w := tabwriter.NewWriter(os.Stdout, 4, 4, 2, ' ', tabwriter.TabIndent|tabwriter.Debug)

	combine := stringCombiner("\t")
	fmt.Fprintln(w, combine("Id", command.Id))
	fmt.Fprintln(w, combine("Uuid", command.Uuid))
	fmt.Fprintln(w, combine("CreateBy", command.CreatedBy))
	fmt.Fprintln(w, combine("CreateTime", command.CreatedAt.String()))
	fmt.Fprintln(w, combine("ModifiedBy", command.UpdatedBy))
	fmt.Fprintln(w, combine("ModifiedTime", command.UpdatedAt.String()))
	fmt.Fprintln(w, combine("Removed", command.DeletedAt))
	fmt.Fprintln(w, combine("Name", command.Name))
	fmt.Fprintln(w, combine("Summary", command.Summary))
	fmt.Fprintln(w, combine("ApiVersion", command.ApiVersion))
	fmt.Fprintln(w, combine("TemplateUuid", command.TemplateUuid))
	fmt.Fprintln(w, combine("Method", command.Method))
	fmt.Fprintln(w, combine("Args(pair)", createKeyValuePairs(command.Args)))
	fmt.Fprintln(w, combine("Args(json)", string(createKeyValueJson(command.Args))))

	w.Flush()
}

func TestTemplateObjectCRUD(t *testing.T) {

	template := NewServiceCommand()

	driver := "mysql"
	dsn := "root:root@tcp(127.0.0.1:3306)/sudory?charset=utf8mb4&parseTime=True&loc=Local"
	engin, err := xorm.NewEngine(driver, dsn)
	if err != nil {
		t.Error(err)
	}
	engin.SetLogLevel(log.LOG_DEBUG)

	session := engin.NewSession()

	insert, err := session.Insert(&template)
	if err != nil {
		t.Error(err)
	}

	t.Log("insert", insert)
	t.Log("template.id", template.Id)

	template_ := new(ServiceCommand)

	get, err := session.
		ID(template.Id).
		// Where("id=?", template.Id).
		Get(template_)

	if err != nil {
		t.Error(err)
	}
	t.Log("get", get)
	LogPrint(os.Stdout, *template_)

	template_ = new(ServiceCommand)
	template_.Name = "abc"

	update, err := session.
		ID(template.Id).
		// Where("id=?", template.Id).
		Cols("name").
		Update(template_)

	if err != nil {
		t.Error(err)
	}
	t.Log("update", update)
	LogPrint(os.Stdout, *template_)

	template_ = new(ServiceCommand)
	get, err = session.
		ID(template.Id).
		// Where("id=?", template.Id).
		Get(template_)

	if err != nil {
		t.Error(err)
	}
	LogPrint(os.Stdout, *template_)

	delete, err := session.
		Where("id=?", template.Id).
		Delete(template_)

	if err != nil {
		t.Error(err)
	}
	t.Log("delete", delete)
	LogPrint(os.Stdout, *template_)

	template_ = new(ServiceCommand)
	get, err = session.
		ID(template.Id).
		// Where("id=?", template.Id).
		Get(template_)

	if err != nil {
		t.Error(err)
	}
	t.Log("get", get)
	LogPrint(os.Stdout, *template_)
}

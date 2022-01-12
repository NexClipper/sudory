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
	serializev1 "github.com/NexClipper/sudory/pkg/server/model/serialize/v1"
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

func NewServiceTemplate() ServiceTemplate {
	return ServiceTemplate{
		DbMeta: metav1.DbMeta{
			// Id:        0,
			Uuid:      serializev1.NewUuidString(),
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
		Origin: "test_defined",
	}
}

func EmptyServiceTemplate() ServiceTemplate {
	return ServiceTemplate{}
}

func LogPrint(out io.Writer, template ServiceTemplate) {
	w := tabwriter.NewWriter(out, 4, 4, 2, ' ', tabwriter.TabIndent|tabwriter.Debug)

	combine := stringCombiner("\t")
	fmt.Fprintln(w, combine("Id", template.Id))
	fmt.Fprintln(w, combine("CreateBy", template.CreatedBy))
	fmt.Fprintln(w, combine("CreateTime", template.CreatedAt.String()))
	fmt.Fprintln(w, combine("ModifiedBy", template.UpdatedBy))
	fmt.Fprintln(w, combine("ModifiedTime", template.UpdatedAt.String()))
	fmt.Fprintln(w, combine("Removed", template.DeletedAt))
	fmt.Fprintln(w, combine("Name", template.Name))
	fmt.Fprintln(w, combine("Summary", template.Summary))
	fmt.Fprintln(w, combine("ApiVersion", template.ApiVersion))
	fmt.Fprintln(w, combine("Origin", template.Origin))

	w.Flush()
}

func TestJsonSerialize(t *testing.T) {

	cmd := NewServiceTemplate()

	cmd.Origin = "test_defined"

	data, err := json.MarshalIndent(cmd, "", " ")
	if err != nil {
		t.Error(err)
	}
	println(string(data))

	cmd_ := new(ServiceTemplate)

	err = json.Unmarshal(data, cmd_)
	if err != nil {
		t.Error(err)
	}

	LogPrint(os.Stdout, cmd)
}

func TestTemplateObjectCRUD(t *testing.T) {

	template := NewServiceTemplate()

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

	template_ := new(ServiceTemplate)

	get, err := session.
		ID(template.Id).
		// Where("id=?", template.Id).
		Get(template_)

	if err != nil {
		t.Error(err)
	}
	t.Log("get", get)
	LogPrint(os.Stdout, *template_)

	template_ = new(ServiceTemplate)
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

	delete, err := session.
		Where("id=?", template.Id).
		Delete(template_)

	if err != nil {
		t.Error(err)
	}
	t.Log("delete", delete)
	LogPrint(os.Stdout, *template_)

	template_ = new(ServiceTemplate)
	get, err = session.
		ID(template.Id).
		// Where("id=?", template.Id).
		Get(template_)

	if err != nil {
		t.Error(err)
	}
	t.Log("get", get)
}

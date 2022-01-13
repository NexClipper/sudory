package database_test

import (
	"encoding/json"
	"testing"

	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
	tcommandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

func TestCreateTemplate(t *testing.T) {

	m := newTemplate()

	cnt, err := newEngine().Insert(m)
	errorHandle(err, func(err error) {
		t.Fatal(err)
	})

	t.Log("create record", cnt)

}

func TestGetTemplate(t *testing.T) {

	m := newTemplate()

	has, err := newEngine().
		Where("uuid = ?", m.Uuid).
		Get(m)

	errorHandle(err, func(err error) {
		t.Fatal(err)
	})

	t.Log("get record", has)
	jsonString(m, t.Log)

}

func TestUpdateTemplate(t *testing.T) {

	m := newTemplate()

	has, err := newEngine().
		Where("uuid = ?", m.Uuid).
		Get(m)

	errorHandle(err, func(err error) {
		t.Fatal(err)
	})

	t.Log("get record", has)

	m.Name = "(updated)" + m.Name

	cnt, err := newEngine().
		ID(m.Id).
		Update(m)

	errorHandle(err, func(err error) {
		t.Fatal(err)
	})

	t.Log("update record", cnt)
	jsonString(m, t.Log)

}

func TestDeleteTemplate(t *testing.T) {

	m := newTemplate()

	has, err := newEngine().
		Where("uuid = ?", m.Uuid).
		Get(m)

	errorHandle(err, func(err error) {
		t.Fatal(err)
	})

	t.Log("get record", has)

	cnt, err := newEngine().
		ID(m.Id).
		Delete()

	errorHandle(err, func(err error) {
		t.Fatal(err)
	})

	t.Log("delete record", cnt)
	jsonString(m, t.Log)

}

const templateUuid = "cda6498a235d4f7eae19661d41bc154c"

func newTemplate() *templatev1.DbSchemaTemplate {
	return &templatev1.DbSchemaTemplate{
		Template: templatev1.Template{
			LabelMeta: metav1.LabelMeta{
				Uuid:       templateUuid,
				Name:       "template_kube_get_pods",
				Summary:    "template_kube_get_pods: ...",
				ApiVersion: "v1",
			},
			TemplateProperty: templatev1.TemplateProperty{
				Origin: "test_defined",
			},
		},
	}
}

const commandUuid = "d1b8af587470407f9b4299b501979e00"

func NewServiceCommand() *tcommandv1.DbSchemaTemplateCommand {
	return &tcommandv1.DbSchemaTemplateCommand{
		DbMeta: metav1.DbMeta{},
		TemplateCommand: tcommandv1.TemplateCommand{
			LabelMeta: metav1.LabelMeta{
				Uuid:       commandUuid,
				Name:       "template_kube_get_pods",
				Summary:    "template_kube_get_pods: ...",
				ApiVersion: "v1",
			},
			TemplateCommandProperty: tcommandv1.TemplateCommandProperty{
				TemplateUuid: templateUuid,
				Method:       "kubernetes.deployment.get.v1",
				Args: map[string]string{
					"--output": "yaml",
				},
			},
		},
	}
}

func EmptyServiceTemplate() templatev1.Template {
	return templatev1.Template{}
}

func newEngine() *xorm.Engine {
	const (
		driver = "mysql"
		dsn    = "root:root@tcp(127.0.0.1:3306)/sudory?charset=utf8mb4&parseTime=True&loc=Local"
	)

	engine, err := xorm.NewEngine(driver, dsn)
	if errorHandle(err) {
		panic(err)
	}

	engine.SetLogLevel(log.LOG_DEBUG)

	return engine
}

func jsonString(m interface{}, writer func(a ...interface{})) {
	j, err := json.MarshalIndent(m, "", " ")

	errorHandle(err, func(err error) {
		panic(err)
	})

	writer(string(j))
}

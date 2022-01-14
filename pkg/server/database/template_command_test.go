package database

import (
	"errors"
	"fmt"
	"testing"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	templatev1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
	tcommandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

//Test TemplateCommand CRUD
//테스트 목적: TemplateCommand 테이블 기능
//테스트 조건:
//		로컬 호스트 mariadb 인스턴스
//		v1.TemplateCommand 테이블
func TestTemplateCommandCRUD(t *testing.T) {

	newEngine := func() *xorm.Engine {
		const (
			driver = "mysql"
			dsn    = "root:root@tcp(127.0.0.1:3306)/sudory?charset=utf8mb4&parseTime=True&loc=Local"
		)

		engine, err := xorm.NewEngine(driver, dsn)
		if ErrorHandle(err) {
			panic(err)
		}

		engine.SetLogLevel(log.LOG_DEBUG)

		return engine
	}

	database := DBManipulator{engine: newEngine()}

	Are := func(expect interface{}, actual interface{}, equal, diff func(string)) {

		ex := fmt.Sprint(expect)
		ac := fmt.Sprint(actual)
		if ex == ac {
			if equal != nil {
				equal(fmt.Sprintf("Equal: expect='%v', actual='%v'", expect, actual))
			}
		} else {
			if equal != nil {
				equal(fmt.Sprintf("NotEq: expect='%v', actual='%v'", expect, actual))
			}
		}
	}

	const templateUuid = "cda6498a235d4f7eae19661d41bc154c"
	create_model := func() *tcommandv1.DbSchemaTemplateCommand {
		return &tcommandv1.DbSchemaTemplateCommand{
			TemplateCommand: tcommandv1.TemplateCommand{
				LabelMeta: metav1.LabelMeta{
					Uuid:       templateUuid,
					Name:       "(NEW) template_kube_get_pods",
					Summary:    "(NEW) template_kube_get_pods: ...",
					ApiVersion: "v1",
				},
				TemplateCommandProperty: tcommandv1.TemplateCommandProperty{
					TemplateUuid: "template_uuid",
					Method:       "kube.pod.get.v1",
					Args: map[string]string{
						"name": "test_name",
						"ns":   "test_namespace",
					},
				},
			},
		}
	}
	updated_model := func() *tcommandv1.DbSchemaTemplateCommand {
		return &tcommandv1.DbSchemaTemplateCommand{
			TemplateCommand: tcommandv1.TemplateCommand{
				LabelMeta: metav1.LabelMeta{
					Uuid:       templateUuid,
					Name:       "(UPDATE) template_kube_get_pods",
					Summary:    "(UPDATE) template_kube_get_pods: ...",
					ApiVersion: "v2",
				},
				TemplateCommandProperty: tcommandv1.TemplateCommandProperty{
					TemplateUuid: "template_uuid",
					Method:       "kube.pod.get.v2",
					Args: map[string]string{
						"name": "update_name",
						"ns":   "update_namespace",
						"args": "test_args",
					},
				},
			},
		}
	}

	empty_model := func() *templatev1.DbSchemaTemplate {
		return new(templatev1.DbSchemaTemplate)
	}
	create := func() error {
		affect, err := database.CreateTemplateCommand(*create_model())
		if !(0 < affect) {
			return errors.New("no record created")
		}
		return err
	} //생성

	read := func() error {
		_, err := database.GetTemplateCommand(create_model().Uuid)
		return err
	} //조회
	validate := func(model *tcommandv1.DbSchemaTemplateCommand) func() error {
		return func() error {
			record, err := database.GetTemplateCommand(model.Uuid)

			are_equl := (func(msg string))(nil)
			are_diff := func(msg string) { t.Error(err) }
			//valied
			Are(model.Uuid, record.Uuid, are_equl, are_diff)
			Are(model.Name, record.Name, are_equl, are_diff)
			Are(model.Summary, record.Summary, are_equl, are_diff)
			Are(model.ApiVersion, record.ApiVersion, are_equl, are_diff)
			Are(model.TemplateUuid, record.TemplateUuid, are_equl, are_diff)
			Are(model.Method, record.Method, are_equl, are_diff)
			Are(model.Args, record.Args, are_equl, are_diff)

			return err
		}
	} //데이터 정합 확인

	update := func() error {
		affect, err := database.UpdateTemplateCommand(*updated_model())

		if !(0 < affect) {
			return errors.New("no record updated")
		}
		return err
	} //데이터 갱신
	delete := func() error {
		affect, err := database.DeleteTemplateCommand(updated_model().Uuid)

		if !(0 < affect) {
			return errors.New("no record deleted")
		}
		return err
	} //삭제
	clear := func() error {
		table := empty_model().TableName()

		sqlrst, err := newEngine().
			Exec(fmt.Sprintf("delete from %s where uuid = ?", table), create_model().Uuid)
		if err != nil {
			return err
		}
		affect, err := sqlrst.RowsAffected()

		if !(0 < affect) {
			return errors.New("no record deleted")
		}
		return err
	} //데이터 정리

	//시나리오 구성
	TestScenarios([]TestChapter{
		{Subject: "create newist record", Method: create},
		{Subject: "vaild created record", Method: validate(create_model())},
		{Subject: "read created record", Method: read},
		{Subject: "update record", Method: update},
		{Subject: "vaild update record", Method: validate(updated_model())},
		{Subject: "delete record", Method: delete},    //테이블에서 지워지는것이 아니라 삭제 플래그가 업데이트 됨 (레코드가 남아있음)
		{Subject: "clear test record", Method: clear}, //남아있는 레코드 정리
	}).Foreach(t)
}

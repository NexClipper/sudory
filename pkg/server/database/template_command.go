package database

import (
	"errors"
	"fmt"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	commandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
)

/* CreateTemplateCommand
   @return []DbSchemaTemplateCommand, error
   @method insert
   @from template_command
   @condition template_command
*/
func (d *DBManipulator) CreateTemplateCommand(model []commandv1.DbSchemaTemplateCommand) (int64, error) {
	var err error
	tx := d.session()
	tx.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	return tx.Insert(&model)
}

/* GetSearchTemplateCommand
   @return []DbSchemaTemplateCommand, error
   @method select
   @from template_command
   @condition uuid%, template_uuid%, %name%, %methods%
*/
func (d *DBManipulator) GetSearchTemplateCommand(param map[string]string) ([]commandv1.DbSchemaTemplateCommand, error) {
	tx := d.session()
	var records = make([]commandv1.DbSchemaTemplateCommand, 0)

	uuid := param["uuid"]
	template_uuid := param["template_uuid"]
	name := param["name"]
	methods := param["methods"]

	uuid = fmt.Sprintf("%s%%", uuid)
	template_uuid = fmt.Sprintf("%s%%", template_uuid)
	name = fmt.Sprintf("%%%s%%", name)
	methods = fmt.Sprintf("%%%s%%", methods)

	//SELECT * FROM template_command WHERE uuid = ? AND template_uuid = ? AND name = ? AND methods = ?
	has, err := tx.Where("uuid = ? AND template_uuid = ? AND name = ? AND methods = ?", uuid, template_uuid, name, methods).
		Get(records)

	if !has {
		return nil, errors.New(ErrorRecordWasNotFound)
	}

	return records, err
}

/* GetTemplateCommandWithTemplateUuid
   @return []DbSchemaTemplateCommand, error
   @method select
   @from template_command
   @condition template_uuid
*/
func (d *DBManipulator) GetTemplateCommandWithTemplateUuid(param map[string]string) ([]commandv1.DbSchemaTemplateCommand, error) {
	tx := d.session()
	var records = make([]commandv1.DbSchemaTemplateCommand, 0)

	template_uuid := param["template_uuid"]

	//SELECT * FROM template_command WHERE template_uuid = ?
	_, err := tx.Where("template_uuid = ?", template_uuid).
		Get(records)

	return records, err
}

/* GetTemplateCommand
   @return DbSchemaTemplateCommand, error
   @method select
   @from template_command
   @condition: uuid
*/
func (d *DBManipulator) GetTemplateCommand(params map[string]string) (*commandv1.DbSchemaTemplateCommand, error) {
	tx := d.session()
	var record = new(commandv1.DbSchemaTemplateCommand)

	uuid := params["uuid"]

	//SELECT * FROM template_command WHERE uuid = ? LIMIT 1
	has, err := tx.Where("uuid = ?", uuid).
		Limit(1).
		Get(record)

	if !has {
		return nil, errors.New(ErrorRecordWasNotFound)
	}

	return record, err
}

/* CreateTemplateCommand
   @return []DbSchemaTemplateCommand, error
   @method insert
   @from template_command
   @condition template_command
*/
func (d *DBManipulator) UpdateTemplateCommand(model commandv1.DbSchemaTemplateCommand) (int64, error) {
	var err error
	tx := d.session()
	tx.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var record = new(commandv1.DbSchemaTemplateCommand)

	uuid := model.Uuid

	//아이디를 가져오기
	has, err := tx.Where("uuid = ?", uuid).
		Get(record)

	if HasError(err) {
		return -1, err
	}
	if !has {
		return -2, errors.New(ErrorRecordWasNotFound)
	}

	//아이디를 조건으로 업데이트
	return tx.ID(record.Id).
		Update(model)
}

//변환 DbSchema -> TemplateCommand
func TransFormDbSchemaTemplateCommand(s []commandv1.DbSchemaTemplateCommand) []commandv1.TemplateCommand {
	var out = make([]commandv1.TemplateCommand, len(s))
	for n, it := range s {
		out[n] = it.TemplateCommand
	}

	return out
}

//변환 TemplateCommand -> DbSchema
func TransFormTemplateCommand(s []commandv1.TemplateCommand) []commandv1.DbSchemaTemplateCommand {
	var out = make([]commandv1.DbSchemaTemplateCommand, len(s))
	for n, it := range s {
		out[n] = commandv1.DbSchemaTemplateCommand{TemplateCommand: it}
	}

	return out
}

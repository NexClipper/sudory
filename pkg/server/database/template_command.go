package database

import (
	"errors"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	tcommandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
)

/* CreateTemplateCommand
   @return int64, error
   @method insert
   @from TemplateCommand
   @condition DbSchemaTemplateCommand
*/
func (d *DBManipulator) CreateTemplateCommand(model tcommandv1.DbSchemaTemplateCommand) (int64, error) {
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

	affect, err := tx.Insert(&model)

	return affect, err
}

/* GetTemplateCommand
   @return DbSchemaTemplateCommand, error
   @method select
   @from TemplateCommand
   @condition: uuid
*/
func (d *DBManipulator) GetTemplateCommand(uuid string) (*tcommandv1.DbSchemaTemplateCommand, error) {
	tx := d.session()

	var record = new(tcommandv1.DbSchemaTemplateCommand)
	//SELECT * FROM template_command WHERE uuid = ? LIMIT 1
	has, err := tx.Where("uuid = ?", uuid).
		Get(record)

	if !has {
		return nil, errors.New(ErrorRecordWasNotFound)
	}

	return record, err
}

/* FindTemplateCommand
   @return []DbSchemaTemplateCommand, error
   @method find
   @from TemplateCommand
   @condition where, args
*/
func (d *DBManipulator) FindTemplateCommand(where string, args ...string) ([]tcommandv1.DbSchemaTemplateCommand, error) {
	tx := d.session()

	var records = make([]tcommandv1.DbSchemaTemplateCommand, 0)
	//SELECT * FROM template_command WHERE [cond]
	err := tx.Where(where, args).
		Find(records)

	return records, err
}

/* UpdateTemplateCommand
   @return int64, error
   @method insert
   @from TemplateCommand
   @condition DbSchemaTemplateCommand
   @comment [panic]
   @comment		message: golang panic hash of unhashable type {noun pointer data struct}
   @comment   	원인: xorm Update or Insert 등의 반환 기능이 있는 메소드 호출 하면서 패닉 발생 값을 포인터로 넘긴다
*/
func (d *DBManipulator) UpdateTemplateCommand(model tcommandv1.DbSchemaTemplateCommand) (int64, error) {
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

	//아이디를 가져오기
	var record = new(tcommandv1.DbSchemaTemplateCommand)
	//SELECT * FROM template_command WHERE uuid = ? LIMIT 1
	has, err := tx.Where("uuid = ?", model.Uuid).
		Get(record)

	if HasError(err) {
		return -1, err
	}
	if !has {
		return -2, errors.New(ErrorRecordWasNotFound)
	}

	//아이디를 조건으로 업데이트
	affect, err := tx.ID(record.Id).
		Update(&model)

	return affect, err
}

/* DeleteTemplateCommand
   @return int64, error
   @method delete
   @from TemplateCommand
   @condition uuid
*/
func (d *DBManipulator) DeleteTemplateCommand(uuid string) (int64, error) {
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

	//아이디를 가져오기
	var record = new(tcommandv1.DbSchemaTemplateCommand)
	//SELECT * FROM template_command WHERE uuid = ? LIMIT 1
	has, err := tx.Where("uuid = ?", uuid).
		Get(record)

	if HasError(err) {
		return -1, err
	}
	if !has {
		return -2, errors.New(ErrorRecordWasNotFound)
	}

	//아이디를 조건으로 업데이트
	affect, err := tx.ID(record.Id).
		Delete(record)

	return affect, err
}

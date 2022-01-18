package database

import (
	tcommandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
)

/* CreateTemplateCommand
   @return int64, error
   @method insert
   @from TemplateCommand
   @condition DbSchemaTemplateCommand
*/
func (d *DBManipulator) CreateTemplateCommand(model tcommandv1.DbSchemaTemplateCommand) error {
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
	if err != nil {
		return err
	}
	if !(0 < affect) {
		return ErrorNoAffecte()
	}
	return nil
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
	//SELECT * FROM {table} WHERE uuid = ? LIMIT 1
	has, err := tx.Where("uuid = ?", uuid).
		Get(record)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrorRecordWasNotFound()
	}

	return record, nil
}

/* FindTemplateCommand
   @return []DbSchemaTemplateCommand, error
   @method find
   @from TemplateCommand
   @condition where, args
*/
func (d *DBManipulator) FindTemplateCommand(where string, args ...interface{}) ([]tcommandv1.DbSchemaTemplateCommand, error) {
	tx := d.session()

	var records = make([]tcommandv1.DbSchemaTemplateCommand, 0)
	//SELECT * FROM {table} WHERE [cond]
	err := tx.Where(where, args...).
		Find(&records)
	if err != nil {
		return nil, err
	}

	return records, nil
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
func (d *DBManipulator) UpdateTemplateCommand(model tcommandv1.DbSchemaTemplateCommand) error {
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

	//아이디를 조건으로 업데이트
	affect, err := tx.Where("uuid = ?", model.Uuid).
		Update(&model)
	if err != nil {
		return err
	}
	if !(0 < affect) {
		return ErrorNoAffecte()
	}
	return nil
}

/* DeleteTemplateCommand
   @return int64, error
   @method delete
   @from TemplateCommand
   @condition uuid
*/
func (d *DBManipulator) DeleteTemplateCommand(uuid string) error {
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

	var model = new(tcommandv1.DbSchemaTemplateCommand)
	//아이디를 조건으로 업데이트
	affect, err := tx.Where("uuid = ?", uuid).
		Delete(model)
	if err != nil {
		return err
	}
	if !(0 < affect) {
		return ErrorNoAffecte()
	}
	return nil
}

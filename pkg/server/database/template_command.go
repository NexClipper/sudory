package database

import (
	tcommandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
)

// CreateTemplateCommand
//  @return error
//  @method insert
//  @from TemplateCommand
//  @condition DbSchemaTemplateCommand
func (ctx Session) CreateTemplateCommand(model tcommandv1.DbSchemaTemplateCommand) error {
	tx := ctx.Tx()

	affect, err := tx.
		AllCols().Insert(&model)
	if err != nil {
		return err
	}
	if !(0 < affect) {
		return ErrorNoAffecte()
	}
	return nil
}

// GetTemplateCommand
//  @return DbSchemaTemplateCommand, error
//  @method get
//  @from TemplateCommand
//  @condition: uuid
func (ctx Session) GetTemplateCommand(uuid string) (*tcommandv1.DbSchemaTemplateCommand, error) {
	tx := ctx.Tx()

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

// FindTemplateCommand
//  @return []DbSchemaTemplateCommand, error
//  @method find
//  @from TemplateCommand
//  @condition where, args
func (ctx Session) FindTemplateCommand(where string, args ...interface{}) ([]tcommandv1.DbSchemaTemplateCommand, error) {
	tx := ctx.Tx()

	var records = make([]tcommandv1.DbSchemaTemplateCommand, 0)
	//SELECT * FROM {table} WHERE [cond]
	err := tx.Where(where, args...).
		Find(&records)
	if err != nil {
		return nil, err
	}

	return records, nil
}

// UpdateTemplateCommand
//  @return error
//  @method insert
//  @from TemplateCommand
//  @condition DbSchemaTemplateCommand
//  @comment [panic] message: golang panic hash of unhashable type {noun pointer data struct}
//     	원인: xorm Update or Insert 등의 반환 기능이 있는 메소드 호출 하면서 패닉 발생 값을 포인터로 넘긴다
func (ctx Session) UpdateTemplateCommand(model tcommandv1.DbSchemaTemplateCommand) error {
	tx := ctx.Tx()

	//아이디를 조건으로 업데이트
	affect, err := tx.Where("uuid = ?", model.Uuid).
		AllCols().Update(&model)
	if err != nil {
		return err
	}
	if !(0 < affect) {
		return ErrorNoAffecte()
	}
	return nil
}

// DeleteTemplateCommand
//  @return error
//  @method delete
//  @from TemplateCommand
//  @condition uuid
func (ctx Session) DeleteTemplateCommand(uuid string) error {
	tx := ctx.Tx()

	var model = new(tcommandv1.DbSchemaTemplateCommand)
	//아이디를 조건으로 업데이트
	affect, err := tx.Where("uuid = ?", uuid).
		Delete(model)
	if err != nil {
		return err
	}
	if !(0 < affect) {
		return nil //idempotent
	}
	return nil
}

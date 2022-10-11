package database

import (
	"reflect"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"

	"github.com/pkg/errors"
	"xorm.io/builder"
	"xorm.io/xorm"
	xorm_names "xorm.io/xorm/names"
)

func XormCreate(tx *xorm.Session, i xorm_names.TableName) error {
	affect, err := tx.InsertOne(i)
	if err != nil {
		return errors.Wrapf(err, "xorm insert%v",
			logs.KVL(
				"type_name", TypeName(i),
				"table", i.TableName(),
			))
	} else if !(0 < affect) {
		return errors.Wrapf(ErrorNoAffected, "xorm insert%v",
			logs.KVL(
				"type_name", TypeName(i),
				"table", i.TableName(),
			))
	}

	return nil
}

func XormGet(tx *xorm.Session, i xorm_names.TableName) error {
	w := builder.NewWriter()
	consd := tx.Conds()
	consd.WriteTo(w)

	has, err := tx.Get(i)
	if err != nil {
		return errors.Wrapf(err, "xorm get%v",
			logs.KVL(
				"type_name", TypeName(i),
				"table", i.TableName(),
				"cond", w.String(),
				"args", w.Args(),
			))
	}
	if !has {
		return errors.Wrapf(ErrorRecordWasNotFound, "xorm get%v",
			logs.KVL(
				"type_name", TypeName(i),
				"table", i.TableName(),
				"cond", w.String(),
				"args", w.Args(),
			))
	}
	return nil
}

// Find
func XormFind(tx *xorm.Session, i interface{}) error {
	w := builder.NewWriter()
	consd := tx.Conds()
	consd.WriteTo(w)

	if err := tx.Find(i); err != nil {
		return errors.Wrapf(err, "xorm find%v",
			logs.KVL(
				"type_name", TypeName(i),
				"cond", w.String(),
				"args", w.Args(),
			))
	}

	return nil
}

// Update
func XormUpdate(tx *xorm.Session, i xorm_names.TableName) error {
	w := builder.NewWriter()
	consd := tx.Conds()
	consd.WriteTo(w)

	//레코드 업데이트
	if _, err := tx.Update(i); err != nil {
		return errors.Wrapf(err, "xorm update%v",
			logs.KVL(
				"type_name", TypeName(i),
				"table", i.TableName(),
				"cond", w.String(),
				"args", w.Args(),
			))
	}

	//입력된 타입의 동일한 복제를 만든다
	t := reflect.TypeOf(i).Elem()
	v := reflect.New(t)
	replica := v.Interface()

	//affect 카운트로 적용 확인 하지 않고
	//Get으로 검사 및 변경 값 가져오기
	has, err := tx.Where(w.String(), w.Args()...).Get(replica)
	if err != nil {
		return errors.Wrapf(err, "xorm update%v",
			logs.KVL(
				"replica_type_name", TypeName(replica),
				"replica_table", replica.(xorm_names.TableName).TableName(),
			))
	}
	if !has {
		return errors.Wrapf(ErrorNoAffected, "xorm update%v",
			logs.KVL(
				"replica_type_name", TypeName(replica),
				"replica_table", replica.(xorm_names.TableName).TableName(),
			))
	}

	// Copy(replica, record)

	reflect.ValueOf(i).Elem().Set(v.Elem())

	return nil
}

func XormDelete(tx *xorm.Session, i xorm_names.TableName) error {
	w := builder.NewWriter()
	consd := tx.Conds()
	consd.WriteTo(w)

	if affect, err := tx.Delete(i); err != nil {
		return errors.Wrapf(err, "xorm delete%v",
			logs.KVL(
				"type_name", TypeName(i),
				"table", i.TableName(),
				"cond", w.String(),
				"args", w.Args(),
			))
	} else if !(0 < affect) {
		return nil //idempotent
	}

	return nil
}

package database

import (
	"reflect"

	"github.com/NexClipper/sudory/pkg/server/database/prepare"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type XormTableNameHolder interface {
	TableName() string
}

// XormContext
type XormContext struct {
	tx        *xorm.Session
	codinator []func() XormContext
}

func (context XormContext) Codinate() XormContext {
	for i := range context.codinator {
		context = context.codinator[i]()
	}

	return context
}

func NewXormContext(session *xorm.Session) Context {
	return &XormContext{tx: session}
}

// Close
//  close session
func (context XormContext) Close() error {
	return context.tx.Close()
}

// Prepared
func (context XormContext) Prepared(preparer prepare.Preparer) Context {
	if context.codinator == nil {
		context.codinator = make([]func() XormContext, 0)
	}

	context.codinator = append(context.codinator, func() XormContext {
		context.tx = preparer.Prepared(context.tx)
		return context
	})

	return context
}

// Where
func (context XormContext) Where(where string, args ...interface{}) Context {
	if context.codinator == nil {
		context.codinator = make([]func() XormContext, 0)
	}

	context.codinator = append(context.codinator, func() XormContext {
		context.tx = context.tx.Where(where, args...)
		return context
	})

	return context
}

// Create
func (context XormContext) Create(record interface{}) error {
	affect, err := context.tx.Insert(record)
	if err != nil {
		return errors.Wrapf(err, "xorm insert%v",
			logs.KVL(
				"type_name", TypeName(record),
			))
	} else if !(0 < affect) {
		return errors.Wrapf(ErrorNoAffected(), "xorm insert%v",
			logs.KVL(
				"type_name", TypeName(record),
			))
	}

	return nil
}

// Count
func (context XormContext) Count(record interface{}) (int64, error) {
	count, err := context.Codinate().tx.Count(record)
	if err != nil {
		return 0, errors.Wrapf(err, "xorm count%v",
			logs.KVL(
				"type_name", TypeName(record),
			))
	}

	return count, nil
}

// Get
func (context XormContext) Get(record interface{}) error {
	if has, err := context.Codinate().tx.Get(record); err != nil {
		return errors.Wrapf(err, "xorm get%v",
			logs.KVL(
				"type_name", TypeName(record),
			))
	} else if !has {
		return errors.Wrapf(ErrorRecordWasNotFound(), "xorm get%v",
			logs.KVL(
				"type_name", TypeName(record),
			))
	}

	return nil
}

// Find
func (context XormContext) Find(records interface{}) error {
	if err := context.Codinate().tx.Find(records); err != nil {
		return errors.Wrapf(err, "xorm find=%v", TypeName(records))
	}

	return nil
}

// Update
func (context XormContext) Update(record interface{}) error {

	//레코드 업데이트
	if _, err := context.Codinate().tx.Update(record); err != nil {
		return errors.Wrapf(err, "xorm update%v",
			logs.KVL(
				"type_name", TypeName(record),
			))
	}

	//입력된 타입의 동일한 복제를 만든다
	t := reflect.TypeOf(record).Elem()
	v := reflect.New(t)
	replica := v.Interface()

	//affect 카운트로 적용 확인 하지 않고
	//Get으로 검사 및 변경 값 가져오기
	has, err := context.Codinate().tx.Get(replica)
	if err != nil {
		return errors.Wrapf(err, "xorm update%v",
			logs.KVL(
				"type_name", TypeName(replica),
			))
	}
	if !has {
		return errors.Wrapf(ErrorNoAffected(), "xorm update%v",
			logs.KVL(
				"type_name", TypeName(replica),
			))
	}

	// Copy(replica, record)

	reflect.ValueOf(record).Elem().Set(v.Elem())

	return nil
}

// Delete
func (context XormContext) Delete(record interface{}) error {
	if affect, err := context.Codinate().tx.Delete(record); err != nil {
		return errors.Wrapf(err, "xorm delete%v",
			logs.KVL(
				"type_name", TypeName(record),
			))
	} else if !(0 < affect) {
		return nil //idempotent
	}

	return nil
}

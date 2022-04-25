package database

import (
	"reflect"

	"github.com/NexClipper/sudory/pkg/server/database/prepare"
)

// Context
//  데이터베이스 컨텍스트
type Context interface {
	//close
	Close() error
	//coordinator
	Prepared(preparer prepare.Preparer) Context
	Where(where string, args ...interface{}) Context
	Limit(limit int, start ...int) Context
	Asc(colNames ...string) Context
	Desc(colNames ...string) Context
	//operator
	Create(record interface{}) error
	Count(records interface{}) (int64, error)
	Get(record interface{}) error
	Find(records interface{}) error
	Update(record interface{}) error
	Delete(record interface{}) error
}

func TypeName(i interface{}) string {

	return reflect.TypeOf(i).String()
}

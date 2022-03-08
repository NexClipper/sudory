package nullable

import (
	"reflect"
	"strconv"
)

type nullInt32 struct {
	int32
	bool
}

func Int32(v interface{}) nullInt32 {
	return *(new(nullInt32).scan(v))
}

func (nullable nullInt32) Value() int32 {
	return nullable.int32
}
func (nullable nullInt32) Ptr() *int32 {
	return &nullable.int32
}
func (nullable nullInt32) Has() bool {
	return nullable.bool
}

func (nullable *nullInt32) scan(v interface{}) *nullInt32 {

	nullable.int32, nullable.bool = 0, false

	if reflect.ValueOf(v).IsNil() {
		return nullable
	}

	switch value := v.(type) {
	case string:
		if i, err := strconv.Atoi(value); err == nil {
			nullable.int32, nullable.bool = int32(i), true
		}
	case *string:
		if i, err := strconv.Atoi(*value); err == nil {
			nullable.int32, nullable.bool = int32(i), true
		}
	case int:
		nullable.int32, nullable.bool = int32(value), true
	case *int:
		nullable.int32, nullable.bool = int32(*value), true
	case int32:
		nullable.int32, nullable.bool = value, true
	case *int32:
		nullable.int32, nullable.bool = *value, true
	case int64:
		nullable.int32, nullable.bool = int32(value), true //(overflow)
	case *int64:
		nullable.int32, nullable.bool = int32(*value), true //(overflow)
	case uint:
		nullable.int32, nullable.bool = int32(value), true //(overflow)
	case *uint:
		nullable.int32, nullable.bool = int32(*value), true //(overflow)
	case uint32:
		nullable.int32, nullable.bool = int32(value), true //(overflow)
	case *uint32:
		nullable.int32, nullable.bool = int32(*value), true //(overflow)
	case uint64:
		nullable.int32, nullable.bool = int32(value), true //(overflow)
	case *uint64:
		nullable.int32, nullable.bool = int32(*value), true //(overflow)
	}

	return nullable
}

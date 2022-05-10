package nullable

import (
	"strconv"
)

type nullInt32 struct {
	int32
	bool
}

func Int32(v *int32) *nullInt32 {
	return new(nullInt32).scan(v)
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

	switch value := v.(type) {
	case string:
		if i, err := strconv.Atoi(value); err == nil {
			nullable.int32, nullable.bool = int32(i), true
		}
	case *string:
		if value != nil {
			if i, err := strconv.Atoi(*value); err == nil {
				nullable.int32, nullable.bool = int32(i), true
			}
		}
	case int:
		nullable.int32, nullable.bool = int32(value), true
	case *int:
		if value != nil {
			nullable.int32, nullable.bool = int32(*value), true
		}
	case int32:
		nullable.int32, nullable.bool = value, true
	case *int32:
		if value != nil {
			nullable.int32, nullable.bool = *value, true
		}
	case int64:
		nullable.int32, nullable.bool = int32(value), true //(overflow)
	case *int64:
		if value != nil {
			nullable.int32, nullable.bool = int32(*value), true //(overflow)
		}
	case uint:
		nullable.int32, nullable.bool = int32(value), true //(overflow)
	case *uint:
		if value != nil {
			nullable.int32, nullable.bool = int32(*value), true //(overflow)
		}
	case uint32:
		nullable.int32, nullable.bool = int32(value), true //(overflow)
	case *uint32:
		if value != nil {
			nullable.int32, nullable.bool = int32(*value), true //(overflow)
		}
	case uint64:
		nullable.int32, nullable.bool = int32(value), true //(overflow)
	case *uint64:
		if value != nil {
			nullable.int32, nullable.bool = int32(*value), true //(overflow)
		}
	}

	return nullable
}

type nullInt8 struct {
	int8
	bool
}

func Int8(v *int8) *nullInt8 {
	return new(nullInt8).scan(v)
}

func (nullable nullInt8) Value() int8 {
	return nullable.int8
}
func (nullable nullInt8) Ptr() *int8 {
	return &nullable.int8
}
func (nullable nullInt8) Has() bool {
	return nullable.bool
}

func (nullable *nullInt8) scan(v interface{}) *nullInt8 {
	nullable.int8, nullable.bool = 0, false

	switch value := v.(type) {
	case int8:
		nullable.int8, nullable.bool = int8(value), true
	case *int8:
		if value != nil {
			nullable.int8, nullable.bool = int8(*value), true
		}

	}

	return nullable
}

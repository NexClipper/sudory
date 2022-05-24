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

type nullInt16 struct {
	int16
	bool
}

func Int16(v *int16) *nullInt16 {
	return new(nullInt16).scan(v)
}

func (nullable nullInt16) Value() int16 {
	return nullable.int16
}
func (nullable nullInt16) Ptr() *int16 {
	return &nullable.int16
}
func (nullable nullInt16) Has() bool {
	return nullable.bool
}

func (nullable *nullInt16) scan(v interface{}) *nullInt16 {
	nullable.int16, nullable.bool = 0, false

	switch value := v.(type) {
	case int16:
		nullable.int16, nullable.bool = int16(value), true
	case *int16:
		if value != nil {
			nullable.int16, nullable.bool = int16(*value), true
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

type nullUint8 struct {
	uint8
	bool
}

func Uint8(v *uint8) *nullUint8 {
	return new(nullUint8).scan(v)
}

func (nullable nullUint8) Value() uint8 {
	return nullable.uint8
}
func (nullable nullUint8) Ptr() *uint8 {
	return &nullable.uint8
}
func (nullable nullUint8) Has() bool {
	return nullable.bool
}

func (nullable *nullUint8) scan(v interface{}) *nullUint8 {
	nullable.uint8, nullable.bool = 0, false

	switch value := v.(type) {
	case uint8:
		nullable.uint8, nullable.bool = uint8(value), true
	case *uint8:
		if value != nil {
			nullable.uint8, nullable.bool = uint8(*value), true
		}

	}

	return nullable
}

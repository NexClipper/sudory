package nullable

import (
	"strconv"
)

type nullString struct {
	string
	bool
}

func String(v interface{}) nullString {
	return *(new(nullString).scan(v))
}

func (nullable nullString) Value() string {
	return nullable.string
}
func (nullable nullString) Ptr() *string {
	return &nullable.string
}
func (nullable nullString) Has() bool {
	return nullable.bool
}

func (nullable *nullString) scan(v interface{}) *nullString {
	nullable.string, nullable.bool = "", false

	switch value := v.(type) {
	case []byte:
		nullable.string, nullable.bool = string(value), true
	case string:
		nullable.string, nullable.bool = value, true
	case *string:
		if value != nil {
			nullable.string, nullable.bool = *value, true
		}
	case int:
		nullable.string, nullable.bool = strconv.FormatInt(int64(value), 10), true
	case *int:
		if value != nil {
			nullable.string, nullable.bool = strconv.FormatInt(int64(*value), 10), true
		}
	case int32:
		nullable.string, nullable.bool = strconv.FormatInt(int64(value), 10), true
	case *int32:
		if value != nil {
			nullable.string, nullable.bool = strconv.FormatInt(int64(*value), 10), true
		}
	case int64:
		nullable.string, nullable.bool = strconv.FormatInt(value, 10), true
	case *int64:
		if value != nil {
			nullable.string, nullable.bool = strconv.FormatInt(*value, 10), true
		}
	case uint:
		nullable.string, nullable.bool = strconv.FormatUint(uint64(value), 10), true
	case *uint:
		if value != nil {
			nullable.string, nullable.bool = strconv.FormatUint(uint64(*value), 10), true
		}
	case uint32:
		nullable.string, nullable.bool = strconv.FormatUint(uint64(value), 10), true
	case *uint32:
		if value != nil {
			nullable.string, nullable.bool = strconv.FormatUint(uint64(*value), 10), true
		}
	case uint64:
		nullable.string, nullable.bool = strconv.FormatUint(value, 10), true
	case *uint64:
		if value != nil {
			nullable.string, nullable.bool = strconv.FormatUint(*value, 10), true
		}
	}

	return nullable
}

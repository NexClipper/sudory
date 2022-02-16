package nullable

import "strconv"

type Int32 struct {
	value int32
	valid bool
}

func NewInt32(v interface{}) *Int32 {
	return new(Int32).scan(v)
}

func (nullable Int32) Int32() int32 {
	return nullable.value
}
func (nullable Int32) Valid() bool {
	return nullable.valid
}

func (nullable *Int32) scan(v interface{}) *Int32 {

	nullable.value, nullable.valid = 0, false

	switch value := v.(type) {
	case string:
		if i, err := strconv.Atoi(value); err == nil {
			nullable.value, nullable.valid = int32(i), true
		}
	case *string:
		if value == nil {
			break
		}
		if i, err := strconv.Atoi(*value); err == nil {
			nullable.value, nullable.valid = int32(i), true
		}
	case int:
		nullable.value, nullable.valid = int32(value), true
	case *int:
		if value == nil {
			break
		}
		nullable.value, nullable.valid = int32(*value), true

	case int32:
		nullable.value, nullable.valid = value, true
	case *int32:
		if value == nil {
			break
		}
		nullable.value, nullable.valid = *value, true
	}

	return nullable
}

package nullable

import "strconv"

type Int32 struct {
	n     int32
	valid bool
}

func NewInt32(v interface{}) *Int32 {
	return new(Int32).scan(v)
}

func (nullable Int32) Int32() int32 {
	return nullable.n
}
func (nullable Int32) Valid() bool {
	return nullable.valid
}

func (nullable *Int32) scan(v interface{}) *Int32 {

	nullable.n, nullable.valid = 0, false

	switch value := v.(type) {
	case string:
		if i, err := strconv.Atoi(value); err == nil {
			nullable.n, nullable.valid = int32(i), true
		}
	case *string:
		if value == nil {
			break
		}
		if i, err := strconv.Atoi(*value); err == nil {
			nullable.n, nullable.valid = int32(i), true
		}
	case int:
		nullable.n, nullable.valid = int32(value), true
	case *int:
		if value == nil {
			break
		}
		nullable.n, nullable.valid = int32(*value), true

	case int32:
		nullable.n, nullable.valid = value, true
	case *int32:
		if value == nil {
			break
		}
		nullable.n, nullable.valid = *value, true
	}

	return nullable
}

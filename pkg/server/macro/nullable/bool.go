package nullable

type nullBool struct {
	value bool
	bool
}

func Bool(v *bool) *nullBool {
	return new(nullBool).scan(v)
}

func (nullable nullBool) Value() bool {
	return nullable.value
}
func (nullable nullBool) Ptr() *bool {
	return &nullable.value
}
func (nullable nullBool) Has() bool {
	return nullable.bool
}

func (nullable *nullBool) scan(v interface{}) *nullBool {
	nullable.value, nullable.bool = false, false

	switch value := v.(type) {
	case bool:
		nullable.value, nullable.bool = value, true
	case *bool:
		if value != nil {
			nullable.value, nullable.bool = *value, true
		}
	}

	return nullable
}

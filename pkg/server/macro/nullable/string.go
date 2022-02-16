package nullable

type String struct {
	value string
	valid bool
}

func NewString(v interface{}) *String {
	return new(String).scan(v)
}

func (n String) String() string {
	return n.value
}
func (n String) Valid() bool {
	return n.valid
}

func (n *String) scan(v interface{}) *String {
	var s string = ""
	var ok bool = true

	switch value := v.(type) {
	case string:
		s = value
	case *string:
		s = *value
	case []byte:
		s = string(value)
	default:
		ok = false
	}

	n.value, n.valid = s, ok
	return n
}

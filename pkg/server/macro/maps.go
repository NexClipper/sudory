package macro

func MapContain(m map[string]interface{}, key string) bool {
	_, ok := m[key]
	return ok
}

func MapValue(m map[string]interface{}, key string) interface{} {
	v, ok := m[key]
	if !ok {
		return nil
	}
	return v
}

func MapMap(m map[string]interface{}, key string) (map[string]interface{}, bool) {
	if v, ok := m[key]; ok {
		a, ok := v.(map[string]interface{})
		return a, ok
	}
	return nil, false
}

func MapString(m map[string]interface{}, key string) (string, bool) {
	if v, ok := m[key]; ok {
		a, ok := v.(string)
		return a, ok
	}
	return "", false
}

func MapInt(m map[string]interface{}, key string) (int, bool) {
	if v, ok := m[key]; ok {
		a, ok := v.(int)
		return a, ok
	}
	return 0, false
}

// WrapArray
func WrapArray(emun ...interface{}) []interface{} {
	return emun
}

// WrapMap
func WrapMap(a string, b interface{}) map[string]interface{} {
	return map[string]interface{}{a: b}
}

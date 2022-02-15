package macro

func MapFound(m map[string]interface{}, key string) bool {
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

func MapString(m map[string]interface{}, key string) (string, bool) {
	if v, ok := m[key]; ok {
		s, ok := v.(string)
		return s, ok
	}
	return "", false
}

func MapInt(m map[string]interface{}, key string) (int, bool) {
	if v, ok := m[key]; ok {
		s, ok := v.(int)
		return s, ok
	}
	return 0, false
}

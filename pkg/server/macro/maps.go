package macro

func MapString(m map[string]interface{}, key string) (string, bool) {

	if v, ok := m[key]; ok {
		s, ok := v.(string)
		return s, ok
	}
	return "", false
}
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

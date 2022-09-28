package grafana

import (
	"fmt"
	"reflect"
)

func FindCastFromMap(m map[string]interface{}, find string, cast interface{}) (bool, error) {
	if m == nil || len(m) <= 0 {
		return false, fmt.Errorf("'%s' not found", find)
	}

	val, ok := m[find]
	if !ok {
		return false, fmt.Errorf("'%s' not found", find)
	}
	found := true

	crv := reflect.ValueOf(cast)
	if crv.Kind() != reflect.Ptr {
		return found, fmt.Errorf("cast value must be pointer")
	}
	crv = crv.Elem()

	vrv := reflect.ValueOf(val)
	if vrv.Type() != crv.Type() {
		return found, fmt.Errorf("type of '%s' must be %s, not %s", find, crv.Type().String(), vrv.Type().String())
	}

	crv.Set(vrv)

	return found, nil
}

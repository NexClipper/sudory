package k8s

import (
	"fmt"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func convertMapToLabelSelector(m map[string]interface{}) (string, error) {
	if len(m) == 0 {
		return "", nil
	}

	labelSelector := &metav1.LabelSelector{MatchLabels: make(map[string]string)}

	for k, v := range m {
		labelSelector.MatchLabels[k] = fmt.Sprintf("%v", v)
	}

	selector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		return "", err
	}

	return selector.String(), nil
}

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

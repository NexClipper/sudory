package alertmanager

import (
	"fmt"
	"reflect"

	"github.com/prometheus/alertmanager/api/v2/models"
	"github.com/prometheus/alertmanager/pkg/labels"
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

func ConvertMathcersToModels(matchers []string) (models.Matchers, error) {
	modelsMatchers := make(models.Matchers, len(matchers))
	for i, mc := range matchers {
		m, err := labels.ParseMatcher(mc)
		if err != nil {
			return nil, err
		}

		modelsMatchers[i] = TypeMatcher(*m)
	}

	return modelsMatchers, nil
}

func TypeMatcher(matcher labels.Matcher) *models.Matcher {
	name := matcher.Name
	value := matcher.Value
	isRegex := (matcher.Type == labels.MatchRegexp) || (matcher.Type == labels.MatchNotRegexp)
	isEqual := (matcher.Type == labels.MatchEqual) || (matcher.Type == labels.MatchRegexp)

	typeMatcher := models.Matcher{
		Name:    &name,
		Value:   &value,
		IsRegex: &isRegex,
		IsEqual: &isEqual,
	}
	return &typeMatcher
}

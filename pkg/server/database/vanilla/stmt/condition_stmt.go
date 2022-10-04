package stmt

import (
	"github.com/pkg/errors"
)

type ConditionStmt map[string]interface{}

func (condition ConditionStmt) Build(engine ConditionBuildEngine) (ConditionResult, error) {
	unwraped := (map[string]interface{})(condition)
	r, err := engine.Build(unwraped)
	if err != nil {
		return nil, errors.Wrapf(err, "build condition statement")
	}
	return r, nil
}

func (condition ConditionStmt) Keys() []string {
	return MapKeys(map[string]interface{}(condition))
}

func MapKeys(m map[string]interface{}) []string {
	var iter func(v interface{})

	isDeadEnd := func(v interface{}) bool {
		switch v := v.(type) {
		case map[string]interface{}:
			return false // is not dead end
		case []interface{}:
			for _, v := range v {
				switch v.(type) {
				case map[string]interface{}:
					return false // is not dead end
				}
			}
			return true // [between|in] method
		default:
			return true // key-value method
		}
	}

	split := func(m map[string]interface{}) (key string, val interface{}) {
		for key, val = range m {
		}
		return
	}

	keys := make(map[string]struct{})
	iter = func(v interface{}) {

		switch v := v.(type) {
		case []interface{}:
			// slice
			for i := range v {
				iter(v[i])
			}
		case map[string]interface{}:
			// key-value
			key, val := split(v)
			if isDeadEnd(val) {
				keys[key] = struct{}{} // save key
				return
			}
			iter(val)
		default:
			// 	fmt.Printf("%#v", v)
		}
	}

	// run inter
	iter(m)

	// change format map -> []
	rst := make([]string, 0, len(keys))
	for k := range keys {
		rst = append(rst, k)
	}

	return rst
}

func And(e ...map[string]interface{}) ConditionStmt {
	v := make([]interface{}, 0, len(e))
	for i := range e {
		if len(e[i]) == 0 {
			continue
		}
		v = append(v, e[i])
	}

	return map[string]interface{}{"and": v}
}

func Or(e ...map[string]interface{}) ConditionStmt {
	v := make([]interface{}, 0, len(e))
	for i := range e {
		if len(e[i]) == 0 {
			continue
		}
		v = append(v, e[i])
	}

	return map[string]interface{}{"or": v}
}

func Not(b interface{}) ConditionStmt {
	return map[string]interface{}{"not": b}
}

func Equal(a string, b interface{}) ConditionStmt {
	return map[string]interface{}{"equal": map[string]interface{}{a: b}}
}

func IsNull(a string) ConditionStmt {
	return map[string]interface{}{"isnull": map[string]interface{}{a: nil}}
}

func In(a string, b ...interface{}) ConditionStmt {
	return map[string]interface{}{"in": map[string]interface{}{a: b}}
}

func GT(a string, b interface{}) ConditionStmt {
	return map[string]interface{}{"gt": map[string]interface{}{a: b}}
}

func GTE(a string, b interface{}) ConditionStmt {
	return map[string]interface{}{"gte": map[string]interface{}{a: b}}
}

func LT(a string, b interface{}) ConditionStmt {
	return map[string]interface{}{"lt": map[string]interface{}{a: b}}
}

func LTE(a string, b interface{}) ConditionStmt {
	return map[string]interface{}{"lte": map[string]interface{}{a: b}}
}

func Like(a string, b interface{}) ConditionStmt {
	return map[string]interface{}{"like": map[string]interface{}{a: b}}
}

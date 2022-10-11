package stmt

import (
	"github.com/pkg/errors"
)

type TypeOrderStmt = map[string][]string

// OrderStmt
//  map[string][]string{"ASC": []string{column1, column2}, "DESC": []string{column3, column4} }
type OrderStmt []TypeOrderStmt

func (order OrderStmt) Build(engine OrderBuildEngine) (OrderResult, error) {
	split := func(order OrderStmt) []interface{} {
		unwraped := ([]TypeOrderStmt)(order)
		r := make([]interface{}, len(unwraped))
		for i := range unwraped {
			r[i] = unwraped[i]
		}
		return r
	}

	o, err := engine.Build(split(order))
	if err != nil {
		return nil, errors.Wrapf(err, "build order statement")
	}

	return o, nil
}

func (order OrderStmt) Keys() []string {
	MapKeys := func(m []map[string][]string) []string {
		var iter func(v interface{})

		// isDeadEnd := func(v interface{}) bool {
		// 	switch v := v.(type) {
		// 	case map[string]interface{}:
		// 		return false // is not dead end
		// 	case []interface{}:
		// 		for _, v := range v {
		// 			switch v.(type) {
		// 			case map[string]interface{}:
		// 				return false // is not dead end
		// 			}
		// 		}
		// 		return true // [between|in] method
		// 	default:
		// 		return true // key-value method
		// 	}
		// }

		split := func(m map[string][]string) (key string, val []string) {
			for key, val = range m {
			}
			return
		}

		keys := make(map[string]struct{})
		iter = func(v interface{}) {

			switch v := v.(type) {
			case []map[string][]string:
				// slice
				for i := range v {
					iter(v[i])
				}
			case []interface{}:
				// slice
				for i := range v {
					iter(v[i])
				}
			case map[string][]string:
				// key-value
				_, val := split(v)
				// if isDeadEnd(val) {
				// 	keys[key] = struct{}{} // save key
				// 	return
				// }
				// iter(val)
				for _, key := range val {
					keys[key] = struct{}{} // save key
				}
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

	return MapKeys([]map[string][]string(order))
}

func (order OrderStmt) Asc(columns ...string) OrderStmt {
	return append(order, Asc(columns...)...)
}

func (order OrderStmt) Desc(columns ...string) OrderStmt {
	return append(order, Desc(columns...)...)
}

func Asc(columns ...string) OrderStmt {
	return []map[string][]string{{"asc": columns}}
}

func Desc(columns ...string) OrderStmt {
	return []map[string][]string{{"desc": columns}}
}

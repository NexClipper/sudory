package stmt

import (
	"fmt"
)

// func quote(exp interface{}) string {
// 	return fmt.Sprintf("%v", exp)
// }

type TypeSlice = []interface{}
type TypeMap = map[string]interface{}

// Slice
func Slice(e interface{}, es ...interface{}) TypeSlice {
	return append([]interface{}{e}, es...)
}

// Map
func Map(a string, b interface{}) TypeMap {
	return TypeMap{a: b}
}

func BackQuote(exp string) string {
	return fmt.Sprintf("`%s`", exp)
}

func Quote(exp string) string {
	return exp
}

func MapQuote(mapper func(string) string) func(src ...string) []string {
	return func(src ...string) []string {
		dst := make([]string, len(src))
		for i := range src {
			dst[i] = mapper(src[i])
		}
		return dst
	}
}

func Repeat(n int, s string) []string {
	ss := make([]string, n)
	for i := 0; i < n; i++ {
		ss[i] = s
	}
	return ss
}

// func StringRepeat(s string, seq string, n int) string {
// 	r := make([]string, n)
// 	for i := 0; i < n; i++ {
// 		r[i] = s
// 	}
// 	return strings.Join(r, seq)
// }

// func NewParser(m map[string]string) (condition *Condition, order *Orders, pagination *Pagination, err error) {
// 	//pagination
// 	if 0 < len(m["p"]) {
// 		pagination, err = NewPagination(m["p"])
// 		err = errors.Wrapf(err, "NewPagination p=%s", m["p"])
// 		if err != nil {
// 			return
// 		}
// 	}
// 	//order
// 	if 0 < len(m["o"]) {
// 		order, err = NewOrder(m["o"])
// 		err = errors.Wrapf(err, "NewOrder o=%s", m["o"])
// 		if err != nil {
// 			return
// 		}
// 	}
// 	//Condition
// 	if 0 < len(m["q"]) {
// 		condition, err = NewCondition(m["q"])
// 		err = errors.Wrapf(err, "NewCondition q=%s", m["q"])
// 		if err != nil {
// 			return
// 		}
// 	}
// 	return
// }

package macro

import (
	"fmt"
	"strings"
)

func StringBuilder() (func(elems ...interface{}), func(sep string) string) {
	buf := make([]string, 0)

	sprint := func(a interface{}) string {
		return fmt.Sprintf("%v", a)
	}

	jointer := func(elems ...interface{}) {
		for _, it := range elems {
			buf = append(buf, sprint(it))
		}
	}

	builder := func(sep string) string {
		return strings.Join(buf, sep)
	}

	return jointer, builder
}

package macro

import (
	"fmt"
	"strings"
)

func StringJoin(sep string) (func(elems ...interface{}), func() string) {
	buf := make([]string, 0)

	sprint := func(a interface{}) string {
		return fmt.Sprintf("%v", a)
	}

	jointer := func(elems ...interface{}) {
		for _, it := range elems {
			buf = append(buf, sprint(it))
		}
	}

	builder := func() string {
		return strings.Join(buf, sep)
	}

	return jointer, builder
}

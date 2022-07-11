package vanilla

import (
	"fmt"
	"strings"
)

const __COLUMN_ALIAS_SEPERATOR__ = "."

func MG(mangling string) func(string) string {
	mangling = strings.TrimSpace(mangling)

	if len(mangling) == 0 {
		return func(s string) string {
			return fmt.Sprintf("`%v`", s)
		}
	}

	manglings := make([]string, 0, 2)
	manglings = append(manglings, mangling)
	return func(s string) string {
		return fmt.Sprintf("`%v`", strings.Join(append(manglings, s), __COLUMN_ALIAS_SEPERATOR__))
	}
}

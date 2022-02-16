package query_parser

import (
	"fmt"

	. "github.com/NexClipper/sudory/pkg/server/macro"
)

type Condition struct {
	where string
	args  []interface{}
}

type ConditionFilter func(key string) (string, string, bool)

func NewCondition(m map[string]interface{}, filter ConditionFilter) *Condition {
	if len(m) == 0 {
		return &Condition{}
	}

	args := make([]interface{}, 0)
	add, build := StringBuilder()

	for key, val := range m {
		operator, format, ok := filter(key)
		if ok {
			args = append(args, fmt.Sprintf(format, val))
			add(fmt.Sprintf("%s %s ?", key, operator)) //조건문 만들기
		}
	}

	return &Condition{where: build(" AND "), args: args}
}

func (cond Condition) Where() string {
	return cond.where
}
func (cond Condition) Args() []interface{} {
	return cond.args
}

package prepare

import (
	"github.com/pkg/errors"
)

// Array
func Array(emun ...interface{}) []interface{} {
	return emun
}

// Map
func Map(a string, b interface{}) map[string]interface{} {
	return map[string]interface{}{a: b}
}

func NewParser(m map[string]string) (condition *Condition, order *Orders, pagination *Pagination, err error) {
	//pagination
	if 0 < len(m["p"]) {
		pagination, err = NewPagination(m["p"])
		err = errors.Wrapf(err, "NewPagination p=%s", m["p"])
		if err != nil {
			return
		}
	}
	//order
	if 0 < len(m["o"]) {
		order, err = NewOrder(m["o"])
		err = errors.Wrapf(err, "NewOrder o=%s", m["o"])
		if err != nil {
			return
		}
	}
	//Condition
	if 0 < len(m["q"]) {
		condition, err = NewCondition(m["q"])
		err = errors.Wrapf(err, "NewCondition q=%s", m["q"])
		if err != nil {
			return
		}
	}
	return
}

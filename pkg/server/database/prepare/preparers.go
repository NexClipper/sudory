package prepare

import (
	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type Preparer interface {
	Prepared(*xorm.Session) *xorm.Session
}

type Preparers []Preparer

func NewParser(m map[string]string) (Preparer, error) {

	preparers := make([]Preparer, 0, 3)

	//pagination
	if pagination, err := NewPagination(m["p"]); err != nil {
		if !macro.Eqaul(ErrorInvalidArgumentEmptyString(), err) {
			return nil, errors.Wrapf(err, "NewPagination json=%s", m["p"])
		}
	} else {
		preparers = append(preparers, pagination)
	}
	//order
	if order, err := NewOrder(m["o"]); err != nil {
		if !macro.Eqaul(ErrorInvalidArgumentEmptyString(), err) {
			return nil, errors.Wrapf(err, "NewOrder json=%s", m["o"])
		}
	} else {
		preparers = append(preparers, order)
	}
	//Condition
	if condition, err := NewCondition(m["q"]); err != nil {
		if !macro.Eqaul(ErrorInvalidArgumentEmptyString(), err) {
			return nil, errors.Wrapf(err, "NewCondition json=%s", m["q"])
		}
	} else {
		preparers = append(preparers, condition)
	}

	return (*Preparers)(&preparers), nil
}

func (preparers Preparers) Prepared(tx *xorm.Session) *xorm.Session {
	for _, preparer := range preparers {
		tx = preparer.Prepared(tx)
	}
	return tx
}

// WrapArray
func WrapArray(emun ...interface{}) []interface{} {
	return emun
}

// WrapMap
func WrapMap(a string, b interface{}) map[string]interface{} {
	return map[string]interface{}{a: b}
}

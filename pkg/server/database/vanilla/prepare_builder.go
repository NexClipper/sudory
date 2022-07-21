package vanilla

import (
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/prepare"
	"github.com/pkg/errors"
)

type PrepareCondition map[string]interface{}

func (m PrepareCondition) Parse() *prepare.Condition {
	q, err := prepare.NewConditionMap(m)
	if err != nil {
		err = errors.Wrapf(err, "parse map to new condition")
		panic(err)
	}

	return q
}

func And(a map[string]interface{}, b ...map[string]interface{}) PrepareCondition {
	b_ := make([]interface{}, len(b))
	for i := range b {
		b_[i] = b[i]
	}

	return map[string]interface{}{"and": append([]interface{}{a}, b_...)}
}

func Or(a map[string]interface{}, b ...map[string]interface{}) PrepareCondition {
	b_ := make([]interface{}, len(b))
	for i := range b {
		b_[i] = b[i]
	}
	return map[string]interface{}{"or": append([]interface{}{a}, b_...)}
}

func Equal(a string, b interface{}) PrepareCondition {
	return map[string]interface{}{"equal": map[string]interface{}{a: b}}
}

func IsNull(a string) PrepareCondition {
	return map[string]interface{}{"isnull": a}
}

func In(a string, b ...interface{}) PrepareCondition {
	return map[string]interface{}{"in": map[string]interface{}{a: b}}
}

func GreaterThan(a string, b interface{}) PrepareCondition {
	return map[string]interface{}{"gt": map[string]interface{}{a: b}}
}

func GreaterThanEqual(a string, b interface{}) PrepareCondition {
	return map[string]interface{}{"gte": map[string]interface{}{a: b}}
}

func LessThan(a string, b interface{}) PrepareCondition {
	return map[string]interface{}{"lt": map[string]interface{}{a: b}}
}

func LessThanEqual(a string, b interface{}) PrepareCondition {
	return map[string]interface{}{"lte": map[string]interface{}{a: b}}
}

func Like(a string, b interface{}) PrepareCondition {
	return map[string]interface{}{"like": map[string]interface{}{a: b}}
}

type PrepareOrder []interface{}

func (sl PrepareOrder) Parse() *prepare.Orders {
	o, err := prepare.NewOrderSlice(sl)
	if err != nil {
		err = errors.Wrapf(err, "parse slice to new order")
		panic(err)
	}

	return o
}

func (a PrepareOrder) Combine(b ...PrepareOrder) PrepareOrder {
	b_ := make([]interface{}, 0, len(b))

	for i := range b {
		b_ = append(b_, ([]interface{})(b[i])...)
	}

	return append(a, b_...)
}

func Asc(a ...string) PrepareOrder {
	a_ := make([]interface{}, len(a))
	for i := range a {
		a_[i] = a[i]
	}
	return append(a_, "asc")
}

func Desc(a ...string) PrepareOrder {
	a_ := make([]interface{}, len(a))
	for i := range a {
		a_[i] = a[i]
	}
	return append(a_, "desc")
}

type PreparePagination map[string]interface{}

func (m PreparePagination) Parse() *prepare.Pagination {
	o, err := prepare.NewPaginationMap(m)
	if err != nil {
		err = errors.Wrapf(err, "parse map to new pagination")
		panic(err)
	}

	return o
}

func (m *PreparePagination) Limit(l int) *PreparePagination {
	if l < 0 {
		panic("limit is greater then equal zero")
	}
	(*m)["limit"] = l
	return m
}

func (m *PreparePagination) Page(p int) *PreparePagination {
	if p <= 0 {
		panic("page is greater then zero")
	}
	(*m)["page"] = p
	return m
}

func Limit(limit int, page ...int) *PreparePagination {
	page_ := 1
	if 0 < len(page) {
		if 0 < page[0] {
			page_ = page[0]
		}
	}

	m := PreparePagination(map[string]interface{}{
		"limit": limit,
		"page":  page_,
	})
	return &m
}

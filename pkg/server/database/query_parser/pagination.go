package query_parser

import (
	"math"
	"strconv"
)

type Pagination struct {
	limit int
	page  int
	// order string

	// query map[string]interface{}
}

const (
	PaginationLimit = math.MaxUint16 //65535
	PaginationPage  = 1
)

func NewPagination(m map[string]interface{}) (*Pagination, map[string]interface{}) {

	var (
		found bool                   = false
		query map[string]interface{} = make(map[string]interface{})
		limit int                    = PaginationLimit
		page  int                    = PaginationPage
		// order string                 = ""
	)

	for n := range m {
		query[n] = m[n]
	}

	if _, ok := m[__PAGINATION_LIMIT__]; ok {
		limit, _ = strconv.Atoi(m[__PAGINATION_LIMIT__].(string))

		found = true
	}
	if _, ok := m[__PAGINATION_PAGE__]; ok {
		page, _ = strconv.Atoi(m[__PAGINATION_PAGE__].(string))
		if page < 1 {
			page = 1
		}
		found = true
	}
	// if _, ok := m[__PAGINATION_ORDER__]; ok {
	// 	order, _ = m[__PAGINATION_ORDER__].(string)
	// }

	delete(query, __PAGINATION_LIMIT__)
	delete(query, __PAGINATION_PAGE__)
	// delete(query, __PAGINATION_ORDER__)

	if !found {
		return nil, query
	}

	return &Pagination{
		// query: query,
		limit: limit,
		page:  page,
		// order: order,
	}, query
}

// // Order
// func (page Pagination) Order() string {
// 	return page.order
// }

// Limit
func (page Pagination) Limit() int {
	return page.limit
}

// Page
func (page Pagination) Page() int {
	return page.page
}

// // Query
// func (page Pagination) Query() map[string]interface{} {
// 	return page.query
// }

// Offset
func (page Pagination) Offset() int {
	return (page.page - 1) * page.limit
}

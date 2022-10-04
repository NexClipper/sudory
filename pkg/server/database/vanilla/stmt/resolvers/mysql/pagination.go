package flavor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/pkg/errors"
)

type Pagination struct {
	limit int
	page  int
}

func (page *Pagination) SetLimit(n int) {
	page.limit = n
}

func (page *Pagination) SetPage(n int) {
	page.page = n
}

// Limit
func (page Pagination) Limit() (int, bool) {
	return page.limit, page.limit != 0
}

// Page
func (page Pagination) Page() (int, bool) {
	return page.page, page.page != 0
}

// Offset
func (page Pagination) Offset() int {
	page = pagination(page)
	return offset(page)
}

// pagination
func pagination(page_ Pagination) Pagination {

	if !(0 < page_.limit) {
		//  (limit == 0) | 127
		//  (limit < 0)  | 127
		//  else         | limit
		page_.limit = stmt.DEFAULT_PAGINATION_LIMIT
	}
	if !(0 < page_.page) {
		//  (page == 0) | 1
		//  (page < 0)  | 1
		//  else        | page
		page_.page = stmt.DEFAULT_PAGINATION_PAGE
	}
	return page_
}

// offset
//  else | (page - 1) * limit
func offset(page_ Pagination) int {
	return (page_.page - 1) * page_.limit
}

// stringer
func (page Pagination) String() string {
	page = pagination(page)

	return fmt.Sprintf("%v, %v", offset(page), page.limit)
}

type MysqlPagination struct {
	stmt.PaginationStmtBuilder
}

func NewMysqlPagination() *MysqlPagination {

	flavor := new(MysqlPagination)
	// flavor.Limit = MakePaginationFunc(&page, "limit")
	// flavor.Page = MakePaginationFunc(&page, "page")

	flavor.PaginationStmtBuilder = stmt.NewPaginationBuildEngine(__DIALECT__, map[string]stmt.PaginationFunctor{
		"limit": flavor.Limit,
		"page":  flavor.Page,
	})

	return flavor
}

var (
	pagination_limit = MakePaginationFunc("limit")
	pagination_page  = MakePaginationFunc("page")
)

func (flavor *MysqlPagination) Limit(v interface{}) (stmt.PaginationResult, error) {
	return pagination_limit(v)
}
func (flavor *MysqlPagination) Page(v interface{}) (stmt.PaginationResult, error) {
	return pagination_page(v)
}

func MakePaginationFunc(method string) func(v interface{}) (stmt.PaginationResult, error) {
	scan := func(v interface{}) (int, error) {
		switch value := v.(type) {
		case string:
			i, err := strconv.Atoi(value)
			return i, err
		case float64:
			return int(value), nil
		case int:
			return value, nil
		default:
			return 0, errors.WithStack(stmt.ErrorUnsupportedType)
		}
	}

	switch strings.ToLower(method) {
	case "limit":
		return func(v interface{}) (stmt.PaginationResult, error) {
			page := new(Pagination)
			var err error
			page.limit, err = scan(v)
			return page, err
		}
	case "page":
		return func(v interface{}) (stmt.PaginationResult, error) {
			page := new(Pagination)
			var err error
			page.page, err = scan(v)
			return page, err
		}
	default:
		panic("unknown handler")
	}
}

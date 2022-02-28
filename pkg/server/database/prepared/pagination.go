package prepared

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type Pagination struct {
	limit int
	page  int
}

type PaginationFunctor func(v interface{}) error

// NewPagination
func NewPagination(s string) (*Pagination, error) {

	if len(s) == 0 {
		return nil, ErrorInvalidArgumentEmptyString()
	}

	v := new(interface{})
	if err := json.Unmarshal([]byte(s), v); err != nil {
		return nil, errors.Wrapf(err, "json.Unmarshal q=%s", s)
	}

	page := &Pagination{
		limit: math.MaxUint8, //255
		page:  1,
	}

	if err := newPaginationBuilder(page).Parse(*v); err != nil {
		return nil, err
	}

	return page, nil
}

// NewPaginationMap
func NewPaginationMap(m map[string]interface{}) (*Pagination, error) {
	page := &Pagination{
		limit: math.MaxUint8, //255
		page:  1,
	}

	if err := newPaginationBuilder(page).Parse(m); err != nil {
		return nil, err
	}

	return page, nil
}

// Limit
func (page Pagination) Limit() int {
	return page.limit
}

// Page
func (page Pagination) Page() int {
	return page.page
}

// Offset
func (page Pagination) Offset() int {
	return (page.Page() - 1) * page.Limit()
}

// Prepared
func (page Pagination) Prepared(tx *xorm.Session) *xorm.Session {
	if page.Offset() < 0 {
		return tx.Limit(page.Limit())
	}
	return tx.Limit(page.Limit(), page.Offset())
}

type paginationBuild struct {
	engine map[string]PaginationFunctor
}

func newPaginationBuilder(page *Pagination) *paginationBuild {
	builder := paginationBuild{}

	engine := map[string]PaginationFunctor{
		"limit": builder.MakePaginationFunc(page),
		"page":  builder.MakePaginationFunc(page),
	}

	builder.engine = engine
	return &builder
}

func (builder *paginationBuild) Parse(v interface{}) error {

	switch value := v.(type) {
	case map[string]interface{}:
		for key := range value {
			functor, ok := builder.engine[strings.ToLower(key)]
			if !ok {
				return ErrorNotFoundHandler(key)
			}
			err := functor(value)
			if err != nil {
				return errors.Wrapf(err, "functor key=%s value=%+v", key, value[key])
			}
		}
	case []interface{}:
		for n := range value {
			err := builder.Parse(value[n]) //(recursion)
			if err != nil {
				return err
			}
		}
	default:
		return ErrorUnsupportedType(value)
	}

	return nil
}
func (builder *paginationBuild) MakePaginationFunc(page *Pagination) func(v interface{}) error {

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
			return 0, ErrorUnsupportedType(value)
		}
	}

	return func(v interface{}) error {
		switch value := v.(type) {
		case map[string]interface{}:
			var err error

			for key := range value {
				switch key {
				case "limit":
					page.limit, err = scan(value[key])
					if err != nil {
						return errors.Wrapf(err, "scan key=%s value=%s", key, value[key])
					}
				case "page":
					page.page, err = scan(value[key])
					if err != nil {
						return errors.Wrapf(err, "scan key=%s value=%s", key, value[key])
					}
				}
			}
		default:
			return ErrorUnsupportedType(value)
		}
		return nil
	}
}

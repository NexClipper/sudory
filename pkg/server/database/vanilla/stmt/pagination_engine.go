package stmt

import (
	"math"
	"strings"

	"github.com/pkg/errors"
)

const (
	DEFAULT_PAGINATION_LIMIT = math.MaxInt8
	DEFAULT_PAGINATION_PAGE  = 1
)

type PaginationFunctor = func(v interface{}) (PaginationResult, error)

// PaginationBuildEngine
//  implement of PaginationStmtResolver
type PaginationBuildEngine struct {
	builtin map[string]PaginationFunctor
	dialect string
}

func NewPaginationBuildEngine(dialect string, builtinSet map[string]PaginationFunctor) *PaginationBuildEngine {
	engine := new(PaginationBuildEngine)
	engine.dialect = dialect
	engine.builtin = builtinSet
	return engine
}

func (engine PaginationBuildEngine) Dialect() string {
	return engine.dialect
}
func (engine PaginationBuildEngine) Build(v interface{}) (PaginationResult, error) {

	switch value := v.(type) {
	case PaginationStmt:
		return value.Build(engine)
	case []interface{}:
		if len(value) == 0 {
			return nil, errors.WithStack(ErrorInvalidArgumentEmptyObject)
		}

		var result PaginationResult
		for _, iter := range value {

			order_result, err := engine.Build(iter)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			result = func(pr PaginationResult) PaginationResult {
				if result == nil {
					return pr
				}

				if limit, ok := pr.Limit(); ok {
					result.SetLimit(limit)
				}
				if page, ok := pr.Page(); ok {
					result.SetPage(page)
				}

				return result
			}(order_result)
		}
		return result, nil
	case map[string]int:
		if len(value) == 0 {
			return nil, errors.WithStack(ErrorInvalidArgumentEmptyObject)
		}

		var result PaginationResult
		for key, val := range value {
			functor, ok := engine.builtin[strings.ToLower(key)]
			if !ok {
				return nil, errors.Cause(ErrorNotFoundHandler)
			}

			pr, err := functor(val)
			if err != nil {
				return nil, errors.Wrapf(err, "functor key=%s value=%v value_type=%T", key, val, val)
			}

			result = func(pr PaginationResult) PaginationResult {
				if result == nil {
					return pr
				}

				if limit, ok := pr.Limit(); ok {
					result.SetLimit(limit)
				}
				if page, ok := pr.Page(); ok {
					result.SetPage(page)
				}

				return result
			}(pr)
		}
		return result, nil
	default:
		return nil, errors.Wrapf(ErrorUnsupportedType, "value=%v type=%T", value, value)
	}
}

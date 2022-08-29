package stmt

import (
	"strings"

	"github.com/pkg/errors"
)

type OrderFunctor = func(v interface{}) (OrderResult, error)

// OrderBuildEngine
//  implement of OrderStmtResolver
type OrderBuildEngine struct {
	builtin map[string]OrderFunctor
	dialect string
}

func NewOrderBuildEngine(dialect string, builtinSet map[string]OrderFunctor) *OrderBuildEngine {
	engine := new(OrderBuildEngine)
	engine.dialect = dialect
	engine.builtin = builtinSet
	return engine
}

func (engine OrderBuildEngine) Dialect() string {
	return engine.dialect
}

func (engine OrderBuildEngine) Build(v interface{}) (OrderResult, error) {
	switch value := v.(type) {
	case OrderStmt:
		return value.Build(engine)
	case []interface{}:
		var result OrderResult

		for _, iter := range value {

			order_result, err := engine.Build(iter)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			result = func(or OrderResult) OrderResult {
				if result == nil {
					return or
				}
				return result.Combine(or)
			}(order_result)
		}
		return result, nil
	case map[string][]string:
		var result OrderResult
		var err error

		if len(value) == 0 {
			return nil, errors.WithStack(ErrorInvalidArgumentEmptyObject)
		}

		for key, val := range value {
			functor, ok := engine.builtin[strings.ToLower(key)]
			if !ok {
				return nil, errors.Cause(ErrorNotFoundHandler)
			}

			result, err = functor(val)
			if err != nil {
				return nil, errors.Wrapf(err, "functor key=%s value=%v value_type=%T", key, val, val)
			}
		}
		return result, nil
	default:
		return nil, errors.Wrapf(ErrorUnsupportedType, "value=%v type=%T", value, value)
	}
}

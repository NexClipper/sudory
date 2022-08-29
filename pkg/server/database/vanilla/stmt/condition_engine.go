package stmt

import (
	"strings"

	"github.com/pkg/errors"
)

type ConditionFunctor = func(v interface{}) (ConditionResult, error)

// ConditionBuildEngine
//  implement of ConditionStmtResolver
type ConditionBuildEngine struct {
	builtin map[string]ConditionFunctor
	dialect string
}

func NewConditionBuildEngine(dialect string, builtinSet map[string]ConditionFunctor) *ConditionBuildEngine {
	engine := new(ConditionBuildEngine)
	engine.dialect = dialect
	engine.builtin = builtinSet
	return engine
}

func (engine ConditionBuildEngine) Dialect() string {
	return engine.dialect
}

func (engine ConditionBuildEngine) Build(v interface{}) (ConditionResult, error) {

	switch value := v.(type) {
	case ConditionStmt:
		return value.Build(engine)
	case map[string]interface{}:
		var result ConditionResult
		var err error

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

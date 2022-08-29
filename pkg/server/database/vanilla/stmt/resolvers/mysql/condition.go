package flavor

import (
	"fmt"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/pkg/errors"
)

type Condition struct {
	query string
	args  []interface{}
}

func (condition Condition) Query() string {
	return condition.query
}

func (condition Condition) Args() []interface{} {
	return condition.args
}

type MysqlCondition struct {
	stmt.ConditionStmtBuilder
	Quote func(exp string) string // column name decorator
}

func NewMysqlCondition() *MysqlCondition {
	flavor := new(MysqlCondition)
	flavor.Quote = stmt.BackQuote
	flavor.ConditionStmtBuilder = stmt.NewConditionBuildEngine(__DIALECT__, map[string]stmt.ConditionFunctor{
		"and":     flavor.And,
		"or":      flavor.Or,
		"not":     flavor.Not,
		"equal":   flavor.Equal,
		"eq":      flavor.Equal, //addition: equal
		"gt":      flavor.GreaterThan,
		"lt":      flavor.LessThan,
		"gte":     flavor.GreaterThanOrEqual,
		"ge":      flavor.GreaterThanOrEqual, //addition: gte
		"lte":     flavor.LessThanOrEqual,
		"le":      flavor.LessThanOrEqual, //addition: lte
		"like":    flavor.Like,
		"isnull":  flavor.IsNull,
		"in":      flavor.In,
		"between": flavor.Between,
	})

	return flavor
}

func (flavor *MysqlCondition) And(v interface{}) (stmt.ConditionResult, error) {
	switch value := v.(type) {
	case []interface{}:
		var str = make([]string, 0, len(value))
		var args = make([]interface{}, 0, len(value))

		for i := range value {
			cond, err := flavor.Build(value[i])
			if err != nil {
				return nil, err
			}

			str = append(str, cond.Query())
			args = append(args, cond.Args()...)
		}

		return &Condition{query: "(" + strings.Join(str, " AND ") + ")", args: args}, nil
	case map[string]interface{}:
		var str = make([]string, 0, len(value))
		var args = make([]interface{}, 0, len(value))

		cond, err := flavor.Build(value)
		if err != nil {
			return nil, err
		}

		str = append(str, cond.Query())
		args = append(args, cond.Args()...)

		return &Condition{query: "(" + strings.Join(str, " AND ") + ")", args: args}, nil
	default:
		return nil, errors.WithStack(stmt.ErrorUnsupportedType)
	}
}
func (flavor *MysqlCondition) Or(v interface{}) (stmt.ConditionResult, error) {
	switch value := v.(type) {
	case []interface{}:
		var str = make([]string, 0, len(value))
		var args = make([]interface{}, 0, len(value))

		for i := range value {
			cond, err := flavor.Build(value[i])
			if err != nil {
				return nil, err
			}

			str = append(str, cond.Query())
			args = append(args, cond.Args()...)
		}
		return &Condition{query: "(" + strings.Join(str, " OR ") + ")", args: args}, nil
	case map[string]interface{}:
		var str = make([]string, 0, len(value))
		var args = make([]interface{}, 0, len(value))

		cond, err := flavor.Build(value)
		if err != nil {
			return nil, err
		}

		str = append(str, cond.Query())
		args = append(args, cond.Args()...)

		return &Condition{query: "(" + strings.Join(str, " OR ") + ")", args: args}, nil
	default:
		return nil, errors.WithStack(stmt.ErrorUnsupportedType)
	}
}
func (flavor *MysqlCondition) Not(v interface{}) (stmt.ConditionResult, error) {
	switch value := v.(type) {
	case map[string]interface{}:
		cond, err := flavor.Build(value)
		if err != nil {
			return nil, err
		}

		return &Condition{query: fmt.Sprintf("NOT %s", cond.Query()), args: cond.Args()}, nil
	default:
		return nil, errors.WithStack(stmt.ErrorUnsupportedType)
	}
}
func (flavor *MysqlCondition) Equal(v interface{}) (stmt.ConditionResult, error) {
	switch value := v.(type) {
	case map[string]interface{}:
		var exp string
		var arg interface{}

		for exp = range value {
			arg = value[exp]
		}
		return &Condition{query: fmt.Sprintf("%s = ?", flavor.Quote(exp)), args: []interface{}{arg}}, nil
	default:
		return nil, errors.WithStack(stmt.ErrorUnsupportedType)
	}
}
func (flavor *MysqlCondition) GreaterThan(v interface{}) (stmt.ConditionResult, error) {
	switch value := v.(type) {
	case map[string]interface{}:
		var exp string
		var arg interface{}

		for exp = range value {
			arg = value[exp]
		}
		return &Condition{query: fmt.Sprintf("%s > ?", flavor.Quote(exp)), args: []interface{}{arg}}, nil
	default:
		return nil, errors.WithStack(stmt.ErrorUnsupportedType)
	}
}
func (flavor *MysqlCondition) LessThan(v interface{}) (stmt.ConditionResult, error) {
	switch value := v.(type) {
	case map[string]interface{}:
		var exp string
		var arg interface{}

		for exp = range value {
			arg = value[exp]
		}
		return &Condition{query: fmt.Sprintf("%s < ?", flavor.Quote(exp)), args: []interface{}{arg}}, nil
	default:
		return nil, errors.WithStack(stmt.ErrorUnsupportedType)
	}
}

func (flavor *MysqlCondition) GreaterThanOrEqual(v interface{}) (stmt.ConditionResult, error) {
	switch value := v.(type) {
	case map[string]interface{}:
		var exp string
		var arg interface{}

		for exp = range value {
			arg = value[exp]
		}
		return &Condition{query: fmt.Sprintf("%s >= ?", flavor.Quote(exp)), args: []interface{}{arg}}, nil
	default:
		return nil, errors.WithStack(stmt.ErrorUnsupportedType)
	}
}

func (flavor *MysqlCondition) LessThanOrEqual(v interface{}) (stmt.ConditionResult, error) {
	switch value := v.(type) {
	case map[string]interface{}:
		var exp string
		var arg interface{}

		for exp = range value {
			arg = value[exp]
		}
		return &Condition{query: fmt.Sprintf("%s <= ?", flavor.Quote(exp)), args: []interface{}{arg}}, nil
	default:
		return nil, errors.WithStack(stmt.ErrorUnsupportedType)
	}
}
func (flavor *MysqlCondition) Like(v interface{}) (stmt.ConditionResult, error) {
	switch value := v.(type) {
	case map[string]interface{}:
		var exp string
		var arg interface{}

		for exp = range value {
			arg = value[exp]
		}
		return &Condition{query: fmt.Sprintf("%s LIKE ?", flavor.Quote(exp)), args: []interface{}{arg}}, nil
	default:
		return nil, errors.WithStack(stmt.ErrorUnsupportedType)
	}
}
func (flavor *MysqlCondition) IsNull(v interface{}) (stmt.ConditionResult, error) {
	switch value := v.(type) {
	// case string:
	// 	return &Condition{query: fmt.Sprintf("%s IS NULL", flavor.Quote(value)), args: []interface{}{}}, nil
	case map[string]interface{}:
		var exp string
		// var arg interface{}

		for exp = range value {
			// arg = value[exp]
		}
		return &Condition{query: fmt.Sprintf("%s IS NULL", flavor.Quote(exp)), args: []interface{}{}}, nil
	default:
		return nil, errors.WithStack(stmt.ErrorUnsupportedType)
	}
}
func (flavor *MysqlCondition) In(v interface{}) (stmt.ConditionResult, error) {
	opt_map := func(m map[string]interface{}) (string, []interface{}, error) {
		var key string
		var args []interface{}

		for key = range m {
			switch value := m[key].(type) {
			case []interface{}:
				if len(value) == 0 {
					return key, args, errors.Errorf("len(value) == 0")
				}
				args = append(args, value...)
			case interface{}:
				if value == nil {
					return key, args, errors.Errorf("value == nil")
				}
				args = append(args, value)
			default:
				return key, args, errors.WithStack(stmt.ErrorUnsupportedType)
			}
		}
		return key, args, nil
	}

	switch value := v.(type) {
	case map[string]interface{}:
		exp, args, err := opt_map(value)
		if err != nil {
			return nil, err
		}

		return &Condition{query: fmt.Sprintf("%s IN (%s)", flavor.Quote(exp), strings.Join(stmt.Repeat(len(args), "?"), ", ")), args: args}, nil
	default:
		return nil, errors.WithStack(stmt.ErrorUnsupportedType)
	}
}
func (flavor *MysqlCondition) Between(v interface{}) (stmt.ConditionResult, error) {
	opt_map := func(m map[string]interface{}) (string, []interface{}, error) {
		var key string
		var args []interface{}

		for key = range m {
			switch value := m[key].(type) {
			case []interface{}:
				if len(value) != 2 {
					return key, args, errors.Errorf("len(value) != 2")
				}

				args = append(args, value...)
			default:
				return key, args, errors.WithStack(stmt.ErrorUnsupportedType)
			}
		}
		return key, args, nil
	}

	switch value := v.(type) {
	case map[string]interface{}:
		exp, args, err := opt_map(value)
		if err != nil {
			return nil, err
		}

		return &Condition{query: fmt.Sprintf("%s BETWEEN ? AND ?", flavor.Quote(exp)), args: args}, nil
	default:
		return nil, errors.WithStack(stmt.ErrorUnsupportedType)
	}
}

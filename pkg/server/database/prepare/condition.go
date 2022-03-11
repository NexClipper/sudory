package prepare

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/database/prepare/sexp"
	"github.com/pkg/errors"
	"xorm.io/xorm"
)

type Condition struct {
	query string
	args  []interface{}
}

type ConditionFunctor func(v interface{}) (*Condition, error)

// NewCondition
func NewCondition(s string) (*Condition, error) {
	if len(s) == 0 {
		return nil, ErrorInvalidArgumentEmptyString()
	}

	s = strings.TrimSpace(s)

	if strings.Index(s, "(") == 0 {
		val, err := sexp.EvalString(s)
		if err != nil {
			return nil, errors.Wrapf(err, "sexp.EvalString q=%s", s)
		}

		v, ok := val.Val().(sexp.ArgsValueHolder)
		if !ok {
			return nil, ErrorUnsupportedType(val.Val())
		}
		return &Condition{query: v.String(), args: v.Args()}, nil
	} else {
		m := make(map[string]interface{})
		if err := json.Unmarshal([]byte(s), &m); err != nil {
			return nil, errors.Wrapf(err, "json.Unmarshal q=%s", s)
		}
		return newConditionBuilder().Parse(m)
	}
}

// NewConditionMap
func NewConditionMap(m map[string]interface{}) (*Condition, error) {
	return newConditionBuilder().Parse(m)
}

// Query
func (cond Condition) Query() string {
	return cond.query
}

// Args
func (cond Condition) Args() []interface{} {
	return cond.args
}

// Prepared
func (cond Condition) Prepared(tx *xorm.Session) *xorm.Session {
	return tx.Where(cond.Query(), cond.Args()...) //condition
}

type conditionBuild struct {
	engine map[string]ConditionFunctor
}

func newConditionBuilder() *conditionBuild {
	builder := conditionBuild{}

	engine := map[string]ConditionFunctor{
		"and":     builder.And,
		"or":      builder.Or,
		"not":     builder.Not,
		"equal":   builder.Equal,
		"eq":      builder.Equal, //addition: equal
		"gt":      builder.GreaterThan,
		"lt":      builder.LessThan,
		"gte":     builder.GreaterThanOrEqual,
		"ge":      builder.GreaterThanOrEqual,
		"lte":     builder.LessThanOrEqual,
		"le":      builder.LessThanOrEqual,
		"like":    builder.Like,
		"isnull":  builder.IsNull,
		"in":      builder.In,
		"between": builder.Between,
	}

	builder.engine = engine
	return &builder
}

func (builder *conditionBuild) Parse(v interface{}) (*Condition, error) {

	switch value := v.(type) {
	case map[string]interface{}:
		var cond *Condition
		var err error

		for key := range value {
			functor, ok := builder.engine[strings.ToLower(key)]
			if !ok {
				return nil, ErrorNotFoundHandler(key)
			}

			cond, err = functor(value[key])
			if err != nil {
				return nil, errors.Wrapf(err, "functor key=%s value=%+v", key, value[key])
			}
		}
		return cond, nil
	default:
		return nil, ErrorUnsupportedType(value)
	}
}

func (builder *conditionBuild) And(v interface{}) (*Condition, error) {
	switch value := v.(type) {
	case []interface{}:
		var str = make([]string, 0)
		var args = make([]interface{}, 0)

		for n := range value {
			cond, err := builder.Parse(value[n])
			if err != nil {
				return nil, err
			}

			str = append(str, cond.query)
			args = append(args, cond.args...)
		}
		return &Condition{query: "(" + strings.Join(str, " AND ") + ")", args: args}, nil
	case map[string]interface{}:
		var str = make([]string, 0)
		var args = make([]interface{}, 0)

		cond, err := builder.Parse(value)
		if err != nil {
			return nil, err
		}

		str = append(str, cond.query)
		args = append(args, cond.args...)

		return &Condition{query: "(" + strings.Join(str, " AND ") + ")", args: args}, nil
	default:
		return nil, ErrorUnsupportedType(value)
	}
}

func (builder *conditionBuild) Or(v interface{}) (*Condition, error) {
	switch value := v.(type) {
	case []interface{}:
		var str = make([]string, 0)
		var args = make([]interface{}, 0)

		for n := range value {
			cond, err := builder.Parse(value[n])
			if err != nil {
				return nil, err
			}

			str = append(str, cond.query)
			args = append(args, cond.args...)
		}
		return &Condition{query: "(" + strings.Join(str, " OR ") + ")", args: args}, nil
	case map[string]interface{}:
		var str = make([]string, 0)
		var args = make([]interface{}, 0)

		cond, err := builder.Parse(value)
		if err != nil {
			return nil, err
		}

		str = append(str, cond.query)
		args = append(args, cond.args...)

		return &Condition{query: "(" + strings.Join(str, " OR ") + ")", args: args}, nil
	default:
		return nil, ErrorUnsupportedType(value)
	}
}

func (builder *conditionBuild) Not(v interface{}) (*Condition, error) {
	switch value := v.(type) {
	case map[string]interface{}:
		cond, err := builder.Parse(value)
		if err != nil {
			return nil, err
		}

		return &Condition{query: fmt.Sprintf("NOT %s", cond.query), args: cond.args}, nil
	default:
		return nil, ErrorUnsupportedType(value)
	}
}

func (builder *conditionBuild) Equal(v interface{}) (*Condition, error) {
	switch value := v.(type) {
	case map[string]interface{}:
		var exp string
		var arg interface{}

		for exp = range value {
			arg = value[exp]
		}
		return &Condition{query: fmt.Sprintf("%s = ?", exp), args: []interface{}{arg}}, nil
	default:
		return nil, ErrorUnsupportedType(value)
	}
}

func (builder *conditionBuild) GreaterThan(v interface{}) (*Condition, error) {
	switch value := v.(type) {
	case map[string]interface{}:
		var exp string
		var arg interface{}

		for exp = range value {
			arg = value[exp]
		}
		return &Condition{query: fmt.Sprintf("%s > ?", exp), args: []interface{}{arg}}, nil
	default:
		return nil, ErrorUnsupportedType(value)
	}
}
func (builder *conditionBuild) LessThan(v interface{}) (*Condition, error) {
	switch value := v.(type) {
	case map[string]interface{}:
		var exp string
		var arg interface{}

		for exp = range value {
			arg = value[exp]
		}
		return &Condition{query: fmt.Sprintf("%s < ?", exp), args: []interface{}{arg}}, nil
	default:
		return nil, ErrorUnsupportedType(value)
	}
}

func (builder *conditionBuild) GreaterThanOrEqual(v interface{}) (*Condition, error) {
	switch value := v.(type) {
	case map[string]interface{}:
		var exp string
		var arg interface{}

		for exp = range value {
			arg = value[exp]
		}
		return &Condition{query: fmt.Sprintf("%s >= ?", exp), args: []interface{}{arg}}, nil
	default:
		return nil, ErrorUnsupportedType(value)
	}
}

func (builder *conditionBuild) LessThanOrEqual(v interface{}) (*Condition, error) {
	switch value := v.(type) {
	case map[string]interface{}:
		var exp string
		var arg interface{}

		for exp = range value {
			arg = value[exp]
		}
		return &Condition{query: fmt.Sprintf("%s <= ?", exp), args: []interface{}{arg}}, nil
	default:
		return nil, ErrorUnsupportedType(value)
	}
}

func (builder *conditionBuild) Like(v interface{}) (*Condition, error) {
	switch value := v.(type) {
	case map[string]interface{}:
		var exp string
		var arg interface{}

		for exp = range value {
			arg = value[exp]
		}
		return &Condition{query: fmt.Sprintf("%s LIKE ?", exp), args: []interface{}{arg}}, nil
	default:
		return nil, ErrorUnsupportedType(value)
	}
}

func (builder *conditionBuild) IsNull(v interface{}) (*Condition, error) {
	switch value := v.(type) {
	case string:
		return &Condition{query: fmt.Sprintf("%s IS ?", value), args: []interface{}{nil}}, nil
	default:
		return nil, ErrorUnsupportedType(value)
	}
}

func (builder *conditionBuild) In(v interface{}) (*Condition, error) {

	opt_map := func(m map[string]interface{}) (string, []interface{}, error) {
		var key string
		var args []interface{}

		for key = range m {
			switch value := m[key].(type) {
			case []interface{}:
				if len(value) == 0 {
					return key, args, fmt.Errorf("len(args) == 0")
				}
				args = append(args, value...)
			case interface{}:
				if value == nil {
					return key, args, fmt.Errorf("args == nil")
				}
				args = append(args, value)
			default:
				return key, args, ErrorUnsupportedType(value)
			}
		}
		return key, args, nil
	}

	makeQ := func(n int) string {
		s := make([]string, n)
		for i := 0; i < n; i++ {
			s[i] = "?"
		}
		return strings.Join(s, ", ")
	}

	switch value := v.(type) {
	case map[string]interface{}:
		exp, args, err := opt_map(value)
		if err != nil {
			return nil, errors.Wrapf(err, "operator map")
		}

		return &Condition{query: fmt.Sprintf("%s IN (%s)", exp, makeQ(len(args))), args: args}, nil
	default:
		return nil, ErrorUnsupportedType(value)
	}
}

func (builder *conditionBuild) Between(v interface{}) (*Condition, error) {

	opt_map := func(m map[string]interface{}) (string, []interface{}, error) {
		var key string
		var args []interface{}

		for key = range m {
			switch value := m[key].(type) {
			case []interface{}:
				if len(value) != 2 {
					return key, args, fmt.Errorf("args length != 2")
				}

				args = append(args, value...)
			default:
				return key, args, ErrorUnsupportedType(m)
			}
		}
		return key, args, nil
	}

	switch value := v.(type) {
	case map[string]interface{}:
		exp, args, err := opt_map(value)
		if err != nil {
			return nil, errors.Wrapf(err, "operator map")
		}

		return &Condition{query: fmt.Sprintf("%s BETWEEN ? AND ?", exp), args: args}, nil
	default:
		return nil, ErrorUnsupportedType(value)
	}
}

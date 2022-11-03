package sexp

import (
	"fmt"

	"github.com/pkg/errors"
)

// var quote = BackQuote

// func BackQuote(exp interface{}) string {
// 	return strings.ReplaceAll(fmt.Sprintf("`%s`", exp), "``", "`")
// }

func quote(exp interface{}) string {
	return fmt.Sprintf("%v", exp)
}

type TypeSlice = []interface{}
type TypeMap = map[string]interface{}

// Slice
func Slice(e interface{}, es ...interface{}) TypeSlice {
	return append([]interface{}{e}, es...)
}

// Map
func Map(a string, b interface{}) TypeMap {
	return TypeMap{a: b}
}

// NewArgsValue
func NewArgsValue(v interface{}) Value {
	return Value{argsValue, v}
}

func CheckValueType(a valueType) func(b valueType) bool {
	return func(b valueType) bool {
		return a == b
	}
}

var CheckArgsValueType = CheckValueType(argsValue)

type Builtin struct{}

var builtin = Builtin{}

// var builtin_commands = map[string]string{
// 	"+":       "Add",
// 	"-":       "Sub",
// 	"*":       "Mul",
// 	">":       "Gt",
// 	"<":       "Lt",
// 	">=":      "Gte",
// 	"<=":      "Lte",
// 	"display": "Display",
// 	"cons":    "Cons",
// 	"car":     "Car",
// 	"cdr":     "Cdr",
// }

// func (Builtin) Display(vars ...Value) (Value, error) {
// 	if len(vars) == 1 {
// 		fmt.Println(vars[0])
// 	} else {
// 		return Nil, fmt.Errorf("Badly formatted arguments: %v", vars)
// 	}
// 	return Nil, nil
// }

// func (Builtin) Cons(vars ...Value) (Value, error) {
// 	if len(vars) == 2 {
// 		cons := Cons{&vars[0], &vars[1]}
// 		return Value{consValue, &cons}, nil
// 	} else {
// 		return Nil, fmt.Errorf("Badly formatted arguments: %v", vars)
// 	}
// }

// func (Builtin) Car(vars ...Value) (Value, error) {
// 	if len(vars) == 1 && vars[0].typ == consValue {
// 		cons := vars[0].Cons()
// 		return *cons.car, nil
// 	} else {
// 		return Nil, fmt.Errorf("Badly formatted arguments: %v", vars)
// 	}
// }

// func (Builtin) Cdr(vars ...Value) (Value, error) {
// 	if len(vars) == 1 && vars[0].typ == consValue {
// 		cons := vars[0].Cons()
// 		return *cons.cdr, nil
// 	} else {
// 		return Nil, fmt.Errorf("Badly formatted arguments: %v", vars)
// 	}
// }

// func (Builtin) Add(vars ...Value) (Value, error) {
// 	var sum float64
// 	for _, v := range vars {
// 		if v.typ == numberValue {
// 			sum += v.Number()
// 		} else {
// 			return Nil, fmt.Errorf("Badly formatted arguments: %v", vars)
// 		}
// 	}
// 	return Value{numberValue, sum}, nil
// }

// func (Builtin) Sub(vars ...Value) (Value, error) {
// 	if vars[0].typ != numberValue {
// 		return Nil, fmt.Errorf("Badly formatted arguments: %v", vars)
// 	}
// 	sum := vars[0].Number()
// 	for _, v := range vars[1:] {
// 		if v.typ == numberValue {
// 			sum -= v.Number()
// 		} else {
// 			return Nil, fmt.Errorf("Badly formatted arguments: %v", vars)
// 		}
// 	}
// 	return Value{numberValue, sum}, nil
// }

// func (Builtin) Mul(vars ...Value) (Value, error) {
// 	if vars[0].typ != numberValue {
// 		return Nil, fmt.Errorf("Badly formatted arguments: %v", vars)
// 	}
// 	sum := vars[0].Number()
// 	for _, v := range vars[1:] {
// 		if v.typ == numberValue {
// 			sum *= v.Number()
// 		} else {
// 			return Nil, fmt.Errorf("Badly formatted arguments: %v", vars)
// 		}
// 	}
// 	return Value{numberValue, sum}, nil
// }

// func (Builtin) Gt(vars ...Value) (Value, error) {
// 	for i := 1; i < len(vars); i++ {
// 		v1 := vars[i-1]
// 		v2 := vars[i]
// 		if v1.typ != numberValue || v2.typ != numberValue {
// 			return Nil, fmt.Errorf("Badly formatted arguments: %v", vars)
// 		} else if !(v1.Number() > v2.Number()) {
// 			return False, nil
// 		}
// 	}
// 	return True, nil
// }

// func (Builtin) Lt(vars ...Value) (Value, error) {
// 	for i := 1; i < len(vars); i++ {
// 		v1 := vars[i-1]
// 		v2 := vars[i]
// 		if v1.typ != numberValue || v2.typ != numberValue {
// 			return Nil, fmt.Errorf("Badly formatted arguments: %v", vars)
// 		} else if !(v1.Number() < v2.Number()) {
// 			return False, nil
// 		}
// 	}
// 	return True, nil
// }

// func (Builtin) Gte(vars ...Value) (Value, error) {
// 	for i := 1; i < len(vars); i++ {
// 		v1 := vars[i-1]
// 		v2 := vars[i]
// 		if v1.typ != numberValue || v2.typ != numberValue {
// 			return Nil, fmt.Errorf("Badly formatted arguments: %v", vars)
// 		} else if !(v1.Number() >= v2.Number()) {
// 			return False, nil
// 		}
// 	}
// 	return True, nil
// }

// func (Builtin) Lte(vars ...Value) (Value, error) {
// 	for i := 1; i < len(vars); i++ {
// 		v1 := vars[i-1]
// 		v2 := vars[i]
// 		if v1.typ != numberValue || v2.typ != numberValue {
// 			return Nil, fmt.Errorf("Badly formatted arguments: %v", vars)
// 		} else if !(v1.Number() <= v2.Number()) {
// 			return False, nil
// 		}
// 	}
// 	return True, nil
// }

type funcBuiltin func(vars ...Value) (Value, error)

var builtin_commands = map[string]funcBuiltin{
	"and":     builtin.And,
	"or":      builtin.Or,
	"not":     builtin.Not,
	"equal":   builtin.Equal,
	"eq":      builtin.Equal, // addition: equal
	"gt":      builtin.GreaterThan,
	"lt":      builtin.LessThan,
	"gte":     builtin.GreaterThanOrEqual,
	"ge":      builtin.GreaterThanOrEqual, // addition: gte
	"lte":     builtin.LessThanOrEqual,
	"le":      builtin.LessThanOrEqual, // addition: lte
	"like":    builtin.Like,
	"isnull":  builtin.IsNull,
	"in":      builtin.In,
	"between": builtin.Between,
}

var ErrorEmptyArgument = fmt.Errorf("empty argument")
var ErrorBablyFormattedArguments = fmt.Errorf("badly formatted arguments")
var ErrorUnsupportedType = fmt.Errorf("unsupported type")

func ErrorArgumentFormat(v interface{}) string {
	return fmt.Sprintf("argument=%T(%v)", v, v)
}

func (Builtin) Not(vars ...Value) (Value, error) {
	if len(vars) != 1 {
		return Nil, errors.Wrapf(ErrorBablyFormattedArguments, ErrorArgumentFormat(vars))
	}

	iter := vars[0]
	if !CheckArgsValueType(iter.typ) {
		return Nil, errors.Wrapf(ErrorUnsupportedType, ErrorArgumentFormat(iter))
	}

	return NewArgsValue(Map("not", iter.val)), nil
}

func (Builtin) And(vars ...Value) (Value, error) {
	if len(vars) == 0 {
		return Nil, errors.Wrapf(ErrorEmptyArgument, ErrorArgumentFormat(vars))
	}

	values := make([]interface{}, 0, len(vars))
	for _, iter := range vars {
		if !CheckArgsValueType(iter.typ) {
			return Nil, errors.Wrapf(ErrorUnsupportedType, ErrorArgumentFormat(iter))
		}

		values = append(values, iter.val)
	}

	return NewArgsValue(Map("and", values)), nil
}

func (Builtin) Or(vars ...Value) (Value, error) {
	if len(vars) == 0 {
		return Nil, errors.Wrapf(ErrorEmptyArgument, ErrorArgumentFormat(vars))
	}

	values := make([]interface{}, 0, len(vars))
	for _, iter := range vars {
		if !CheckArgsValueType(iter.typ) {
			return Nil, errors.Wrapf(ErrorUnsupportedType, ErrorArgumentFormat(iter))
		}

		values = append(values, iter.val)
	}

	return NewArgsValue(Map("or", values)), nil
}

func (Builtin) Equal(vars ...Value) (Value, error) {
	if len(vars) != 2 {
		return Nil, errors.Wrapf(ErrorBablyFormattedArguments, ErrorArgumentFormat(vars))
	}

	return NewArgsValue(Map("equal", Map(quote(vars[0].val), vars[1].val))), nil
}

func (Builtin) GreaterThan(vars ...Value) (Value, error) {
	if len(vars) != 2 {
		return Nil, errors.Wrapf(ErrorBablyFormattedArguments, ErrorArgumentFormat(vars))
	}

	return NewArgsValue(Map("gt", Map(quote(vars[0].val), vars[1].val))), nil
}

func (Builtin) LessThan(vars ...Value) (Value, error) {
	if len(vars) != 2 {
		return Nil, errors.Wrapf(ErrorBablyFormattedArguments, ErrorArgumentFormat(vars))
	}

	return NewArgsValue(Map("lt", Map(quote(vars[0].val), vars[1].val))), nil
}

func (Builtin) GreaterThanOrEqual(vars ...Value) (Value, error) {
	if len(vars) != 2 {
		return Nil, errors.Wrapf(ErrorBablyFormattedArguments, ErrorArgumentFormat(vars))
	}

	return NewArgsValue(Map("gte", Map(quote(vars[0].val), vars[1].val))), nil
}

func (Builtin) LessThanOrEqual(vars ...Value) (Value, error) {
	if len(vars) != 2 {
		return Nil, errors.Wrapf(ErrorBablyFormattedArguments, ErrorArgumentFormat(vars))
	}

	return NewArgsValue(Map("lte", Map(quote(vars[0].val), vars[1].val))), nil
}

func (Builtin) Like(vars ...Value) (Value, error) {
	if len(vars) != 2 {
		return Nil, errors.Wrapf(ErrorBablyFormattedArguments, ErrorArgumentFormat(vars))
	}

	return NewArgsValue(Map("like", Map(quote(vars[0].val), vars[1].val))), nil
}

func (Builtin) IsNull(vars ...Value) (Value, error) {
	if len(vars) != 1 {
		return Nil, errors.Wrapf(ErrorBablyFormattedArguments, ErrorArgumentFormat(vars))
	}

	return NewArgsValue(Map("isnull", Map(quote(vars[0].val), nil))), nil
}

func (Builtin) In(vars ...Value) (Value, error) {
	if !(1 < len(vars)) {
		return Nil, errors.Wrapf(ErrorBablyFormattedArguments, ErrorArgumentFormat(vars))
	}

	var values []interface{}
	for _, iter := range vars[1:] {
		switch iter.typ {
		case consValue:
			if values == nil {
				values = make([]interface{}, 0, iter.Cons().Len())
			}
			iter.Cons().Map(func(v Value) (Value, error) {
				switch v.typ {
				case stringValue, numberValue:
					values = append(values, v.val)
				default:
					return Nil, errors.Wrapf(ErrorUnsupportedType, ErrorArgumentFormat(v.val))
				}
				return Nil, nil
			})
		default:
			if values == nil {
				values = make([]interface{}, 0, len(vars))
			}
			values = append(values, iter.val)
		}
	}

	return NewArgsValue(Map("in", Map(quote(vars[0].val), values))), nil
}

func (Builtin) Between(vars ...Value) (Value, error) {
	if !(1 < len(vars)) {
		return Nil, errors.Wrapf(ErrorBablyFormattedArguments, ErrorArgumentFormat(vars))
	}

	var values []interface{}
	for _, iter := range vars[1:] {

		switch iter.typ {
		case consValue:
			if values == nil {
				values = make([]interface{}, 0, iter.Cons().Len())
			}
			iter.Cons().Map(func(v Value) (Value, error) {
				switch v.typ {
				case stringValue, numberValue:
					values = append(values, v.val)
				default:
					return Nil, errors.Wrapf(ErrorUnsupportedType, ErrorArgumentFormat(v.val))
				}
				return Nil, nil
			})
		default:
			if values == nil {
				values = make([]interface{}, 0, len(vars))
			}
			values = append(values, iter.val)
		}
	}

	return NewArgsValue(Map("between", Map(quote(vars[0].val), values))), nil
}

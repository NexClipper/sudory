package sexp

import (
	"fmt"
	"strings"
)

var quote = BackQuote

func BackQuote(exp interface{}) string {
	return strings.ReplaceAll(fmt.Sprintf("`%s`", exp), "``", "`")
}

func NoneQuote(exp interface{}) string {
	return fmt.Sprintf("%s", exp)
}

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
	"eq":      builtin.Equal,
	"gt":      builtin.GreaterThan,
	"lt":      builtin.LessThan,
	"gte":     builtin.GreaterThanOrEqual,
	"ge":      builtin.GreaterThanOrEqual,
	"lte":     builtin.LessThanOrEqual,
	"le":      builtin.LessThanOrEqual,
	"like":    builtin.Like,
	"isnull":  builtin.IsNull,
	"in":      builtin.In,
	"between": builtin.Between,
}

func (Builtin) Not(vars ...Value) (Value, error) {
	if len(vars) == 1 {
		iter := vars[0]
		val, ok := iter.val.(ArgsValueHolder)
		if !ok {
			return Nil, fmt.Errorf("unsupported type value=%#v", iter)
		}
		return NewArgsValue(fmt.Sprintf("NOT %s", val.String()), val.Args()...), nil
	} else {
		return Nil, fmt.Errorf("badly formatted arguments=%v", vars)
	}
}

func (Builtin) And(vars ...Value) (Value, error) {
	if len(vars) == 0 {
		return Nil, fmt.Errorf("empty argument")
	} else {
		args := make([]interface{}, 0, len(vars))
		s := make([]string, 0, len(vars))

		for n := range vars {
			iter := vars[n]

			switch iter.typ {
			case argsValue:
				val, ok := iter.val.(ArgsValueHolder)
				if !ok {
					return Nil, fmt.Errorf("unsupported type value=%#v", iter)
				}

				s = append(s, val.String())
				args = append(args, val.Args()...)
			default:
				return Nil, fmt.Errorf("unsupported type value=%#v", iter)
			}
		}
		return NewArgsValue(fmt.Sprintf("(%s)", strings.Join(s, " AND ")), args...), nil
	}
}

func (Builtin) Or(vars ...Value) (Value, error) {
	if len(vars) == 0 {
		return Nil, fmt.Errorf("empty argument")
	} else {
		args := make([]interface{}, 0, len(vars))
		s := make([]string, 0, len(vars))

		for n := range vars {
			iter := vars[n]

			switch iter.typ {
			case argsValue:
				val, ok := iter.val.(ArgsValueHolder)
				if !ok {
					return Nil, fmt.Errorf("unsupported type value=%#v", iter)
				}

				s = append(s, val.String())
				args = append(args, val.Args()...)
			default:
				return Nil, fmt.Errorf("unsupported type value=%#v", iter)
			}
		}
		return NewArgsValue(fmt.Sprintf("(%s)", strings.Join(s, " OR ")), args...), nil
	}
}

func (Builtin) Equal(vars ...Value) (Value, error) {
	if len(vars) == 2 {
		return NewArgsValue(fmt.Sprintf("%s = ?", quote(vars[0].val)), vars[1].val), nil
	} else {
		return Nil, fmt.Errorf("badly formatted arguments=%v", vars)
	}
}

func (Builtin) GreaterThan(vars ...Value) (Value, error) {
	if len(vars) == 2 {
		return NewArgsValue(fmt.Sprintf("%s > ?", quote(vars[0].val)), vars[1].val), nil
	} else {
		return Nil, fmt.Errorf("badly formatted arguments=%v", vars)
	}
}

func (Builtin) LessThan(vars ...Value) (Value, error) {
	if len(vars) == 2 {
		return NewArgsValue(fmt.Sprintf("%s < ?", quote(vars[0].val)), vars[1].val), nil
	} else {
		return Nil, fmt.Errorf("badly formatted arguments=%v", vars)
	}
}

func (Builtin) GreaterThanOrEqual(vars ...Value) (Value, error) {
	if len(vars) == 2 {
		return NewArgsValue(fmt.Sprintf("%s >= ?", quote(vars[0].val)), vars[1].val), nil
	} else {
		return Nil, fmt.Errorf("badly formatted arguments=%v", vars)
	}
}

func (Builtin) LessThanOrEqual(vars ...Value) (Value, error) {
	if len(vars) == 2 {
		return NewArgsValue(fmt.Sprintf("%s <= ?", quote(vars[0].val)), vars[1].val), nil
	} else {
		return Nil, fmt.Errorf("badly formatted arguments=%v", vars)
	}
}

func (Builtin) Like(vars ...Value) (Value, error) {
	if len(vars) == 2 {
		return NewArgsValue(fmt.Sprintf("%s LIKE ?", quote(vars[0].val)), vars[1].val), nil
	} else {
		return Nil, fmt.Errorf("badly formatted arguments=%v", vars)
	}
}

func (Builtin) IsNull(vars ...Value) (Value, error) {
	if len(vars) == 1 {
		return NewArgsValue(fmt.Sprintf("%s IS ?", quote(vars[0].val)), nil), nil
	} else {
		return Nil, fmt.Errorf("badly formatted arguments=%v", vars)
	}
}

func (Builtin) In(vars ...Value) (Value, error) {

	makeQ := func(n int) string {
		s := make([]string, n)
		for i := 0; i < n; i++ {
			s[i] = "?"
		}
		return strings.Join(s, ", ")
	}

	if 1 < len(vars) {
		var args []interface{}
		for n := range vars[1:] {
			iter := vars[n+1]

			switch iter.typ {
			case consValue:
				if args == nil {
					args = make([]interface{}, 0, iter.Cons().Len())
				}
				iter.Cons().Map(func(v Value) (Value, error) {
					switch v.typ {
					case stringValue, numberValue:
						args = append(args, v.val)
					default:
						return Nil, fmt.Errorf("unsupported type value=%#v", v.val)
					}
					return Nil, nil
				})
			default:
				if args == nil {
					args = make([]interface{}, 0, len(vars))
				}
				args = append(args, iter.val)
			}
		}
		return NewArgsValue(fmt.Sprintf("%s IN (%s)", quote(vars[0].val), makeQ(len(args))), args...), nil
	} else {
		return Nil, fmt.Errorf("badly formatted arguments=%v", vars)
	}
}

func (Builtin) Between(vars ...Value) (Value, error) {
	if 1 < len(vars) {
		var args []interface{}
		for n := range vars[1:] {
			iter := vars[n+1]

			switch iter.typ {
			case consValue:
				if args == nil {
					args = make([]interface{}, 0, iter.Cons().Len())
				}
				iter.Cons().Map(func(v Value) (Value, error) {
					switch v.typ {
					case stringValue, numberValue:
						args = append(args, v.val)
					default:
						return Nil, fmt.Errorf("unsupported type value=%#v", v.val)
					}
					return Nil, nil
				})
			default:
				if args == nil {
					args = make([]interface{}, 0, len(vars))
				}
				args = append(args, iter.val)
			}
		}
		return NewArgsValue(fmt.Sprintf("%s BETWEEN ? AND ?", quote(vars[0].val)), args...), nil
	} else {
		return Nil, fmt.Errorf("badly formatted arguments=%v", vars)
	}
}

package stmt

import (
	"strings"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/pkg/errors"
)

type Order struct {
	order   string
	columns []string
}

// Order
func (order Order) Order() string {
	quote := stmt.MapQuote(stmt.BackQuote)

	//컬럼이 없으면 empty string 리턴
	//값이 없으면 Prepare에서 거른다
	if len(order.columns) == 0 {
		return ""
	}
	return strings.Join([]string{strings.Join(quote(order.columns...), ", "), order.order}, " ")
}

func (order Order) Combine(other stmt.OrderResult) stmt.OrderResult {
	return CombineOrder([]stmt.OrderResult{order, other})
}

type CombineOrder []stmt.OrderResult

func (orders CombineOrder) Order() string {

	o := make([]string, 0, len(orders))

	for i := range orders {
		o = append(o, orders[i].Order())
	}

	return strings.Join(o, ", ")
}

func (order CombineOrder) Combine(other stmt.OrderResult) stmt.OrderResult {
	return append(order, other)
}

type MysqlOrder struct {
	stmt.OrderStmtBuilder
	// Asc  func(v interface{}) (stmt.OrderResult, error)
	// Desc func(v interface{}) (stmt.OrderResult, error)
}

func NewMysqlOrder() *MysqlOrder {
	flavor := new(MysqlOrder)
	// flavor.Asc = MakeOrderFunc("ASC")
	// flavor.Desc = MakeOrderFunc("DESC")
	flavor.OrderStmtBuilder = stmt.NewOrderBuildEngine(__DIALECT__, map[string]stmt.OrderFunctor{
		"asc":  flavor.Asc,
		"desc": flavor.Desc,
	})

	return flavor
}

var (
	order_asc  = MakeOrderFunc("ASC")
	order_desc = MakeOrderFunc("DESC")
)

func (flavor *MysqlOrder) Asc(v interface{}) (stmt.OrderResult, error) {
	return order_asc(v)
}
func (flavor *MysqlOrder) Desc(v interface{}) (stmt.OrderResult, error) {
	return order_desc(v)
}

func MakeOrderFunc(order string) func(v interface{}) (stmt.OrderResult, error) {
	scan := func(emun []interface{}) ([]string, error) {
		s := make([]string, len(emun))

		for n := range emun {
			switch value := emun[n].(type) {
			case string:
				s[n] = value
			default:
				return nil, errors.WithStack(stmt.ErrorUnsupportedType)
			}
		}
		return s, nil
	}

	return func(v interface{}) (stmt.OrderResult, error) {
		switch value := v.(type) {
		case []interface{}:
			columns, err := scan(value)
			if err != nil {
				return nil, errors.Wrapf(err, "scan %T", value)
			}
			return &Order{order: order, columns: columns}, nil
		case []string:
			return &Order{order: order, columns: value}, nil
		case string:
			return &Order{order: order, columns: []string{value}}, nil
		default:
			return nil, errors.WithStack(stmt.ErrorUnsupportedType)
		}
	}
}

package prepare

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

type Order struct {
	order   string
	columns []string
}

type Orders []Order

func (orders Orders) Order() string {

	o := make([]string, 0, 2)

	for i := range orders {
		o = append(o, orders[i].Order())
	}

	return strings.Join(o, ", ")
}

type OrdersFunctor = func(v interface{}) (*Orders, error)

// NewOrder
func NewOrder(s string) (*Orders, error) {
	if len(s) == 0 {
		return nil, ErrorInvalidArgumentEmptyString()
	}

	v := new(interface{})
	if err := json.Unmarshal([]byte(s), v); err != nil {
		return nil, errors.Wrapf(err, "json.Unmarshal q=%s", s)
	}

	return newOrderBuilder().Parse(*v)
}

// NewOrderMap
func NewOrderMap(m map[string]interface{}) (*Orders, error) {
	return newOrderBuilder().parseMap(m)
}

func NewOrderSlice(s []interface{}) (*Orders, error) {
	return newOrderBuilder().parseSlice(s)
}

// Order
func (order Order) Order() string {
	//컬럼이 없으면 empty string 리턴
	//값이 없으면 Prepare에서 거른다
	if len(order.columns) == 0 {
		return ""
	}
	return strings.Join([]string{strings.Join(order.columns, ", "), order.order}, " ")
}

type orderEngine struct {
	builtin map[string]OrdersFunctor
}

func newOrderBuilder() *orderEngine {
	builder := orderEngine{}

	builtin := map[string]OrdersFunctor{
		"asc":  builder.MakeOrderFunc("ASC"),
		"desc": builder.MakeOrderFunc("DESC"),
	}

	builder.builtin = builtin
	return &builder
}

func (builder *orderEngine) parseMap(m map[string]interface{}) (*Orders, error) {
	orderslice := make([]Order, 0)

	for key := range m {
		functor, ok := builder.builtin[strings.ToLower(key)]
		if !ok {
			return nil, ErrorNotFoundHandler(key)
		}
		orders, err := functor(m[key])
		if err != nil {
			return nil, errors.Wrapf(err, "functor key=%s value=%+v", key, m[key])
		}
		orderslice = append(orderslice, ([]Order)(*orders)...)
	}
	return (*Orders)(&orderslice), nil
}

func (builder *orderEngine) parseSlice(enum []interface{}) (*Orders, error) {
	orderslice := make([]Order, 0)

	for n := range enum {
		switch elem := enum[n].(type) {
		case string:
			//incase; ["foo","foobar", "DESC", "bar", "barfoo", "ASC"]
			//스트링 배열입력
			//
			//슬라이스 초기화
			//미리 만들어 둔다
			//따라오는 코드는 미리 만들어둔 슬라이스의 마지막 주소의 객체를 갱신
			if len(orderslice) == 0 {
				orderslice = append(orderslice, Order{})
			}

			// is this element functor?
			functor, ok := builder.builtin[strings.ToLower(elem)]
			if !ok {
				//append columns
				orderslice[len(orderslice)-1].columns = append(orderslice[len(orderslice)-1].columns, elem)
			} else {
				//펑터를 찾으면, 펑터의 위치가 가장 마지막에 온다는 가정으로 동작
				//펑터를 실행하면, 다름 객체를 append하여 다음 배열값에 대비
				//
				//펑터를 실행하고 결과값으로 슬라이스의 마지막 주소 객체에 갱신

				//exec functor
				orders, err := functor(orderslice[len(orderslice)-1].columns)
				if err != nil {
					return nil, errors.Wrapf(err, "functor key=%s value=%+v", elem, orderslice[len(orderslice)-1].columns)
				}
				//over write last order element
				orders_ := ([]Order)(*orders)
				for n := range orders_ {
					orderslice[len(orderslice)-1].order = orders_[n].order
					orderslice[len(orderslice)-1].columns = orders_[n].columns
				}
				//check last
				//if exist more element add tail
				if n < len(enum)-1 {
					orderslice = append(orderslice, Order{}) //add new one
				}
			}
		default:
			//incase; element is array map[string]interface{}
			orders, err := builder.Parse(elem)
			if err != nil {
				return nil, err
			}
			orderslice = append(orderslice, ([]Order)(*orders)...)
		}

	}
	return (*Orders)(&orderslice), nil
}

func (builder *orderEngine) Parse(v interface{}) (*Orders, error) {

	switch value := v.(type) {
	case map[string]interface{}:
		orders, err := builder.parseMap(value)
		if err != nil {
			return nil, errors.Wrapf(err, "operator map")
		}
		return orders, nil
	case []interface{}:
		orders, err := builder.parseSlice(value)
		if err != nil {
			return nil, errors.Wrapf(err, "operator slice")
		}
		return orders, nil
	default:
		return nil, ErrorUnsupportedType(value)
	}

}

func (builder *orderEngine) MakeOrderFunc(order string) func(v interface{}) (*Orders, error) {

	scan := func(emun []interface{}) ([]string, error) {
		s := make([]string, len(emun))

		for n := range emun {
			switch value := emun[n].(type) {
			case string:
				s[n] = value
			default:
				return nil, ErrorUnsupportedType(value)
			}
		}
		return s, nil
	}

	return func(v interface{}) (*Orders, error) {
		switch value := v.(type) {
		case []interface{}:
			columns, err := scan(value)
			if err != nil {
				return nil, errors.Wrapf(err, "scan value=%+v", value)
			}
			return (*Orders)(&[]Order{{order: order, columns: columns}}), nil
		case []string:
			return (*Orders)(&[]Order{{order: order, columns: value}}), nil
		case string:
			return (*Orders)(&[]Order{{order: order, columns: []string{value}}}), nil
		default:
			return nil, ErrorUnsupportedType(value)
		}
	}
}

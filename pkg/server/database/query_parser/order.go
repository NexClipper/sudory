package query_parser

type Order struct {
	order string

	// query map[string]interface{}
}

func NewOrder(m map[string]interface{}) (*Order, map[string]interface{}) {

	var (
		found bool                   = false
		query map[string]interface{} = make(map[string]interface{})
		order string                 = ""
	)

	for n := range m {
		query[n] = m[n]
	}

	if _, ok := m[__ORDER_ORDER__]; ok {
		order, _ = m[__ORDER_ORDER__].(string)

		found = true
	}

	delete(query, __ORDER_ORDER__)

	if !found {
		return nil, query
	}
	return &Order{
		// query: query,
		order: order,
	}, query
}

// Order
func (order Order) Order() string {
	return order.order
}

// // Query
// func (order Order) Query() map[string]interface{} {
// 	return order.query
// }

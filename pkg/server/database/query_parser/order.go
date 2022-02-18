package query_parser

type Order struct {
	order string

	// query map[string]interface{}
}

func NewOrder(m map[string]string) (*Order, map[string]string) {

	var (
		found bool   = false
		query        = make(map[string]string)
		order string = ""
	)

	for n := range m {
		query[n] = m[n]
	}

	if _, ok := m[__ORDER_ORDER__]; ok {
		order = m[__ORDER_ORDER__]

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

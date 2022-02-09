package query_parser

import "xorm.io/xorm"

type QueryParser struct {
	condition  *Condition
	pagination *Pagination
	order      *Order
}

func NewQueryParser(m map[string]interface{}, filter ConditionFilter) *QueryParser {

	//pagination
	pagination, m := NewPagination(m)
	//order
	order, m := NewOrder(m)
	//Condition
	condition := NewCondition(m, filter)

	return &QueryParser{
		condition:  condition,
		pagination: pagination,
		order:      order,
	}
}

func (query QueryParser) Condition() *Condition {
	return query.condition
}

func (query QueryParser) Pagination() *Pagination {
	return query.pagination
}
func (query QueryParser) Order() *Order {
	return query.order
}

func (query QueryParser) Prepare(tx *xorm.Session) *xorm.Session {
	if query.Condition() != nil {
		tx = tx.Where(query.Condition().Where(), query.Condition().Args()...) //condition
	}
	if query.Pagination() != nil {
		tx = tx.Limit(query.Pagination().Limit(), query.Pagination().Offset()) //pagination
	}
	if query.Order() != nil {
		tx = tx.OrderBy(query.Order().Order()) //pagination
	}
	return tx
}

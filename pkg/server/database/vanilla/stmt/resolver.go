package stmt

import (
	"sync"

	"github.com/pkg/errors"
)

// type mysql struct{}

// func (mysql) ParseCondition(stmt interface{}) (ConditionResult, error) {
// 	builder := resolvers.ConditionStmtBuilder("mysql")
// 	return builder.Build(stmt)
// }

// func (mysql) ParseOrder(stmt interface{}) (OrderResult, error) {
// 	builder := resolvers.OrderStmtBuilder("mysql")
// 	return builder.Build(stmt)
// }

// func (mysql) ParsePagination(stmt interface{}) (PaginationResult, error) {
// 	builder := resolvers.PaginationStmtBuilder("mysql")
// 	return builder.Build(stmt)
// }

// var (
// 	MySql = new(mysql)
// )

var (
	//
	resolvers = new(resolverCollection).init()
	// export for resolver's SetConditionStmtBuilder
	SetConditionStmtBuilder = resolvers.SetConditionStmtBuilder
	// export for resolver's SetOrderStmtBuilder
	SetOrderStmtBuilder = resolvers.SetOrderStmtBuilder
	// export for resolver's SetPaginationStmtBuilder
	SetPaginationStmtBuilder = resolvers.SetPaginationStmtBuilder

	// get condition resolver
	GetConditionStmtResolver = resolvers.ConditionStmtResolver
	// get order resolver
	GetOrderStmtResolver = resolvers.OrderStmtResolver
	// get pagination resolver
	GetPaginationStmtResolver = resolvers.PaginationStmtResolver
)

type ConditionStmtResolverSet = map[string]ConditionStmtResolver

type OrderStmtResolverSet = map[string]OrderStmtResolver

type PaginationStmtResolverSet = map[string]PaginationStmtResolver

type resolverCollection struct {
	mux sync.Mutex

	condition  ConditionStmtResolverSet
	order      OrderStmtResolverSet
	pagination PaginationStmtResolverSet
}

func (op *resolverCollection) init() *resolverCollection {
	op.mux.Lock()
	defer op.mux.Unlock()

	op.condition = make(ConditionStmtResolverSet)
	op.order = make(OrderStmtResolverSet)
	op.pagination = make(PaginationStmtResolverSet)
	return op
}

func (collection *resolverCollection) SetConditionStmtBuilder(dialect string, resolver ConditionStmtResolver) {
	collection.mux.Lock()
	defer collection.mux.Unlock()

	if resolver == nil {
		panic(errors.New("resolver is nil"))
	}

	if _, dup := collection.condition[dialect]; dup {
		panic(errors.New("register called twice"))
	}

	collection.condition[dialect] = resolver
}

func (collection *resolverCollection) SetOrderStmtBuilder(dialect string, resolver OrderStmtResolver) {
	collection.mux.Lock()
	defer collection.mux.Unlock()

	if resolver == nil {
		panic(errors.New("resolver is nil"))
	}

	if _, dup := collection.order[dialect]; dup {
		panic(errors.New("register called twice"))
	}

	collection.order[dialect] = resolver
}

func (collection *resolverCollection) SetPaginationStmtBuilder(dialect string, resolver PaginationStmtResolver) {
	collection.mux.Lock()
	defer collection.mux.Unlock()

	if resolver == nil {
		panic(errors.New("resolver is nil"))
	}

	if _, dup := collection.pagination[dialect]; dup {
		panic(errors.New("register called twice"))
	}

	collection.pagination[dialect] = resolver
}

func (collection *resolverCollection) ConditionStmtResolver(dialect string) ConditionStmtBuilder {
	collection.mux.Lock()
	defer collection.mux.Unlock()

	resolver, ok := collection.condition[dialect]
	if !ok {
		return &FakeConditionStmtBuildEngine{
			dialect: dialect,
			err:     errors.Errorf(`missing dialect="%v"`, dialect),
		}
	}

	return resolver
}

func (collection *resolverCollection) OrderStmtResolver(dialect string) OrderStmtBuilder {
	collection.mux.Lock()
	defer collection.mux.Unlock()

	resolver, ok := collection.order[dialect]
	if !ok {
		return &FakeOrderStmtBuildEngine{
			dialect: dialect,
			err:     errors.Errorf(`missing dialect="%v"`, dialect),
		}
	}

	return resolver
}

func (collection *resolverCollection) PaginationStmtResolver(dialect string) PaginationStmtBuilder {
	collection.mux.Lock()
	defer collection.mux.Unlock()

	resolver, ok := collection.pagination[dialect]
	if !ok {
		return &FakePaginationStmtBuildEngine{
			dialect: dialect,
			err:     errors.Errorf(`missing dialect="%v"`, dialect),
		}
	}

	return resolver
}

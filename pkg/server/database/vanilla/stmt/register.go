package stmt

import (
	"sync"

	"github.com/pkg/errors"
)

var (
	//
	collection = new(collectionResolver).init()
	// export for resolver's SetConditionStmtBuilder
	SetConditionStmtBuilder = collection.SetConditionStmtBuilder
	// export for resolver's SetOrderStmtBuilder
	SetOrderStmtBuilder = collection.SetOrderStmtBuilder
	// export for resolver's SetPaginationStmtBuilder
	SetPaginationStmtBuilder = collection.SetPaginationStmtBuilder

	// get condition resolver
	GetConditionStmtBuilder = collection.ConditionStmtBuilder
	// get order resolver
	GetOrderStmtBuilder = collection.OrderStmtBuilder
	// get pagination resolver
	GetPaginationStmtBuilder = collection.PaginationStmtBuilder
)

type ConditionStmtResolverSet = map[string]ConditionStmtResolver

type OrderStmtResolverSet = map[string]OrderStmtResolver

type PaginationStmtResolverSet = map[string]PaginationStmtResolver

type collectionResolver struct {
	mux sync.Mutex

	condition  ConditionStmtResolverSet
	order      OrderStmtResolverSet
	pagination PaginationStmtResolverSet
}

func (op *collectionResolver) init() *collectionResolver {
	op.mux.Lock()
	defer op.mux.Unlock()

	op.condition = make(ConditionStmtResolverSet)
	op.order = make(OrderStmtResolverSet)
	op.pagination = make(PaginationStmtResolverSet)
	return op
}

func (collection *collectionResolver) SetConditionStmtBuilder(dialect string, resolver ConditionStmtResolver) {
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

func (collection *collectionResolver) SetOrderStmtBuilder(dialect string, resolver OrderStmtResolver) {
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

func (collection *collectionResolver) SetPaginationStmtBuilder(dialect string, resolver PaginationStmtResolver) {
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

func (collection *collectionResolver) ConditionStmtBuilder(dialect string) ConditionStmtBuilder {
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

func (collection *collectionResolver) OrderStmtBuilder(dialect string) OrderStmtBuilder {
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

func (collection *collectionResolver) PaginationStmtBuilder(dialect string) PaginationStmtBuilder {
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

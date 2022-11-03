package excute

import (
	"sync"

	"github.com/pkg/errors"
)

var (
	//
	collection = new(collectionSqlExcutor).init()
	// export for resolver's SetSqlExcutor
	SetSqlExcutor = collection.SetSqlExcutor

	// get condition resolver
	GetSqlExcutor = collection.SqlExcutor
)

type ConditionStmtResolverSet = map[string]SqlExcutor

type collectionSqlExcutor struct {
	mux sync.Mutex

	condition ConditionStmtResolverSet
}

func (op *collectionSqlExcutor) init() *collectionSqlExcutor {
	op.mux.Lock()
	defer op.mux.Unlock()

	op.condition = make(ConditionStmtResolverSet)
	return op
}

func (collection *collectionSqlExcutor) SetSqlExcutor(dialect string, resolver SqlExcutor) {
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

func (collection *collectionSqlExcutor) SqlExcutor(dialect string) SqlExcutor {
	collection.mux.Lock()
	defer collection.mux.Unlock()

	resolver, ok := collection.condition[dialect]
	if !ok {
		return &FakeResolver{
			err: errors.Errorf(`missing dialect="%v"`, dialect),
		}
	}

	return resolver
}

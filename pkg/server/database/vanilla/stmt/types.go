package stmt

type ResolverSelector interface {
	Dialect() string
}

type ConditionResult interface {
	Query() string
	Args() []interface{}
}

type PaginationResult interface {
	String() string
	Limit() (int, bool)
	Page() (int, bool)
	Offset() int
	SetLimit(int)
	SetPage(int)
}

type OrderResult interface {
	Order() string
	// Order() string
	Combine(OrderResult) OrderResult
}

type ConditoinContent interface {
	And(v interface{}) (ConditionResult, error)
	Or(v interface{}) (ConditionResult, error)
	Not(v interface{}) (ConditionResult, error)
	Equal(v interface{}) (ConditionResult, error)
	GreaterThan(v interface{}) (ConditionResult, error)
	LessThan(v interface{}) (ConditionResult, error)
	GreaterThanOrEqual(v interface{}) (ConditionResult, error)
	LessThanOrEqual(v interface{}) (ConditionResult, error)
	Like(v interface{}) (ConditionResult, error)
	IsNull(v interface{}) (ConditionResult, error)
	In(v interface{}) (ConditionResult, error)
	Between(v interface{}) (ConditionResult, error)
}

type OrderContent interface {
	Asc(v interface{}) (OrderResult, error)
	Desc(v interface{}) (OrderResult, error)
}

type PaginationContent interface {
	Limit(v interface{}) (PaginationResult, error)
	Page(v interface{}) (PaginationResult, error)
}

type ConditionStmtBuilder interface {
	Build(v interface{}) (ConditionResult, error)
	Dialect() string
}

type OrderStmtBuilder interface {
	Build(v interface{}) (OrderResult, error)
	Dialect() string
}

type PaginationStmtBuilder interface {
	Build(v interface{}) (PaginationResult, error)
	Dialect() string
}

type ConditionStmtResolver = interface {
	ConditionStmtBuilder
	ConditoinContent
}

type OrderStmtResolver = interface {
	OrderStmtBuilder
	OrderContent
}

type PaginationStmtResolver = interface {
	PaginationStmtBuilder
	PaginationContent
}

type FakeConditionStmtBuildEngine struct {
	dialect string
	err     error
}

func (engine FakeConditionStmtBuildEngine) Build(v interface{}) (ConditionResult, error) {
	return nil, engine.err
}
func (engine FakeConditionStmtBuildEngine) Dialect() string {
	return engine.dialect
}

type FakeOrderStmtBuildEngine struct {
	dialect string
	err     error
}

func (engine FakeOrderStmtBuildEngine) Build(v interface{}) (OrderResult, error) {
	return nil, engine.err
}
func (engine FakeOrderStmtBuildEngine) Dialect() string {
	return engine.dialect
}

type FakePaginationStmtBuildEngine struct {
	dialect string
	err     error
}

func (engine FakePaginationStmtBuildEngine) Build(v interface{}) (PaginationResult, error) {
	return nil, engine.err
}
func (engine FakePaginationStmtBuildEngine) Dialect() string {
	return engine.dialect
}

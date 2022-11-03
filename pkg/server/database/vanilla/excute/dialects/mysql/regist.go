package excute

import (
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	_ "github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/dialects/mysql"
)

const (
	__DIALECT__ = "mysql"
)

func NewPlaceHolderBuilder() func() string {
	const __SQL_PREPARED_STMT_PLACEHOLDER__ = "?"
	return func() string {
		return __SQL_PREPARED_STMT_PLACEHOLDER__
	}
}

func Dialect() string {
	return __DIALECT__
}

func init() {
	excute.SetSqlExcutor(__DIALECT__, &MySql{
		conditionStmtBuilder:  stmt.GetConditionStmtBuilder(Dialect()),
		orderStmtBuilder:      stmt.GetOrderStmtBuilder(Dialect()),
		paginationStmtBuilder: stmt.GetPaginationStmtBuilder(Dialect()),
	})
}

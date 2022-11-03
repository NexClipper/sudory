package excute

import (
	"context"
	"database/sql"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
)

type ConditionStmt = stmt.ConditionStmt
type OrderStmt = stmt.OrderStmt
type PaginationStmt = stmt.PaginationStmt

type Preparer = interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Prepare(query string) (*sql.Stmt, error)
}

type Scanner = interface {
	Scan(dest ...interface{}) error
}
type CallbackScanner = func(scan Scanner) error
type CallbackScannerWithIndex = func(scan Scanner, _ int) error

// SqlExcutor
type SqlExcutor interface {
	// Insert
	Insert(table_name string, columns []string, values ...[]interface{}) func(ctx context.Context, tx Preparer) (affected int64, lastid int64, err error)
	// InsertOrUpdate
	InsertOrUpdate(table_name string, insert_columns []string, update_columns []string, values ...[]interface{}) func(ctx context.Context, tx Preparer) (affected int64, lastid int64, err error)
	// Update
	Update(table_name string, keys_values map[string]interface{}, q ConditionStmt) func(ctx context.Context, tx Preparer) (affected int64, err error)
	// QueryRow
	QueryRow(table_name string, columns []string, q ConditionStmt, o OrderStmt, p PaginationStmt) func(ctx context.Context, tx Preparer) func(CallbackScanner) error
	// QueryRows
	QueryRows(table_name string, columns []string, q ConditionStmt, o OrderStmt, p PaginationStmt) func(ctx context.Context, tx Preparer) func(CallbackScannerWithIndex) error
	// Delete
	Delete(table_name string, q ConditionStmt) func(ctx context.Context, tx Preparer) (affected int64, err error)
	// Count
	Count(tableName string, cond stmt.ConditionStmt, page stmt.PaginationStmt) func(ctx context.Context, tx Preparer) (count int, err error)
	// Exist
	Exist(tableName string, cond stmt.ConditionStmt) func(ctx context.Context, tx Preparer) (exist bool, err error)
}

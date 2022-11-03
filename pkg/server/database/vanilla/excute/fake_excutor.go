package excute

import (
	"context"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
)

type FakeResolver struct {
	err error
}

func (flavor FakeResolver) Update(table_name string, keys_values map[string]interface{}, q ConditionStmt) func(ctx context.Context, tx Preparer) (affected int64, err error) {
	return func(ctx context.Context, tx Preparer) (affected int64, err error) {
		return 0, flavor.err
	}
}
func (flavor FakeResolver) Insert(table_name string, columns []string, values ...[]interface{}) func(ctx context.Context, tx Preparer) (affected int64, lastid int64, err error) {
	return func(ctx context.Context, tx Preparer) (affected int64, lastid int64, err error) {
		return 0, 0, flavor.err
	}
}

func (flavor FakeResolver) InsertOrUpdate(table_name string, insert_columns []string, update_columns []string, values ...[]interface{}) func(ctx context.Context, tx Preparer) (affected int64, lastid int64, err error) {
	return func(ctx context.Context, tx Preparer) (affected int64, lastid int64, err error) {
		return 0, 0, flavor.err
	}
}

func (flavor FakeResolver) Delete(table_name string, q ConditionStmt) func(ctx context.Context, tx Preparer) (affected int64, err error) {
	return func(ctx context.Context, tx Preparer) (affected int64, err error) {
		return 0, flavor.err
	}
}

func (flavor FakeResolver) QueryRow(table_name string, columns []string, q ConditionStmt, o OrderStmt, p PaginationStmt) func(ctx context.Context, tx Preparer) func(CallbackScanner) error {
	return func(ctx context.Context, tx Preparer) func(CallbackScanner) error {
		return func(cs CallbackScanner) error {
			return flavor.err
		}
	}
}

func (flavor FakeResolver) QueryRows(table_name string, columns []string, q ConditionStmt, o OrderStmt, p PaginationStmt) func(ctx context.Context, tx Preparer) func(CallbackScannerWithIndex) error {
	return func(ctx context.Context, tx Preparer) func(CallbackScannerWithIndex) error {
		return func(cs CallbackScannerWithIndex) error {
			return flavor.err
		}
	}
}

func (flavor FakeResolver) Count(tableName string, cond stmt.ConditionStmt, page stmt.PaginationStmt) func(ctx context.Context, tx Preparer) (count int, err error) {
	return func(ctx context.Context, tx Preparer) (count int, err error) {
		return 0, flavor.err
	}
}

func (flavor FakeResolver) Exist(tableName string, cond stmt.ConditionStmt) func(ctx context.Context, tx Preparer) (exist bool, err error) {
	return func(ctx context.Context, tx Preparer) (exist bool, err error) {
		return false, flavor.err
	}
}

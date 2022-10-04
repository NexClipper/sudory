package stmtex

import (
	"context"

	vanilla "github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/pkg/errors"
)

type ErrorExec struct {
	err error
}

func (exec ErrorExec) Exec(tx Preparer, dialect string) (affected int64, lastid int64, err error) {
	err = exec.err
	return
}

func (exec ErrorExec) ExecContext(ctx context.Context, tx Preparer, dialect string) (affected int64, lastid int64, err error) {
	err = exec.err
	return
}

func NewErrorExec(err error) *ErrorExec {
	return &ErrorExec{err: err}
}

type SqlInsertExec struct {
	table_name string
	columns    []string
	values     []interface{}
}

func Insert(table_name string, columns []string, values ...[]interface{}) SqlExecutor {
	if len(columns) == 0 {
		return NewErrorExec(errors.New("column_names length is not 0"))
	}
	if len(values) == 0 {
		return NewErrorExec(errors.New("column_values length is not 0"))
	}

	// flat column_values
	var args []interface{}
	for _, column_value := range values {
		if len(column_value)%len(columns) != 0 {
			return NewErrorExec(errors.New("column_names and column_values do not match in length"))
		}

		if args == nil {
			args = make([]interface{}, 0, len(values)*len(column_value))
		}
		args = append(args, column_value...)
	}

	return &SqlInsertExec{
		table_name: table_name,
		columns:    columns,
		values:     args,
	}
}

func (exec SqlInsertExec) Exec(tx Preparer, dialect string) (affected int64, lastid int64, err error) {
	return InsertContext(context.Background(), tx, exec.table_name, exec.columns, exec.values)
}

func (exec SqlInsertExec) ExecContext(ctx context.Context, tx Preparer, dialect string) (affected int64, lastid int64, err error) {
	return InsertContext(ctx, tx, exec.table_name, exec.columns, exec.values)
}

type SqlUpdateExec struct {
	table_name string
	columns    []string
	values     []interface{}
	cond       vanilla.ConditionStmt
}

func Update(table_name string, keys_values map[string]interface{}, q vanilla.ConditionStmt) SqlUpdator {
	keys := make([]string, 0, len(keys_values))
	values := make([]interface{}, 0, len(keys_values))
	for k, v := range keys_values {
		keys = append(keys, k)
		values = append(values, v)
	}

	return &SqlUpdateExec{
		table_name: table_name,
		columns:    keys,
		values:     values,
		cond:       q,
	}
}

func (exec SqlUpdateExec) Exec(tx Preparer, dialect string) (affected int64, err error) {
	var cond vanilla.ConditionResult
	cond, err = vanilla.GetConditionStmtResolver(dialect).Build(exec.cond)
	if err != nil {
		return
	}
	affected, _, err = UpdateContext(context.Background(), tx, exec.table_name, exec.columns, exec.values, cond)
	return
}

func (exec SqlUpdateExec) ExecContext(ctx context.Context, tx Preparer, dialect string) (affected int64, err error) {
	var cond vanilla.ConditionResult
	cond, err = vanilla.GetConditionStmtResolver(dialect).Build(exec.cond)
	if err != nil {
		return
	}

	affected, _, err = UpdateContext(ctx, tx, exec.table_name, exec.columns, exec.values, cond)
	return
}

type SqlInsertOrUpdateExec struct {
	table_name     string
	insert_columns []string
	update_columns []string
	values         []interface{}
}

func InsertOrUpdate(table_name string, insert_columns []string, update_columns []string, values ...[]interface{}) SqlExecutor {
	if len(insert_columns) == 0 {
		return NewErrorExec(errors.New("column_names length is not 0"))
	}
	if len(values) == 0 {
		return NewErrorExec(errors.New("column_values length is not 0"))
	}

	// flat column_values
	var args []interface{}
	for _, value := range values {
		if len(value)%len(insert_columns) != 0 {
			return NewErrorExec(errors.New("column_names and column_values do not match in length"))
		}

		if args == nil {
			args = make([]interface{}, 0, len(values)*len(value))
		}
		args = append(args, value...)
	}

	return SqlInsertOrUpdateExec{
		table_name:     table_name,
		insert_columns: insert_columns,
		update_columns: update_columns,
		values:         args,
	}
}

func (exec SqlInsertOrUpdateExec) Exec(tx Preparer, dialect string) (affected int64, lastid int64, err error) {
	return InsertOrUpdateContext(context.Background(), tx, exec.table_name, exec.insert_columns, exec.update_columns, exec.values)
}

func (exec SqlInsertOrUpdateExec) ExecContext(ctx context.Context, tx Preparer, dialect string) (affected int64, lastid int64, err error) {
	return InsertOrUpdateContext(ctx, tx, exec.table_name, exec.insert_columns, exec.update_columns, exec.values)
}

type SqlDeleteExec struct {
	table_name string
	cond       vanilla.ConditionStmt
}

func Delete(table_name string, q vanilla.ConditionStmt) SqlRemover {
	return SqlDeleteExec{
		table_name: table_name,
		cond:       q,
	}
}

func (exec SqlDeleteExec) Exec(tx Preparer, dialect string) (affected int64, err error) {
	var cond vanilla.ConditionResult
	cond, err = vanilla.GetConditionStmtResolver(dialect).Build(exec.cond)
	if err != nil {
		return
	}
	return DeleteContext(context.Background(), tx, exec.table_name, cond)
}

func (exec SqlDeleteExec) ExecContext(ctx context.Context, tx Preparer, dialect string) (affected int64, err error) {
	var cond vanilla.ConditionResult
	cond, err = vanilla.GetConditionStmtResolver(dialect).Build(exec.cond)
	if err != nil {
		return
	}
	return DeleteContext(ctx, tx, exec.table_name, cond)
}

package excute

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/excute"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/pkg/errors"
)

// type ConditionStmt = stmt.ConditionStmt
// type OrderStmt = stmt.OrderStmt
// type PaginationStmt = stmt.PaginationStmt

// type ConditionResult = stmt.ConditionResult
// type OrderResult = stmt.OrderResult
// type PaginationResult = stmt.PaginationResult

type Preparer = excute.Preparer

var Repeat = excute.Repeat

type MySql struct {
	conditionStmtBuilder  stmt.ConditionStmtBuilder
	orderStmtBuilder      stmt.OrderStmtBuilder
	paginationStmtBuilder stmt.PaginationStmtBuilder
}

var _ excute.SqlExcutor = (*MySql)(nil)

func (flavor MySql) Update(table_name string, keys_values map[string]interface{}, q stmt.ConditionStmt) func(ctx context.Context, tx Preparer) (affected int64, err error) {
	buildPlaceHolder := NewPlaceHolderBuilder()
	buildPlaceHolder_ := func(column_names []string) []string {
		ss := make([]string, len(column_names))
		for i := range column_names {
			ss[i] = column_names[i] + "=" + buildPlaceHolder()
		}
		return ss
	}

	keys := make([]string, 0, len(keys_values))
	values := make([]interface{}, 0, len(keys_values))
	for k, v := range keys_values {
		keys = append(keys, k)
		values = append(values, v)
	}

	cond, err := flavor.conditionStmtBuilder.Build(q)
	if err != nil {
		return func(ctx context.Context, tx Preparer) (affected int64, err error) {
			return 0, err
		}
	}

	var args []interface{} = values
	if cond != nil {
		args = append(values, cond.Args()...)
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, `UPDATE %v SET %v`,
		table_name,
		strings.Join(buildPlaceHolder_(keys), ","),
	)
	if cond != nil {
		fmt.Fprintf(&buf, "\nWHERE %v", cond.Query())
	}

	return func(ctx context.Context, tx Preparer) (affected int64, err error) {
		affected, _, err = excute.ExecContext(ctx, tx, buf.String(), args)
		return
	}

}

func (flavor MySql) Insert(table_name string, columns []string, values ...[]interface{}) func(ctx context.Context, tx Preparer) (affected int64, lastid int64, err error) {
	if len(columns) == 0 {
		return func(ctx context.Context, tx Preparer) (affected int64, lastid int64, err error) {
			return 0, 0, errors.New("column_names length is not 0")
		}
	}
	if len(values) == 0 {
		return func(ctx context.Context, tx Preparer) (affected int64, lastid int64, err error) {
			return 0, 0, errors.New("column_values length is not 0")
		}
	}

	// flat column_values
	var args []interface{}
	for _, column_value := range values {
		if len(column_value)%len(columns) != 0 {
			return func(ctx context.Context, tx Preparer) (affected int64, lastid int64, err error) {
				return 0, 0, errors.New("column_names and column_values do not match in length")
			}
		}

		if args == nil {
			args = make([]interface{}, 0, len(values)*len(column_value))
		}
		args = append(args, column_value...)
	}

	buildPlaceHolder := NewPlaceHolderBuilder()
	var buf bytes.Buffer
	fmt.Fprintf(&buf, `INSERT INTO %v (%v) VALUES %v`,
		table_name,
		strings.Join(columns, ","),
		strings.Join(
			Repeat(len(values),
				fmt.Sprintf("(%v)", strings.Join(
					Repeat(len(columns),
						buildPlaceHolder()), ","))), ","),
	)

	return func(ctx context.Context, tx Preparer) (affected int64, lastid int64, err error) {
		affected, lastid, err = excute.ExecContext(ctx, tx, buf.String(), args)
		return
	}
}

func (flavor MySql) InsertOrUpdate(table_name string, insert_columns []string, update_columns []string, values ...[]interface{}) func(ctx context.Context, tx Preparer) (affected int64, lastid int64, err error) {
	if len(insert_columns) == 0 {
		return func(ctx context.Context, tx Preparer) (affected int64, lastid int64, err error) {
			return 0, 0, errors.New("column_names length is not 0")
		}
	}
	if len(values) == 0 {
		return func(ctx context.Context, tx Preparer) (affected int64, lastid int64, err error) {
			return 0, 0, errors.New("column_values length is not 0")
		}
	}

	// flat column_values
	var args []interface{}
	for _, value := range values {
		if len(value)%len(insert_columns) != 0 {
			return func(ctx context.Context, tx Preparer) (affected int64, lastid int64, err error) {
				return 0, 0, errors.New("column_names and column_values do not match in length")
			}
		}

		if args == nil {
			args = make([]interface{}, 0, len(values)*len(value))
		}
		args = append(args, value...)
	}

	buildPlaceHolder := NewPlaceHolderBuilder()
	var buf bytes.Buffer
	fmt.Fprintf(&buf, `INSERT INTO %v (%v) VALUES %v`,
		table_name,
		strings.Join(insert_columns, ","),
		strings.Join(
			Repeat(len(values),
				fmt.Sprintf("(%v)", strings.Join(
					Repeat(len(insert_columns),
						buildPlaceHolder()), ","))), ","),
	)

	if 0 < len(update_columns) {
		s := make([]string, 0, len(update_columns))
		for i := range update_columns {
			update_column := update_columns[i]
			s = append(s, fmt.Sprintf("%v = VALUES(%v)", update_column, update_column))
		}
		fmt.Fprintf(&buf, "\nON DUPLICATE KEY UPDATE %v", strings.Join(s, ", "))
	}

	return func(ctx context.Context, tx Preparer) (affected int64, lastid int64, err error) {
		affected, lastid, err = excute.ExecContext(ctx, tx, buf.String(), args)
		return
	}
}

func (flavor MySql) Delete(table_name string, q stmt.ConditionStmt) func(ctx context.Context, tx Preparer) (affected int64, err error) {
	cond, err := flavor.conditionStmtBuilder.Build(q)
	if err != nil {
		return func(ctx context.Context, tx Preparer) (affected int64, err error) {
			return 0, err
		}
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, `DELETE FROM %v`, table_name)

	var args []interface{}
	if cond != nil {
		fmt.Fprintf(&buf, "\nWHERE %v", cond.Query())

		args = cond.Args()
	}
	return func(ctx context.Context, tx Preparer) (affected int64, err error) {
		affected, _, err = excute.ExecContext(ctx, tx, buf.String(), args)
		return
	}
}

type CallbackScanner = excute.CallbackScanner
type CallbackScannerWithIndex = excute.CallbackScannerWithIndex

func (flavor MySql) QueryRow(table_name string, columns []string, q stmt.ConditionStmt, o stmt.OrderStmt, p stmt.PaginationStmt) func(ctx context.Context, tx Preparer) func(CallbackScanner) error {
	build := func() (cond stmt.ConditionResult, order stmt.OrderResult, page stmt.PaginationResult, err error) {
		if q != nil {
			cond, err = flavor.conditionStmtBuilder.Build(q)
			if err != nil {
				err = errors.Wrapf(err, "%v=%+v",
					"cond", q,
				)
				return
			}

		}
		if o != nil {
			order, err = flavor.orderStmtBuilder.Build(o)
			if err != nil {
				err = errors.Wrapf(err, "%v=%+v",
					"order", o,
				)
				return
			}

		}
		if p != nil {
			page, err = flavor.paginationStmtBuilder.Build(p)
			if err != nil {
				err = errors.Wrapf(err, "%v=%+v",
					"page", p,
				)
				return
			}
		}
		return
	}

	cond, order, page, err := build()
	if err != nil {
		return func(ctx context.Context, tx Preparer) func(CallbackScanner) error {
			return func(cs CallbackScanner) error {
				return err
			}
		}
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, `SELECT %v FROM %v`,
		strings.Join(columns, ","),
		table_name,
	)
	var args []interface{}
	if cond != nil {
		fmt.Fprintf(&buf, "\nWHERE %v", cond.Query())
		args = cond.Args()
	}
	if order != nil {
		fmt.Fprintf(&buf, "\nORDER BY %v", order.Order())
	}
	if page != nil {
		limit := func(fn func() (int, bool)) int {
			a, _ := fn()
			return a
		}
		fmt.Fprintf(&buf, "\nLIMIT %v, %v", page.Offset(), limit(page.Limit))
	}

	return func(ctx context.Context, tx Preparer) func(CallbackScanner) error {
		return excute.QueryRowContext(ctx, tx, buf.String(), args)
	}
}

func (flavor MySql) QueryRows(table_name string, columns []string, q stmt.ConditionStmt, o stmt.OrderStmt, p stmt.PaginationStmt) func(ctx context.Context, tx Preparer) func(CallbackScannerWithIndex) error {
	build := func() (cond stmt.ConditionResult, order stmt.OrderResult, page stmt.PaginationResult, err error) {
		if q != nil {
			cond, err = flavor.conditionStmtBuilder.Build(q)
			if err != nil {
				err = errors.Wrapf(err, "%v=%+v",
					"cond", q,
				)
				return
			}

		}
		if o != nil {
			order, err = flavor.orderStmtBuilder.Build(o)
			if err != nil {
				err = errors.Wrapf(err, "%v=%+v",
					"order", o,
				)
				return
			}

		}
		if p != nil {
			page, err = flavor.paginationStmtBuilder.Build(p)
			if err != nil {
				err = errors.Wrapf(err, "%v=%+v",
					"page", p,
				)
				return
			}
		}
		return
	}

	cond, order, page, err := build()
	if err != nil {
		return func(ctx context.Context, tx Preparer) func(CallbackScannerWithIndex) error {
			return func(cs CallbackScannerWithIndex) error {
				return err
			}
		}
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, `SELECT %v FROM %v`,
		strings.Join(columns, ","),
		table_name,
	)
	var args []interface{}
	if cond != nil {
		fmt.Fprintf(&buf, "\nWHERE %v", cond.Query())
		args = cond.Args()
	}
	if order != nil {
		fmt.Fprintf(&buf, "\nORDER BY %v", order.Order())
	}
	if page != nil {
		limit := func(fn func() (int, bool)) int {
			a, _ := fn()
			return a
		}
		fmt.Fprintf(&buf, "\nLIMIT %v, %v", page.Offset(), limit(page.Limit))
	}

	return func(ctx context.Context, tx Preparer) func(CallbackScannerWithIndex) error {
		return excute.QueryRowsContext(ctx, tx, buf.String(), args)
	}
}

// Count
func (flavor MySql) Count(tableName string, cond stmt.ConditionStmt, page stmt.PaginationStmt) func(ctx context.Context, tx Preparer) (count int, err error) {
	return excute.Count(flavor, tableName, cond, page)
}

// Exist
func (flavor MySql) Exist(tableName string, cond stmt.ConditionStmt) func(ctx context.Context, tx Preparer) (exist bool, err error) {
	return excute.Exist(flavor, tableName, cond)
}

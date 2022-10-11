package stmtex

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	vanilla "github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
)

type SqlQuery struct {
	table_name string
	columns    []string
	cond       vanilla.ConditionStmt
	order      vanilla.OrderStmt
	page       vanilla.PaginationStmt
}

func Select(table_name string, columns []string, q vanilla.ConditionStmt, o vanilla.OrderStmt, p vanilla.PaginationStmt) SqlScanner {
	return &SqlQuery{
		table_name: table_name,
		columns:    columns,
		cond:       q,
		order:      o,
		page:       p,
	}
}

func (exec SqlQuery) Parse(dialect string) (s string, args []interface{}, err error) {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, `SELECT %v FROM %v`,
		strings.Join(exec.columns, ","),
		exec.table_name,
	)

	args = make([]interface{}, 0, __DEFAULT_ARGS_CAPACITY__)
	if exec.cond != nil {
		var cond vanilla.ConditionResult
		cond, err = vanilla.GetConditionStmtResolver(dialect).Build(exec.cond)
		if err != nil {
			err = errors.Wrapf(err, logs.KVL(
				"cond", exec.cond,
			))
			return
		}
		fmt.Fprintf(&buf, "\nWHERE %v", cond.Query())

		args = append(args, cond.Args()...)
	}
	if exec.order != nil {
		var order vanilla.OrderResult
		order, err = vanilla.GetOrderStmtResolver(dialect).Build(exec.order)
		if err != nil {
			err = errors.Wrapf(err, logs.KVL(
				"order", exec.order,
			))
			return
		}

		fmt.Fprintf(&buf, "\nORDER BY %v", order.Order())
	}
	if exec.page != nil {
		var page vanilla.PaginationResult
		page, err = vanilla.GetPaginationStmtResolver(dialect).Build(exec.page)
		if err != nil {
			err = errors.Wrapf(err, logs.KVL(
				"page", exec.page,
			))
			return
		}

		limit := func(fn func() (int, bool)) int {
			a, _ := fn()
			return a
		}
		fmt.Fprintf(&buf, "\nLIMIT %v, %v", page.Offset(), limit(page.Limit))
	}

	s = buf.String()

	return
}

func (exec SqlQuery) QueryRow(tx Preparer, dialect string) func(CallbackScanner) error {
	query, args, err := exec.Parse(dialect)
	if err != nil {
		return func(CallbackScanner) error { return err }
	}
	return QueryRowContext(context.Background(), tx, query, args)
}
func (exec SqlQuery) QueryRowContext(ctx context.Context, tx Preparer, dialect string) func(CallbackScanner) error {
	query, args, err := exec.Parse(dialect)
	if err != nil {
		return func(CallbackScanner) error { return err }
	}
	return QueryRowContext(ctx, tx, query, args)
}
func (exec SqlQuery) QueryRows(tx Preparer, dialect string) func(CallbackScannerWithIndex) error {
	query, args, err := exec.Parse(dialect)
	if err != nil {
		return func(CallbackScannerWithIndex) error { return err }
	}
	return QueryRowsContext(context.Background(), tx, query, args)
}
func (exec SqlQuery) QueryRowsContext(ctx context.Context, tx Preparer, dialect string) func(CallbackScannerWithIndex) error {
	query, args, err := exec.Parse(dialect)
	if err != nil {
		return func(CallbackScannerWithIndex) error { return err }
	}
	return QueryRowsContext(ctx, tx, query, args)
}

func Count(table_name string, q vanilla.ConditionStmt, p vanilla.PaginationStmt) func(ctx context.Context, tx Preparer, dialect string) (count int, err error) {
	builder := SqlQuery{
		table_name: table_name,
		columns:    []string{"COUNT(1)"},
		cond:       q,
		page:       p,
	}

	return func(ctx context.Context, tx Preparer, dialect string) (count int, err error) {
		err = builder.QueryRowContext(ctx, tx, dialect)(func(scan Scanner) error {
			err := scan.Scan(&count)
			err = errors.Wrapf(err, fmt.Sprintf("scan %v", builder.table_name))
			return err
		})

		return
	}
}

func ExistContext(table_name string, q vanilla.ConditionStmt) func(ctx context.Context, tx Preparer, dialect string) (exist bool, err error) {
	builder := SqlQuery{
		table_name: table_name,
		columns:    []string{"COUNT(1)"},
		cond:       q,
		page:       vanilla.Limit(1),
	}

	return func(ctx context.Context, tx Preparer, dialect string) (exist bool, err error) {
		var count int
		err = builder.QueryRowContext(ctx, tx, dialect)(func(scan Scanner) error {
			err := scan.Scan(&count)
			err = errors.Wrapf(err, fmt.Sprintf("scan %v", builder.table_name))
			return err
		})
		exist = 0 < count

		return
	}

}

func Exist(table_name string, q vanilla.ConditionStmt) func(tx Preparer, dialect string) (exist bool, err error) {
	builder := SqlQuery{
		table_name: table_name,
		columns:    []string{"COUNT(1)"},
		cond:       q,
		page:       vanilla.Limit(1),
	}

	return func(tx Preparer, dialect string) (exist bool, err error) {
		var count int
		err = builder.QueryRowContext(context.Background(), tx, dialect)(func(scan Scanner) error {
			err := scan.Scan(&count)
			err = errors.Wrapf(err, fmt.Sprintf("scan %v", builder.table_name))
			return err
		})
		exist = 0 < count

		return
	}

}

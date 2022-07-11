package vanilla

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/error_compose"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/prepare"
	"github.com/pkg/errors"
)

var __VANILLA_DEBUG_PRINT_STATMENT__ = func() func() bool {
	ok, _ := strconv.ParseBool(os.Getenv("VANILLA_DEBUG_PRINT_STATMENT"))
	return func() bool {
		return ok
	}
}()

type Preparer interface {
	Prepare(query string) (*sql.Stmt, error)
}

type Queryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func Exec(tx Preparer, query string, args []interface{}) (affected int64, err error) {
	if __VANILLA_DEBUG_PRINT_STATMENT__() {
		println(fmt.Sprintf("query=%v", query))
		println(fmt.Sprintf("args=%+v", args))
	}

	stmt, err := tx.Prepare(query)
	err = errors.Wrapf(err, "sql.Tx.Prepare")
	if err != nil {
		return
	}
	defer func() {
		err = error_compose.Composef(err, stmt.Close(), "sql.Stmt.Close")
	}()

	result, err := stmt.Exec(args...)
	err = errors.Wrapf(err, "sql.Stmt.Exec")
	if err != nil {
		return
	}

	affected, err = result.RowsAffected()
	err = errors.Wrapf(err, "sql.Result.RowsAffected")
	return
}

type Scanner interface {
	Scan(dest ...interface{}) error
}
type CallbackScanner = func(scan Scanner) error
type CallbackScannerWithIndex = func(scan Scanner, _ int) error

func QueryRow(tx Preparer, query string, args []interface{}) func(CallbackScanner) error {
	if __VANILLA_DEBUG_PRINT_STATMENT__() {
		println(fmt.Sprintf("query=%v", query))
		println(fmt.Sprintf("args=%+v", args))
	}

	return func(scan CallbackScanner) (err error) {
		stmt, err := tx.Prepare(query)
		err = errors.Wrapf(err, "sql.Tx.Prepare")
		if err != nil {
			return
		}
		defer func() {
			err = error_compose.Composef(err, stmt.Close(), "sql.Stmt.Close")
		}()

		row := stmt.QueryRow(args...)
		err = scan(row)
		err = errors.Wrapf(err, "sql.Row.Scan")
		err = error_compose.Composef(err, row.Err(), "sql.Row; during scan")

		err = errors.Wrapf(err, "faild to query row\nquery=\"%v\"\nargs=\"%+v\"",
			query,
			args,
		)
		return
	}
}

func QueryRows(tx Preparer, query string, args []interface{}) func(CallbackScannerWithIndex) error {
	if __VANILLA_DEBUG_PRINT_STATMENT__() {
		println(fmt.Sprintf("query=%v", query))
		println(fmt.Sprintf("args=%+v", args))
	}

	return func(scan CallbackScannerWithIndex) (err error) {
		stmt, err := tx.Prepare(query)
		err = errors.Wrapf(err, "sql.Tx.Prepare")
		if err != nil {
			return
		}
		defer func() {
			err = error_compose.Composef(err, stmt.Close(), "sql.Stmt.Close")
		}()

		var rows *sql.Rows
		rows, err = stmt.Query(args...)
		err = errors.Wrapf(err, "sql.Stmt.Query")
		if err != nil {
			return
		}

		defer func() {
			err = error_compose.Composef(err, rows.Close(), "sql.Rows.Close")
		}()
		i := 0
		for rows.Next() {
			err = scan(rows, i)
			err = errors.Wrapf(err, "sql.Row.Scan")
			if err != nil {
				break
			}
			i++
		}
		err = error_compose.Composef(err, rows.Err(), "sql.Rows; during scan")

		err = errors.Wrapf(err, "faild to query rows\nquery=\"%v\"\nargs=\"%+v\"",
			query,
			args,
		)
		return
	}
}

func QueryRows2(tx Queryer, query string, args []interface{}) func(CallbackScannerWithIndex) error {
	if __VANILLA_DEBUG_PRINT_STATMENT__() {
		println(fmt.Sprintf("query=%v", query))
		println(fmt.Sprintf("args=%+v", args))
	}

	return func(scan CallbackScannerWithIndex) (err error) {
		// stmt, err := tx.Prepare(query)
		// err = errors.Wrapf(err, "sql.Tx.Prepare")
		// if err != nil {
		// 	return
		// }
		// defer func() {
		// 	err = error_compose.Composef(err, stmt.Close(), "sql.Stmt.Close")
		// }()

		var rows *sql.Rows
		rows, err = tx.Query(query, args...)
		err = errors.Wrapf(err, "sql.Stmt.Query")
		if err != nil {
			return
		}

		defer func() {
			err = error_compose.Composef(err, rows.Close(), "sql.Rows.Close")
		}()
		i := 0
		for rows.Next() {
			err = scan(rows, i)
			err = errors.Wrapf(err, "sql.Row.Scan")
			if err != nil {
				break
			}
			i++
		}
		err = error_compose.Composef(err, rows.Err(), "sql.Rows; during scan")

		err = errors.Wrapf(err, "faild to query rows\nquery=\"%v\"\nargs=\"%+v\"",
			query,
			args,
		)
		return
	}
}

func Count(tx Preparer, table_name string, q *prepare.Condition) (count int, err error) {
	column_names := []string{
		"COUNT(1)",
	}

	err = Stmt.Select(table_name, column_names, q, nil, nil).
		QueryRow(tx)(func(s Scanner) (err error) {
		err = s.Scan(&count)
		err = errors.Wrapf(err, "scan count")
		return
	})

	err = errors.Wrapf(err, "faild to count table=\"%v\"", table_name)

	return
}

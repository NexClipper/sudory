package stmtex

import (
	"bytes"
	"context"
	"database/sql"
	"strings"
	"sync"

	"fmt"

	"github.com/NexClipper/sudory/pkg/server/database"
	vanilla "github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/pkg/errors"
)

const __DEFAULT_ARGS_CAPACITY__ = 10
const __SQL_PREPARED_STMT_PLACEHOLDER__ = "?"

func Repeat(n int, s string) []string {
	ss := make([]string, n)
	for i := 0; i < n; i++ {
		ss[i] = s
	}
	return ss
}

func Flat(e ...[]interface{}) ([]interface{}, error) {
	var f []interface{}

	for _, iter := range e {
		if len(e[0]) == len(iter) {
			return nil, errors.New("diff column length")
		}
		if f == nil {
			f = make([]interface{}, 0, len(e[0])*len(e))
		}

		f = append(f, iter...)
	}

	return f, nil
}

// type Preparer interface {
// 	stmt.ResolverSelector
// 	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
// }

// type sqlDbWithDialect struct {
// 	*sql.DB
// 	dialect string
// }

// func (sqldb *sqlDbWithDialect) Dialect() string {
// 	return sqldb.dialect
// }

// func NewSqlDB(db *sql.DB, dialect string) *sqlDbWithDialect {
// 	return &sqlDbWithDialect{
// 		DB:      db,
// 		dialect: dialect,
// 	}
// }

// type ScopePreparer interface {
// 	stmt.ResolverSelector
// 	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
// 	Commit() error
// 	Rollback() error
// }

// type sqlTxWithDialect struct {
// 	*sql.Tx
// 	dialect string
// }

// func (sqldb *sqlTxWithDialect) Dialect() string {
// 	return sqldb.dialect
// }

// func NewSqlTx(tx *sql.Tx, dialect string) *sqlTxWithDialect {
// 	return &sqlTxWithDialect{
// 		Tx:      tx,
// 		dialect: dialect,
// 	}
// }

type Preparer interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Prepare(query string) (*sql.Stmt, error)
}

type SqlUpdator interface {
	Exec(tx Preparer, dialect string) (affected int64, err error)
	ExecContext(ctx context.Context, tx Preparer, dialect string) (affected int64, err error)
}

type SqlExecutor interface {
	Exec(tx Preparer, dialect string) (affected int64, lastid int64, err error)
	ExecContext(ctx context.Context, tx Preparer, dialect string) (affected int64, lastid int64, err error)
}

type SqlRemover interface {
	Exec(tx Preparer, dialect string) (affected int64, err error)
	ExecContext(ctx context.Context, tx Preparer, dialect string) (affected int64, err error)
}

func ExecContext(ctx context.Context, tx Preparer, query string, args []interface{}) (affected int64, lastid int64, err error) {
	if __VANILLA_DEBUG_PRINT_STATMENT__() {
		println(fmt.Sprintf("query=%v", query))
		println(fmt.Sprintf("args=%+v", args))
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		err = errors.Wrapf(err, "sql.Tx.Prepare")
		return
	}
	defer func() {
		err = vanilla.ErrorCompose(err, errors.Wrapf(stmt.Close(), "sql.Stmt.Close"))
	}()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		err = errors.Wrapf(err, "sql.Stmt.Exec")
		return
	}

	affected, err = result.RowsAffected()
	if err != nil {
		err = errors.Wrapf(err, "sql.Result.RowsAffected")
		return
	}

	lastid, err = result.LastInsertId()
	if err != nil {
		err = errors.Wrapf(err, "sql.Result.LastInsertId")
		return
	}

	return
}

func InsertContext(ctx context.Context, tx Preparer, table_name string, columns []string, values []interface{}) (affected int64, lastid int64, err error) {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, `INSERT INTO %v (%v) VALUES %v`,
		table_name,
		strings.Join(columns, ","),
		strings.Join(
			Repeat(len(values)/len(columns),
				fmt.Sprintf("(%v)", strings.Join(
					Repeat(len(columns),
						__SQL_PREPARED_STMT_PLACEHOLDER__ /*placeholder=?*/), ","))), ","),
	)

	return ExecContext(ctx, tx, buf.String(), values)
}

func InsertOrUpdateContext(ctx context.Context, tx Preparer, table_name string, insert_columns, update_columns []string, values []interface{}) (affected int64, lastid int64, err error) {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, `INSERT INTO %v (%v) VALUES %v`,
		table_name,
		strings.Join(insert_columns, ","),
		strings.Join(
			Repeat(len(values)/len(insert_columns),
				fmt.Sprintf("(%v)", strings.Join(
					Repeat(len(insert_columns),
						__SQL_PREPARED_STMT_PLACEHOLDER__ /*placeholder=?*/), ","))), ","),
	)

	if 0 < len(update_columns) {
		s := make([]string, 0, len(update_columns))
		for i := range update_columns {
			update_column := update_columns[i]
			s = append(s, fmt.Sprintf("%v = VALUES(%v)", update_column, update_column))
		}
		fmt.Fprintf(&buf, "\nON DUPLICATE KEY UPDATE %v", strings.Join(s, ", "))
	}

	return ExecContext(ctx, tx, buf.String(), values)
}

func UpdateContext(ctx context.Context, tx Preparer, table_name string, columns []string, values []interface{}, cond vanilla.ConditionResult) (affected int64, lastid int64, err error) {
	set_prepare_placeholder := func(column_names []string) []string {
		ss := make([]string, len(column_names))
		for i := range column_names {
			ss[i] = column_names[i] + "=" + __SQL_PREPARED_STMT_PLACEHOLDER__ /*placeholder=?*/
		}
		return ss
	}

	args := make([]interface{}, 0, __DEFAULT_ARGS_CAPACITY__)
	args = append(args, values...)

	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, `UPDATE %v SET %v`,
		table_name,
		strings.Join(set_prepare_placeholder(columns), ","),
	)
	if cond != nil {
		fmt.Fprintf(&buf, "\nWHERE %v", cond.Query())

		args = append(args, cond.Args()...)
	}

	return ExecContext(ctx, tx, buf.String(), args)
}

func DeleteContext(ctx context.Context, tx Preparer, table_name string, cond vanilla.ConditionResult) (affected int64, err error) {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, `DELETE FROM %v`, table_name)

	var args []interface{}
	if cond != nil {
		fmt.Fprintf(&buf, "\nWHERE %v", cond.Query())

		args = cond.Args()
	}

	affected, _, err = ExecContext(ctx, tx, buf.String(), args)

	return
}

type Scanner = interface {
	Scan(dest ...interface{}) error
}
type CallbackScanner = func(scan Scanner) error
type CallbackScannerWithIndex = func(scan Scanner, _ int) error

type SqlScanner interface {
	QueryRow(tx Preparer, dialect string) func(CallbackScanner) error
	QueryRowContext(ctx context.Context, tx Preparer, dialect string) func(CallbackScanner) error
	QueryRows(tx Preparer, dialect string) func(CallbackScannerWithIndex) error
	QueryRowsContext(ctx context.Context, tx Preparer, dialect string) func(CallbackScannerWithIndex) error
}

func QueryRowContext(ctx context.Context, tx Preparer, query string, args []interface{}) func(CallbackScanner) error {
	if __VANILLA_DEBUG_PRINT_STATMENT__() {
		println(fmt.Sprintf("query=%v", query))
		println(fmt.Sprintf("args=%+v", args))
	}

	return func(scan CallbackScanner) (err error) {
		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			err = errors.Wrapf(err, "sql.Tx.Prepare")
			return
		}
		defer func() {
			err = vanilla.ErrorCompose(err, errors.Wrapf(stmt.Close(), "sql.Stmt.Close"))
		}()

		row := stmt.QueryRowContext(ctx, args...)
		if row.Err() != nil {
			err = errors.Wrapf(row.Err(), "sql.Row.Err")
			return
		}

		err = scan(row)

		if err != nil {
			err = errors.Wrapf(err, "sql.Row.Scan")
			err = errors.Wrapf(err, "faild to query row\nquery=\"%v\"\nargs=\"%+v\"",
				query,
				args,
			)

			var once sync.Once
			vanilla.CauseIter(err, func(er error) {
				if er == sql.ErrNoRows {
					once.Do(func() {
						err = errors.Wrapf(err, database.ErrorRecordWasNotFound.Error())
					})
				}
			})

			return
		}

		return
	}
}

func QueryRowsContext(ctx context.Context, tx Preparer, query string, args []interface{}) func(CallbackScannerWithIndex) error {
	if __VANILLA_DEBUG_PRINT_STATMENT__() {
		println(fmt.Sprintf("query=%v", query))
		println(fmt.Sprintf("args=%+v", args))
	}

	return func(scan CallbackScannerWithIndex) (err error) {
		stmt, err := tx.PrepareContext(ctx, query)
		err = errors.Wrapf(err, "sql.Tx.Prepare")
		if err != nil {
			return
		}
		defer func() {
			err = vanilla.ErrorCompose(err, errors.Wrapf(stmt.Close(), "sql.Stmt.Close"))
		}()

		var rows *sql.Rows
		rows, err = stmt.QueryContext(ctx, args...)
		if err != nil {
			err = errors.Wrapf(err, "sql.Stmt.Query")
			return
		}
		defer func() {
			err = vanilla.ErrorCompose(err, errors.Wrapf(rows.Close(), "sql.Rows.Close"))
		}()

		if rows.Err() != nil {
			err = errors.Wrapf(rows.Err(), "sql.Rows.Err")
			return
		}

		i := 0
		for rows.Next() {
			err = scan(rows, i)

			if err != nil {
				err = errors.Wrapf(err, "sql.Row.Scan")
				err = errors.Wrapf(err, "faild to query rows\nquery=\"%v\"\nargs=\"%+v\"",
					query,
					args,
				)

				var once sync.Once
				vanilla.CauseIter(err, func(er error) {
					if er == sql.ErrNoRows {
						once.Do(func() {
							err = errors.Wrapf(err, database.ErrorRecordWasNotFound.Error())
						})
					}
				})

				break
			}
			i++
		}

		return
	}
}

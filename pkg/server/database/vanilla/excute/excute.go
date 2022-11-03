package excute

import (
	"context"
	"database/sql"
	"sync"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/pkg/errors"
)

var ErrorCompose = stmt.ErrorCompose
var CauseIter = stmt.CauseIter

func ExecContext(ctx context.Context, tx Preparer, query string, args []interface{}) (affected int64, lastid int64, err error) {
	VANILLA_DEBUG_PRINT(query, args)

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		err = errors.Wrapf(err, "sql.Tx.Prepare")
		return
	}
	defer func() {
		err = ErrorCompose(err, errors.Wrapf(stmt.Close(), "sql.Stmt.Close"))
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

func QueryRowContext(ctx context.Context, tx Preparer, query string, args []interface{}) func(CallbackScanner) error {
	VANILLA_DEBUG_PRINT(query, args)

	return func(scan CallbackScanner) (err error) {
		stmt, err := tx.PrepareContext(ctx, query)
		if err != nil {
			err = errors.Wrapf(err, "sql.Tx.Prepare")
			return
		}
		defer func() {
			err = ErrorCompose(err, errors.Wrapf(stmt.Close(), "sql.Stmt.Close"))
		}()

		row := stmt.QueryRowContext(ctx, args...)
		if row.Err() != nil {
			err = errors.Wrapf(row.Err(), "sql.Row.Err")
			return
		}

		err = scan(row)

		if err != nil {
			err = errors.Wrapf(err, "sql.Row.Scan")
			err = errors.Wrapf(err, "failed to query row\nquery=\"%v\"\nargs=\"%+v\"",
				query,
				args,
			)

			var once sync.Once
			CauseIter(err, func(er error) {
				if er == sql.ErrNoRows {
					once.Do(func() {
						err = errors.Wrapf(err, vanilla.ErrorRecordWasNotFound.Error())
					})
				}
			})

			return
		}

		return
	}
}

func QueryRowsContext(ctx context.Context, tx Preparer, query string, args []interface{}) func(CallbackScannerWithIndex) error {
	VANILLA_DEBUG_PRINT(query, args)

	return func(scan CallbackScannerWithIndex) (err error) {
		stmt, err := tx.PrepareContext(ctx, query)
		err = errors.Wrapf(err, "sql.Tx.Prepare")
		if err != nil {
			return
		}
		defer func() {
			err = ErrorCompose(err, errors.Wrapf(stmt.Close(), "sql.Stmt.Close"))
		}()

		var rows *sql.Rows
		rows, err = stmt.QueryContext(ctx, args...)
		if err != nil {
			err = errors.Wrapf(err, "sql.Stmt.Query")
			return
		}
		defer func() {
			err = ErrorCompose(err, errors.Wrapf(rows.Close(), "sql.Rows.Close"))
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
				err = errors.Wrapf(err, "failed to query rows\nquery=\"%v\"\nargs=\"%+v\"",
					query,
					args,
				)

				var once sync.Once
				CauseIter(err, func(er error) {
					if er == sql.ErrNoRows {
						once.Do(func() {
							err = errors.Wrapf(err, vanilla.ErrorRecordWasNotFound.Error())
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

func Repeat(n int, s string) []string {
	ss := make([]string, n)
	for i := 0; i < n; i++ {
		ss[i] = s
	}
	return ss
}

func Count(dialect SqlExcutor, tableName string, cond stmt.ConditionStmt, page stmt.PaginationStmt) func(ctx context.Context, tx Preparer) (count int, err error) {
	var (
		columns = []string{"COUNT(1)"}
	)

	return func(ctx context.Context, tx Preparer) (int, error) {
		var count int
		err := dialect.QueryRow(tableName, columns, cond, nil, page)(
			ctx, tx)(
			func(scan Scanner) error {
				err := scan.Scan(&count)
				err = errors.WithStack(err)
				return err
			})

		return count, err
	}
}

func Exist(dialect SqlExcutor, tableName string, cond stmt.ConditionStmt) func(ctx context.Context, tx Preparer) (exist bool, err error) {
	var (
		columns = []string{"COUNT(1)"}
		page    = stmt.Limit(1)
	)

	return func(ctx context.Context, tx Preparer) (bool, error) {
		var count int
		err := dialect.QueryRow(tableName, columns, cond, nil, page)(
			ctx, tx)(
			func(scan Scanner) error {
				err := scan.Scan(&count)
				err = errors.WithStack(err)
				return err
			})

		return 0 < count, err
	}

}

package vanilla

import (
	"context"
	"database/sql"
	"time"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/block"
	"github.com/pkg/errors"
)

type SqlDbEx struct {
	*sql.DB
	timeout time.Duration
}

func NewSqlDbEx(db *sql.DB, timeout ...time.Duration) *SqlDbEx {
	if len(timeout) == 0 {
		timeout = append(timeout, 3*time.Second)
	}

	return &SqlDbEx{
		DB:      db,
		timeout: timeout[0],
	}
}

// func (ctl *SqlDbEx) DB() *sql.DB {
// 	return ctl.db
// }

// Scope
func (ctl SqlDbEx) Scope(fn func(*sql.Tx) error) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), ctl.timeout) //timeout;
	defer cancel()

	return ScopeTx(ctl.DB, ctx, fn)
}

// ScopeTx
func (ctl SqlDbEx) ScopeTx(ctx context.Context, fn func(*sql.Tx) error) (err error) {
	return ScopeTx(ctl.DB, ctx, fn)
}

func Scope(db *sql.DB, fn func(*sql.Tx) error) (err error) {
	return ScopeTx(db, nil, fn)
}

func ScopeTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) (err error) {
	var cancel = func() {}
	if ctx == nil {
		ctx, cancel = context.WithTimeout(context.Background(), 60*time.Second) //timeout;
	}
	defer cancel()

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	err = errors.Wrapf(err, "sql begin tran")
	if err != nil {
		return err
	}

	block.Block{
		Try: func() {
			err = fn(tx)
			err = errors.Wrapf(err, "tx exec")
		},
		Catch: func(ex error) {
			err = errors.Wrapf(ex, "catch")
		},
		Finally: func() {
			switch err == nil {
			case true:
				err = tx.Commit()
				err = errors.Wrapf(err, "tx commit")
			default:
				err = ErrorCompose(err, errors.Wrapf(tx.Rollback(), "tx rollback"))
			}
		},
	}.Do()

	return
}

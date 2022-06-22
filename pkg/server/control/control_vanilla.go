package control

import (
	"context"
	"database/sql"
	"time"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/NexClipper/sudory/pkg/server/macro/block"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

const (
	__INIT_RECORD_CAPACITY__        = 5
	__INTERNAL_TRANSACTIN_TIMEOUT__ = 60 * time.Second
)

type ControlVanilla struct {
	db *sql.DB
}

func NewVanilla(db *sql.DB) *ControlVanilla {
	return &ControlVanilla{
		db: db,
	}
}

func (ctl *ControlVanilla) DB() *sql.DB {
	return ctl.db
}

// Scope
func (ctl ControlVanilla) Scope(fn func(*sql.Tx) error) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), __INTERNAL_TRANSACTIN_TIMEOUT__) //timeout;
	defer cancel()

	return ctl.ScopeTx(ctx, fn)
}

// ScopeTx
func (ctl ControlVanilla) ScopeTx(ctx context.Context, fn func(*sql.Tx) error) (err error) {
	var cancel = func() {}
	if ctx == nil {
		ctx, cancel = context.WithTimeout(ctx, __INTERNAL_TRANSACTIN_TIMEOUT__) //timeout;
	}
	defer cancel()

	var tx *sql.Tx
	tx, err = ctl.db.BeginTx(ctx, &sql.TxOptions{})
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

// HttpError
func HttpError(err error, code int) error {
	if err == nil {
		return nil
	}
	return echo.NewHTTPError(code).SetInternal(err)
}

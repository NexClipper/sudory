package stmtex

import (
	"context"
	"database/sql"
	"time"

	vanilla "github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	"github.com/NexClipper/sudory/pkg/server/macro/block"
	"github.com/pkg/errors"
)

type BeginTx interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

func ScopeTx(ctx context.Context, db BeginTx, fn func(*sql.Tx) error) (err error) {
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
				err = vanilla.ErrorCompose(err, errors.Wrapf(tx.Rollback(), "tx rollback"))
			}
		},
	}.Do()

	return
}

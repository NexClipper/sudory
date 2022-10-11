package vanilla

// import (
// 	"context"
// 	"database/sql"

// 	. "github.com/NexClipper/sudory/pkg/server/macro"
// 	"github.com/NexClipper/sudory/pkg/server/macro/block"
// 	"github.com/pkg/errors"
// )

// type SqlDbEx struct {
// 	*sql.DB
// }

// // func (ctl *SqlDbEx) DB() *sql.DB {
// // 	return ctl.db
// // }

// // Scope
// func (ctl *SqlDbEx) Scope(fn func(*sql.Tx) error) (err error) {
// 	return ScopeTx(ctl.DB, nil, fn)
// }

// // ScopeTx
// func (ctl *SqlDbEx) ScopeTx(ctx context.Context, fn func(*sql.Tx) error) (err error) {
// 	return ScopeTx(ctl.DB, ctx, fn)
// }

// func Scope(db *sql.DB, fn func(*sql.Tx) error) (err error) {
// 	return ScopeTx(db, nil, fn)
// }

// func ScopeTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) (err error) {
// 	if ctx == nil {
// 		ctx = context.Background()
// 	}

// 	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
// 	err = errors.Wrapf(err, "sql begin tran")
// 	if err != nil {
// 		return err
// 	}

// 	block.Block{
// 		Try: func() {
// 			err = fn(tx)
// 			err = errors.Wrapf(err, "tx exec")
// 		},
// 		Catch: func(ex error) {
// 			err = errors.Wrapf(ex, "catch")
// 		},
// 		Finally: func() {
// 			switch err == nil {
// 			case true:
// 				err = tx.Commit()
// 				err = errors.Wrapf(err, "tx commit")
// 			default:
// 				err = ErrorCompose(err, errors.Wrapf(tx.Rollback(), "tx rollback"))
// 			}
// 		},
// 	}.Do()

// 	return
// }

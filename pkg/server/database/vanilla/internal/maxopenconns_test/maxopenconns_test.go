package maxopenconns_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmtex"
	_ "github.com/go-sql-driver/mysql"
)

func TestMaxOpenConns(t *testing.T) {
	const (
		dialect    = "mysql"
		connstring = "sudory:sudory@tcp(127.0.0.1:3306)/sudory?charset=utf8mb4&parseTime=True&loc=Local"
	)
	var db *sql.DB
	db, err := sql.Open(dialect, connstring)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// engine, err := xorm.NewEngine(dialect, connstring)
	// if err != nil {
	// 	panic(err)
	// }
	// defer engine.Close()

	// db = engine.DB().DB

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(1 * time.Second)

	fn_odd := func(db *sql.DB) {

		var ctx = context.Background()
		// var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		// defer cancel()

		stmtex.ScopeTx(ctx, db, func(tx *sql.Tx) error {

			_, _, err = stmtex.ExecContext(ctx, tx, "SELECT 1", []interface{}{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return err
			}

			_, _, err = stmtex.ExecContext(ctx, tx, "SELECT 2", []interface{}{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				return err
			}

			return nil
		})
	}

	fn_even := func(db *sql.DB) {
		var v int
		var ctx = context.Background()
		// var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
		// defer cancel()

		err = stmtex.QueryRowContext(ctx, db, "SELECT 1", []interface{}{})(func(scan stmtex.Scanner) error {
			return scan.Scan(&v)
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}

		err = stmtex.QueryRowContext(ctx, db, "SELECT 2", []interface{}{})(func(scan stmtex.Scanner) error {
			return scan.Scan(&v)
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return
		}
	}

	var current int32
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {

		wg.Add(1)
		go func(i int, db *sql.DB) {
			defer wg.Done()

			fmt.Println("gorutine", atomic.AddInt32(&current, 1))
			if i%2 == 0 {
				fn_even(db)
			} else {
				fn_odd(db)
			}

			fmt.Println("gorutine", atomic.AddInt32(&current, -1)*-1)
		}(i, db)
	}

	wg.Wait()

	// fmt.Println("query #2")
	// go func() {
	// 	time.Sleep(2 * time.Second)
	// 	panic(nil)
	// }()
	// err = db.QueryRow("SELECT 2").Scan(&v)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("value #2: ", v)
}

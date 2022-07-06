package vanilla

// import (
// 	"database/sql"
// 	"fmt"
// 	"strings"

// 	"github.com/NexClipper/sudory/pkg/server/database/vanilla/error_compose"
// 	"github.com/pkg/errors"
// )

// const __DEFAULT_ARGS_CAPACITY__ = 10
// const __DEBUG_PRINT_STATMENT__ = false
// const __SQL_PREPARED_STMT_PLACEHOLDER__ = "?"

// type Preparer interface {
// 	Prepare(query string) (*sql.Stmt, error)
// }

// // type QueryMuxer interface {
// // 	Query() string
// // 	// Args() []interface{}
// // 	Combine(...QueryMuxer) QueryMuxer
// // 	Prepare(tx Preparer) (args []interface{}, stmt *sql.Stmt, err error)
// // }

// // type QueryMux struct {
// // 	query string
// // 	args  []interface{}
// // }

// // func (q QueryMux) Query() string {
// // 	return q.query
// // }

// // func (q QueryMux) Args() []interface{} {
// // 	return q.args
// // }

// // func (q QueryMux) Combine(querirs ...QueryMuxer) QueryMuxer {
// // 	mq := MultiQueryMux{}
// // 	mq.querirs = append([]QueryMuxer{&q}, querirs...)
// // 	return &mq
// // }

// // func (q QueryMux) Prepare(tx Preparer) (args []interface{}, stmt *sql.Stmt, err error) {
// // 	stmt, err = tx.Prepare(q.Query())
// // 	return q.Args(), stmt, err
// // }

// type MultiQueryMux struct {
// 	querirs []QueryMuxer
// }

// func (mq MultiQueryMux) Query() string {
// 	const sep = ";"
// 	queries := make([]string, 0, len(mq.querirs))

// 	for _, item := range mq.querirs {
// 		//split and append
// 		for _, q := range strings.Split(item.Query(), sep) {
// 			q = strings.TrimSpace(q)
// 			queries = append(queries, q)
// 		}
// 	}
// 	s := strings.Join(queries, sep+"\n") + sep + "\n" // "string_1" + ";\n" + "string_2" + ";\n"

// 	if __DEBUG_PRINT_STATMENT__ {
// 		println(s)
// 	}

// 	return s
// }

// func (mq MultiQueryMux) Args() []interface{} {
// 	args := make([]interface{}, 0, __DEFAULT_ARGS_CAPACITY__*len(mq.querirs)*2)
// 	for _, item := range mq.querirs {
// 		args = append(args, item.Args()...)
// 	}
// 	return args
// }

// func (mq MultiQueryMux) Combine(querirs ...QueryMuxer) QueryMuxer {
// 	mq.querirs = append(mq.querirs, querirs...)
// 	return &mq
// }

// func (mq MultiQueryMux) Prepare(tx Preparer) (args []interface{}, stmt *sql.Stmt, err error) {
// 	stmt, err = tx.Prepare(mq.Query())
// 	return mq.Args(), stmt, err
// }

// type Condition struct {
// 	Condition string
// 	Args      []interface{}
// }

// func NewCond(cond string, args ...interface{}) *Condition {
// 	return &Condition{
// 		Condition: cond,
// 		Args:      args,
// 	}
// }

// func Insert(tablename string, columns []string) QueryMuxer {
// 	const query = `INSERT INTO %v (%v) VALUES (%v)`
// 	s := fmt.Sprintf(query,
// 		tablename,
// 		strings.Join(columns, ","),
// 		strings.Join(Repeat(len(columns), __SQL_PREPARED_STMT_PLACEHOLDER__ /*placeholder=?*/), ","),
// 	)

// 	if __DEBUG_PRINT_STATMENT__ {
// 		println(s)
// 	}

// 	return &QueryMux{
// 		query: s,
// 		args:  []interface{}{},
// 	}
// }

// func Select(tablename string, columns []string, conditions ...Condition) QueryMuxer {
// 	args := make([]interface{}, 0, __DEFAULT_ARGS_CAPACITY__)

// 	query := `SELECT %v FROM %v`
// 	for i := range conditions {
// 		query = query + "\n" + conditions[i].Condition
// 		args = append(args, conditions[i].Args...)
// 	}

// 	s := fmt.Sprintf(query,
// 		strings.Join(columns, ","),
// 		tablename,
// 	)

// 	if __DEBUG_PRINT_STATMENT__ {
// 		println(s)
// 	}

// 	return &QueryMux{
// 		query: s,
// 		args:  args,
// 	}
// }

// func Update(tablename string, column_names []string, column_values []interface{}, conditions ...Condition) QueryMuxer {
// 	set_prepare_placeholder := func() []string {
// 		ss := make([]string, len(column_names))
// 		for i := range column_names {
// 			ss[i] = column_names[i] + "=" + __SQL_PREPARED_STMT_PLACEHOLDER__ /*placeholder=?*/
// 		}
// 		return ss
// 	}

// 	args := make([]interface{}, 0, __DEFAULT_ARGS_CAPACITY__)
// 	args = append(args, column_values...)

// 	query := `UPDATE %v SET %v`
// 	for i := range conditions {
// 		query = query + "\n" + conditions[i].Condition
// 		args = append(args, conditions[i].Args...)
// 	}

// 	s := fmt.Sprintf(query,
// 		tablename,
// 		strings.Join(set_prepare_placeholder(), ","),
// 	)

// 	if __DEBUG_PRINT_STATMENT__ {
// 		println(s)
// 	}

// 	return &QueryMux{
// 		query: s,
// 		args:  args,
// 	}
// }

// func Delete(tablename string, conditions ...Condition) QueryMuxer {
// 	args := make([]interface{}, 0, __DEFAULT_ARGS_CAPACITY__)

// 	query := `DELETE FROM %v`
// 	for i := range conditions {
// 		query = query + "\n" + conditions[i].Condition
// 		args = append(args, conditions[i].Args...)
// 	}
// 	s := fmt.Sprintf(query,
// 		tablename,
// 	)

// 	if __DEBUG_PRINT_STATMENT__ {
// 		println(s)
// 	}

// 	return &QueryMux{
// 		query: s,
// 		args:  args,
// 	}
// }

// func Repeat(n int, s string) []string {
// 	ss := make([]string, n)
// 	for i := 0; i < n; i++ {
// 		ss[i] = s
// 	}
// 	return ss
// }

// type Scanner interface {
// 	Scan(dest ...interface{}) error
// }
// type CallbackScanner = func(Scanner) error

// func QueryRow(tx Preparer, table_name string, column_names []string, conditions ...Condition) func(CallbackScanner) error {
// 	return func(scan CallbackScanner) (err error) {
// 		args, stmt, err := Select(table_name, column_names, conditions...).
// 			Prepare(tx)
// 		err = errors.Wrapf(err, "sql.DB.Prepare")
// 		if err != nil {
// 			return err
// 		}

// 		defer func() {
// 			err = error_compose.Composef(err, stmt.Close(), "sql.Stmt.Close")
// 		}()

// 		row := stmt.QueryRow(args...)
// 		err = scan(row)
// 		err = errors.Wrapf(err, "sql.Row.Scan")
// 		err = error_compose.Composef(err, row.Err(), "sql.Row; during scan")

// 		err = errors.Wrapf(err, "faild to query row table=\"%v\"",
// 			table_name,
// 		)
// 		return
// 	}
// }

// func QueryRows(tx Preparer, table_name string, column_names []string, conditions ...Condition) func(CallbackScanner) error {
// 	return func(scan CallbackScanner) (err error) {
// 		args, stmt, err := Select(table_name, column_names, conditions...).
// 			Prepare(tx)
// 		err = errors.Wrapf(err, "sql.DB.Prepare")
// 		if err != nil {
// 			return
// 		}

// 		defer func() {
// 			err = error_compose.Composef(err, stmt.Close(), "sql.Stmt.Close")
// 		}()

// 		var rows *sql.Rows
// 		rows, err = stmt.Query(args...)
// 		err = errors.Wrapf(err, "sql.Stmt.Query")
// 		if err != nil {
// 			return
// 		}

// 		defer func() {
// 			err = error_compose.Composef(err, rows.Close(), "sql.Rows.Close")
// 		}()
// 		for rows.Next() {
// 			err = scan(rows)
// 			err = errors.Wrapf(err, "sql.Row.Scan")
// 			if err != nil {
// 				break
// 			}
// 		}
// 		err = error_compose.Composef(err, rows.Err(), "sql.Rows; during iteration")

// 		err = errors.Wrapf(err, "faild to query rows table=\"%v\"",
// 			table_name,
// 		)
// 		return
// 	}
// }

// type Executor interface {
// 	Exec(args ...interface{}) (sql.Result, error)
// }
// type CallbackExecutor = func(Executor) (sql.Result, error)

// func InsertRow(tx Preparer, table_name string, column_names []string) func(CallbackExecutor) error {
// 	return func(exec CallbackExecutor) (err error) {
// 		_, stmt, err := Insert(table_name, column_names).
// 			Prepare(tx)
// 		err = errors.Wrapf(err, "sql.Tx.Prepare")
// 		if err != nil {
// 			return
// 		}

// 		defer func() {
// 			err = error_compose.Composef(err, stmt.Close(), "sql.Stmt.Close")
// 		}()

// 		var result sql.Result
// 		result, err = exec(stmt)
// 		err = errors.Wrapf(err, "sql.Stmt.Exec")
// 		if err != nil {
// 			return
// 		}

// 		affected, err := result.RowsAffected()
// 		err = errors.Wrapf(err, "sql.Result.RowsAffected")
// 		if affected == 0 {
// 			err = error_compose.Compose(err, errors.New("no affected"))
// 		}

// 		err = errors.Wrapf(err, "faild to insert rows table=\"%v\"",
// 			table_name,
// 		)
// 		return
// 	}
// }

// func InsertRows(tx Preparer, table_name string, column_names []string) func(func(int) ([]interface{}, bool)) error {
// 	return func(argsIter func(int) ([]interface{}, bool)) (err error) {
// 		_, stmt, err := Insert(table_name, column_names).
// 			Prepare(tx)
// 		err = errors.Wrapf(err, "sql.Tx.Prepare")
// 		if err != nil {
// 			return
// 		}

// 		defer func() {
// 			err = error_compose.Composef(err, stmt.Close(), "sql.Stmt.Close")
// 		}()

// 		i := 0
// 		for args, ok := argsIter(i); ok; args, ok = argsIter(i) {
// 			var result sql.Result
// 			result, err = stmt.Exec(args...)
// 			err = errors.Wrapf(err, "sql.Stmt.Exec")
// 			if err != nil {
// 				return
// 			}
// 			var n int64
// 			n, err = result.RowsAffected()
// 			err = errors.Wrapf(err, "sql.Result.RowsAffected")
// 			if n == 0 {
// 				err = error_compose.Compose(err, errors.New("no affected"))
// 				break
// 			}

// 			i++
// 		}

// 		err = errors.Wrapf(err, "faild to insert rows table=\"%v\"",
// 			table_name,
// 		)
// 		return
// 	}
// }

// func UpdateRow(tx Preparer, table_name string, keys_values map[string]interface{}, conditions ...Condition) (affected int64, err error) {
// 	keys := make([]string, 0, len(keys_values))
// 	values := make([]interface{}, 0, len(keys_values))
// 	for k, v := range keys_values {
// 		keys = append(keys, k)
// 		values = append(values, v)
// 	}

// 	args, stmt, err := Update(table_name, keys, values, conditions...).
// 		Prepare(tx)
// 	err = errors.Wrapf(err, "sql.Tx.Prepare")
// 	if err != nil {
// 		return affected, err
// 	}

// 	defer func() {
// 		err = error_compose.Composef(err, stmt.Close(), "sql.Stmt.Close")
// 	}()

// 	result, err := stmt.Exec(args...)
// 	err = errors.Wrapf(err, "sql.Stmt.Exec")
// 	if err != nil {
// 		return affected, err
// 	}

// 	affected, err = result.RowsAffected()
// 	err = errors.Wrapf(err, "sql.Result.RowsAffected")
// 	// if affected == 0 {
// 	// 	err = ErrorCompose(err, errors.New("no affected"))
// 	// }

// 	err = errors.Wrapf(err, "faild to insert rows table=\"%v\"",
// 		table_name,
// 	)
// 	return affected, err
// }

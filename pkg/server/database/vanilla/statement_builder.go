package vanilla

// import (
// 	"bytes"
// 	"context"
// 	"fmt"
// 	"strings"

// 	"github.com/NexClipper/sudory/pkg/server/database/vanilla/prepare"
// 	"github.com/pkg/errors"
// )

// const __DEFAULT_ARGS_CAPACITY__ = 10
// const __SQL_PREPARED_STMT_PLACEHOLDER__ = "?"

// var Stmt = stmt{}

// type StmtBuild struct {
// 	query string
// 	args  []interface{}
// }

// func (sb StmtBuild) Print() string {
// 	buf := bytes.Buffer{}
// 	fmt.Fprintf(&buf, "query=\"%v\"\n", sb.query)
// 	fmt.Fprintf(&buf, "args=%+v\n", sb.args...)

// 	return buf.String()
// }

// func (sb StmtBuild) Query() string {
// 	return sb.query
// }

// func (sb StmtBuild) Args() []interface{} {
// 	return sb.args
// }

// func (sb StmtBuild) QueryRow(tx Preparer) func(CallbackScanner) error {
// 	return QueryRow(tx, sb.Query(), sb.Args())
// }

// func (sb StmtBuild) QueryRowContext(ctx context.Context, tx Preparer) func(CallbackScanner) error {
// 	return QueryRowContext(ctx, tx, sb.Query(), sb.Args())
// }

// func (sb StmtBuild) QueryRows(tx Preparer) func(CallbackScannerWithIndex) error {
// 	return QueryRows(tx, sb.Query(), sb.Args())
// }

// func (sb StmtBuild) QueryRowsContext(ctx context.Context, tx Preparer) func(CallbackScannerWithIndex) error {
// 	return QueryRowsContext(ctx, tx, sb.Query(), sb.Args())
// }

// func (sb StmtBuild) Exec(tx Preparer) (affected int64, err error) {
// 	return Exec(tx, sb.Query(), sb.Args())
// }

// func (sb StmtBuild) ExecContext(ctx context.Context, tx Preparer) (affected int64, err error) {
// 	return ExecContext(ctx, tx, sb.Query(), sb.Args())
// }

// type stmt struct{}

// func (stmt) Insert(table_name string, column_names []string, column_values ...[]interface{}) (stmt *StmtBuild, err error) {
// 	if len(column_names) == 0 {
// 		err = errors.New("column_names length is not 0")
// 		return
// 	}
// 	if len(column_values) == 0 {
// 		err = errors.New("column_values length is not 0")
// 		return
// 	}

// 	// flat column_values
// 	var args []interface{}
// 	for _, column_value := range column_values {
// 		if len(column_value)%len(column_names) != 0 {
// 			return stmt, errors.New("column_names and column_values do not match in length")
// 		}

// 		if args == nil {
// 			args = make([]interface{}, 0, len(column_values)*len(column_value))
// 		}
// 		args = append(args, column_value...)
// 	}

// 	buf := bytes.Buffer{}
// 	fmt.Fprintf(&buf, `INSERT INTO %v (%v) VALUES %v`,
// 		table_name,
// 		strings.Join(column_names, ","),
// 		strings.Join(
// 			Repeat(len(args)/len(column_names),
// 				fmt.Sprintf("(%v)", strings.Join(
// 					Repeat(len(column_names),
// 						__SQL_PREPARED_STMT_PLACEHOLDER__ /*placeholder=?*/), ","))), ","),
// 	)

// 	stmt = &StmtBuild{
// 		query: buf.String(),
// 		args:  args,
// 	}

// 	return
// }

// func (stmt) InsertOrUpdate(table_name string, column_names []string, update_column_names []string, column_values ...[]interface{}) (stmt *StmtBuild, err error) {
// 	if len(column_names) == 0 {
// 		err = errors.New("column_names length is not 0")
// 		return
// 	}
// 	if len(column_values) == 0 {
// 		err = errors.New("column_values length is not 0")
// 		return
// 	}

// 	// flat column_values
// 	var args []interface{}
// 	for _, column_value := range column_values {
// 		if len(column_value)%len(column_names) != 0 {
// 			return stmt, errors.New("column_names and column_values do not match in length")
// 		}

// 		if args == nil {
// 			args = make([]interface{}, 0, len(column_values)*len(column_value))
// 		}
// 		args = append(args, column_value...)
// 	}

// 	buf := bytes.Buffer{}
// 	fmt.Fprintf(&buf, `INSERT INTO %v (%v) VALUES %v`,
// 		table_name,
// 		strings.Join(column_names, ","),
// 		strings.Join(
// 			Repeat(len(args)/len(column_names),
// 				fmt.Sprintf("(%v)", strings.Join(
// 					Repeat(len(column_names),
// 						__SQL_PREPARED_STMT_PLACEHOLDER__ /*placeholder=?*/), ","))), ","),
// 	)

// 	if 0 < len(update_column_names) {
// 		s := make([]string, 0, len(update_column_names))
// 		for i := range update_column_names {
// 			column := update_column_names[i]
// 			s = append(s, fmt.Sprintf("%v = VALUES(%v)", column, column))
// 		}
// 		fmt.Fprintf(&buf, "\nON DUPLICATE KEY UPDATE %v", strings.Join(s, ", "))
// 	}

// 	stmt = &StmtBuild{
// 		query: buf.String(),
// 		args:  args,
// 	}

// 	return
// }

// func (stmt) Select(table_name string, column_names []string, q *prepare.Condition, o *prepare.Orders, p *prepare.Pagination) (stmt *StmtBuild) {
// 	buf := bytes.Buffer{}
// 	fmt.Fprintf(&buf, `SELECT %v FROM %v`,
// 		strings.Join(column_names, ","),
// 		table_name,
// 	)

// 	args := make([]interface{}, 0, __DEFAULT_ARGS_CAPACITY__)
// 	if q != nil {
// 		fmt.Fprintf(&buf, "\nWHERE %v", q.Query())

// 		args = append(args, q.Args()...)
// 	}
// 	if o != nil {
// 		fmt.Fprintf(&buf, "\nORDER BY %v", o.Order())
// 	}
// 	if p != nil {
// 		fmt.Fprintf(&buf, "\nLIMIT %v, %v", p.Offset(), p.Limit())
// 	}

// 	stmt = &StmtBuild{
// 		query: buf.String(),
// 		args:  args,
// 	}

// 	return
// }

// func (stmt) Count(table_name string, q *prepare.Condition, p *prepare.Pagination) func(ctx context.Context, tx Preparer) (int, error) {
// 	buf := bytes.Buffer{}
// 	fmt.Fprintf(&buf, `SELECT COUNT(1) FROM %v`,
// 		table_name,
// 	)

// 	args := make([]interface{}, 0, __DEFAULT_ARGS_CAPACITY__)
// 	if q != nil {
// 		fmt.Fprintf(&buf, "\nWHERE %v", q.Query())

// 		args = append(args, q.Args()...)
// 	}

// 	if p != nil {
// 		fmt.Fprintf(&buf, "\nLIMIT %v, %v", p.Offset(), p.Limit())
// 	}

// 	stmt := &StmtBuild{
// 		query: buf.String(),
// 		args:  args,
// 	}

// 	return func(ctx context.Context, tx Preparer) (count int, err error) {
// 		err = stmt.QueryRowContext(ctx, tx)(func(scan Scanner) (err error) {
// 			err = scan.Scan(&count)
// 			return
// 		})
// 		return
// 	}
// }

// func (stmt) Exist(table_name string, q *prepare.Condition) func(ctx context.Context, tx Preparer) (bool, error) {
// 	buf := bytes.Buffer{}
// 	fmt.Fprintf(&buf, `SELECT COUNT(1) FROM %v`,
// 		table_name,
// 	)

// 	args := make([]interface{}, 0, __DEFAULT_ARGS_CAPACITY__)
// 	if q != nil {
// 		fmt.Fprintf(&buf, "\nWHERE %v", q.Query())

// 		args = append(args, q.Args()...)
// 	}

// 	fmt.Fprintf(&buf, "\nLIMIT %v", 1)

// 	stmt := &StmtBuild{
// 		query: buf.String(),
// 		args:  args,
// 	}

// 	return func(ctx context.Context, tx Preparer) (exist bool, err error) {
// 		var count int
// 		err = stmt.QueryRowContext(ctx, tx)(func(scan Scanner) (err error) {
// 			err = scan.Scan(&count)
// 			return
// 		})

// 		exist = 0 < count

// 		return
// 	}
// }

// func (stmt) Update(table_name string, keys_values map[string]interface{}, q *prepare.Condition) (stmt *StmtBuild) {
// 	set_prepare_placeholder := func(column_names []string) []string {
// 		ss := make([]string, len(column_names))
// 		for i := range column_names {
// 			ss[i] = column_names[i] + "=" + __SQL_PREPARED_STMT_PLACEHOLDER__ /*placeholder=?*/
// 		}
// 		return ss
// 	}

// 	keys := make([]string, 0, len(keys_values))
// 	values := make([]interface{}, 0, len(keys_values))
// 	for k, v := range keys_values {
// 		keys = append(keys, k)
// 		values = append(values, v)
// 	}

// 	args := make([]interface{}, 0, __DEFAULT_ARGS_CAPACITY__)
// 	args = append(args, values...)

// 	// for i := range conditions {
// 	// 	query = query + "\n" + conditions[i].Condition
// 	// 	args = append(args, conditions[i].Args...)
// 	// }

// 	buf := bytes.Buffer{}
// 	fmt.Fprintf(&buf, `UPDATE %v SET %v`,
// 		table_name,
// 		strings.Join(set_prepare_placeholder(keys), ","),
// 	)
// 	if q != nil {
// 		fmt.Fprintf(&buf, "\nWHERE %v", q.Query())

// 		args = append(args, q.Args()...)
// 	}

// 	stmt = &StmtBuild{
// 		query: buf.String(),
// 		args:  args,
// 	}

// 	return
// }

// func (stmt) Delete(table_name string, q *prepare.Condition) (stmt *StmtBuild) {
// 	buf := bytes.Buffer{}

// 	fmt.Fprintf(&buf, `DELETE FROM %v`,
// 		table_name)

// 	var args []interface{}
// 	if q != nil {
// 		fmt.Fprintf(&buf, "\nWHERE %v", q.Query())

// 		args = q.Args()
// 	}

// 	stmt = &StmtBuild{
// 		query: buf.String(),
// 		args:  args,
// 	}

// 	return
// }

// func Repeat(n int, s string) []string {
// 	ss := make([]string, n)
// 	for i := 0; i < n; i++ {
// 		ss[i] = s
// 	}
// 	return ss
// }

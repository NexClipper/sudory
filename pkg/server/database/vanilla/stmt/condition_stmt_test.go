package stmt_test

import (
	"reflect"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
	_ "github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/resolvers/mysql"
)

func TestConditionStmt(t *testing.T) {

	var scenarios = []struct {
		dialect     string             // dialect
		name        string             // name
		args        stmt.ConditionStmt // input args
		want_query  interface{}        // expected output
		want_args   interface{}        // expected output
		expectedErr error              // expected error
	}{
		{
			dialect:    "mysql",
			name:       `mysql: Equal`,
			args:       stmt.Equal("foo", "bar"),
			want_query: "`foo` = ?",
			want_args:  []interface{}{"bar"},
		},
		{
			dialect:    "mysql",
			name:       `mysql: IsNull`,
			args:       stmt.IsNull("foo"),
			want_query: "`foo` IS NULL",
			want_args:  []interface{}{},
		},
		{
			dialect:    "mysql",
			name:       `mysql: In`,
			args:       stmt.In("foo", 1, 2, 3, 4),
			want_query: "`foo` IN (?, ?, ?, ?)",
			want_args:  []interface{}{1, 2, 3, 4},
		},
		{
			dialect:    "mysql",
			name:       `mysql: And+Equal+Equal `,
			args:       stmt.And(stmt.Equal("c1", 1), stmt.Equal("c2", 2)),
			want_query: "(`c1` = ? AND `c2` = ?)",
			want_args:  []interface{}{1, 2},
		},
	}

	for _, test_case := range scenarios {
		t.Run(test_case.name, func(t *testing.T) {

			// build sql query
			o, err := stmt.GetConditionStmtResolver(test_case.dialect).Build(test_case.args)
			if err != nil {
				t.Fatal(err)
			}

			// test
			var expected interface{} = test_case.want_query
			var actual interface{} = o.Query()
			// test; query string
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected=%v actual=%v", expected, actual)
			}
			// test; query args
			expected = test_case.want_args
			actual = o.Args()
			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected=%v actual=%v", expected, actual)
			}
		})
	}

}

func TestConditionStmt_and(t *testing.T) {
	var q = stmt.ConditionStmt{}

	q = stmt.And(q, stmt.Equal(
		"foo", "bar",
	))

	o, err := stmt.GetConditionStmtResolver("mysql").Build(q)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(o.Query())
	t.Log(o.Args()...)
}

// func Test_GreaterThan(t *testing.T) {
// 	c := vanilla.GreaterThan("foo", 1).Parse()
// 	t.Log(c.Query(), c.Args())
// }

// func Test_GreaterThanEqual(t *testing.T) {
// 	c := vanilla.GreaterThanEqual("foo", 1).Parse()
// 	t.Log(c.Query(), c.Args())
// }

// func Test_LessThan(t *testing.T) {
// 	c := vanilla.LessThan("foo", 1).Parse()
// 	t.Log(c.Query(), c.Args())
// }

// func Test_LessThanEqual(t *testing.T) {
// 	c := vanilla.LessThanEqual("foo", 1).Parse()
// 	t.Log(c.Query(), c.Args())
// }

// func Test_Like(t *testing.T) {
// 	c := vanilla.Like("foo", "bar%").Parse()
// 	t.Log(c.Query(), c.Args())
// }

// func Test_Or(t *testing.T) {
// 	c := vanilla.Or(vanilla.Equal("c1", 1), vanilla.Equal("c2", 2)).Parse()
// 	t.Log(c.Query(), c.Args())
// }

// func Test_AndOr(t *testing.T) {
// 	c := vanilla.And(vanilla.Equal("c1", 1), vanilla.Equal("c2", 2),
// 		vanilla.Or(vanilla.Equal("c3", 3), vanilla.Equal("c4", 4))).Parse()
// 	t.Log(c.Query(), c.Args())
// }

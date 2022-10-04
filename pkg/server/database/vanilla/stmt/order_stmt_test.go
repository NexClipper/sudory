package stmt_test

import (
	"reflect"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
)

func TestOrderStmt(t *testing.T) {
	var scenarios = []struct {
		dialect     string         // dialect
		name        string         // name
		args        stmt.OrderStmt // input args
		want        interface{}    // expected output
		expectedErr error          // expected error
	}{
		{
			dialect: "mysql",
			name:    `mysql: ASC`,
			args:    stmt.Asc("a1", "b1", "c1"),
			want:    "`a1`, `b1`, `c1` ASC",
		},
		{
			dialect: "mysql",
			name:    `mysql: DESC`,
			args:    stmt.Desc("a2", "b2", "c2"),
			want:    "`a2`, `b2`, `c2` DESC",
		},
		{
			dialect: "mysql",
			name:    `mysql: ASC+DESC+ASC`,
			args:    stmt.Asc("a1", "b1", "c1").Desc("a2", "b2", "c2").Asc("a3", "b3", "c3"),
			want:    "`a1`, `b1`, `c1` ASC, `a2`, `b2`, `c2` DESC, `a3`, `b3`, `c3` ASC",
		},
	}

	for _, test_case := range scenarios {
		t.Run(test_case.name, func(t *testing.T) {
			o, err := stmt.GetOrderStmtResolver(test_case.dialect).Build(test_case.args)
			if err != nil {
				panic(err)
			}

			var expected interface{} = test_case.want
			var actual interface{} = o.Order()

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected=%v actual=%v", expected, actual)
			}
		})
	}
}

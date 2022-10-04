package stmt_test

import (
	"reflect"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
)

func TestPaginationStmt(t *testing.T) {
	var scenarios = []struct {
		dialect     string              // dialect
		name        string              // name
		args        stmt.PaginationStmt // input args
		want        interface{}         // expected output
		expectedErr error               // expected error
	}{
		{
			dialect: "mysql",
			name:    `mysql: limit: 24`,
			args:    stmt.PaginationStmt{"limit": 24},
			want:    "0, 24",
		},
		{
			dialect: "mysql",
			name:    `mysql: page: 42`,
			args:    stmt.PaginationStmt{"page": 42},
			want:    "5207, 127",
		},
		{
			dialect: "mysql",
			name:    `mysql: Limit(24, 42)`,
			args:    stmt.Limit(24, 42),
			want:    "984, 24",
		},
		{
			dialect: "mysql",
			name:    `mysql: limit: 24 page: 42`,
			args:    stmt.PaginationStmt{"limit": 24, "page": 42},
			want:    "984, 24",
		},
	}

	for _, test_case := range scenarios {
		t.Run(test_case.name, func(t *testing.T) {
			o, err := stmt.GetPaginationStmtResolver(test_case.dialect).Build(test_case.args)
			if err != nil {
				panic(err)
			}

			var expected interface{} = test_case.want
			var actual interface{} = o.String()

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected=%v actual=%v", expected, actual)
			}
		})
	}
}

package stmt_test

import (
	"reflect"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
)

func TestPaginationLexer(t *testing.T) {
	var scenarios = []struct {
		name        string      // name
		args        string      // input args
		want        interface{} // expected output
		expectedErr error       // expected error
	}{
		{
			name: `limit: 24 page: 42`,
			args: `{"limit": 24, "page": 42}`,
			want: stmt.TypeMap{"limit": 24, "page": 42},
		},
		{
			name: `limit: "24" page: "42"`,
			args: `{"limit": "24", "page": "42"}`,
			want: stmt.TypeMap{"limit": 24, "page": 42},
		},
		{
			name: `limit: 24`,
			args: `{"limit": 24}`,
			want: stmt.TypeMap{"limit": 24},
		},
		{
			name: `page: 42`,
			args: `{"page": 42}`,
			want: stmt.TypeMap{"page": 42},
		},
		{
			name: `[limit: 24, page: 42]`,
			args: `[{"limit": 24}, {"page": 42}]`,
			want: stmt.TypeMap{"limit": 24, "page": 42},
		},
	}

	for _, test_case := range scenarios {
		t.Run(test_case.name, func(t *testing.T) {
			lex, err := stmt.PaginationLexer.Parse(test_case.args)
			if err != nil {
				t.Error(err)
				return
			}
			var expected interface{} = test_case.want
			var actual interface{} = (map[string]int)(lex)

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected=%v actual=%v", expected, actual)
			}
		})
	}
}

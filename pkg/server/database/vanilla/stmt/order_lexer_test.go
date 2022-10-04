package stmt_test

import (
	"reflect"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
)

func TestOrderLexer(t *testing.T) {
	var scenarios = []struct {
		name        string      // name
		args        string      // input args
		want        interface{} // expected output
		expectedErr error       // expected error
	}{
		{
			name: `{"asc": "foo" }`,
			args: `{"asc": "foo" }`,
			want: []stmt.TypeOrderStmt{{"asc": []string{"foo"}}},
		},
		{
			name: `{"asc": ["foo", "bar"] }`,
			args: `{"asc": ["foo", "bar"] }`,
			want: []stmt.TypeOrderStmt{{"asc": []string{"foo", "bar"}}},
		},
		{
			name: `{"desc": ["foo", "bar"] }`,
			args: `{"desc": ["foo", "bar"] }`,
			want: []stmt.TypeOrderStmt{{"desc": []string{"foo", "bar"}}},
		},
		{
			name: `[{"asc": ["foo", "bar"] }, {"desc": ["foobar", "baz"] }]`,
			args: `[{"asc": ["foo", "bar"] }, {"desc": ["foobar", "baz"] }]`,
			want: []stmt.TypeOrderStmt{
				{"asc": []string{"foo", "bar"}},
				{"desc": []string{"foobar", "baz"}}},
		},
	}
	for _, test_case := range scenarios {
		t.Run(test_case.name, func(t *testing.T) {
			exp := test_case.args
			lex, err := stmt.OrderLexer.Parse(exp)

			if err != nil {
				t.Fatal(err)
			}

			var expected interface{} = test_case.want
			var actual interface{} = ([]stmt.TypeOrderStmt)(lex)

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected=%v actual=%v", expected, actual)
			}
		})
	}
}

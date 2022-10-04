package stmt_test

import (
	"reflect"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
)

func TestConditionLexer(t *testing.T) {

	var scenarios = []struct {
		name        string      // name
		args        string      // input args
		want        interface{} // expected output
		expectedErr error       // expected error
	}{
		{
			name: `s-exp`,
			args: `(and (eq foo "abc") (eq bar "def"))`,
			want: stmt.TypeMap{"and": []interface{}{
				stmt.TypeMap{"equal": stmt.TypeMap{"foo": "abc"}},
				stmt.TypeMap{"equal": stmt.TypeMap{"bar": "def"}},
			}},
		},
		{
			name: `json`,
			args: `{"and": [{"eq": {"foo": "abc"}}, {"eq": {"bar": "def"}}]}`,
			want: stmt.TypeMap{"and": []interface{}{
				stmt.TypeMap{"eq": stmt.TypeMap{"foo": "abc"}},
				stmt.TypeMap{"eq": stmt.TypeMap{"bar": "def"}},
			}},
		},
	}
	for _, test_case := range scenarios {
		t.Run(test_case.name, func(t *testing.T) {
			exp := test_case.args
			lex, err := stmt.ConditionLexer.Parse(exp)

			if err != nil {
				t.Fatal(err)
			}

			var expected interface{} = test_case.want
			var actual interface{} = (map[string]interface{})(lex)

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected=%v actual=%v", expected, actual)
			}
		})
	}
}

package flavor

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
)

func TestNewCondition(t *testing.T) {

	var scenarios = []struct {
		name    string
		args    map[string]interface{}
		want    *Condition
		wantErr error
	}{
		{name: "take 1",
			args: stmt.Map("EQ", stmt.Map("foo", 123)),
			want: &Condition{query: "`foo` = ?", args: stmt.Slice(123)}, wantErr: nil},
		{name: "take 2",
			args: stmt.Map("and", stmt.Slice(stmt.Map("EQ", stmt.Map("foo", 123)))),
			want: &Condition{query: "(`foo` = ?)", args: stmt.Slice(123)}, wantErr: nil},
		{name: "take 3",
			args: stmt.Map("and", stmt.Map("EQ", stmt.Map("foo", 123))),
			want: &Condition{query: "(`foo` = ?)", args: stmt.Slice(123)}, wantErr: nil},
		{name: "take 4",
			args: stmt.Map("and", stmt.Slice(stmt.Map("EQ", stmt.Map("foo", 123)), stmt.Map("EQ", stmt.Map("bar", 456)), stmt.Map("EQ", stmt.Map("foobar", 123456)))),
			want: &Condition{query: "(`foo` = ? AND `bar` = ? AND `foobar` = ?)", args: []interface{}{123, 456, 123456}}, wantErr: nil},
		{name: "take 5",
			args: stmt.Map("or", stmt.Map("EQ", stmt.Map("foo", 123))),
			want: &Condition{query: "(`foo` = ?)", args: []interface{}{123}}, wantErr: nil},
		{name: "take 6",
			args: stmt.Map("or", stmt.Slice(stmt.Map("EQ", stmt.Map("foo", 123)), stmt.Map("EQ", stmt.Map("bar", 456)))),
			want: &Condition{query: "(`foo` = ? OR `bar` = ?)", args: []interface{}{123, 456}}, wantErr: nil},
		{name: "take 7",
			args: stmt.Map("or", stmt.Slice(stmt.Map("and", stmt.Slice(stmt.Map("EQ", stmt.Map("foo", 123)), stmt.Map("EQ", stmt.Map("bar", 456)))), stmt.Map("EQ", stmt.Map("foobar", 123456)))),
			want: &Condition{query: "((`foo` = ? AND `bar` = ?) OR `foobar` = ?)", args: []interface{}{123, 456, 123456}}, wantErr: nil},
		{name: "take 8",
			args: stmt.Map("not", stmt.Map("and", stmt.Slice(stmt.Map("EQ", stmt.Map("foo", 123)), stmt.Map("EQ", stmt.Map("bar", 456))))),
			want: &Condition{query: "NOT (`foo` = ? AND `bar` = ?)", args: []interface{}{123, 456}}, wantErr: nil},
		{name: "take 9",
			args: stmt.Map("gt", stmt.Map("foo", 123)),
			want: &Condition{query: "`foo` > ?", args: []interface{}{123}}, wantErr: nil},
		{name: "take 10",
			args: stmt.Map("lt", stmt.Map("foo", 123)),
			want: &Condition{query: "`foo` < ?", args: []interface{}{123}}, wantErr: nil},
		{name: "take 11",
			args: stmt.Map("gte", stmt.Map("foo", 123)),
			want: &Condition{query: "`foo` >= ?", args: []interface{}{123}}, wantErr: nil},
		{name: "take 12",
			args: stmt.Map("lte", stmt.Map("foo", 123)),
			want: &Condition{query: "`foo` <= ?", args: []interface{}{123}}, wantErr: nil},
		{name: "take 13",
			args: stmt.Map("like", stmt.Map("foo", "string%")),
			want: &Condition{query: "`foo` LIKE ?", args: stmt.Slice("string%")}, wantErr: nil},
		{name: "take 14",
			args: stmt.Map("ISNULL", "foo"),
			want: &Condition{query: "`foo` IS NULL", args: []interface{}{}}, wantErr: nil},
		{name: "take 15",
			args: stmt.Map("in", stmt.Map("foo", stmt.Slice(123, 456, 789))),
			want: &Condition{query: "`foo` IN (?, ?, ?)", args: stmt.Slice(123, 456, 789)}, wantErr: nil},
		{name: "take 16",
			args: stmt.Map("in", stmt.Map("foo", 123)),
			want: &Condition{query: "`foo` IN (?)", args: stmt.Slice(123)}, wantErr: nil},
		{name: "take 16",
			args: stmt.Map("in", stmt.Map("foo", stmt.Slice(123))),
			want: &Condition{query: "`foo` IN (?)", args: stmt.Slice(123)}, wantErr: nil},
		{name: "take 16-error",
			args: stmt.Map("in", stmt.Map("foo", []interface{}{})),
			want: &Condition{query: "`foo` IN ()", args: []interface{}{}}, wantErr: fmt.Errorf("functor key=in value=map[foo:[]] value_type=map[string]interface {}: len(value) == 0")},
		{name: "take 16-error",
			args: stmt.Map("in", "foo"),
			want: &Condition{query: "`foo` IN ()", args: []interface{}{}}, wantErr: fmt.Errorf("functor key=in value=foo value_type=string: unsupported type")},
		{name: "take 17",
			// args: args{`{"BETWEEN":{"foo":[123,456]}}`},
			args: stmt.Map("between", stmt.Map("foo", stmt.Slice(123, 456))),
			want: &Condition{query: "`foo` BETWEEN ? AND ?", args: stmt.Slice(123, 456)}, wantErr: nil},
		{name: "take 17-error",
			args: stmt.Map("between", stmt.Map("foo", stmt.Slice(123, 456, 789))),
			want: &Condition{query: "`foo` BETWEEN ? AND ?", args: stmt.Slice(123, 456, 789)}, wantErr: fmt.Errorf("functor key=between value=map[foo:[123 456 789]] value_type=map[string]interface {}: len(value) != 2")},

		// TODO: Add test cases.
	}

	for _, test_case := range scenarios {

		builder := NewMysqlCondition()

		t.Run(test_case.name, func(t *testing.T) {
			// println(test_case.name)

			got, err := builder.Build(test_case.args)
			if test_case.wantErr != nil {
				if ErrorToString(test_case.wantErr) != ErrorToString(err) {
					t.Errorf("\nexpected=%v,\nactual  =%v", ErrorToString(test_case.wantErr), ErrorToString(err))
					return
				}
				return
			}

			if err != nil {
				t.Error(err)
				return
			}

			if !reflect.DeepEqual(got.Query(), test_case.want.Query()) {
				t.Errorf("\nexpected=%#v,\nactual  =%#v", test_case.want.Query(), got.Query())
			}

			if !reflect.DeepEqual(got.Args(), test_case.want.Args()) {
				t.Errorf("\nexpected=%#v,\nactual  =%#v", test_case.want.Args(), got.Args())
			}
		})
	}
}

func ErrorToString(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

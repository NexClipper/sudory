package flavor

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt"
)

func TestNewOrder(t *testing.T) {

	var scenarios = []struct {
		name    string
		args    interface{}
		want    stmt.OrderResult
		wantErr error
	}{
		{name: "take 1",
			args: M("asc", "foo"),
			want: &Order{order: "ASC", columns: []string{"foo"}}, wantErr: nil},
		{name: "take 1",
			args: M("asc", A("foo")),
			want: &Order{order: "ASC", columns: []string{"foo"}}, wantErr: nil},
		{name: "take 2",
			args: M("desc", "foo"),
			want: &Order{order: "DESC", columns: []string{"foo"}}, wantErr: nil},
		{name: "take 3",
			args: A(M("asc", "foo"), M("desc", "bar")),
			want: &CombineOrder{&Order{order: "ASC", columns: []string{"foo"}}, &Order{order: "DESC", columns: []string{"bar"}}}, wantErr: nil},
		{name: "take 3",
			args: A(M("asc", A("foo")), M("desc", A("bar"))),
			want: &CombineOrder{&Order{order: "ASC", columns: []string{"foo"}}, &Order{order: "DESC", columns: []string{"bar"}}}, wantErr: nil},
		{name: "take 4",
			args: M("asc", A("foo", "bar")),
			want: &Order{order: "ASC", columns: []string{"foo", "bar"}}, wantErr: nil},
		{name: "take 5",
			args: M("desc", A("foo", "bar")),
			want: &Order{order: "DESC", columns: []string{"foo", "bar"}}, wantErr: nil},
		{name: "take 6",
			args: A(M("asc", A("foo", "bar")), M("desc", A("foobar"))),
			want: &CombineOrder{&Order{order: "ASC", columns: []string{"foo", "bar"}}, &Order{order: "DESC", columns: []string{"foobar"}}}, wantErr: nil},
		{name: "take 7",
			args: A(M("asc", A("foo", "bar")), M("desc", A("foobar")), M("asc", A("barfoo"))),
			want: &CombineOrder{&Order{order: "ASC", columns: []string{"foo", "bar"}}, &Order{order: "DESC", columns: []string{"foobar"}}, &Order{order: "ASC", columns: []string{"barfoo"}}}, wantErr: nil},
		{name: "take 8",
			args: A("foo", "bar", "foobar"),
			want: &Order{order: "", columns: []string{"foo", "bar", "foobar"}}, wantErr: fmt.Errorf("value=foo type=string: unsupported type")},
		{name: "take 9",
			args: A("foo", "bar", "foobar", "DeSc"),
			want: &Order{order: "DESC", columns: []string{"foo", "bar", "foobar"}}, wantErr: fmt.Errorf("value=foo type=string: unsupported type")},
		{name: "take 10",
			args: A("foo", "foobar", "DESC", "bar", "barfoo", "ASC"),
			want: &CombineOrder{
				&Order{order: "DESC", columns: []string{"foo", "foobar"}},
				&Order{order: "ASC", columns: []string{"bar", "barfoo"}}}, wantErr: fmt.Errorf("value=foo type=string: unsupported type")},

		// TODO: Add test cases.
	}
	for _, test_case := range scenarios {

		builder := NewMysqlOrder()

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

			if !reflect.DeepEqual(got.Order(), test_case.want.Order()) {
				t.Errorf("\nexpected=%#v,\nactual  =%#v", test_case.want.Order(), got.Order())
			}
		})
	}
}

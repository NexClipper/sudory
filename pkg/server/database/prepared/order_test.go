package prepared

import (
	"reflect"
	"testing"
)

func TestNewOrder(t *testing.T) {
	type args struct {
		q string
	}
	tests := []struct {
		name    string
		args    args
		want    *Orders
		wantErr bool
	}{
		{name: "take 1",
			args: args{`{"ASC" : "foo"}`},
			want: (*Orders)(&[]Order{{order: "ASC", columns: []string{"foo"}}}), wantErr: false},
		{name: "take 2",
			args: args{`{"DESC" : "foo"}`},
			want: (*Orders)(&[]Order{{order: "DESC", columns: []string{"foo"}}}), wantErr: false},
		{name: "take 3",
			args: args{`{"ASC" : "foo", "DESC" : "bar"}`},
			want: (*Orders)(&[]Order{{order: "ASC", columns: []string{"foo"}}, {order: "DESC", columns: []string{"bar"}}}), wantErr: false},
		{name: "take 4",
			args: args{`{"ASC" : ["foo","bar"]}`},
			want: (*Orders)(&[]Order{{order: "ASC", columns: []string{"foo", "bar"}}}), wantErr: false},
		{name: "take 5",
			args: args{`{"DESC" : ["foo","bar"]}`},
			want: (*Orders)(&[]Order{{order: "DESC", columns: []string{"foo", "bar"}}}), wantErr: false},
		{name: "take 6",
			args: args{`[{"ASC" : ["foo","bar"]},{"DESC" : "foobar"}]`},
			want: (*Orders)(&[]Order{{order: "ASC", columns: []string{"foo", "bar"}}, {order: "DESC", columns: []string{"foobar"}}}), wantErr: false},
		{name: "take 7",
			args: args{`[{"ASC" : ["foo","bar"]},{"DESC" : "foobar"},{"ASC" : ["barfoo"]}]`},
			want: (*Orders)(&[]Order{{order: "ASC", columns: []string{"foo", "bar"}}, {order: "DESC", columns: []string{"foobar"}}, {order: "ASC", columns: []string{"barfoo"}}}), wantErr: false},
		{name: "take 8",
			args: args{`["foo","bar","foobar"]`},
			want: (*Orders)(&[]Order{{order: "", columns: []string{"foo", "bar", "foobar"}}}), wantErr: false},
		{name: "take 9",
			args: args{`["foo","bar","foobar", "DeSc"]`},
			want: (*Orders)(&[]Order{{order: "DESC", columns: []string{"foo", "bar", "foobar"}}}), wantErr: false},
		{name: "take 10",
			args: args{`["foo","foobar", "DESC", "bar", "barfoo", "ASC"]`},
			want: (*Orders)(&[]Order{{order: "DESC", columns: []string{"foo", "foobar"}}, {order: "ASC", columns: []string{"bar", "barfoo"}}}), wantErr: false},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			println(tt.name)

			got, err := NewOrder(tt.args.q)
			if (err != nil) && tt.wantErr {
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

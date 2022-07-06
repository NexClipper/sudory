package prepare

import (
	"reflect"
	"testing"
)

func TestNewCondition(t *testing.T) {
	type args struct {
		q string
	}
	tests := []struct {
		name    string
		args    args
		want    *Condition
		wantErr bool
	}{
		{name: "take 1",
			args: args{`{"EQ":{"foo":123}}`},
			want: &Condition{query: "foo = ?", args: []interface{}{123.0}}, wantErr: false},
		{name: "take 2",
			args: args{`{"AND":[{"EQ":{"foo":123}}]}`},
			want: &Condition{query: "(foo = ?)", args: []interface{}{123.0}}, wantErr: false},
		{name: "take 3",
			args: args{`{"AND":{"EQ":{"foo":123}}}`},
			want: &Condition{query: "(foo = ?)", args: []interface{}{123.0}}, wantErr: false},
		{name: "take 4",
			args: args{`{"AND":[{"EQ":{"foo":123}},{"EQ":{"bar":456}}, {"EQ":{"foobar":123456}}]}`},
			want: &Condition{query: "(foo = ? AND bar = ? AND foobar = ?)", args: []interface{}{123.0, 456.0, 123456.0}}, wantErr: false},
		{name: "take 5",
			args: args{`{"OR":{"EQ":{"foo":123}}}`},
			want: &Condition{query: "(foo = ?)", args: []interface{}{123.0}}, wantErr: false},
		{name: "take 6",
			args: args{`{"OR":[{"EQ":{"foo":123}},{"EQ":{"bar":456}}]}`},
			want: &Condition{query: "(foo = ? OR bar = ?)", args: []interface{}{123.0, 456.0}}, wantErr: false},
		{name: "take 7",
			args: args{`{"OR":[{"AND":[{"EQ":{"foo":123}},{"EQ":{"bar":456}}]}, {"EQ":{"foobar":123456}}]}`},
			want: &Condition{query: "((foo = ? AND bar = ?) OR foobar = ?)", args: []interface{}{123.0, 456.0, 123456.0}}, wantErr: false},
		{name: "take 8",
			args: args{`{"NOT":{"AND":[{"EQ":{"foo":123}},{"EQ":{"bar":456}}]}}`},
			want: &Condition{query: "NOT (foo = ? AND bar = ?)", args: []interface{}{123.0, 456.0}}, wantErr: false},
		{name: "take 9",
			args: args{`{"GT":{"foo": 123}}`},
			want: &Condition{query: "foo > ?", args: []interface{}{123.0}}, wantErr: false},
		{name: "take 10",
			args: args{`{"LT":{"foo": 123}}`},
			want: &Condition{query: "foo < ?", args: []interface{}{123.0}}, wantErr: false},
		{name: "take 11",
			args: args{`{"GE":{"foo": 123}}`},
			want: &Condition{query: "foo >= ?", args: []interface{}{123.0}}, wantErr: false},
		{name: "take 12",
			args: args{`{"LE":{"foo": 123}}`},
			want: &Condition{query: "foo <= ?", args: []interface{}{123.0}}, wantErr: false},
		{name: "take 13",
			args: args{`{"LIKE":{"foo": "string%"}}`},
			want: &Condition{query: "foo LIKE ?", args: []interface{}{"string%"}}, wantErr: false},
		{name: "take 14",
			args: args{`{"ISNULL":"foo"}`},
			want: &Condition{query: "foo IS ?", args: []interface{}{nil}}, wantErr: false},
		{name: "take 15",
			args: args{`{"IN":{"foo":[123,456,789]}}`},
			want: &Condition{query: "foo IN (?, ?, ?)", args: []interface{}{123.0, 456.0, 789.0}}, wantErr: false},
		{name: "take 16",
			args: args{`{"IN":{"foo":123}}`},
			want: &Condition{query: "foo IN (?)", args: []interface{}{123.0}}, wantErr: false},
		{name: "take 16-error",
			args: args{`{"IN":{"foo":[]}}`},
			want: &Condition{query: "foo IN ()", args: []interface{}{}}, wantErr: true},
		{name: "take 17",
			args: args{`{"BETWEEN":{"foo":[123,456]}}`},
			want: &Condition{query: "foo BETWEEN ? AND ?", args: []interface{}{123.0, 456.0}}, wantErr: false},
		{name: "take 17-error",
			args: args{`{"BETWEEN":{"foo":[123,456,789]}}`},
			want: &Condition{query: "foo BETWEEN ? AND ?", args: []interface{}{123.0, 456.0, 789.0}}, wantErr: true},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			println(tt.name)

			got, err := NewCondition(tt.args.q)
			if (err != nil) && tt.wantErr {
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCondition() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

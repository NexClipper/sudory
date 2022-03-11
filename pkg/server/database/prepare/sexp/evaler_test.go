package sexp

import (
	"reflect"
	"testing"
)

func TestEvalString(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name    string
		args    args
		want    Value
		wantErr bool
	}{
		{"take 1",
			args{`(eq foo 1)`},
			Value{argsValue, ArgsValue{"foo = ?", []interface{}{1.0}}}, false},
		{"take 2",
			args{`(and (eq foo 1))`},
			Value{argsValue, ArgsValue{"(foo = ?)", []interface{}{1.0}}}, false},
		{"take 3",
			args{`(and (eq foo 1) (eq bar 2))`},
			Value{argsValue, ArgsValue{"(foo = ? AND bar = ?)", []interface{}{1.0, 2.0}}}, false},
		{"take 4",
			args{`(and (eq foo 1) (eq bar 2) (eq foobar "foo bar"))`},
			Value{argsValue, ArgsValue{"(foo = ? AND bar = ? AND foobar = ?)", []interface{}{1.0, 2.0, "foo bar"}}}, false},
		{"take 5",
			args{`(or (eq foo 1))`},
			Value{argsValue, ArgsValue{"(foo = ?)", []interface{}{1.0}}}, false},
		{"take 6",
			args{`(or (eq foo 1) (eq bar 2))`},
			Value{argsValue, ArgsValue{"(foo = ? OR bar = ?)", []interface{}{1.0, 2.0}}}, false},
		{"take 7",
			args{`(not (eq foo 1))`},
			Value{argsValue, ArgsValue{"NOT foo = ?", []interface{}{1.0}}}, false},
		{"take 8",
			args{`(not (and (eq foo 1) (eq bar 2)))`},
			Value{argsValue, ArgsValue{"NOT (foo = ? AND bar = ?)", []interface{}{1.0, 2.0}}}, false},
		{"take 9",
			args{`(gt foo 1)`},
			Value{argsValue, ArgsValue{"foo > ?", []interface{}{1.0}}}, false},
		{"take 10",
			args{`(lt foo 1)`},
			Value{argsValue, ArgsValue{"foo < ?", []interface{}{1.0}}}, false},
		{"take 11",
			args{`(ge foo 1)`},
			Value{argsValue, ArgsValue{"foo >= ?", []interface{}{1.0}}}, false},
		{"take 12",
			args{`(le foo 1)`},
			Value{argsValue, ArgsValue{"foo <= ?", []interface{}{1.0}}}, false},
		{"take 13",
			args{`(like foo "foo%")`},
			Value{argsValue, ArgsValue{"foo LIKE ?", []interface{}{"foo%"}}}, false},
		{"take 14",
			args{`(isnull foo)`},
			Value{argsValue, ArgsValue{"foo IS ?", []interface{}{nil}}}, false},
		{"take 15",
			args{`(in foo (quote 1))`},
			Value{argsValue, ArgsValue{"foo IN (?)", []interface{}{1.0}}}, false},
		{"take 15",
			args{`(in foo '(1))`},
			Value{argsValue, ArgsValue{"foo IN (?)", []interface{}{1.0}}}, false},
		{"take 16",
			args{`(in foo '(1 2 3))`},
			Value{argsValue, ArgsValue{"foo IN (?, ?, ?)", []interface{}{1.0, 2.0, 3.0}}}, false},
		{"take 17",
			args{`(in foo 1 2 3)`},
			Value{argsValue, ArgsValue{"foo IN (?, ?, ?)", []interface{}{1.0, 2.0, 3.0}}}, false},
		{"take 18",
			args{`(between foo '(1 2))`},
			Value{argsValue, ArgsValue{"foo BETWEEN ? AND ?", []interface{}{1.0, 2.0}}}, false},
		{"take 19",
			args{`(between foo 1 2)`},
			Value{argsValue, ArgsValue{"foo BETWEEN ? AND ?", []interface{}{1.0, 2.0}}}, false},
		{"take 20",
			args{`(eq foobar "foo bar")`},
			Value{argsValue, ArgsValue{"foobar = ?", []interface{}{"foo bar"}}}, false},
		{"take 21",
			args{`(EQ foobar "foo bar")`},
			Value{argsValue, ArgsValue{"foobar = ?", []interface{}{"foo bar"}}}, false},
		{"take 22",
			args{`(EQ foobar foobar)`},
			Value{argsValue, ArgsValue{"foobar = ?", []interface{}{"foobar"}}}, false},
		{"take 23",
			args{`(in foo '("1" "2" "3"))`},
			Value{argsValue, ArgsValue{"foo IN (?, ?, ?)", []interface{}{"1", "2", "3"}}}, false},
		{"take 24",
			args{`(in foo "1" "2" "3")`},
			Value{argsValue, ArgsValue{"foo IN (?, ?, ?)", []interface{}{"1", "2", "3"}}}, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EvalString(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("EvalString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EvalString() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}

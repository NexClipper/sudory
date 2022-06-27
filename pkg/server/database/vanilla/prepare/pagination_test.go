package prepare

import (
	"reflect"
	"testing"
)

func TestNewPagination(t *testing.T) {
	type args struct {
		q string
	}
	tests := []struct {
		name    string
		args    args
		want    *Pagination
		wantErr bool
	}{

		{name: "take 1", args: args{`{}`}, want: &Pagination{limit: 255, page: 1}, wantErr: false},
		{name: "take 2", args: args{`{"limit":10,"page":2}`}, want: &Pagination{limit: 10, page: 2}, wantErr: false},
		{name: "take 3", args: args{`[{"limit":10},{"page":2}]`}, want: &Pagination{limit: 10, page: 2}, wantErr: false},
		{name: "take 4", args: args{`[{"limit":11},{"limit":10},{"page":3},{"page":2}]`}, want: &Pagination{limit: 10, page: 2}, wantErr: false},
		{name: "take 5", args: args{`{"limit":"10","page":"2"}`}, want: &Pagination{limit: 10, page: 2}, wantErr: false},
		{name: "take 6", args: args{`[{"limit":"10"},{"page":"2"}]`}, want: &Pagination{limit: 10, page: 2}, wantErr: false},
		{name: "take 7", args: args{`[{"limit":"11"},{"limit":"10"},{"page":"3"},{"page":"2"}]`}, want: &Pagination{limit: 10, page: 2}, wantErr: false},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPagination(tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPagination() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPagination() = %v, want %v", got, tt.want)
			}
		})
	}
}

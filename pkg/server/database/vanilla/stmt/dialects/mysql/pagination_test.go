package stmt

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestNewPagination(t *testing.T) {

	var scenarios = []struct {
		name    string
		args    interface{}
		want    *Pagination
		wantErr error
	}{

		{name: "take 1",
			args: map[string]interface{}{},
			want: &Pagination{limit: math.MaxInt8, page: 1}, wantErr: fmt.Errorf("empty object")},
		{name: "take 2",
			args: map[string]interface{}{"limit": 10, "page": 2},
			want: &Pagination{limit: 10, page: 2}, wantErr: nil},
		{name: "take 3",
			args: A(M("limit", 10), M("page", 2)),
			want: &Pagination{limit: 10, page: 2}, wantErr: nil},
		{name: "take 4",
			args: A(M("limit", 11), M("limit", 10), M("page", 3), M("page", 2)),
			want: &Pagination{limit: 10, page: 2}, wantErr: nil},
		{name: "take 5",
			args: map[string]interface{}{"limit": "10", "page": "2"},
			want: &Pagination{limit: 10, page: 2}, wantErr: nil},
		{name: "take 6",
			args: A(M("limit", "10"), M("page", "2")),
			want: &Pagination{limit: 10, page: 2}, wantErr: nil},
		{name: "take 7",
			args: A(M("limit", "11"), M("limit", "10"), M("page", "3"), M("page", "2")),
			want: &Pagination{limit: 10, page: 2}, wantErr: nil},

		// TODO: Add test cases.
	}

	for _, test_case := range scenarios {

		builder := NewMysqlPagination()

		t.Run(test_case.name, func(t *testing.T) {
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

			if !reflect.DeepEqual(got.String(), test_case.want.String()) {
				t.Errorf("\nexpected=%#v,\nactual  =%#v", test_case.want.String(), got.String())
			}
		})
	}
}

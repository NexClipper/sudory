package macro

import (
	"sort"
	"testing"
)

func TestSortSlice(t *testing.T) {

	var vv = []int{}
	for i := 0; i < 10; i++ {
		vv = append(vv, i)
	}

	t.Logf("%+v", vv)

	sort.Slice(vv, func(i, j int) bool {
		return vv[i] > vv[j]
	})

	t.Logf("%+v", vv)

}

func TestSortNilSlice(t *testing.T) {

	var vv []int

	t.Logf("%+v", vv)

	sort.Slice(vv, func(i, j int) bool {
		return vv[i] > vv[j]
	})

	t.Logf("%+v", vv)

}

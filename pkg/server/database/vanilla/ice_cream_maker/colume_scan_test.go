package ice_cream_maker_test

import (
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/ice_cream_maker"
)

func TestColumnScan(t *testing.T) {

	objs := []interface{}{
		ServiceStep_essential{},
		ServiceStep{},
	}

	s, err := ice_cream_maker.ColumnScan(objs...)
	if err != nil {
		t.Fatal(err)
	}

	println(s)
}

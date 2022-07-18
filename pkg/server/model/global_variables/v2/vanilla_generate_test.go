package v2_test

import (
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/ice_cream_maker"
	v2 "github.com/NexClipper/sudory/pkg/server/model/global_variables/v2"
)

var objs = []interface{}{
	v2.GlobalVariables_property{},
	v2.GlobalVariables{},
}

func TestNoXormColumns(t *testing.T) {
	s, err := ice_cream_maker.GenerateParts(objs, ice_cream_maker.Ingredients)
	if err != nil {
		t.Fatal(err)
	}

	println(s)
}

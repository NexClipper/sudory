package ice_cream_maker_test

import (
	"testing"

	vanilla "github.com/NexClipper/sudory/pkg/server/database/vanilla/ice_cream_maker"
)

func TestGenerateParts(t *testing.T) {
	objs := []interface{}{
		ServiceStep_essential{},
		ServiceStep{},
	}

	s, err := vanilla.GenerateParts(objs, vanilla.AllParts)
	if err != nil {
		t.Fatal(err)
	}

	println(s)
}

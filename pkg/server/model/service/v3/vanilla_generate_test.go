package v3_test

import (
	"os"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/ice_cream_maker"
	v3 "github.com/NexClipper/sudory/pkg/server/model/service/v3"
)

var objs = []interface{}{
	v3.Service_create{},
	v3.Service{},

	v3.ServiceResult_create{},
	v3.ServiceResult{},

	v3.ServiceStep_create{},
	v3.ServiceStep{},

	v3.Service_polling{},
}

func TestNoXormColumns(t *testing.T) {

	s, err := ice_cream_maker.GenerateParts(objs, ice_cream_maker.Ingredients)
	if err != nil {
		t.Fatal(err)
	}

	println(s)

	if true {
		filename := "vanilla_generated.go"
		fd, err := os.Create(filename)
		if err != nil {
			t.Fatal(err)
		}

		if _, err = fd.WriteString(s); err != nil {
			t.Fatal(err)
		}
	}
}

package service_test

import (
	"os"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/ice_cream_maker"
	"github.com/NexClipper/sudory/pkg/server/model/service/v3"
)

var objs = []interface{}{
	service.Service_create{},
	service.Service{},

	service.ServiceResult_create{},
	service.ServiceResult{},

	service.ServiceStep_create{},
	service.ServiceStep{},

	service.Service_polling{},
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

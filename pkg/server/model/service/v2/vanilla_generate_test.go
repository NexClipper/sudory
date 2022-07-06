package v2_test

import (
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/ice_cream_maker"
	v2 "github.com/NexClipper/sudory/pkg/server/model/service/v2"
)

var objs = []interface{}{
	v2.Service_essential{},
	v2.Service{},
	v2.ServiceStatus_essential{},
	v2.ServiceStatus{},
	v2.ServiceResults_essential{},
	v2.ServiceResult{},
	v2.Service_tangled{},
	v2.Service_status{},

	v2.ServiceStep_essential{},
	v2.ServiceStep{},
	v2.ServiceStepStatus_essential{},
	v2.ServiceStepStatus{},
	v2.ServiceStep_tangled{},
}

func TestNoXormColumns(t *testing.T) {

	s, err := ice_cream_maker.GenerateParts(objs, ice_cream_maker.AllParts)
	if err != nil {
		t.Fatal(err)
	}

	println(s)
}

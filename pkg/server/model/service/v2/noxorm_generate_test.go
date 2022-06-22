package v2_test

import (
	"fmt"
	"testing"

	noxorm "github.com/NexClipper/sudory/pkg/server/model/noxorm/v2"
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

	v2.ServiceStep_essential{},
	v2.ServiceStep{},
	v2.ServiceStepStatus_essential{},
	v2.ServiceStepStatus{},
	v2.ServiceStep_tangled{},
}

func TestNoXormColumns(t *testing.T) {
	{
		s, err := noxorm.ColumnPackage(objs...)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(s)
	}
	{
		s, err := noxorm.ColumnNames(objs...)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(s)
	}
	{
		s, err := noxorm.ColumnValues(objs...)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(s)
	}
	{
		s, err := noxorm.ColumnScan(objs...)
		if err != nil {
			t.Fatal(err)
		}

		fmt.Println(s)
	}
}

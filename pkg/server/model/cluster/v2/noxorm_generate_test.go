package v2_test

import (
	"fmt"
	"testing"

	v2 "github.com/NexClipper/sudory/pkg/server/model/cluster/v2"
	noxorm "github.com/NexClipper/sudory/pkg/server/model/noxorm/v2"
)

var objs = []interface{}{
	v2.Cluster_essential{},
	v2.Cluster{},
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

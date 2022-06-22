package v2_test

import (
	"fmt"
	"testing"

	noxorm "github.com/NexClipper/sudory/pkg/server/model/noxorm/v2"
	v2 "github.com/NexClipper/sudory/pkg/server/model/session/v2"
)

var objs = []interface{}{
	v2.Session_essential{},
	v2.Session{},
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

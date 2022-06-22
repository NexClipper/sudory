package v2_test

import (
	"testing"

	v2 "github.com/NexClipper/sudory/pkg/server/model/noxorm/v2"
)

func TestColumnValue(t *testing.T) {

	objs := []interface{}{
		ServiceStep_essential{},
		ServiceStep{},
	}

	v2.ColumnValues(objs)
}

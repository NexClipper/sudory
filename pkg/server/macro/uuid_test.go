package macro_test

import (
	"testing"

	"github.com/NexClipper/sudory/pkg/server/macro"
	"github.com/google/uuid"
)

func TestNewUuid(t *testing.T) {

	for i := 0; i < 10; i++ {
		println(macro.NewUuidString())
	}
}

func TestUuidParse(t *testing.T) {

	u, err := uuid.Parse("ab6e82680f79457d8ca67843fbe6ce2e")
	if err != nil {
		t.Error(err)
	}
	t.Log(u.Domain().String())
	t.Log(u.Variant().String())
	t.Log(u.Time().UnixTime())
	t.Log(u.String())
}

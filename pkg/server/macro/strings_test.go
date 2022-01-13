package macro_test

import (
	"testing"

	. "github.com/NexClipper/sudory/pkg/server/macro"
)

func TestStringJoin(t *testing.T) {

	jointer, builder := StringJoin(",")

	for i := 0; i < 10; i++ {
		jointer(i)
	}

	t.Log(builder())

}

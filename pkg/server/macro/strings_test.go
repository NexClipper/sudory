package macro_test

import (
	"testing"

	. "github.com/NexClipper/sudory/pkg/server/macro"
)

func TestStringJoin(t *testing.T) {

	adder, builder := StringBuilder()

	for i := 0; i < 10; i++ {
		adder(i)
	}

	t.Log(builder(","))

}

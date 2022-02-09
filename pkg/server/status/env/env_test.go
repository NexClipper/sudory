package env_test

import (
	"testing"

	"github.com/NexClipper/sudory/pkg/server/macro"
)

func TestGenerateUuid(t *testing.T) {

	for i := 0; i < 100; i++ {

		println(macro.NewUuidString())
	}

}

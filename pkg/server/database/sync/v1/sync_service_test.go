package v1

import (
	"testing"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	servicev1 "github.com/NexClipper/sudory/pkg/server/model/service/v1"
)

func TestServiceSync(t *testing.T) {
	sync := func() {

		engine := newEngine()

		model := new(servicev1.DbSchemaService)

		err := engine.Sync(model)

		if ErrorWithHandler(err) {
			t.Fatal(err)
		}
	}
	sync()
}

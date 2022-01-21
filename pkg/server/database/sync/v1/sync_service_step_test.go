package v1

import (
	"testing"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	servstepv1 "github.com/NexClipper/sudory/pkg/server/model/service_step/v1"
)

func TestServiceStepSync(t *testing.T) {
	sync := func() {

		engine := newEngine()

		model := new(servstepv1.DbSchemaServiceStep)

		err := engine.Sync(model)

		if ErrorWithHandler(err) {
			t.Fatal(err)
		}
	}
	sync()
}

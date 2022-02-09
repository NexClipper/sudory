package v1

import (
	"testing"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	clinetv1 "github.com/NexClipper/sudory/pkg/server/model/client/v1"
)

func TestClientSync(t *testing.T) {
	sync := func() {

		engine := newEngine()

		model := new(clinetv1.DbSchemaClient)

		err := engine.Sync(model)

		if ErrorWithHandler(err) {
			t.Fatal(err)
		}
	}
	sync()
}

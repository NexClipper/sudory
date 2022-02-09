package v1

import (
	"testing"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
)

func TestSessionSync(t *testing.T) {
	sync := func() {

		engine := newEngine()

		model := new(sessionv1.DbSchemaSession)

		err := engine.Sync(model)

		if ErrorWithHandler(err) {
			t.Fatal(err)
		}
	}
	sync()
}

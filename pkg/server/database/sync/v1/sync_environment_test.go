package v1

import (
	"testing"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	envv1 "github.com/NexClipper/sudory/pkg/server/model/environment/v1"
)

func TestEnvironmentSync(t *testing.T) {
	sync := func() {

		engine := newEngine()

		model := new(envv1.DbSchemaEnvironment)

		err := engine.Sync(model)

		if ErrorWithHandler(err) {
			t.Fatal(err)
		}
	}
	sync()
}

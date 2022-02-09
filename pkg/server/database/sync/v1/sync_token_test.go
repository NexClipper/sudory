package v1

import (
	"testing"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	tokenv1 "github.com/NexClipper/sudory/pkg/server/model/token/v1"
)

func TestTokenSync(t *testing.T) {
	sync := func() {

		engine := newEngine()

		model := new(tokenv1.DbSchemaToken)

		err := engine.Sync(model)

		if ErrorWithHandler(err) {
			t.Fatal(err)
		}
	}
	sync()
}

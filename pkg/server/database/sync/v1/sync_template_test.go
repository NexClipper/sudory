package v1

import (
	"testing"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	clinetv1 "github.com/NexClipper/sudory/pkg/server/model/template/v1"
)

func TestTemplateSync(t *testing.T) {
	sync := func() {

		engine := newEngine()

		model := new(clinetv1.DbSchemaTemplate)

		err := engine.Sync(model)

		if ErrorWithHandler(err) {
			t.Fatal(err)
		}
	}
	sync()
}

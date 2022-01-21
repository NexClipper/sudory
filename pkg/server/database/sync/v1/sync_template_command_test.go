package v1

import (
	"testing"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	commandv1 "github.com/NexClipper/sudory/pkg/server/model/template_command/v1"
)

func TestTemplateCommandSync(t *testing.T) {
	sync := func() {

		engine := newEngine()

		model := new(commandv1.TemplateCommand)

		err := engine.Sync(model)

		if ErrorWithHandler(err) {
			t.Fatal(err)
		}
	}
	sync()
}

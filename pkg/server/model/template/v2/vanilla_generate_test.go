package v2_test

import (
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/ice_cream_maker"
	v2 "github.com/NexClipper/sudory/pkg/server/model/template/v2"
)

var objs = []interface{}{
	v2.Template_essential{},
	v2.Template{},
	v2.TemplateCommand_essential{},
	v2.TemplateCommand{},
}

func TestNoXormColumns(t *testing.T) {
	s, err := ice_cream_maker.GenerateParts(objs, ice_cream_maker.AllParts)
	if err != nil {
		t.Fatal(err)
	}

	println(s)
}

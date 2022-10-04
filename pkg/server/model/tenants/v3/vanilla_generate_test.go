package tenants_test

import (
	"os"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/ice_cream_maker"
	"github.com/NexClipper/sudory/pkg/server/model/tenants/v3"
)

var objs = []interface{}{
	tenants.Tenant{},
	tenants.TenantClusters{},
	tenants.TenantChannels{},
}

func TestNoXormColumns(t *testing.T) {
	s, err := ice_cream_maker.GenerateParts(objs, ice_cream_maker.Ingredients)
	if err != nil {
		t.Fatal(err)
	}

	println(s)

	if true {
		filename := "vanilla_generated.go"
		fd, err := os.Create(filename)
		if err != nil {
			t.Fatal(err)
		}

		if _, err = fd.WriteString(s); err != nil {
			t.Fatal(err)
		}
	}
}

package vanilla

import (
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/ice_cream_maker"
	clusterv2 "github.com/NexClipper/sudory/pkg/server/model/cluster/v2"
	servicev2 "github.com/NexClipper/sudory/pkg/server/model/service/v2"
	sessionv2 "github.com/NexClipper/sudory/pkg/server/model/session/v2"
	templatev2 "github.com/NexClipper/sudory/pkg/server/model/template/v2"
)

var Schemas = []interface{}{

	servicev2.Service{},
	servicev2.ServiceStatus{},
	servicev2.ServiceResult{},
	servicev2.ServiceStep{},
	servicev2.ServiceStatus{},

	clusterv2.Cluster{},
	sessionv2.Session{},
	templatev2.Template{},
	templatev2.TemplateCommand{},
}

func TestVanilaFlavor(t *testing.T) {

	s, err := ice_cream_maker.GenerateParts(Schemas, ice_cream_maker.VanilaFlavor)
	if err != nil {
		t.Fatal(err)
	}

	println(s)
}

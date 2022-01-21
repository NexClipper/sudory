package v1

import (
	"testing"

	. "github.com/NexClipper/sudory/pkg/server/macro"
	clusterv1 "github.com/NexClipper/sudory/pkg/server/model/cluster/v1"
)

func TestClusterSync(t *testing.T) {
	sync := func() {

		engine := newEngine()

		model := new(clusterv1.DbSchemaCluster)

		err := engine.Sync(model)

		if ErrorWithHandler(err) {
			t.Fatal(err)
		}
	}
	sync()
}

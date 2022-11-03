package cluster_token_test

import (
	"os"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/ice_cream_maker"
	"github.com/NexClipper/sudory/pkg/server/model/cluster_token/v3"
)

var objs = []interface{}{
	cluster_token.ClusterToken{},
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

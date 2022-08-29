package v2_test

import (
	"os"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/ice_cream_maker"
	v2 "github.com/NexClipper/sudory/pkg/server/model/channel/v2"
)

var objs = []interface{}{
	v2.ManagedChannel{},
	v2.ManagedChannel_property{},
	v2.ManagedChannel_option{},
	v2.ManagedChannel_tangled{},

	// notifier edge
	v2.NotifierEdge{},
	v2.NotifierEdge_property{},
	v2.NotifierEdge_option{},

	// notifier
	v2.NotifierConsole{},
	// v2.NotifierConsole_property{},
	v2.NotifierWebhook{},
	v2.NotifierWebhook_property{},
	v2.NotifierRabbitMq{},
	v2.NotifierRabbitMq_property{},
	v2.NotifierSlackhook{},
	v2.NotifierSlackhook_property{},

	// notifier status option
	v2.ChannelStatusOption{},
	// notifier status
	v2.ChannelStatus{},
	// filter
	v2.Filter{},
	// format
	v2.Format{},
	v2.Format_property{},
}

func TestNoXormColumns(t *testing.T) {
	s, err := ice_cream_maker.GenerateParts(objs, append(ice_cream_maker.Ingredients, ice_cream_maker.ColumnNamesWithAlias))
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

package v3_test

import (
	"os"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/ice_cream_maker"
	v3 "github.com/NexClipper/sudory/pkg/server/model/channel/v3"
)

var objs = []interface{}{
	v3.ManagedChannel{},
	// v2.ManagedChannel_property{},
	// v3.ManagedChannel_option{},
	// v3.ManagedChannel_tangled{},

	// notifier edge
	v3.NotifierEdge{},
	v3.NotifierEdge_property{},
	// v3.NotifierEdge_option{},

	// notifier
	v3.NotifierConsole{},
	// v2.NotifierConsole_property{},
	v3.NotifierWebhook{},
	v3.NotifierWebhook_property{},
	v3.NotifierRabbitMq{},
	v3.NotifierRabbitMq_property{},
	v3.NotifierSlackhook{},
	v3.NotifierSlackhook_property{},

	// notifier status option
	v3.ChannelStatusOption_property{},
	v3.ChannelStatusOption{},
	// notifier status
	v3.ChannelStatus{},
	// filter
	v3.Filter{},
	// format
	v3.Format_property{},
	v3.Format{},
	// v3.Format_property{},
}

func TestNoXormColumns(t *testing.T) {
	// s, err := ice_cream_maker.GenerateParts(objs, append(ice_cream_maker.Ingredients, ice_cream_maker.ColumnNamesWithAlias))
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

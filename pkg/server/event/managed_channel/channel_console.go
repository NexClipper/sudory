package managed_channel

import (
	"fmt"
	"os"

	channelv2 "github.com/NexClipper/sudory/pkg/server/model/channel/v2"
	"github.com/pkg/errors"
)

type ChannelConsole struct {
	uuid string
	opt  *channelv2.NotifierConsole_property
	// sub event.EventNotifierMuxer
}

func NewChannelConsole(uuid string, opt channelv2.NotifierConsole_property) *ChannelConsole {
	notifier := &ChannelConsole{}
	notifier.uuid = uuid
	notifier.opt = &opt

	return notifier
}

func (channel ChannelConsole) Type() fmt.Stringer {
	return channel.opt.Type()
}
func (channel ChannelConsole) Uuid() string {
	return channel.uuid
}

func (channel ChannelConsole) Property() map[string]string {
	return map[string]string{
		"type": channel.opt.Type().String(),
		"uuid": channel.uuid,
	}
}

func (channel *ChannelConsole) Close() {}

func (channel ChannelConsole) OnNotify(factory *MarshalFactory) error {
	const content_type = "application/json"

	w := os.Stdout

	b, err := factory.Marshal(content_type)
	if err != nil {
		return errors.Wrapf(err, "marshal factory")
	}

	fmt.Fprintf(w, "%s\n", b)

	return nil
}

package managed_channel

import (
	"fmt"
	"os"

	channelv1 "github.com/NexClipper/sudory/pkg/server/model/channel/v2"
	"github.com/pkg/errors"
)

type ChannelConsole struct {
	uuid string
	opt  *channelv1.NotifierConsole_property
	// sub event.EventNotifierMuxer
}

func NewChannelConsole(uuid string, opt channelv1.NotifierConsole_property) *ChannelConsole {
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

func (channel ChannelConsole) OnNotify(factory MarshalFactoryResult) error {
	w := os.Stdout

	b, err := factory("application/json")
	if err != nil {
		return errors.Wrapf(err, "marshal factory")
	}

	fmt.Fprintf(w, "%s\n", b)

	return nil
}

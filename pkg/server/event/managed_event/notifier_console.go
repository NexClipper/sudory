package managed_event

import (
	"fmt"
	"os"

	channelv1 "github.com/NexClipper/sudory/pkg/server/model/channel/v1"
	"github.com/pkg/errors"
)

var _ Notifier = (*consoleNotifier)(nil)

type consoleNotifier struct {
	opt *channelv1.NotifierConsole
	// sub event.EventNotifierMuxer
}

func NewConsoleNotifier(opt *channelv1.NotifierConsole) *consoleNotifier {
	notifier := &consoleNotifier{}
	notifier.opt = opt

	return notifier
}

func (notifier consoleNotifier) Type() fmt.Stringer {
	return notifier.opt.Type()
}
func (notifier consoleNotifier) Uuid() string {
	return notifier.opt.Uuid
}

func (notifier consoleNotifier) Property() map[string]string {
	return map[string]string{
		"type": notifier.opt.Type().String(),
		"uuid": notifier.opt.Uuid,
	}
}

func (notifier *consoleNotifier) Close() {}

func (notifier consoleNotifier) OnNotify(factory MarshalFactoryResult) error {
	w := os.Stdout

	b, err := factory(notifier.opt.ContentType)
	if err != nil {
		return errors.Wrapf(err, "marshal factory")
	}
	for _, b := range b {
		fmt.Fprintf(w, "%s\n", b)
	}

	return nil
}

package managed_event

import (
	"fmt"
	"os"

	eventv1 "github.com/NexClipper/sudory/pkg/server/model/event/v1"
	"github.com/pkg/errors"
)

var _ Notifier = (*consoleNotifier)(nil)

type consoleNotifier struct {
	opt *eventv1.EventNotifierConsole
	// sub event.EventNotifierMuxer
}

func NewConsoleNotifier(opt *eventv1.EventNotifierConsole) *consoleNotifier {
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

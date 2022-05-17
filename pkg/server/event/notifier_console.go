package event

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
)

type consoleNotifier struct {
	sub EventNotifierMuxer
}

func NewConsoleNotifier() *consoleNotifier {
	notifier := &consoleNotifier{}

	return notifier
}
func (notifier consoleNotifier) Type() fmt.Stringer {
	return NotifierTypeConsole
}

func (notifier consoleNotifier) Property() map[string]string {
	return map[string]string{
		"name": notifier.sub.(EventNotifiMuxConfigHolder).Config().Name,
		"type": notifier.Type().String(),
	}
}

func (notifier *consoleNotifier) Regist(sub EventNotifierMuxer) {
	//Subscribe
	if !(sub == nil && notifier.sub != nil) {
		notifier.sub = sub
		notifier.sub.Notifiers().Add(notifier)
	}
}

func (notifier *consoleNotifier) Close() {
	//Unsubscribe
	if notifier.sub != nil {
		notifier.sub.Notifiers().Remove(notifier)
		notifier.sub = nil
	}
}

func (notifier consoleNotifier) OnNotify(factory MarshalFactoryResult) error {
	w := os.Stdout

	b, err := factory("application/json")
	if err != nil {
		return errors.Wrapf(err, "marshal factory")
	}
	for _, b := range b {
		fmt.Fprintf(w, "%s\n", b)
	}

	return nil
}

package event

import (
	"bytes"
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

type consoleNotifier struct {
	sub EventSubscriber
}

func NewConsoleNotifier() *consoleNotifier {
	notifier := &consoleNotifier{}

	return notifier
}
func (notifier consoleNotifier) Type() string {
	return NotifierTypeConsole.String()
}

func (notifier consoleNotifier) Property() map[string]string {
	return map[string]string{
		"name": notifier.sub.Config().Name,
		"type": notifier.Type(),
	}
}

func (notifier consoleNotifier) PropertyString() string {
	buff := bytes.Buffer{}
	for key, value := range notifier.Property() {
		if 0 < buff.Len() {
			buff.WriteString(" ")
		}
		buff.WriteString(key)
		buff.WriteString("=")
		buff.WriteString(strconv.Quote(value))
	}
	return buff.String()
}

func (notifier *consoleNotifier) Regist(sub EventSubscriber) {
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

func (notifier consoleNotifier) OnNotify(factory MarshalFactory) error {
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

func (notifier consoleNotifier) OnNotifyAsync(factory MarshalFactory) <-chan NotifierFuture {
	future := make(chan NotifierFuture)
	go func() {
		defer close(future)

		future <- NotifierFuture{&notifier, notifier.OnNotify(factory)}
	}()

	return future
}

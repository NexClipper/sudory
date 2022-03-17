package event

import (
	"bytes"
	"strconv"

	"github.com/NexClipper/sudory/pkg/server/event/filepool"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
)

type fileNotifier struct {
	opt FileNotifierConfig
	sub EventSubscriber
}

func NewFileNotifier(opt FileNotifierConfig) (*fileNotifier, error) {
	if _, err := filepool.OpenFile(opt.Path); err != nil {
		return nil, errors.Wrapf(err, "open filepool %s",
			logs.KVL(
				"file", opt.Path,
			))
	}

	notifier := &fileNotifier{}
	notifier.opt = opt

	return notifier, nil
}
func (notifier fileNotifier) Type() string {
	return NotifierTypeFile.String()
}

func (notifier fileNotifier) Property() map[string]string {
	return map[string]string{
		"name": notifier.sub.Config().Name,
		"type": notifier.Type(),
		"path": notifier.opt.Path,
	}
}

func (notifier fileNotifier) PropertyString() string {
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

func (notifier *fileNotifier) Regist(sub EventSubscriber) {
	//Subscribe
	if !(sub == nil && notifier.sub != nil) {
		notifier.sub = sub
		notifier.sub.Notifiers().Add(notifier)
	}
}

func (notifier *fileNotifier) Close() {
	//Unsubscribe
	if notifier.sub != nil {
		notifier.sub.Notifiers().Remove(notifier)
		notifier.sub = nil
	}

	filepath := notifier.opt.Path
	filepool.Close(filepath)
}

func (notifier fileNotifier) OnNotify(factory MarshalFactory) error {
	filepath := notifier.opt.Path

	buff := bytes.Buffer{}
	b, err := factory("application/json")
	if err != nil {
		return errors.Wrapf(err, "marshal factory")
	}
	for _, b := range b {
		buff.Write(b)
		buff.WriteByte('\n')
	}

	if _, err := filepool.OpenFile(filepath); err != nil {
		return errors.Wrapf(err, "open file %s",
			logs.KVL(
				"file", filepath,
			))
	}
	if err := filepool.Write(filepath, buff.Bytes()); err != nil {
		return errors.Wrapf(err, "write to file %s",
			logs.KVL(
				"file", filepath,
			))
	}

	return nil
}

func (notifier fileNotifier) OnNotifyAsync(factory MarshalFactory) <-chan NotifierFuture {
	future := make(chan NotifierFuture)
	go func() {
		defer close(future)

		future <- NotifierFuture{&notifier, notifier.OnNotify(factory)}
	}()

	return future
}

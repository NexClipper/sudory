package event

import (
	"bytes"
	"fmt"

	"github.com/NexClipper/sudory/pkg/server/event/filepool"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
)

type fileNotifier struct {
	opt FileNotifierConfig
	sub EventNotifierMultiplexer
}

func NewFileNotifier(opt FileNotifierConfig) (*fileNotifier, error) {
	if _, err := filepool.OpenFile(opt.Path); err != nil {
		return nil, errors.Wrapf(err, "open filepool%s",
			logs.KVL(
				"file", opt.Path,
			))
	}

	notifier := &fileNotifier{}
	notifier.opt = opt

	return notifier, nil
}
func (notifier fileNotifier) Type() fmt.Stringer {
	return NotifierTypeFile
}

func (notifier fileNotifier) Property() map[string]string {
	return map[string]string{
		"name": notifier.sub.(EventNotifiMuxConfigHolder).Config().Name,
		"type": notifier.Type().String(),
		"path": notifier.opt.Path,
	}
}

func (notifier *fileNotifier) Regist(sub EventNotifierMultiplexer) {
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

func (notifier fileNotifier) OnNotify(factory MarshalFactoryResult) error {
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
		return errors.Wrapf(err, "open file%s",
			logs.KVL(
				"file", filepath,
			))
	}
	if err := filepool.Write(filepath, buff.Bytes()); err != nil {
		return errors.Wrapf(err, "write to file%s",
			logs.KVL(
				"file", filepath,
			))
	}

	return nil
}

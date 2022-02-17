package logs

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/macro/logs/internal/serialize"
	"github.com/pkg/errors"
)

func put(keysAndValues ...interface{}) string {

	buf := bytes.Buffer{}

	serialize.KVListFormat(&buf, keysAndValues...)

	return buf.String()
}

type sinker interface {
	Id() uint64
	Names() []string
	Values() []interface{}
	Error() error

	WithId(id uint64) sinker
	WithName(name string) sinker
	WithError(err error) sinker
	WithValue(keysAndValues ...interface{}) sinker

	String() string
}

type sink struct {
	id            uint64
	names         []string
	keysAndValues []interface{}
	err           error
}

func (sink sink) String() string {
	var a [4]bool

	id := sink.Id()
	name := sink.Names()
	err := sink.Error()
	values := sink.Values()

	// if id != nil {
	a[0] = true
	// }
	if name != nil {
		a[1] = true
	}
	if err != nil {
		a[2] = true
	}

	if values != nil {
		a[3] = true
	}

	var s = bytes.Buffer{}
	for i := 0; i < len(a); i++ {

		if !a[i] {
			continue
		}
		if 0 < s.Len() {
			s.WriteString(" ")
		}
		switch i {
		case 0:
			s.WriteString(strconv.FormatUint(uint64(id), 10))
		case 1:
			s.WriteString("\"" + strings.Join(name, ",") + "\"")
		case 2:
			s.WriteString("err=" + err.Error())
		case 3:
			s.WriteString(put(values...))
		}

	}
	return s.String()
}

func (s sink) Id() uint64 {
	return s.id
}
func (s sink) Names() []string {
	return s.names
}
func (s sink) Values() []interface{} {
	return s.keysAndValues
}
func (s sink) Error() error {
	return s.err
}
func (s sink) WithId(id uint64) sinker {

	return &sink{id: id, names: s.names, keysAndValues: s.keysAndValues, err: s.err}
}
func (s sink) WithName(name string) sinker {

	names := []string{name}
	if s.names != nil {
		names = append(s.names, name)
	}

	return &sink{id: s.id, names: names, keysAndValues: s.keysAndValues, err: s.err}
}
func (s sink) WithError(err error) sinker {

	if s.err != nil {
		err = errors.WithMessage(err, s.err.Error())
	}

	return &sink{id: s.id, names: s.names, keysAndValues: s.keysAndValues, err: err}
}
func (s sink) WithValue(keysAndValues ...interface{}) sinker {

	if s.keysAndValues != nil {
		keysAndValues = append(s.keysAndValues, keysAndValues...)
	}

	return &sink{id: s.id, names: s.names, keysAndValues: keysAndValues, err: s.err}
}

func WithId(id uint64) sinker {
	return &sink{id: id}
}

func WithName(name string) sinker {
	return &sink{names: []string{name}}
}

func WithValue(keysAndValues ...interface{}) sinker {
	return &sink{keysAndValues: keysAndValues}
}
func WithError(err error) sinker {
	return &sink{err: err}
}

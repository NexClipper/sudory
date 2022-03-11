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
	errors        []error
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

	var buf = bytes.Buffer{}
	for i := 0; i < len(a); i++ {

		if !a[i] {
			continue
		}
		if 0 < buf.Len() {
			buf.WriteString(" ")
		}
		switch i {
		case 0:
			buf.WriteString(strconv.FormatUint(uint64(id), 10))
		case 1:
			buf.WriteString("\"" + strings.Join(name, "<") + "\"")
		case 2:
			buf.WriteString("err=" + err.Error())
		case 3:
			buf.WriteString(put(values...))
		}

	}
	return buf.String()
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

	var err error
	if 0 < len(s.errors) {
		err = s.errors[0]
		s.errors = s.errors[1:]
	}

	for n := range s.errors {
		err = errors.Wrapf(err, s.errors[n].Error())
	}

	return err
}
func (s sink) WithId(id uint64) sinker {

	return &sink{id: id, names: s.names, keysAndValues: s.keysAndValues, errors: s.errors}
}
func (s sink) WithName(name string) sinker {

	names := []string{name}
	if s.names != nil {
		names = append(s.names, names...)
	}

	return &sink{id: s.id, names: names, keysAndValues: s.keysAndValues, errors: s.errors}
}
func (s sink) WithError(err error) sinker {

	errors := []error{err}
	if s.errors != nil {
		errors = append(s.errors, errors...)
	}

	return &sink{id: s.id, names: s.names, keysAndValues: s.keysAndValues, errors: errors}
}
func (s sink) WithValue(keysAndValues ...interface{}) sinker {

	if s.keysAndValues != nil {
		keysAndValues = append(s.keysAndValues, keysAndValues...)
	}

	return &sink{id: s.id, names: s.names, keysAndValues: keysAndValues, errors: s.errors}
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
	return &sink{errors: []error{err}}
}

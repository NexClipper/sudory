package macro

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func ErrorWithHandler(err error, handler ...func(err error)) bool {
	if err == nil {
		return false
	}
	for n := range handler {
		func(fn func(err error)) {
			defer func() {
				_ = recover()
			}()
			fn(err)
		}(handler[n])
	}
	return true
}

func HasError(err error) bool {
	return err != nil
}

func Eqaul(a error, b ...error) bool {
	eqaul := func(a, b error) bool {
		if a == nil && b == nil {
			return true
		}
		if a == nil || b == nil {
			return false
		}

		x := strings.Split(a.Error(), "=")
		y := strings.Split(b.Error(), "=")

		return strings.Compare(x[0], y[0]) == 0
	}

	var ok bool = true
	for _, e := range b {
		ok = ok && eqaul(a, e)
	}
	return ok
}

func Stack(err error) string {

	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	buf := &bytes.Buffer{}

	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			fmt.Fprintf(buf, "%+s:%d\n", f, f)
		}
	}

	return buf.String()
}

func ErrorComposef(a, b error, format string, args ...interface{}) error {
	if b == nil {
		return a
	}
	if a != nil {
		return errors.Wrapf(a, errors.Wrapf(b, format, args...).Error())
	}

	return errors.Wrapf(b, format, args)
}
func ErrorCompose(a, b error) error {
	if b == nil {
		return a
	}
	if a != nil {
		return errors.Wrapf(a, b.Error())
	}

	return b
}

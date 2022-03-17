package logs

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
)

func CauseIter(err error, fn func(error)) {

	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
		fn(err)
	}
}

func StackIter(err error, fn func(string)) {

	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	if err, ok := err.(stackTracer); ok {
		buff := &bytes.Buffer{}
		for _, f := range err.StackTrace() {
			fmt.Fprintf(buff, "%+s:%d\n", f, f)
		}
		fn(buff.String())
	}
}

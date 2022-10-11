package stmt

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrorInvalidArgumentEmptyString = fmt.Errorf("empty string")
	ErrorInvalidArgumentEmptyObject = fmt.Errorf("empty object")
	ErrorNotFoundHandler            = fmt.Errorf("not found handler")
	ErrorUnsupportedType            = fmt.Errorf("unsupported type")
	ErrorUnsupportedPaginationKeys  = fmt.Errorf("unsupported pagination keys")
	ErrorUnsupportedOrderKeys       = fmt.Errorf("unsupported order keys")
)

// func ErrorComposef(a, b error, format string, args ...interface{}) error {
// 	if b == nil {
// 		return a
// 	}
// 	if a != nil {
// 		return errors.Wrapf(a, errors.Wrapf(b, format, args...).Error())
// 	}

// 	return errors.Wrapf(b, format, args)
// }

func ErrorCompose(a, b error) error {
	if b == nil {
		return a
	}
	if a != nil {
		return errors.Wrapf(a, b.Error())
	}

	return b
}

func CauseIter(err error, iter func(er error)) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		// call iter
		iter(err)

		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

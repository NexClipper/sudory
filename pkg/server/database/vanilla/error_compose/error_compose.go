package error_compose

import (
	"github.com/pkg/errors"
)

func Composef(a, b error, format string, args ...interface{}) error {
	if b == nil {
		return a
	}
	if a != nil {
		return errors.Wrapf(a, errors.Wrapf(b, format, args...).Error())
	}

	return errors.Wrapf(b, format, args)
}

func Compose(a, b error) error {
	if b == nil {
		return a
	}
	if a != nil {
		return errors.Wrapf(a, b.Error())
	}

	return b
}

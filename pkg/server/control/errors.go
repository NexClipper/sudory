package control

import (
	"fmt"

	"github.com/pkg/errors"
)

func ErrorInvaliedRequestParameter() error {
	return fmt.Errorf("invalid request parameter")
}
func ErrorInvaliedRequestParameterName(name string) error {
	return fmt.Errorf("invalid request parameter name='%s'", name)
}

// func ErrorInvaliedRequestParameterError(err error) error {
// 	return errors.WithMessage(err, "invalied request parameter")
// }

func ErrorBindRequestObject(err error) error {
	return errors.Wrapf(err, "cannot bind request")
}
func ErrorFailedCast() error {
	return fmt.Errorf("failed cast")
}

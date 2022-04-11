package control

import (
	"fmt"
)

func ErrorInvalidRequestParameter() error {
	return fmt.Errorf("invalid request parameter")
}

// func ErrorInvaliedRequestParameterName() error {
// 	return fmt.Errorf("invalid request parameter")
// }

// func ErrorInvaliedRequestParameterError(err error) error {
// 	return errors.WithMessage(err, "invalid request parameter")
// }

func ErrorBindRequestObject() error {
	return fmt.Errorf("cannot bind request")
}
func ErrorFailedCast() error {
	return fmt.Errorf("failed cast")
}

package control

import (
	"fmt"
)

func ErrorInvaliedRequestParameter() error {
	return fmt.Errorf("invalied request params")
}
func ErrorInvaliedRequestParameterName(name string) error {
	return fmt.Errorf("invalied request param name='%s'", name)
}
func ErrorBindRequestObject(err error) error {
	return fmt.Errorf("cannot bind request object; %w", err)
}
func ErrorFailedCast() error {
	return fmt.Errorf("failed cast")
}

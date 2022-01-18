package control

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb/errors"
)

func ErrorInvaliedRequestParameter() error {
	return errors.New("invalied request params")
}
func ErrorBindRequestObject(err error) error {
	return fmt.Errorf("cannot bind request object; %w", err)
}
func ErrorFailedCast() error {
	return fmt.Errorf("failed cast")
}

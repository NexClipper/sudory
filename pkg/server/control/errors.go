package control

import (
	"fmt"

	"github.com/pkg/errors"
)

func ErrorInvaliedRequestParameter() error {
	return fmt.Errorf("invalied request parameter")
}
func ErrorInvaliedRequestParameterName(name string) error {
	return fmt.Errorf("invalied request parameter name='%s'", name)
}

// func ErrorInvaliedRequestParameterError(err error) error {
// 	return errors.WithMessage(err, "invalied request parameter")
// }

func ErrorBindRequestObject(err error) error {
	return errors.WithMessage(err, "cannot bind request")
}
func ErrorFailedCast() error {
	return fmt.Errorf("failed cast")
}

type withMessage struct {
	cause error
	msg   string
}

func WithMessage(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   msg,
	}
}

func (e withMessage) Error() string {
	return fmt.Sprintf("%s: %s", e.msg, e.cause.Error())
}

func (e withMessage) Cause() error { return e.cause }

type withCode struct {
	error
	code int
}

func WithCode(err error, msg string, code int) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   msg,
	}
	return &withCode{
		error: err,
		code:  code,
	}
}

func (e withCode) Error() string {
	return fmt.Sprintf("code=%d: %s", e.code, e.error.Error())
}

func (e withCode) Cause() error { return e.error }

func (e withCode) Code() int { return e.code }

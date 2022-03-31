package control

import (
	"fmt"

	"github.com/labstack/echo/v4"
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

func HttpError(code int, err error) *echo.HTTPError {
	// msg := make([]string, 0, 3)
	// msg = append(msg, http.StatusText(code))
	// msg = append(msg, err.Error())
	// logs.CauseIter(err, func(err error) {
	// 	logs.StackIter(err, func(stack string) {
	// 		msg = append(msg, logs.KVL(
	// 			"stack", stack,
	// 		))
	// 	})
	// })

	return echo.NewHTTPError(code).SetInternal(err)
}

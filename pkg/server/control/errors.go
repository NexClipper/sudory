package control

import (
	"fmt"
)

var (
	ErrorInvalidRequestParameter = fmt.Errorf("invalid request parameter")
	ErrorBindRequestObject       = fmt.Errorf("could not bind request")
	ErrorFailedCast              = fmt.Errorf("failed to cast")
)

package vanilla

import (
	"os"
	"strconv"
)

// __VANILLA_DEBUG_PRINT_STATMENT__
//  print SQL statement for debug
//  env:
//    VANILLA_DEBUG_PRINT_STATMENT=false
var __VANILLA_DEBUG_PRINT_STATMENT__ = func() func() bool {
	ok, _ := strconv.ParseBool(os.Getenv("VANILLA_DEBUG_PRINT_STATMENT"))
	return func() bool {
		return ok
	}
}()

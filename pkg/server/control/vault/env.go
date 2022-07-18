package vault

import (
	"os"
	"strconv"
)

// __INIT_SLICE_CAPACITY__
//  init slice capacity
//  usage:
//    make([]interface{}, 0, __INIT_SLICE_CAPACITY__())
//  env:
//    INIT_SLICE_CAPACITY=5
var __INIT_SLICE_CAPACITY__ = func() func() int {
	n, err := strconv.Atoi(os.Getenv("INIT_SLICE_CAPACITY"))
	if err != nil {
		return func() int { return 5 }
	}

	if n < 0 {
		n = 1
	}
	return func() int { return n }
}()

package state

import (
	"os"
	"strconv"
	"time"
)

// __INIT_SLICE_CAPACITY__
//  init slice capacity
//  env:
//    INIT_SLICE_CAPACITY=5
var ENV__INIT_SLICE_CAPACITY__ = func() func() int {
	n, err := strconv.Atoi(os.Getenv("INIT_SLICE_CAPACITY"))
	if err != nil {
		return func() int { return 5 }
	}

	if n < 0 {
		n = 1
	}
	return func() int { return n }
}()

// __CONTROL_TRANSACTION_TIMEOUT__
//  transaction timeout for control package
//  env:
//    CONTROL_TRANSACTION_TIMEOUT=3s
var ENV__CONTROL_TRANSACTION_TIMEOUT__ = func() func() time.Duration {
	td, err := time.ParseDuration(os.Getenv("CONTROL_TRANSACTION_TIMEOUT"))
	if err != nil {
		return func() time.Duration { return 3 * time.Second }
	}
	return func() time.Duration { return td }
}()

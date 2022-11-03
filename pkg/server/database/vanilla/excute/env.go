package excute

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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

func VANILLA_DEBUG_PRINT(query string, args []interface{}) {
	if !__VANILLA_DEBUG_PRINT_STATMENT__() {
		return
	}

	aggregate := func(args []interface{}, seed string, accum func(a string, b interface{}) string) string {
		for i := range args {
			seed = accum(seed, args[i])
		}
		return seed
	}
	toString := func(v interface{}) string {
		switch v := v.(type) {
		case time.Time:
			return v.Format(time.RFC3339Nano)
		case fmt.Stringer:
			return v.String()
		case interface{ Print() string }:
			return v.Print()
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	accum := func(a string, b interface{}) string {
		if len(a) == 0 {
			return toString(b)
		}

		return strings.Join([]string{a, toString(b)}, ", ")
	}

	println("query:", query)
	println("args:", aggregate(args, "", accum))
}

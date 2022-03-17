package reflected

import (
	"reflect"
	"strconv"
)

// TypeName
//  reflect.TypeOf(v).String
func TypeName(i interface{}, quote ...bool) string {
	for i := range quote {
		if quote[i] {
			return strconv.Quote(reflect.TypeOf(i).String())
		}
	}

	return reflect.TypeOf(i).String()
}

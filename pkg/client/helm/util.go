package helm

import (
	"fmt"
	"reflect"
	"strings"
)

func convertArgsToStruct(args map[string]interface{}, obj interface{}) error {
	v := reflect.ValueOf(obj)

	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("non-pointer %v", v.Type().String())
	}
	if v.IsNil() {
		return fmt.Errorf("nil %v", v.Type().String())
	}

	v = v.Elem()

	for i := 0; i < v.NumField(); i++ {
		var tagName string
		var isOptional bool
		f := v.Type().Field(i)
		if paramTag := f.Tag.Get("param"); paramTag != "" {
			tags := strings.Split(paramTag, ",")
			tagName = tags[0]
			if len(tags) > 1 && tags[1] == "optional" {
				isOptional = true
			}
		} else {
			tagName = camelCaseToSnakeCase(f.Name)
		}

		if value, ok := args[tagName]; !ok {
			if isOptional {
				continue
			} else {
				return fmt.Errorf("argument(%s) does not exists", tagName)
			}
		} else {
			argValue := reflect.ValueOf(value)
			if v.Field(i).Type() == argValue.Type() {
				if v.Field(i).CanSet() {
					v.Field(i).Set(argValue)
				}
			} else {
				return fmt.Errorf("argument(%s)'s type is invalid: want(%v), got(%v)", tagName, v.Field(i).Type(), argValue.Type())
			}
		}
	}
	return nil
}

func camelCaseToSnakeCase(input string) string {
	var res string
	for k, v := range input {
		if k == 0 {
			res = strings.ToLower(string(v))
		} else {
			if v >= 'A' && v <= 'Z' {
				res += "_" + strings.ToLower(string(v))
			} else {
				res += string(v)
			}
		}
	}

	return res
}

package ice_cream_maker

import (
	"bytes"
	"reflect"
	"runtime"
)

var (
	Ingredients = []FuncPart{
		PrintWarning,
		ColumnPackage,
		ColumnNames,
		// ColumnNamesWithAlias,
		ColumnValues,
		ColumnScan,
		ColumnPtrs,
	}
)

type FuncPart = func(...interface{}) (string, error)

func GenerateParts(objs []interface{}, parts []FuncPart) (string, error) {
	buf := new(bytes.Buffer)
	appendString := func(s string) {
		buf.WriteString(s)
	}

	for _, part := range parts {

		part_name := GetFunctionName(part)

		s, err := part(objs...)
		if err != nil {
			return part_name, err
		}

		// appendString(fmt.Sprintf("// part of: %v\n", path.Base(part_name)))
		appendString(s)
	}

	return buf.String(), nil
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

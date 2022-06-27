package ice_cream_maker

import (
	"bytes"
	"fmt"
	"path"
	"reflect"
	"strings"
	"text/tabwriter"
)

func PkgNameDir(_type reflect.Type) string {
	return path.Base(path.Dir(_type.PkgPath()))
}

func PkgNameBase(_type reflect.Type) string {
	return path.Base(_type.PkgPath())
}

func FieldNames(_type reflect.Type) []string {
	names := make([]string, 0, 2*_type.NumField())
	for i := 0; i < _type.NumField(); i++ {
		field := _type.Field(i)
		names = append(names, field.Name)
	}

	return names
}

// ParseColumnTag
//  discover ColumnInfo by reflect
//  default path == ""
func ParseColumnTag(_type reflect.Type, rootpath string) []ColumnInfo {
	joinPath := func(path, base string) string {
		if len(path) == 0 {
			return base
		}
		return strings.Join([]string{path, base}, ".")
	}

	columnTag := func(sf reflect.StructField) (columntag ColumnTag, ok bool) {
		tag := sf.Tag.Get("column")
		if 0 < len(tag) {
			ok = true
		}
		if tag == "-" {
			ok = false
		}

		columntag = ColumnTag{
			Name: tag,
		}
		return
	}

	switch _type.Kind() {
	case reflect.Ptr:
		// is pointer type;
		elementType := _type.Elem()                      // get element type
		subTags := ParseColumnTag(elementType, rootpath) // recursive
		return subTags
	case reflect.Struct:
		// that type is want to
	default:
		panic("it is not a proper type")
	} // !switch

	if _type.NumField() == 0 {
		panic("no field")
	}

	tags := make([]ColumnInfo, 0, 2*_type.NumField())
	for i := 0; i < _type.NumField(); i++ {
		field := _type.Field(i)
		newbase := joinPath(rootpath, field.Name)

		switch field.Type.Kind() {
		case reflect.Ptr:
			if tag, ok := columnTag(field); !ok {
				// pointer with no tag
				elementType := field.Type.Elem()                 // get element type
				subTags := ParseColumnTag(elementType, rootpath) // recursive
				tags = append(tags, subTags...)
			} else {
				// pointer with tag
				tags = append(tags, ColumnInfo{
					Name:    tag.Name,
					Pointer: true,
					Path:    newbase,
					Type:    field.Type.String(),
				})
			}
		case reflect.Struct:
			if tag, ok := columnTag(field); !ok {
				// struct with no tag
				tags = append(tags, ParseColumnTag(field.Type, newbase)...) // recursive
			} else {
				// struct with tag
				tags = append(tags, ColumnInfo{
					Name:    tag.Name,
					Pointer: false,
					Path:    newbase,
					Type:    field.Type.String(),
				})
			}
		default:
			// default
			if tag, ok := columnTag(field); ok {
				// other type with tag
				tags = append(tags, ColumnInfo{
					Name:    tag.Name,
					Pointer: false,
					Path:    newbase,
					Type:    field.Type.String(),
				})
			}
		} // !switch
	} // !for

	return tags
}

// type MethodDiscriptor struct {
// 	StructName  string
// 	MethodName  string
// 	InputTypes  string
// 	OutputTypes string
// }

type ColumnInfo struct {
	Name    string
	Path    string
	Type    string
	Pointer bool
}

type ColumnTag struct {
	Name string
}

func (columninfo ColumnInfo) String() string {
	buf := bytes.Buffer{}
	func() {
		w := tabwriter.NewWriter(&buf, 1, 1, 2, ' ', 0)
		defer w.Flush()
		fmt.Fprintf(w, "%v:\t\n", columninfo.Name)
		fmt.Fprintf(w, "- %v:\t%v\n", "name", columninfo.Name)
		fmt.Fprintf(w, "- %v:\t%v\n", "path", columninfo.Path)
		fmt.Fprintf(w, "- %v:\t%v\n", "type", columninfo.Type)
		fmt.Fprintf(w, "- %v:\t%v\n", "pointer", columninfo.Pointer)
	}()

	return buf.String()
}

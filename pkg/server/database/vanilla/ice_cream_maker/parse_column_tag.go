package ice_cream_maker

import (
	"bytes"
	"database/sql"
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

type ParseColumnTag_opt struct {
	rootpath     string
	jointabletag *TableAliasTag
}

// ParseColumnTag
//  discover ColumnInfo by reflect
//  default path == ""
func ParseColumnTag(_type reflect.Type, opt ParseColumnTag_opt) []ColumnInfo {
	joinPath := func(path, base string, embed bool) string {
		if embed {
			return path
		}
		if len(path) == 0 {
			return base
		}
		return strings.Join([]string{path, base}, ".")
	}

	columnTag := func(sf reflect.StructField) (tag ColumnTag, ok bool) {
		column_tag := sf.Tag.Get("column")
		if 0 < len(column_tag) {
			ok = true
		}
		if column_tag == "-" {
			ok = false
		}

		tag_split := strings.Split(column_tag, ",")

		if 0 < len(tag_split) {
			tag = ColumnTag{
				Name: tag_split[0],
			}
		}
		if 1 < len(tag_split) {
			for _, tag_split := range tag_split {
				begin_tag := "default("
				begin := strings.Index(tag_split, begin_tag)
				end := strings.Index(tag_split, ")")
				if -1 < begin && -1 < end {
					tag.Default = tag_split[begin+len(begin_tag) : end]
				}
			}
		}

		return
	}

	jointableTag := func(sf reflect.StructField) (tag TableAliasTag, ok bool) {
		alias_tag := sf.Tag.Get("alias")
		if 0 < len(alias_tag) {
			ok = true
		}
		if alias_tag == "-" {
			ok = false
		}

		if tag_split := strings.Split(alias_tag, ","); 1 < len(tag_split) {
			tag = TableAliasTag{
				Alias:    tag_split[0],
				Mangling: sql.NullString{String: tag_split[1], Valid: true},
			}
			return
		}

		tag = TableAliasTag{
			Alias: alias_tag,
		}
		return
	}

	switch _type.Kind() {
	case reflect.Ptr:
		// is pointer type;
		elementType := _type.Elem()                                                                                        // get element type
		subTags := ParseColumnTag(elementType, ParseColumnTag_opt{rootpath: opt.rootpath, jointabletag: opt.jointabletag}) // recursive
		return subTags
	case reflect.Struct:
		// that type is want to
	default:
		panic("it is not a proper type")
	} // !switch

	if _type.NumField() == 0 {
		// panic("no field")
		return []ColumnInfo{}
	}

	tags := make([]ColumnInfo, 0, 2*_type.NumField())
	for i := 0; i < _type.NumField(); i++ {
		field := _type.Field(i)
		newbase := joinPath(opt.rootpath, field.Name, field.Anonymous)

		if tableAlias, ok := jointableTag(field); ok {
			opt.jointabletag = &tableAlias
		}

		switch field.Type.Kind() {
		case reflect.Ptr:
			if tag, ok := columnTag(field); !ok {
				// pointer with no tag
				elementType := field.Type.Elem()                                                                              // get element type
				subTags := ParseColumnTag(elementType, ParseColumnTag_opt{rootpath: newbase, jointabletag: opt.jointabletag}) // recursive
				tags = append(tags, subTags...)
			} else {
				// pointer with tag
				tags = append(tags, ColumnInfo{
					Name:          tag.Name,
					Default:       tag.Default,
					Pointer:       true,
					Path:          newbase,
					Type:          field.Type.String(),
					TableAliasTag: opt.jointabletag,
				})
			}
		case reflect.Struct:
			if tag, ok := columnTag(field); !ok {
				// struct with no tag
				tags = append(tags, ParseColumnTag(field.Type, ParseColumnTag_opt{rootpath: newbase, jointabletag: opt.jointabletag})...) // recursive
			} else {
				// struct with tag
				tags = append(tags, ColumnInfo{
					Name:          tag.Name,
					Default:       tag.Default,
					Pointer:       false,
					Path:          newbase,
					Type:          field.Type.String(),
					TableAliasTag: opt.jointabletag,
				})
			}
		default:
			// default
			if tag, ok := columnTag(field); ok {
				// other type with tag
				tags = append(tags, ColumnInfo{
					Name:          tag.Name,
					Default:       tag.Default,
					Pointer:       false,
					Path:          newbase,
					Type:          field.Type.String(),
					TableAliasTag: opt.jointabletag,
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
	Name          string
	Default       string
	Path          string
	Type          string
	Pointer       bool
	TableAliasTag *TableAliasTag
}

type ColumnTag struct {
	Name    string
	Default string
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

type TableAliasTag struct {
	Alias    string
	Mangling sql.NullString
}

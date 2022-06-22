package v2

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"
)

func ColumnValues(objs ...interface{}) (s string, err error) {
	tmpl := TemplateColumnValue{
		Discriptors: make([]TemplateColumnValueDiscriptor, len(objs)),
	}

	for i := range objs {
		if len(tmpl.Package) == 0 {
			tmpl.Package = PkgName(reflect.TypeOf(objs[i]))
		}

		columninfos := ColumnInfos(reflect.TypeOf(objs[i]), "")

		vals := make([]string, 0)
		for j := range columninfos {
			ref := columninfos[j].Path
			ref = "row." + ref
			if columninfos[j].Pointer {
				ref = "*" + ref
			}
			vals = append(vals, ref)
		}

		discriptor := TemplateColumnValueDiscriptor{
			StructName: reflect.TypeOf(objs[i]).Name(),
			Values:     strings.Join(vals, ", "),
		}

		tmpl.Discriptors[i] = discriptor
	}

	column_names, err := template.New("column_value.go").Parse(strings.TrimSpace(tmpl.Text()))
	if err != nil {
		return
	}

	buf := bytes.Buffer{}
	if err = column_names.Execute(&buf, tmpl); err != nil {
		return
	}
	s = buf.String()

	// filename := "colume_names.go"
	// fd, err := os.Create(filename)
	// if err != nil {
	// 	return
	// }
	// defer fd.Close()

	// if err = column_names.Execute(fd, tmpl); err != nil {
	// 	return
	// }

	return
}

type TemplateColumnValue struct {
	Package     string
	Discriptors []TemplateColumnValueDiscriptor
}

type TemplateColumnValueDiscriptor struct {
	StructName string
	Values     string // scan.Foo, scan.Bar, *scan.Baz
}

func (TemplateColumnValue) Text() string {
	return `
{{ range .Discriptors }} 
func (row {{ .StructName }}) Values() []interface{} {
	return []interface{}{ {{ .Values }} }
}
{{ end }}
`
}

package v2

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"
)

func ColumnPackage(objs ...interface{}) (s string, err error) {
	tmpl := TemplateColumnPackage{}

	for i := range objs {
		if len(tmpl.Package) == 0 {
			tmpl.Package = PkgName(reflect.TypeOf(objs[i]))
		}
	}

	column_names, err := template.New("column_package.go").Parse(strings.TrimSpace(tmpl.Text()))
	if err != nil {
		return
	}

	buf := bytes.Buffer{}
	if err = column_names.Execute(&buf, tmpl); err != nil {
		return
	}
	s = buf.String()

	return
}

type TemplateColumnPackage struct {
	Package string
	// Discriptors []TemplateColumnPackageDiscriptor
}

type TemplateColumnPackageDiscriptor struct {
	StructName string
	References string // &scan.Foo, &scan.Bar, scan.Baz
}

func (TemplateColumnPackage) Text() string {
	return `
package {{ .Package }}
`
}

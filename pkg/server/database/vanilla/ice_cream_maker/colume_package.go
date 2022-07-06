package ice_cream_maker

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"
)

func ColumnPackage(objs ...interface{}) (s string, err error) {
	tmpl := TemplateColumnPackage{Package: make([]string, 0, 1)}

	for i := range objs {
		if len(tmpl.Package) == 0 {
			tmpl.Package = append(tmpl.Package, PkgNameBase(reflect.TypeOf(objs[i])))
		}
	}

	tmpl_, err := template.New("column_package.go").Parse(strings.TrimSpace(tmpl.Text()))
	if err != nil {
		return
	}

	buf := bytes.Buffer{}
	if err = tmpl_.Execute(&buf, tmpl); err != nil {
		return
	}

	s = buf.String()
	return
}

type TemplateColumnPackage struct {
	Package []string
	// Discriptors []TemplateColumnPackageDiscriptor
}

func (TemplateColumnPackage) Text() string {
	return `
{{ range $index, $pkg := .Package -}} 
package {{ $pkg }}
{{ end }}
`
}

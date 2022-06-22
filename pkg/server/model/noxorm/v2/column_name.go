package v2

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

func ColumnNames(objs ...interface{}) (s string, err error) {
	tmpl := TemplateColumnNames{
		Discriptors: make([]TemplateColumnNamesDiscriptor, len(objs)),
	}

	for i := range objs {
		if len(tmpl.Package) == 0 {
			tmpl.Package = PkgName(reflect.TypeOf(objs[i]))
		}

		columninfos := ColumnInfos(reflect.TypeOf(objs[i]), "")

		names := make([]string, 0)
		for j := range columninfos {
			name := columninfos[j].Name
			names = append(names, strconv.Quote(name))
		}

		discriptor := TemplateColumnNamesDiscriptor{
			StructName:  reflect.TypeOf(objs[i]).Name(),
			ColumnNames: strings.Join(names, ", "),
		}

		tmpl.Discriptors[i] = discriptor
	}

	column_names, err := template.New("column_names.go").Parse(strings.TrimSpace(tmpl.Text()))
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

type TemplateColumnNames struct {
	Package     string
	Discriptors []TemplateColumnNamesDiscriptor
}

type TemplateColumnNamesDiscriptor struct {
	StructName  string
	ColumnNames string
}

func (TemplateColumnNames) Text() string {
	return `
{{ range .Discriptors }} 
func ({{ .StructName }}) ColumnNames() []string {
	return []string{ {{ .ColumnNames }} }
}
{{ end }}
`
}

package ice_cream_maker

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"
)

func ColumnNames(objs ...interface{}) (s string, err error) {
	tmpl := TemplateColumnNames{
		Discriptors: make([]TemplateColumnNamesDiscriptor, len(objs)),
	}

	for i := range objs {
		columninfos := ParseColumnTag(reflect.TypeOf(objs[i]), "")

		names := make([]string, 0)
		for j := range columninfos {
			name := columninfos[j].Name
			// names = append(names, strconv.Quote(name))
			names = append(names, name)
		}

		discriptor := TemplateColumnNamesDiscriptor{
			StructName:  reflect.TypeOf(objs[i]).Name(),
			ColumnNames: names,
		}

		tmpl.Discriptors[i] = discriptor
	}

	tmpl_, err := template.New("column_names.go").Parse(strings.TrimSpace(tmpl.Text()))
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

type TemplateColumnNames struct {
	Discriptors []TemplateColumnNamesDiscriptor
}

type TemplateColumnNamesDiscriptor struct {
	StructName  string
	ColumnNames []string
}

func (TemplateColumnNames) Text() string {
	return `
{{ $space := " " }}
{{- range .Discriptors }} 
func ({{ .StructName }}) ColumnNames() []string {
	return []string{ 
		{{ range $index, $columnName := .ColumnNames -}}
 		"{{ $columnName }}",{{ $space }}
		{{- end }}
	}
}
{{ end }}
`
}

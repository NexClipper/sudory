package ice_cream_maker

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
		columninfos := ParseColumnTag(reflect.TypeOf(objs[i]), "")

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
			Values:     vals,
		}

		tmpl.Discriptors[i] = discriptor
	}

	tmpl_, err := template.New("column_value.go").Parse(strings.TrimSpace(tmpl.Text()))
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

type TemplateColumnValue struct {
	Discriptors []TemplateColumnValueDiscriptor
}

type TemplateColumnValueDiscriptor struct {
	StructName string
	Values     []string // scan.Foo, scan.Bar, *scan.Baz
}

func (TemplateColumnValue) Text() string {
	return `
{{ $space := " " }}
{{- range .Discriptors }} 
func (row {{ .StructName }}) Values() []interface{} {
	return []interface{}{ 
		{{ range $index, $value := .Values -}}
		{{ $value }},{{ $space }}
		{{- end }}
	}
}
{{ end -}}
`
}

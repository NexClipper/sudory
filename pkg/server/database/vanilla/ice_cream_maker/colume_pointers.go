package ice_cream_maker

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"
)

func ColumnPtrs(objs ...interface{}) (s string, err error) {
	tmpl := TemplateColumnPtr{
		Discriptors: make([]TemplateColumnPtrDiscriptor, len(objs)),
	}

	for i := range objs {

		columninfos := ParseColumnTag(reflect.TypeOf(objs[i]), ParseColumnTag_opt{})

		vals := make([]string, 0)
		for j := range columninfos {
			ref := columninfos[j].Path
			ref = "row." + ref
			if !columninfos[j].Pointer {
				ref = "&" + ref
			}
			vals = append(vals, ref)
		}

		discriptor := TemplateColumnPtrDiscriptor{
			StructName: reflect.TypeOf(objs[i]).Name(),
			Ptrs:       vals,
		}

		tmpl.Discriptors[i] = discriptor
	}

	tmpl_, err := template.New("colume_ptrs.go").Parse(strings.TrimSpace(tmpl.Text()))
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

type TemplateColumnPtr struct {
	Discriptors []TemplateColumnPtrDiscriptor
}

type TemplateColumnPtrDiscriptor struct {
	StructName string
	Ptrs       []string // &scan.Foo, &scan.Bar, scan.Baz
}

func (TemplateColumnPtr) Text() string {
	return `
{{ $space := "" }}
{{- range .Discriptors }} 
func (row *{{ .StructName }}) Ptrs() []interface{} {
	return []interface{}{ 
		{{- range $index, $value := .Ptrs }}
		{{ $value }},{{ $space }}
		{{- end }}
	}
}
{{ end -}}
`
}

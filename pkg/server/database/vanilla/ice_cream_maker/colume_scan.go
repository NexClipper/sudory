package ice_cream_maker

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"
)

func ColumnScan(objs ...interface{}) (s string, err error) {
	tmpl := TemplateColumnScan{
		Discriptors: make([]TemplateColumnScanDiscriptor, len(objs)),
	}

	for i := range objs {
		columninfos := ParseColumnTag(reflect.TypeOf(objs[i]), ParseColumnTag_opt{})

		refs := make([]string, 0)
		for j := range columninfos {
			ref := columninfos[j].Path
			ref = "row." + ref
			if !columninfos[j].Pointer {
				ref = "&" + ref
			}
			refs = append(refs, ref)
		}

		discriptor := TemplateColumnScanDiscriptor{
			StructName: reflect.TypeOf(objs[i]).Name(),
			References: refs,
		}

		tmpl.Discriptors[i] = discriptor
	}

	tmpl_, err := template.New("column_scans.go").Parse(strings.TrimSpace(tmpl.Text()))
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

type TemplateColumnScan struct {
	Discriptors []TemplateColumnScanDiscriptor
}

type TemplateColumnScanDiscriptor struct {
	StructName string
	References []string // &scan.Foo, &scan.Bar, scan.Baz
}

func (TemplateColumnScan) Text() string {
	return `
{{ $space := "" }}
type Scanner interface {
	Scan(dest ...interface{}) error
}
{{ range .Discriptors }} 
func (row *{{ .StructName }}) Scan(scanner Scanner) error {
	return scanner.Scan(
		{{- range $index, $ref := .References }}
		{{ $ref }},{{ $space }}
		{{- end }}
	)
}
{{ end -}}
`
}

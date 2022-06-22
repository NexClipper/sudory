package v2

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
		if len(tmpl.Package) == 0 {
			tmpl.Package = PkgName(reflect.TypeOf(objs[i]))
		}

		columninfos := ColumnInfos(reflect.TypeOf(objs[i]), "")

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
			References: strings.Join(refs, ", "),
		}

		tmpl.Discriptors[i] = discriptor
	}

	column_names, err := template.New("column_scans.go").Parse(strings.TrimSpace(tmpl.Text()))
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

type TemplateColumnScan struct {
	Package     string
	Discriptors []TemplateColumnScanDiscriptor
}

type TemplateColumnScanDiscriptor struct {
	StructName string
	References string // &scan.Foo, &scan.Bar, scan.Baz
}

func (TemplateColumnScan) Text() string {
	return `	
type Scanner interface {
	Scan(dest ...interface{}) error
}
{{ range .Discriptors }} 
func (row *{{ .StructName }}) Scan(scanner Scanner) error {
	return scanner.Scan({{ .References }})
}
{{ end }}
`
}

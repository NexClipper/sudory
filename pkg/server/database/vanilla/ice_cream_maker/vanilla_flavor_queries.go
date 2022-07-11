package ice_cream_maker

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"
)

func VanillaFlavorQueries(objs ...interface{}) (s string, err error) {
	tmpl := TemplateVanilaFlavorQueries{
		Discriptors: make([]TemplateVanilaFlavorQueriesDiscriptor, len(objs)),
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

		discriptor := TemplateVanilaFlavorQueriesDiscriptor{
			PackageName: PkgNameDir(reflect.TypeOf(objs[i])),
			StructName:  reflect.TypeOf(objs[i]).Name(),
			References:  refs,
		}

		tmpl.Discriptors[i] = discriptor
	}

	tmpl_, err := template.New("vanilla_flavor_queries.go").Parse(strings.TrimSpace(tmpl.Text()))
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

type TemplateVanilaFlavorQueries struct {
	// Package     string
	Discriptors []TemplateVanilaFlavorQueriesDiscriptor
}

type TemplateVanilaFlavorQueriesDiscriptor struct {
	PackageName string
	StructName  string
	References  []string // &scan.Foo, &scan.Bar, scan.Baz
}

func (TemplateVanilaFlavorQueries) Text() string {
	return `
type Preparer interface {
	Prepare(query string) (*sql.Stmt, error)
}

type Condition interface {
	Query() string
	Args() []interface{}
}

type Order interface {
	Order() string
}

type Pagination interface {
	Offset() int
	Limit() int
}
{{ range .Discriptors }} 
func {{ .StructName }}_QueryRow(tx Preparer, q Condition, o Order, p Pagination) (r *{{ .PackageName }}.{{ .StructName }}, err error) {
	r = new({{ .PackageName }}.{{ .StructName }})
	column_names := strings.Join({{ .StructName }}_ColumnNames(), ", ")
	table_name := {{ .StructName }}_TableName()
	args := make([]interface{}, 0, 10)
	s := fmt.Sprintf("SELECT %v FROM %v", column_names, table_name)
	if q != nil {
		s += "\nWHERE " + q.Query()
		args = append(args, q.Args()...)
	}
	if o != nil {
		s += "\nORDER BY " + o.Order()
	}
	if p != nil {
		s += fmt.Sprintf("\nLIMIT %v, %v", p.Offset(), p.Limit())
	}

	stmt, err := tx.Prepare(s)
	err = errors.Wrapf(err, "sql.DB.Prepare")
	if err != nil {
		return
	}
	defer func() {
		err = error_compose.Composef(err, stmt.Close(), "sql.Stmt.Close")
	}()

	row := stmt.QueryRow(args...)

	err = row.Scan(r.Dests()...)
	err = errors.Wrapf(err, "sql.Row.Scan")
	err = error_compose.Composef(err, row.Err(), "sql.Row; during scan")

	err = errors.Wrapf(err, "faild to query row table=\"%v\"",
		table_name,
	)

	return
}
{{ end -}}
`
}

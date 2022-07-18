package ice_cream_maker

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

const __TABLE_ALIAS_SEPERATOR__ = "."
const __COLUMN_ALIAS_SEPERATOR__ = "."

func ColumnNames(objs ...interface{}) (s string, err error) {
	tmpl := TemplateColumnNames{
		Discriptors: make([]TemplateColumnNamesDiscriptor, len(objs)),
	}

	for i := range objs {
		columninfos := ParseColumnTag(reflect.TypeOf(objs[i]), ParseColumnTag_opt{})

		names := make([]string, 0)
		for j := range columninfos {
			columninfo := columninfos[j]

			name := columninfo.Name
			names = append(names, name)

			if columninfo.TableAliasTag != nil {
				// aliases := make([]string, 0, 2)
				manglings := make([]string, 0, 2)
				// if 0 < len(columninfo.TableAliasTag.Alias) {
				// 	aliases = append(aliases, columninfo.TableAliasTag.Alias)
				// }

				if 0 < len(columninfo.TableAliasTag.Mangling.String) {
					manglings = append(manglings, columninfo.TableAliasTag.Mangling.String)
				}

				// aliases = append(aliases, name)
				manglings = append(manglings, name)

				// names[j] = fmt.Sprintf("%v AS `%v`",
				// 	strings.Join(aliases, __TABLE_ALIAS_SEPERATOR__),
				// 	strings.Join(manglings, __COLUMN_ALIAS_SEPERATOR__))
				names[j] = fmt.Sprintf("`%v`",
					strings.Join(manglings, __COLUMN_ALIAS_SEPERATOR__))

				if 0 < len(columninfo.Default) {
					names[j] = fmt.Sprintf("IFNULL(%v, %v) AS %v",
						names[j],
						columninfo.Default,
						names[j])
				}
			}
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
{{ $space := "" }}
{{- range .Discriptors }} 
func ({{ .StructName }}) ColumnNames() []string {
	return []string{ 
		{{- range $index, $columnName := .ColumnNames }}
 		"{{ $columnName }}",{{ $space }}
		{{- end }}
	}
}
{{ end }}
`
}

func ColumnNamesWithAlias(objs ...interface{}) (s string, err error) {
	tmpl := TemplateColumnNamesWithAlias{
		Discriptors: make([]TemplateColumnNamesDiscriptorWithAlias, len(objs)),
	}

	for i := range objs {
		columninfos := ParseColumnTag(reflect.TypeOf(objs[i]), ParseColumnTag_opt{})

		names := make([]string, 0)
		for j := range columninfos {
			columninfo := columninfos[j]

			name := columninfo.Name
			names = append(names, name)

			if columninfo.TableAliasTag != nil {
				aliases := make([]string, 0, 2)
				manglings := make([]string, 0, 2)
				if 0 < len(columninfo.TableAliasTag.Alias) {
					aliases = append(aliases, columninfo.TableAliasTag.Alias)
				}

				if 0 < len(columninfo.TableAliasTag.Mangling.String) {
					manglings = append(manglings, columninfo.TableAliasTag.Mangling.String)
				}

				aliases = append(aliases, name)
				manglings = append(manglings, name)

				if 0 < len(columninfo.Default) {
					names[j] = fmt.Sprintf("IFNULL(%v, %v) AS `%v`",
						strings.Join(aliases, __TABLE_ALIAS_SEPERATOR__),
						columninfo.Default,
						strings.Join(manglings, __COLUMN_ALIAS_SEPERATOR__))
				} else {
					names[j] = fmt.Sprintf("%v AS `%v`",
						strings.Join(aliases, __TABLE_ALIAS_SEPERATOR__),
						strings.Join(manglings, __COLUMN_ALIAS_SEPERATOR__))
				}
			}
		}

		discriptor := TemplateColumnNamesDiscriptorWithAlias{
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

type TemplateColumnNamesWithAlias struct {
	Discriptors []TemplateColumnNamesDiscriptorWithAlias
}

type TemplateColumnNamesDiscriptorWithAlias struct {
	StructName  string
	ColumnNames []string
}

func (TemplateColumnNamesWithAlias) Text() string {
	return `
{{ $space := "" }}
{{- range .Discriptors }} 
func ({{ .StructName }}) ColumnNamesWithAlias() []string {
	return []string{ 
		{{- range $index, $columnName := .ColumnNames }}
 		"{{ $columnName }}",{{ $space }}
		{{- end }}
	}
}
{{ end }}
`
}

package ice_cream_maker

import (
	"bytes"
	"strings"
	"text/template"
)

func VanillaFlavorPackage(pkg string) func(objs ...interface{}) (s string, err error) {
	return func(objs ...interface{}) (s string, err error) {
		tmpl := TemplateVanillaFlavorPackagePackage{Packages: []string{
			pkg,
		}}

		tmpl_, err := template.New("vanilla_flavor_package.go").Parse(strings.TrimSpace(tmpl.Text()))
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
}

type TemplateVanillaFlavorPackagePackage struct {
	Packages []string
	// Discriptors []TemplateColumnPackageDiscriptor
}

func (TemplateVanillaFlavorPackagePackage) Text() string {
	return `
{{ range $index, $pkg := .Packages -}} 
package {{ $pkg }}
{{ end }}
`
}

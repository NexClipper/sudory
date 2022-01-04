package model

type Catalogues struct {
	Items []*Catalogue
}

func (m *Catalogues) GetType() string {
	return "CATALOGUES"
}

type Catalogue struct {
	Name string
}

func (m *Catalogue) GetType() string {
	return "CATALOGUE"
}

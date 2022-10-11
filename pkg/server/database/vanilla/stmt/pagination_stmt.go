package stmt

import (
	"github.com/pkg/errors"
)

type PaginationStmt map[string]int

func (page PaginationStmt) Build(engine PaginationBuildEngine) (PaginationResult, error) {
	unwraped := (map[string]int)(page)
	r, err := engine.Build(unwraped)
	if err != nil {
		return nil, errors.Wrapf(err, "build pagination statement")
	}

	return r, nil
}

func (page PaginationStmt) Limit() (int, bool) {
	v, ok := (page)["limit"]
	if !ok {
		return 0, ok
	}

	return v, ok
}

func (page PaginationStmt) Page() (int, bool) {
	v, ok := (page)["page"]
	if !ok {
		return 0, ok
	}

	return v, ok
}

func (page PaginationStmt) SetLimit(l int) (PaginationStmt, error) {
	if l < 0 {
		return page, errors.New("limit is greater then equal zero")
	}
	(page)["limit"] = l
	return page, nil
}

func (page PaginationStmt) SetPage(p int) (PaginationStmt, error) {
	if p <= 0 {
		return page, errors.New("page is greater then equal zero")
	}

	page["page"] = p
	return page, nil
}

func Limit(limit int, page ...int) PaginationStmt {
	option := func() int {
		if len(page) == 0 {
			return 1
		}
		if !(0 < page[0]) {
			return 1
		}

		return page[0]
	}

	m := map[string]int{
		"limit": limit,
		"page":  option(),
	}
	return m
}

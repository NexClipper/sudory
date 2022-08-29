package stmt

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var PaginationLexer = new(paginationLexer)

type paginationLexer struct{}

func (lexer paginationLexer) Parse(s string) (PaginationStmt, error) {
	s = strings.TrimSpace(s)

	if len(s) == 0 {
		return nil, errors.WithStack(ErrorInvalidArgumentEmptyString)
	}

	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return nil, errors.Wrapf(err, "json.Unmarshal q=%s", s)
	}

	return parsePagination(v)
}

func parsePagination(v interface{}) (PaginationStmt, error) {
	r := PaginationStmt{}

	switch value := v.(type) {
	case map[string]interface{}:
		for k, v := range value {

			f, err := scanJsonNumber(v)
			if err != nil {
				return nil, errors.Wrapf(err, "scan %T", v)
			}

			switch strings.ToLower(k) {
			case "limit":
				r, err = r.SetLimit((int)(f))
				if err != nil {
					return r, errors.Wrapf(err, "set limit l=%v", f)
				}
			case "page":
				r, err = r.SetPage((int)(f))
				if err != nil {
					return r, errors.Wrapf(err, "set page p=%v", f)
				}
			default:
				return nil, errors.Wrapf(ErrorUnsupportedPaginationKeys, "key=%v", k)
			}

		}
		return r, nil

	case []interface{}:
		for i := range value {
			v, err := parsePagination(value[i])
			if err != nil {
				return nil, errors.Wrapf(err, "parseOrderMap %T", value[i])
			}
			if limit, ok := v.Limit(); ok {
				r, err = r.SetLimit(limit)
				if err != nil {
					return r, errors.WithStack(err)
				}
			}

			if page, ok := v.Page(); ok {
				r, err = r.SetPage(page)
				if err != nil {
					return r, errors.WithStack(err)
				}
			}
		}
	default:
		return nil, errors.Wrapf(ErrorUnsupportedType, "value=%v type=%T", value, value)
	}

	return r, nil
}

func scanJsonNumber(v interface{}) (float64, error) {
	switch value := v.(type) {
	case float64:
		return value, nil
	case string:
		f, err := strconv.ParseFloat(value, 64)
		return f, errors.Wrapf(err, "strconv.ParseFloat s=%v", value)
	default:
		return 0, errors.Wrapf(ErrorUnsupportedType, "value=%v type=%T", value, value)
	}
}

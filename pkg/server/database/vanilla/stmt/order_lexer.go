package stmt

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

var OrderLexer = new(orderLexer)

type orderLexer struct{}

func (lexer orderLexer) Parse(s string) (OrderStmt, error) {
	s = strings.TrimSpace(s)

	if len(s) == 0 {
		return nil, errors.WithStack(ErrorInvalidArgumentEmptyString)
	}

	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return nil, errors.Wrapf(err, "json.Unmarshal q=%s", s)
	}

	return parseOrder(v)
}

func parseOrder(v interface{}) (OrderStmt, error) {

	r := make([]map[string][]string, 0)
	switch value := v.(type) {
	case map[string]interface{}:

		var key string
		var val []string
		var err error
		for k, v := range value {
			switch strings.ToLower(k) {
			case "asc":
			case "desc":
			default:
				return nil, errors.Wrapf(ErrorUnsupportedOrderKeys, "key=%v", k)
			}

			key = k
			val, err = scanStrings(v)
			if err != nil {
				return nil, errors.Wrapf(err, "scan %T", v)
			}
		}
		r = append(r, map[string][]string{key: val})

	case []interface{}:
		for i := range value {
			v, err := parseOrder(value[i])
			if err != nil {
				return nil, errors.Wrapf(err, "parseOrderMap %T", value[i])
			}
			r = append(r, v...)
		}
	default:
		return nil, errors.Wrapf(ErrorUnsupportedType, "value=%v type=%T", value, value)
	}

	return OrderStmt(r), nil
}

func scanStrings(v interface{}) ([]string, error) {
	switch value := v.(type) {
	case string:
		return []string{value}, nil
	case []string:
		return value, nil
	case []interface{}:
		r := make([]string, 0, len(value))

		for i := range value {
			ss, err := scanStrings(value[i])
			if err != nil {
				return nil, errors.Wrapf(err, "scanStrings value=%T", value[i])
			}
			r = append(r, ss...)
		}
		return r, nil
	default:
		return nil, errors.Wrapf(ErrorUnsupportedType, "value=%v type=%T", value, value)
	}
}

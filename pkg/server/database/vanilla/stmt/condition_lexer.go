package stmt

import (
	"encoding/json"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla/stmt/internal/sexp"
	"github.com/pkg/errors"
)

var ConditionLexer = new(conditionLexer)

type conditionLexer struct{}

func (lexer conditionLexer) Parse(s string) (ConditionStmt, error) {
	s = strings.TrimSpace(s)

	if len(s) == 0 {
		return nil, errors.WithStack(ErrorInvalidArgumentEmptyString)
	}

	var fn func() (ConditionStmt, error) = func() (ConditionStmt, error) { return nil, errors.New("there is no implementation") }

	if strings.Index(s, "(") == 0 {
		// s-exp
		fn = func() (ConditionStmt, error) {
			val, err := sexp.EvalString(s)
			if err != nil {
				return nil, errors.Wrapf(err, "sexp.EvalString exp=\"%v\"", s)
			}
			value := val.Val()
			v, ok := value.(map[string]interface{})
			if !ok {
				return nil, errors.Errorf("convert from=%T to=%T", value, v)
			}

			return ConditionStmt(v), err
		}
	} else {
		// default json
		fn = func() (ConditionStmt, error) {
			v := make(map[string]interface{})
			if err := json.Unmarshal([]byte(s), &v); err != nil {
				return nil, errors.Wrapf(err, "json.Unmarshal q=%s", s)
			}

			return ConditionStmt(v), nil
		}
	}

	return fn()
}

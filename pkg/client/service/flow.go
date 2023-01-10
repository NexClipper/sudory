package service

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/itchyny/gojq"
)

const (
	predefinedFlowExprStr_INPUTS = "$inputs"
)

func hasFlowInputsExpr(s string) bool {
	if len(s) <= 0 {
		return false
	}

	length := len(predefinedFlowExprStr_INPUTS)

	if strings.HasPrefix(s, predefinedFlowExprStr_INPUTS) {
		if len(s) > length {
			if s[length] != '.' {
				return false
			}
		}
		return true
	}
	return false
}

type FlowStepInput struct {
	initKey string
	data    map[string]interface{}
}

func findReplaceDeferredInput(key interface{}, specified map[string]interface{}) (bool, interface{}, error) {
	switch kt := key.(type) {
	case string:
		if hasFlowInputsExpr(kt) {
			path := strings.TrimPrefix(kt, predefinedFlowExprStr_INPUTS)
			query, err := gojq.Parse(path)
			if err != nil {
				return true, nil, err
			}

			iter := query.Run(specified)
			var res interface{}
			for {
				v, ok := iter.Next()
				if !ok {
					break
				}
				if err, ok := v.(error); ok {
					return true, nil, err
				}

				if v != nil {
					res = v
					break
				}
			}

			if res == nil {
				return true, nil, fmt.Errorf("not found key %q", kt)
			}

			return true, res, nil
		}
	case map[string]interface{}:
		for k, v := range kt {
			found, value, err := findReplaceDeferredInput(v, specified)
			if err != nil {
				return false, nil, err
			}
			if found {
				kt[k] = value
			}
		}
	case []interface{}:
		for i, v := range kt {
			found, value, err := findReplaceDeferredInput(v, specified)
			if err != nil {
				return false, nil, err
			}
			if found {
				kt[i] = value
			}
		}
	}

	return false, nil, nil
}

func (f *FlowStepInput) FindReplaceDeferredInputsFrom(in map[string]interface{}) error {
	if hasFlowInputsExpr(f.initKey) {
		path := strings.TrimPrefix(f.initKey, predefinedFlowExprStr_INPUTS)
		query, err := gojq.Parse(path)
		if err != nil {
			return err
		}

		iter := query.Run(in)
		var res interface{}
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}
			if err, ok := v.(error); ok {
				return err
			}

			if v != nil {
				res = v
				break
			}
		}
		m, ok := res.(map[string]interface{})
		if !ok {
			return fmt.Errorf("inputs(%s) must be json object", f.initKey)
		}
		f.data = m
		return nil
	}

	_, _, err := findReplaceDeferredInput(f.data, in)
	if err != nil {
		return err
	}

	return nil
}

func (f *FlowStepInput) GetInputs() map[string]interface{} {
	return f.data
}

func findReplacePassedInput(key interface{}, dataset map[string]interface{}) (bool, interface{}, error) {
	switch kt := key.(type) {
	case string:
		if strings.HasPrefix(kt, "$") {
			fullPath := strings.TrimPrefix(kt, "$")
			cnt := 0
			ind := strings.IndexFunc(fullPath, func(r rune) bool {
				if r == '.' {
					cnt++
					if cnt == 2 {
						return true
					}
				}
				return false
			})
			var stepInOutKey, keyPath string
			if ind < 0 {
				stepInOutKey = fullPath
			} else {
				stepInOutKey = fullPath[:ind]
				keyPath = fullPath[ind:]
			}

			val, ok := dataset[stepInOutKey]
			if !ok {
				return true, nil, fmt.Errorf("not found key %q", stepInOutKey)
			}

			var jv interface{}
			switch vv := val.(type) {
			case string:
				x := strings.TrimLeft(vv, " \t\r\n")
				var temp interface{}
				if len(x) > 0 {
					if x[0] == '{' {
						temp = &map[string]interface{}{}
					} else if x[0] == '[' {
						temp = &[]interface{}{}
					} else {
						// basic string
						return true, x, nil
					}
				}
				if err := json.Unmarshal([]byte(vv), temp); err != nil {
					return true, nil, err
				}
				jv = reflect.ValueOf(temp).Elem().Interface()
			default:
				jv = vv
			}

			query, err := gojq.Parse(keyPath)
			if err != nil {
				return true, nil, err
			}

			iter := query.Run(jv)
			var res interface{}
			for {
				v, ok := iter.Next()
				if !ok {
					break
				}
				if err, ok := v.(error); ok {
					return true, nil, err
				}

				if v != nil {
					res = v
					break
				}
			}

			if res == nil {
				return true, nil, fmt.Errorf("not found key %q", kt)
			}

			return true, res, nil
		}
	case map[string]interface{}:
		for k, v := range kt {
			found, value, err := findReplacePassedInput(v, dataset)
			if err != nil {
				return false, nil, err
			}
			if found {
				kt[k] = value
			}
		}
	case []interface{}:
		for i, v := range kt {
			found, value, err := findReplacePassedInput(v, dataset)
			if err != nil {
				return false, nil, err
			}
			if found {
				kt[i] = value
			}
		}
	}

	return false, nil, nil
}

func (f *FlowStepInput) FindReplacePassedInputsFrom(in map[string]interface{}) error {
	_, _, err := findReplacePassedInput(f.data, in)
	if err != nil {
		return err
	}
	return nil
}

type FlowStep struct {
	Id      string        `json:"$id"`
	Command string        `json:"$command"`
	Inputs  FlowStepInput `json:"inputs,omitempty"`
	Outputs interface{}   `json:"outputs,omitempty"`
	Error   error         `json:"error,omitempty"`
}

func (s *FlowStep) UnmarshalJSON(b []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	idInf, found := m["$id"]
	if !found {
		return fmt.Errorf("'$id' is required")
	} else {
		id, ok := idInf.(string)
		if !ok {
			return fmt.Errorf("'$id' type must be 'string'. not '%T'", idInf)
		}
		s.Id = id
	}

	cmdInf, found := m["$command"]
	if !found {
		return fmt.Errorf("'$command' is required")
	} else {
		cmd, ok := cmdInf.(string)
		if !ok {
			return fmt.Errorf("'$command' type must be 'string'. not '%T'", cmdInf)
		}
		s.Command = cmd
	}

	inputsInf, found := m["inputs"]
	if found {
		switch t := inputsInf.(type) {
		case string:
			if !hasFlowInputsExpr(t) {
				return fmt.Errorf("'string' type of 'inputs' must start with %q. not %q", predefinedFlowExprStr_INPUTS, t)
			}
			s.Inputs = FlowStepInput{initKey: t}
		case map[string]interface{}:
			s.Inputs = FlowStepInput{data: t}
		default:
			return fmt.Errorf("'inputs' type must be '$string' or 'map[string]interface{}'. not '%T'", t)
		}
	}

	return nil
}

type Flow []*FlowStep

func (s *Flow) UnmarshalJSON(b []byte) error {
	var l []json.RawMessage

	if err := json.Unmarshal(b, &l); err != nil {
		return err
	}

	ids := make(map[string]bool)

	for _, e := range l {
		step := new(FlowStep)
		if err := json.Unmarshal(e, step); err != nil {
			return err
		}

		if _, ok := ids[step.Id]; ok {
			return fmt.Errorf("step's $id(%q) already exists", step.Id)
		}
		ids[step.Id] = true
		*s = append(*s, step)
	}

	return nil
}

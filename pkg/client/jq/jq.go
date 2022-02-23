package jq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/itchyny/gojq"
)

func ProcessJson(input map[string]interface{}, filter string) (string, error) {
	if input == nil || filter == "" {
		return "", fmt.Errorf("input or filter value is empty")
	}

	query, err := gojq.Parse(filter)
	if err != nil {
		return "", err
	}

	res := &JqResult{Filter: filter}

	iter := query.RunWithContext(context.TODO(), input)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return "", err
		}

		if v != nil {
			res.Results = append(res.Results, v)
		}
	}

	b, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

type JqResult struct {
	Filter  string        `json:"filter"`
	Results []interface{} `json:"results"`
}

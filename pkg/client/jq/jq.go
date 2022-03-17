package jq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/itchyny/gojq"
)

func Request(input map[string]interface{}, filter string) (string, error) {
	res := &JqResult{Filter: filter}
	var err error

	res.Results, err = Process(input, filter)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(res)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

type JqResult struct {
	Filter  string      `json:"filter"`
	Results interface{} `json:"results"`
}

func Process(input map[string]interface{}, filter string) (interface{}, error) {
	if input == nil || filter == "" {
		return "", fmt.Errorf("input or filter value is empty")
	}

	query, err := gojq.Parse(filter)
	if err != nil {
		return "", err
	}

	var rootRes interface{}
	var res []interface{}

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
			res = append(res, v)
		}
	}

	if len(res) > 1 {
		rootRes = res
	} else if len(res) == 1 {
		rootRes = res[0]
	}

	return rootRes, nil
}

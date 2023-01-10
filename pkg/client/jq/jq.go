package jq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/itchyny/gojq"
)

const defaultJqTimeout = 10 * time.Second

func Request(input interface{}, filter string) (string, error) {
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

func Process(input interface{}, filter string) (interface{}, error) {
	query, err := gojq.Parse(filter)
	if err != nil {
		return "", err
	}

	var rootRes interface{}
	var res []interface{}

	ctx, cancel := context.WithTimeout(context.Background(), defaultJqTimeout)
	defer cancel()

	iter := query.RunWithContext(ctx, input)
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

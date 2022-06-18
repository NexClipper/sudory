package p8s

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/prometheus/common/model"
)

func (c *Client) Query(params map[string]interface{}) (string, []string, error) {
	type queryParams struct {
		Query string    `json:"query,omitempty"`
		Time  time.Time `json:"time,omitempty"`
	}

	qp := &queryParams{}
	if err := mapToStruct(params, qp); err != nil {
		return "", nil, err
	}

	urlValues := url.Values{}
	urlValues.Set("query", qp.Query)
	if !qp.Time.IsZero() {
		urlValues.Set("time", formatTime(qp.Time))
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	body, err := c.client.PostForm(ctx, "/api/v1/query", nil, []byte(urlValues.Encode()))
	if err != nil {
		return "", nil, err
	}

	var apiresp apiResponse
	if err := json.Unmarshal(body, &apiresp); err != nil {
		return "", nil, err
	}

	if apiresp.Status != "success" {
		return "", apiresp.Warnings, fmt.Errorf(apiresp.Error)
	}

	return string(apiresp.Data), nil, nil
}

func (c *Client) QueryRange(params map[string]interface{}) (string, []string, error) {
	type queryParams struct {
		Query string    `json:"query,omitempty"`
		Start time.Time `json:"start,omitempty"`
		End   time.Time `json:"end,omitempty"`
		Step  Duration  `json:"step,omitempty"`
	}

	qp := &queryParams{}

	if err := mapToStruct(params, qp); err != nil {
		return "", nil, err
	}

	urlValues := url.Values{}
	urlValues.Set("query", qp.Query)
	urlValues.Set("start", formatTime(qp.Start))
	urlValues.Set("end", formatTime(qp.End))
	urlValues.Set("step", strconv.FormatFloat(time.Duration(qp.Step).Seconds(), 'f', -1, 64))

	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	body, err := c.client.PostForm(ctx, "/api/v1/query_range", nil, []byte(urlValues.Encode()))
	if err != nil {
		return "", nil, err
	}

	var apiresp apiResponse
	if err := json.Unmarshal(body, &apiresp); err != nil {
		return "", nil, err
	}

	if apiresp.Status != "success" {
		return "", apiresp.Warnings, fmt.Errorf(apiresp.Error)
	}

	return string(apiresp.Data), nil, nil
}

type QueryResult struct {
	ResultType string      `json:"resultType"`
	Result     model.Value `json:"result"`
}

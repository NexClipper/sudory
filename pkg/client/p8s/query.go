package p8s

import (
	"context"
	"encoding/json"
	"time"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
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

	v1api := v1.NewAPI(c.client)
	ctx, cancel := context.WithTimeout(context.TODO(), defaultQueryTimeout)
	defer cancel()

	data, warnings, err := v1api.Query(ctx, qp.Query, qp.Time)
	if err != nil {
		return "", warnings, err
	}

	qr := &QueryResult{ResultType: data.Type().String(), Result: data}

	b, err := json.Marshal(qr)
	if err != nil {
		return "", warnings, err
	}

	return string(b), warnings, nil
	// return fmt.Sprintf(`{"resultType":"%s","result":%s}`, data.Type().String(), string(b)), warnings, nil
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

	r := v1.Range{
		Start: qp.Start,
		End:   qp.End,
		Step:  time.Duration(qp.Step),
	}

	v1api := v1.NewAPI(c.client)
	ctx, cancel := context.WithTimeout(context.TODO(), defaultQueryTimeout)
	defer cancel()

	data, warnings, err := v1api.QueryRange(ctx, qp.Query, r)
	if err != nil {
		return "", nil, err
	}

	qr := &QueryResult{ResultType: data.Type().String(), Result: data}

	b, err := json.Marshal(qr)
	if err != nil {
		return "", nil, err
	}

	return string(b), warnings, nil
	// return fmt.Sprintf(`{"resultType":"%s","result":%s}`, data.Type().String(), string(b)), warnings, nil
}

type QueryResult struct {
	ResultType string      `json:"resultType"`
	Result     model.Value `json:"result"`
}

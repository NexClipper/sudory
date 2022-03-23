package p8s

import (
	"context"
	"encoding/json"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

func (c *Client) Targets() (string, error) {
	v1api := v1.NewAPI(c.client)
	ctx, cancel := context.WithTimeout(context.TODO(), defaultQueryTimeout)
	defer cancel()

	data, err := v1api.Targets(ctx)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c *Client) TargetsMetadata(params map[string]interface{}) (string, error) {
	type targetsMetadataParams struct {
		MatchTarget string `json:"match_target,omitempty"`
		Metric      string `json:"metric,omitempty"`
		Limit       string `json:"limit,omitempty"`
	}

	tmParams := &targetsMetadataParams{}
	if err := mapToStruct(params, tmParams); err != nil {
		return "", err
	}

	v1api := v1.NewAPI(c.client)
	ctx, cancel := context.WithTimeout(context.TODO(), defaultQueryTimeout)
	defer cancel()

	data, err := v1api.TargetsMetadata(ctx, tmParams.MatchTarget, tmParams.Metric, tmParams.Limit)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

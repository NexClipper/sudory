package p8s

import (
	"context"
	"encoding/json"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

func (c *Client) Rules() (string, error) {
	v1api := v1.NewAPI(c.client)
	ctx, cancel := context.WithTimeout(context.TODO(), defaultQueryTimeout)
	defer cancel()

	data, err := v1api.Rules(ctx)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

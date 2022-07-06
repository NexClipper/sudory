package p8s

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) Targets() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	body, err := c.client.Get("/api/v1/targets").Do(ctx).Raw()
	if err != nil {
		return "", err
	}

	var apiresp apiResponse
	if err := json.Unmarshal(body, &apiresp); err != nil {
		return "", err
	}

	if apiresp.Status != "success" {
		return "", fmt.Errorf(apiresp.Error)
	}

	return string(apiresp.Data), nil
}

func (c *Client) TargetsMetadata(params map[string]interface{}) (string, error) {
	m := make(map[string][]string)

	for k, v := range params {
		str, ok := v.(string)
		if !ok {
			return "", fmt.Errorf("params['%s']'s type must be string, not %T", k, v)
		}
		m[k] = []string{str}
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	body, err := c.client.Get("/api/v1/targets/metadata").Do(ctx).Raw()
	if err != nil {
		return "", err
	}

	var apiresp apiResponse
	if err := json.Unmarshal(body, &apiresp); err != nil {
		return "", err
	}

	if apiresp.Status != "success" {
		return "", fmt.Errorf(apiresp.Error)
	}

	return string(apiresp.Data), nil
}

package p8s

import (
	"context"
	"encoding/json"
	"fmt"
)

func (c *Client) Rules() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()

	body, err := c.client.Get("/api/v1/rules").Do(ctx).Raw()
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

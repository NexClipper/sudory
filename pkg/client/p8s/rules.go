package p8s

import (
	"context"
	"encoding/json"
	"fmt"

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
	for _, group := range data.Groups {
		for i, r := range group.Rules {
			switch v := r.(type) {
			case v1.RecordingRule:
				group.Rules[i] = RecordingRule{
					Type:          string(v1.RuleTypeRecording),
					RecordingRule: v,
				}
			case v1.AlertingRule:
				group.Rules[i] = AlertingRule{
					Type:         string(v1.RuleTypeAlerting),
					AlertingRule: v,
				}
			default:
				fmt.Printf("unknown rule type %s", v)
			}
		}
	}

	b, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// AlertingRule models a alerting rule.
type AlertingRule struct {
	Type            string `json:"type"`
	v1.AlertingRule `json:",inline"`
}

// RecordingRule models a recording rule.
type RecordingRule struct {
	Type             string `json:"type"`
	v1.RecordingRule `json:",inline"`
}

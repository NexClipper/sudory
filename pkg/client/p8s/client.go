package p8s

import (
	"fmt"

	"github.com/prometheus/client_golang/api"

	"github.com/NexClipper/sudory/pkg/client/log"
)

type Client struct {
	client api.Client
}

func NewClient(url string) (*Client, error) {
	client, err := api.NewClient(api.Config{Address: url})
	if err != nil {
		return nil, err
	}

	return &Client{client: client}, nil
}

func (c *Client) ApiRequest(apiVersion, apiName string, queryParams map[string]interface{}) (string, error) {
	var data string
	var warnings []string
	var err error

	switch apiVersion {
	case "v1":
		switch apiName {
		case "query":
			data, warnings, err = c.Query(queryParams)
		case "query_range":
			data, warnings, err = c.QueryRange(queryParams)
		case "alerts":
			data, err = c.Alerts()
		case "rules":
			data, err = c.Rules()
		case "alertmanagers":
			data, err = c.AlertManagers()
		default:
			return "", fmt.Errorf("unknown api name(%s)", apiName)
		}
	default:
		return "", fmt.Errorf("unknown api version(%s)", apiVersion)
	}

	if len(warnings) > 0 {
		log.Warnf("Prometheus API(%s) Warnings : %v\n", apiName, warnings)
	}

	return data, err
}

package alertmanager

import (
	"fmt"
	"time"

	"github.com/NexClipper/sudory/pkg/client/httpclient"
	"github.com/NexClipper/sudory/pkg/client/log"
)

const defaultQueryTimeout = 10 * time.Second

type Client struct {
	client *httpclient.HttpClient
}

func NewClient(url string) (*Client, error) {
	return &Client{client: httpclient.NewHttpClient(url, "", 0, 0)}, nil
}

func (c *Client) ApiRequest(apiVersion, apiName, verb string, params map[string]interface{}) (string, error) {
	var data string
	var warnings []string
	var err error

	switch apiVersion {
	case "v2":
		apiPath := "/api/v2"
		switch apiName {
		case "silences":
			switch verb {
			case "get":
				data, err = c.GetSilence(apiPath, params)
			case "list":
				data, err = c.GetSilences(apiPath, params)
			case "create":
				data, err = c.CreateSilences(apiPath, params)
			case "delete":
				data, err = c.DeleteSilence(apiPath, params)
			case "update":
				data, err = c.UpdateSilence(apiPath, params)
			default:
				return "", fmt.Errorf("unknown verb name(%s)", verb)
			}
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

package grafana

import (
	"fmt"
	"time"

	"github.com/NexClipper/sudory/pkg/client/httpclient"
)

const defaultApiTimeout = 10 * time.Second

type Client struct {
	client      *httpclient.HttpClient
	getApiKeyFn func() ([]byte, error)
}

func NewClient(url string, getApiKeyFn func() ([]byte, error)) (*Client, error) {
	client, err := httpclient.NewHttpClient(url, false, 0, 0)
	if err != nil {
		return nil, err
	}
	client.SetDisableKeepAlives()
	return &Client{client: client, getApiKeyFn: getApiKeyFn}, nil
}

func (c *Client) ApiRequest(apiName, verb string, params map[string]interface{}) (string, error) {
	var data string
	var err error

	switch apiName {
	case "datasources":
		apiPath := "/api/datasources"
		switch verb {
		case "get":
			data, err = c.GetDatasource(apiPath, params)
		case "list":
			data, err = c.ListDatasources(apiPath)
		case "delete":
			data, err = c.DeleteDatasource(apiPath, params)
		default:
			return "", fmt.Errorf("unknown verb name(%s)", verb)
		}
	default:
		return "", fmt.Errorf("unknown api name(%s)", apiName)
	}

	return data, err
}

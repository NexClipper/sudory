package openstack

import (
	"fmt"
	"time"

	"github.com/NexClipper/sudory/pkg/client/httpclient"
)

const (
	defaultApiTimeout    = 10 * time.Second
	xAuthTokenHeaderName = "X-AUTH-TOKEN"
)

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

func (c *Client) ApiRequest(api, resource, verb string, params map[string]interface{}) (string, error) {
	var data string
	var err error

	switch api {
	case "identity":
		switch resource {
		case "projects":
			switch verb {
			case "get":
				data, err = c.GetIdentityV3Project(api, params)
			case "list":
				data, err = c.ListIdentityV3Projects(api, params)
			default:
				return "", fmt.Errorf("unknown verb name(%s)", verb)
			}
		default:
			return "", fmt.Errorf("unknown resource name(%s)", resource)
		}
	case "compute":
		switch resource {
		case "servers":
			switch verb {
			case "get":
				data, err = c.GetComputeV2_1Server(api, params)
			case "list":
				data, err = c.ListComputeV2_1Servers(api, params)
			default:
				return "", fmt.Errorf("unknown verb name(%s)", verb)
			}
		case "hypervisors":
			switch verb {
			case "get":
				data, err = c.GetComputeV2_1Hypervisors(api, params)
			case "list":
				data, err = c.ListComputeV2_1Hypervisors(api, params)
			default:
				return "", fmt.Errorf("unknown verb name(%s)", verb)
			}
		default:
			return "", fmt.Errorf("unknown resource name(%s)", resource)
		}
	case "networking":
		switch resource {
		case "networks":
			switch verb {
			case "get":
				data, err = c.GetNetworkingV2_0Network(api, params)
			case "list":
				data, err = c.ListNetworkingV2_0Networks(api, params)
			default:
				return "", fmt.Errorf("unknown verb name(%s)", verb)
			}
		case "routers":
			switch verb {
			case "get":
				data, err = c.GetNetworkingV2_0Router(api, params)
			case "list":
				data, err = c.ListNetworkingV2_0Routers(api, params)
			default:
				return "", fmt.Errorf("unknown verb name(%s)", verb)
			}
		case "subnets":
			switch verb {
			case "get":
				data, err = c.GetNetworkingV2_0Subnet(api, params)
			case "list":
				data, err = c.ListNetworkingV2_0Subnets(api, params)
			default:
				return "", fmt.Errorf("unknown verb name(%s)", verb)
			}
		default:
			return "", fmt.Errorf("unknown resource name(%s)", resource)
		}
	default:
		return "", fmt.Errorf("unknown api name(%s)", api)
	}

	return data, err
}

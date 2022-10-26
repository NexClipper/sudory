package openstack

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
)

type Client struct {
	pClient *gophercloud.ProviderClient
}

func NewClient(pClient *gophercloud.ProviderClient) *Client {
	return &Client{pClient: pClient}
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
				data, err = c.GetIdentityV3Project(params)
			case "list":
				data, err = c.ListIdentityV3Projects(params)
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
				data, err = c.GetComputeV2_1Server(params)
			case "list":
				data, err = c.ListComputeV2_1Servers(params)
			default:
				return "", fmt.Errorf("unknown verb name(%s)", verb)
			}
		case "hypervisors":
			switch verb {
			case "get":
				data, err = c.GetComputeV2_1Hypervisors(params)
			case "list":
				data, err = c.ListComputeV2_1Hypervisors(params)
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
				data, err = c.GetNetworkingV2_0Network(params)
			case "list":
				data, err = c.ListNetworkingV2_0Networks(params)
			default:
				return "", fmt.Errorf("unknown verb name(%s)", verb)
			}
		case "routers":
			switch verb {
			case "get":
				data, err = c.GetNetworkingV2_0Router(params)
			case "list":
				data, err = c.ListNetworkingV2_0Routers(params)
			default:
				return "", fmt.Errorf("unknown verb name(%s)", verb)
			}
		case "subnets":
			switch verb {
			case "get":
				data, err = c.GetNetworkingV2_0Subnet(params)
			case "list":
				data, err = c.ListNetworkingV2_0Subnets(params)
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

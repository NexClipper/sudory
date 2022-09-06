package sudoryclient

import (
	"fmt"
	"time"

	"github.com/NexClipper/sudory/pkg/client/k8s"
)

const defaultTimeout = 10 * time.Second

type Client struct {
	k8sClient *k8s.Client
}

func NewClient() (*Client, error) {
	k8sClient, err := k8s.GetClient()
	if err != nil {
		return nil, err
	}

	return &Client{k8sClient: k8sClient}, nil
}

func (c *Client) Request(api, verb string, args map[string]interface{}) (string, error) {
	var result string
	var err error

	switch api {
	case "credential":
		result, err = c.Credential(verb, args)
		if err != nil {
			break
		}
	default:
		err = fmt.Errorf("unknown api(%s)", verb)
	}

	if err != nil {
		return "", err
	}

	return result, nil
}

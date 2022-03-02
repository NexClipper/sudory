package helm

import (
	"fmt"

	"helm.sh/helm/v3/pkg/cli"
)

type Client struct {
	settings *cli.EnvSettings
}

func NewClient() (*Client, error) {
	settings := cli.New()

	return &Client{settings: settings}, nil
}

func (c *Client) Request(cmd string, args map[string]interface{}) (string, error) {
	var result string
	var err error

	switch cmd {
	case "repo_add":
		result, err = c.RepoAdd(args)
	default:
		return "", fmt.Errorf("unknown command(%s)", cmd)
	}

	return result, err
}

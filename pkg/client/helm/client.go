package helm

import (
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"

	"github.com/NexClipper/sudory/pkg/client/log"
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
	case "install":
		result, err = c.Install(args)
	case "uninstall":
		result, err = c.Uninstall(args)
	case "upgrade":
		result, err = c.Upgrade(args)
	case "get_values":
		result, err = c.GetValues(args)
	case "repo_add":
		result, err = c.RepoAdd(args)
	case "repo_list":
		result, err = c.RepoList(args)
	case "repo_update":
		result, err = c.RepoUpdate(args)
	default:
		return "", fmt.Errorf("unknown command(%s)", cmd)
	}

	return result, err
}

func (c *Client) getActionConfig() (*action.Configuration, error) {
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(c.settings.RESTClientGetter(), c.settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Debugf); err != nil {
		return nil, err
	}

	return actionConfig, nil
}

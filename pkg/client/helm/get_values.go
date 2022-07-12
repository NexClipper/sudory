package helm

import (
	"encoding/json"

	"helm.sh/helm/v3/pkg/action"
)

func (c *Client) GetValues(args map[string]interface{}) (string, error) {
	type GetValuesParams struct {
		Namespace string `param:"namespace"`
		Name      string `param:"name"`
		All       bool   `param:"all,optional"`
	}

	params := &GetValuesParams{}

	if err := convertArgsToStruct(args, params); err != nil {
		return "", err
	}
	// set namespace
	c.settings.SetNamespace(params.Namespace)

	// get 'GetValues' action client
	actionConfig, err := c.getActionConfig()
	if err != nil {
		return "", err
	}
	client := action.NewGetValues(actionConfig)

	// set all option
	client.AllValues = params.All

	m, err := client.Run(params.Name)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

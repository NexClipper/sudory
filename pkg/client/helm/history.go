package helm

import (
	"encoding/json"
	"fmt"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/releaseutil"
)

func (c *Client) History(args map[string]interface{}) (string, error) {
	type HistoryParams struct {
		Namespace string `param:"namespace"`
		Name      string `param:"name"`
	}

	params := &HistoryParams{}

	if err := convertArgsToStruct(args, params); err != nil {
		return "", err
	}

	// set namespace
	c.settings.SetNamespace(params.Namespace)

	// get history action client
	actionConfig, err := c.getActionConfig()
	if err != nil {
		return "", err
	}
	client := action.NewHistory(actionConfig)

	rels, err := client.Run(params.Name)
	if err != nil {
		return "", err
	}

	releaseutil.Reverse(rels, releaseutil.SortByRevision)

	b, err := transformHistoryResultToJson(rels)
	if err != nil {
		return fmt.Sprintf("chart(%s) history is success, but failed to transform result to json : %s", params.Name, err.Error()), nil
	}

	return string(b), nil
}

func transformHistoryResultToJson(rels []*release.Release) ([]byte, error) {
	m, err := extractHistoryResultFrom(rels)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&m)
}

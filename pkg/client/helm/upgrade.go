package helm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"

	"github.com/NexClipper/sudory/pkg/client/log"
)

func setDefaultUpgradeSettings(client *action.Upgrade) {
	client.Timeout = 300 * time.Second
	client.Wait = true
	client.WaitForJobs = true
	client.Atomic = true

	// allow deletion of new resources created in this upgrade when upgrade fails
	client.CleanupOnFail = true
}

func (c *Client) Upgrade(args map[string]interface{}) (string, error) {
	type UpgradeParams struct {
		Namespace    string                 `param:"namespace"`
		Name         string                 `param:"name"`
		ChartName    string                 `param:"chart_name"`
		RepoURL      string                 `param:"repo_url"`
		ChartVersion string                 `param:"chart_version,optional"`
		Values       map[string]interface{} `param:"values,optional"`
	}

	params := &UpgradeParams{}

	if err := convertArgsToStruct(args, params); err != nil {
		return "", err
	}

	// set namespace
	c.settings.SetNamespace(params.Namespace)

	// get upgrade action client
	actionConfig, err := c.getActionConfig()
	if err != nil {
		return "", err
	}
	client := action.NewUpgrade(actionConfig)

	// default settings
	setDefaultUpgradeSettings(client)

	client.ChartPathOptions.Version = params.ChartVersion
	client.ChartPathOptions.RepoURL = params.RepoURL
	client.Namespace = c.settings.Namespace()

	// look for chart directory
	chartPath, err := client.ChartPathOptions.LocateChart(params.ChartName, c.settings)
	if err != nil {
		return "", err
	}

	// load chart
	chartLoaded, err := loader.Load(chartPath)
	if err != nil {
		return "", err
	}
	if req := chartLoaded.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(chartLoaded, req); err != nil {
			return "", err
		}
	}

	if chartLoaded.Metadata.Deprecated {
		log.Warnf("chart(%s) is deprecated", params.ChartName)
	}

	rel, err := client.RunWithContext(context.TODO(), params.Name, chartLoaded, params.Values)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(rel)
	if err != nil {
		return fmt.Sprintf("chart(%s) upgrade is success, but failed to json.Marhsal : %s", params.Name, err.Error()), nil
	}

	return string(b), nil
}

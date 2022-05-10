package helm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/release"

	"github.com/NexClipper/sudory/pkg/client/log"
)

func setDefaultInstallSettings(client *action.Install) {
	client.CreateNamespace = true
	client.Timeout = 300 * time.Second
	client.Wait = true
	client.WaitForJobs = true
	client.Atomic = true
}

func (c *Client) Install(args map[string]interface{}) (string, error) {
	type InstallParams struct {
		Namespace    string                 `param:"namespace"`
		Name         string                 `param:"name"`
		ChartName    string                 `param:"chart_name"`
		RepoURL      string                 `param:"repo_url"`
		ChartVersion string                 `param:"chart_version,optional"`
		Values       map[string]interface{} `param:"values,optional"`
	}

	params := &InstallParams{}

	if err := convertArgsToStruct(args, params); err != nil {
		return "", err
	}

	// set namespace
	c.settings.SetNamespace(params.Namespace)

	// get install action client
	actionConfig, err := c.getActionConfig()
	if err != nil {
		return "", err
	}
	client := action.NewInstall(actionConfig)

	// default settings
	setDefaultInstallSettings(client)

	// client.Description =
	client.ChartPathOptions.Version = params.ChartVersion
	client.ChartPathOptions.RepoURL = params.RepoURL

	client.ReleaseName = params.Name

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

	// chart's type("" or "application") is only installable
	if chartLoaded.Metadata.Type != "" && chartLoaded.Metadata.Type != "application" {
		return "", fmt.Errorf("chart's type(%s) are not installable", chartLoaded.Metadata.Type)
	}

	if chartLoaded.Metadata.Deprecated {
		log.Warnf("chart(%s) is deprecated", params.ChartName)
	}

	if reqs := chartLoaded.Metadata.Dependencies; reqs != nil {
		if err := action.CheckDependencies(chartLoaded, reqs); err != nil {
			return "", err
		}
	}

	client.Namespace = c.settings.Namespace()

	rel, err := client.RunWithContext(context.TODO(), chartLoaded, params.Values)
	if err != nil {
		return "", err
	}

	b, err := transformInstallResultToJson(rel)
	if err != nil {
		return fmt.Sprintf("chart(%s) install is success, but failed to transform result to json : %s", params.Name, err.Error()), nil
	}

	return string(b), nil
}

func transformInstallResultToJson(rel *release.Release) ([]byte, error) {
	m, err := extractResultFrom(rel)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&m)
}

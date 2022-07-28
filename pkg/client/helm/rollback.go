package helm

import (
	"fmt"
	"time"

	"helm.sh/helm/v3/pkg/action"
)

func setDefaultRollbackSettings(client *action.Rollback) {
	client.Timeout = 300 * time.Second
	client.Wait = true
	client.WaitForJobs = true
	client.MaxHistory = defaultMaxHistory
	
	// allow deletion of new resources created in this upgrade when upgrade fails
	client.CleanupOnFail = true
}

func (c *Client) Rollback(args map[string]interface{}) (string, error) {
	type InstallParams struct {
		Namespace string  `param:"namespace"`
		Name      string  `param:"name"`
		Revision  float64 `param:"revision,optional"` // encoding/json/decode.go:53
	}

	params := &InstallParams{}

	if err := convertArgsToStruct(args, params); err != nil {
		return "", err
	}

	// set namespace
	c.settings.SetNamespace(params.Namespace)

	// get rollback action client
	actionConfig, err := c.getActionConfig()
	if err != nil {
		return "", err
	}
	client := action.NewRollback(actionConfig)

	// default settings
	setDefaultRollbackSettings(client)

	// set revision
	if params.Revision > 0 {
		client.Version = int(params.Revision)
	}

	if err := client.Run(params.Name); err != nil {
		return "", err
	}

	if params.Revision <= 0 {
		return fmt.Sprintf("successfully rolled back to the previous release(%s)", params.Name), nil
	}

	return fmt.Sprintf("successfully rolled back release(%s) to revision %d", params.Name, int(params.Revision)), nil
}

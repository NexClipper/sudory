package grafana

import (
	"context"
	"fmt"
)

func (c *Client) GetDatasource(apiPath string, params map[string]interface{}) (string, error) {
	var id, uid, name string

	if found, err := FindCastFromMap(params, "datasource_id", &id); found && err != nil {
		return "", err
	}
	if found, err := FindCastFromMap(params, "datasource_uid", &uid); found && err != nil {
		return "", err
	}
	if found, err := FindCastFromMap(params, "datasource_name", &name); found && err != nil {
		return "", err
	}

	path := ""
	if id != "" {
		path = "/" + id
	} else if uid != "" {
		path = "/uid/" + uid
	} else if name != "" {
		path = "/name/" + name
	} else {
		return "", fmt.Errorf("one of datasource_id, datasource_uid, datasource_name must have")
	}

	apikey, err := c.getApiKeyFn()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultApiTimeout)
	defer cancel()

	body, err := c.client.Get(apiPath+path).
		SetHeader("Authorization", string(apikey)).
		SetHeader("Accept", "application/json").
		Do(ctx).Raw()
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) ListDatasources(apiPath string) (string, error) {
	apikey, err := c.getApiKeyFn()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultApiTimeout)
	defer cancel()

	body, err := c.client.Get(apiPath).
		SetHeader("Authorization", string(apikey)).
		SetHeader("Accept", "application/json").
		Do(ctx).Raw()
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) DeleteDatasource(apiPath string, params map[string]interface{}) (string, error) {
	var id, uid, name string

	if found, err := FindCastFromMap(params, "datasource_id", &id); found && err != nil {
		return "", err
	}
	if found, err := FindCastFromMap(params, "datasource_uid", &uid); found && err != nil {
		return "", err
	}
	if found, err := FindCastFromMap(params, "datasource_name", &name); found && err != nil {
		return "", err
	}

	path := ""
	if id != "" {
		path = "/" + id
	} else if uid != "" {
		path = "/uid/" + uid
	} else if name != "" {
		path = "/name/" + name
	} else {
		return "", fmt.Errorf("one of datasource_id, datasource_uid, datasource_name must have")
	}

	apikey, err := c.getApiKeyFn()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultApiTimeout)
	defer cancel()

	body, err := c.client.Delete(apiPath+path).
		SetHeader("Authorization", string(apikey)).
		SetHeader("Accept", "application/json").
		Do(ctx).Raw()
	if err != nil {
		return "", err
	}

	return string(body), nil
}

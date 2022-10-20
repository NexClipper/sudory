package openstack

import (
	"context"
	"fmt"
)

const (
	computeApiV2_1BasePath = "/compute/v2.1"
	computeApiServersPath  = "/servers"
)

func (c *Client) GetComputeV2_1Server(api string, params map[string]interface{}) (string, error) {
	var path = computeApiV2_1BasePath + computeApiServersPath
	var id string
	var query = make(map[string]interface{})

	if found, err := FindCastFromMap(params, "id", &id); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if id == "" {
		return "", fmt.Errorf("server_id is empty")
	}

	if found, err := FindCastFromMap(params, "query", &query); found && err != nil {
		return "", err
	}

	q := convertQuery(query)

	path += "/" + id

	apikey, err := c.getApiKeyFn()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultApiTimeout)
	defer cancel()

	body, err := c.client.Get(path).
		SetHeader(xAuthTokenHeaderName, string(apikey)).
		SetHeader("Accept", "application/json").
		SetParamFromQuery(q).
		Do(ctx).Raw()
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) ListComputeV2_1Servers(apiPath string, params map[string]interface{}) (string, error) {
	var path = computeApiV2_1BasePath + computeApiServersPath
	var query = make(map[string]interface{})

	if found, err := FindCastFromMap(params, "query", &query); found && err != nil {
		return "", err
	}

	q := convertQuery(query)

	apikey, err := c.getApiKeyFn()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultApiTimeout)
	defer cancel()

	body, err := c.client.Get(path).
		SetHeader(xAuthTokenHeaderName, string(apikey)).
		SetHeader("Accept", "application/json").
		SetParamFromQuery(q).
		Do(ctx).Raw()
	if err != nil {
		return "", err
	}

	return string(body), nil
}

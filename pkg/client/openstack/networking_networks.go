package openstack

import (
	"context"
	"fmt"
)

const (
	networkingApiV2_0BasePath = "/v2.0"
	networkingApiNetworksPath = "/networks"
)

func (c *Client) GetNetworkingV2_0Network(api string, params map[string]interface{}) (string, error) {
	var path = networkingApiV2_0BasePath + networkingApiNetworksPath
	var id string
	var query = make(map[string]interface{})

	if found, err := FindCastFromMap(params, "id", &id); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if id == "" {
		return "", fmt.Errorf("network_id is empty")
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

func (c *Client) ListNetworkingV2_0Networks(apiPath string, params map[string]interface{}) (string, error) {
	var path = networkingApiV2_0BasePath + networkingApiNetworksPath
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

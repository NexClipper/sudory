package openstack

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
)

func (c *Client) GetComputeV2_1Hypervisors(params map[string]interface{}) (string, error) {
	var id string

	if found, err := FindCastFromMap(params, "id", &id); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if id == "" {
		return "", fmt.Errorf("server_id is empty")
	}

	client, err := openstack.NewComputeV2(c.pClient, gophercloud.EndpointOpts{})
	if err != nil {
		return "", err
	}

	r := hypervisors.Get(client, id)
	if r.Err != nil {
		return "", r.Err
	}

	b, err := json.Marshal(r.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c *Client) ListComputeV2_1Hypervisors(params map[string]interface{}) (string, error) {
	var query = make(map[string]interface{})

	if found, err := FindCastFromMap(params, "query", &query); found && err != nil {
		return "", err
	}

	lo, err := convertQueryToHypervisorsListOpts(query)
	if err != nil {
		return "", err
	}

	client, err := openstack.NewComputeV2(c.pClient, gophercloud.EndpointOpts{})
	if err != nil {
		return "", err
	}

	allPages, err := hypervisors.List(client, lo).AllPages()
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(allPages.GetBody())
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func convertQueryToHypervisorsListOpts(query map[string]interface{}) (hypervisors.ListOpts, error) {
	if len(query) <= 0 {
		return hypervisors.ListOpts{}, nil
	}

	lo := hypervisors.ListOpts{}

	for k, v := range query {
		switch k {
		case "limit":
			i, err := strconv.ParseInt(fmt.Sprintf("%s", v), 10, 64)
			if err != nil {
				return lo, err
			}
			ii := int(i)
			lo.Limit = &ii
		case "marker":
			s := fmt.Sprintf("%s", v)
			lo.Marker = &s
		case "hypervisor_hostname_pattern":
			s := fmt.Sprintf("%s", v)
			lo.HypervisorHostnamePattern = &s
		case "with_servers":
			b, err := strconv.ParseBool(fmt.Sprintf("%s", v))
			if err != nil {
				return lo, err
			}
			lo.WithServers = &b
		}
	}

	return lo, nil
}

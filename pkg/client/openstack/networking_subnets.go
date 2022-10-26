package openstack

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

func (c *Client) GetNetworkingV2_0Subnet(params map[string]interface{}) (string, error) {
	var id string

	if found, err := FindCastFromMap(params, "id", &id); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if id == "" {
		return "", fmt.Errorf("subnet_id is empty")
	}

	client, err := openstack.NewNetworkV2(c.pClient, gophercloud.EndpointOpts{})
	if err != nil {
		return "", err
	}

	r := subnets.Get(client, id)
	if r.Err != nil {
		return "", r.Err
	}

	b, err := json.Marshal(r.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c *Client) ListNetworkingV2_0Subnets(params map[string]interface{}) (string, error) {
	var query = make(map[string]interface{})

	if found, err := FindCastFromMap(params, "query", &query); found && err != nil {
		return "", err
	}

	lo, err := convertQueryToSubnetsListOpts(query)
	if err != nil {
		return "", err
	}

	client, err := openstack.NewNetworkV2(c.pClient, gophercloud.EndpointOpts{})
	if err != nil {
		return "", err
	}

	allPages, err := subnets.List(client, lo).AllPages()
	if err != nil {
		return "", err
	}

	allSubnets, err := subnets.ExtractSubnets(allPages)
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(allSubnets)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func convertQueryToSubnetsListOpts(query map[string]interface{}) (subnets.ListOpts, error) {
	if len(query) <= 0 {
		return subnets.ListOpts{}, nil
	}

	lo := subnets.ListOpts{}

	for k, v := range query {
		switch k {
		case "name":
			lo.Name = fmt.Sprintf("%s", v)
		case "description":
			lo.Description = fmt.Sprintf("%s", v)
		case "enable_dhcp":
			b, err := strconv.ParseBool(fmt.Sprintf("%s", v))
			if err != nil {
				return lo, err
			}
			lo.EnableDHCP = &b
		case "network_id":
			lo.NetworkID = fmt.Sprintf("%s", v)
		case "tenant_id":
			lo.TenantID = fmt.Sprintf("%s", v)
		case "project_id":
			lo.ProjectID = fmt.Sprintf("%s", v)
		case "ip_version":
			i, err := strconv.ParseInt(fmt.Sprintf("%s", v), 10, 64)
			if err != nil {
				return lo, err
			}
			lo.IPVersion = int(i)
		case "gateway_ip":
			lo.GatewayIP = fmt.Sprintf("%s", v)
		case "cidr":
			lo.CIDR = fmt.Sprintf("%s", v)
		case "ipv6_address_mode":
			lo.IPv6AddressMode = fmt.Sprintf("%s", v)
		case "ipv6_ra_mode":
			lo.IPv6RAMode = fmt.Sprintf("%s", v)
		case "id":
			lo.ID = fmt.Sprintf("%s", v)
		case "subnetpool_id":
			lo.SubnetPoolID = fmt.Sprintf("%s", v)
		case "limit":
			i, err := strconv.ParseInt(fmt.Sprintf("%s", v), 10, 64)
			if err != nil {
				return lo, err
			}
			lo.Limit = int(i)
		case "marker":
			lo.Marker = fmt.Sprintf("%s", v)
		case "sort_key":
			lo.SortKey = fmt.Sprintf("%s", v)
		case "sort_dir":
			lo.SortDir = fmt.Sprintf("%s", v)
		case "tags":
			lo.Tags = fmt.Sprintf("%s", v)
		case "tags-any":
			lo.TagsAny = fmt.Sprintf("%s", v)
		case "not-tags":
			lo.NotTags = fmt.Sprintf("%s", v)
		case "not-tags-any":
			lo.NotTagsAny = fmt.Sprintf("%s", v)
		}
	}

	return lo, nil
}

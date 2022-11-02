package openstack

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func (c *Client) GetComputeV2Server(params map[string]interface{}) (string, error) {
	var id string
	var microversion string

	if found, err := FindCastFromMap(params, "id", &id); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if id == "" {
		return "", fmt.Errorf("server_id is empty")
	}

	if found, err := FindCastFromMap(params, "microversion", &microversion); found && err != nil {
		return "", err
	}

	client, err := openstack.NewComputeV2(c.pClient, gophercloud.EndpointOpts{})
	if err != nil {
		return "", err
	}

	if microversion != "" {
		client.Microversion = microversion
	}

	r := servers.Get(client, id)
	if r.Err != nil {
		return "", r.Err
	}

	b, err := json.Marshal(r.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c *Client) ListComputeV2Servers(params map[string]interface{}) (string, error) {
	var query = make(map[string]interface{})
	var microversion string

	if found, err := FindCastFromMap(params, "query", &query); found && err != nil {
		return "", err
	}

	if found, err := FindCastFromMap(params, "microversion", &microversion); found && err != nil {
		return "", err
	}

	lo, err := convertQueryToServersListOpts(query)
	if err != nil {
		return "", err
	}

	client, err := openstack.NewComputeV2(c.pClient, gophercloud.EndpointOpts{})
	if err != nil {
		return "", err
	}

	if microversion != "" {
		client.Microversion = microversion
	}

	allPages, err := servers.List(client, lo).AllPages()
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(allPages.GetBody())
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func convertQueryToServersListOpts(query map[string]interface{}) (servers.ListOpts, error) {
	if len(query) <= 0 {
		return servers.ListOpts{}, nil
	}

	lo := servers.ListOpts{}

	for k, v := range query {
		switch k {
		case "changes-since":
			lo.ChangesSince = fmt.Sprintf("%s", v)
		case "image":
			lo.Image = fmt.Sprintf("%s", v)
		case "flavor":
			lo.Flavor = fmt.Sprintf("%s", v)
		case "ip":
			lo.IP = fmt.Sprintf("%s", v)
		case "ip6":
			lo.IP6 = fmt.Sprintf("%s", v)
		case "name":
			lo.Name = fmt.Sprintf("%s", v)
		case "status":
			lo.Status = fmt.Sprintf("%s", v)
		case "host":
			lo.Host = fmt.Sprintf("%s", v)
		case "marker":
			lo.Marker = fmt.Sprintf("%s", v)
		case "limit":
			i, err := strconv.ParseInt(fmt.Sprintf("%s", v), 10, 64)
			if err != nil {
				return lo, err
			}
			lo.Limit = int(i)
		case "all_tenants":
			b, err := strconv.ParseBool(fmt.Sprintf("%s", v))
			if err != nil {
				return lo, err
			}
			lo.AllTenants = b
		case "tenant_id":
			lo.TenantID = fmt.Sprintf("%s", v)
		case "user_id":
			lo.UserID = fmt.Sprintf("%s", v)
		case "tags":
			lo.Tags = fmt.Sprintf("%s", v)
		case "tags-any":
			lo.TagsAny = fmt.Sprintf("%s", v)
		case "not-tags":
			lo.NotTags = fmt.Sprintf("%s", v)
		case "not-tags-any":
			lo.NotTagsAny = fmt.Sprintf("%s", v)
		case "availability_zone":
			lo.AvailabilityZone = fmt.Sprintf("%s", v)
		}
	}

	return lo, nil
}

package openstack

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
)

func (c *Client) GetNetworkingV2_0Router(params map[string]interface{}) (string, error) {
	var id string

	if found, err := FindCastFromMap(params, "id", &id); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if id == "" {
		return "", fmt.Errorf("router_id is empty")
	}

	client, err := openstack.NewNetworkV2(c.pClient, gophercloud.EndpointOpts{})
	if err != nil {
		return "", err
	}

	r := routers.Get(client, id)
	if r.Err != nil {
		return "", r.Err
	}

	b, err := json.Marshal(r.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c *Client) ListNetworkingV2_0Routers(params map[string]interface{}) (string, error) {
	var query = make(map[string]interface{})

	if found, err := FindCastFromMap(params, "query", &query); found && err != nil {
		return "", err
	}

	lo, err := convertQueryToRoutersListOpts(query)
	if err != nil {
		return "", err
	}

	client, err := openstack.NewNetworkV2(c.pClient, gophercloud.EndpointOpts{})
	if err != nil {
		return "", err
	}

	allPages, err := routers.List(client, lo).AllPages()
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(allPages.GetBody())
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func convertQueryToRoutersListOpts(query map[string]interface{}) (routers.ListOpts, error) {
	if len(query) <= 0 {
		return routers.ListOpts{}, nil
	}

	lo := routers.ListOpts{}

	for k, v := range query {
		switch k {
		case "id":
			lo.ID = fmt.Sprintf("%s", v)
		case "name":
			lo.Name = fmt.Sprintf("%s", v)
		case "description":
			lo.Description = fmt.Sprintf("%s", v)
		case "admin_state_up":
			b, err := strconv.ParseBool(fmt.Sprintf("%s", v))
			if err != nil {
				return lo, err
			}
			lo.AdminStateUp = &b
		case "distributed":
			b, err := strconv.ParseBool(fmt.Sprintf("%s", v))
			if err != nil {
				return lo, err
			}
			lo.Distributed = &b
		case "status":
			lo.Status = fmt.Sprintf("%s", v)
		case "tenant_id":
			lo.TenantID = fmt.Sprintf("%s", v)
		case "project_id":
			lo.ProjectID = fmt.Sprintf("%s", v)
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

package openstack

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

func (c *Client) GetIdentityV3Project(params map[string]interface{}) (string, error) {
	var id string

	if found, err := FindCastFromMap(params, "id", &id); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if id == "" {
		return "", fmt.Errorf("project_id is empty")
	}

	client, err := openstack.NewIdentityV3(c.pClient, gophercloud.EndpointOpts{})
	if err != nil {
		return "", err
	}

	r := projects.Get(client, id)
	if r.Err != nil {
		return "", r.Err
	}

	b, err := json.Marshal(r.Body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (c *Client) ListIdentityV3Projects(params map[string]interface{}) (string, error) {
	var query = make(map[string]interface{})

	if found, err := FindCastFromMap(params, "query", &query); found && err != nil {
		return "", err
	}

	lo, err := convertQueryToProjectsListOpts(query)
	if err != nil {
		return "", err
	}

	client, err := openstack.NewIdentityV3(c.pClient, gophercloud.EndpointOpts{})
	if err != nil {
		return "", err
	}

	allPages, err := projects.List(client, lo).AllPages()
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(allPages.GetBody())
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func convertQueryToProjectsListOpts(query map[string]interface{}) (projects.ListOpts, error) {
	if len(query) <= 0 {
		return projects.ListOpts{}, nil
	}

	lo := projects.ListOpts{}

	for k, v := range query {
		switch k {
		case "domain_id":
			lo.DomainID = fmt.Sprintf("%s", v)
		case "enabled":
			b, err := strconv.ParseBool(fmt.Sprintf("%s", v))
			if err != nil {
				return lo, err
			}
			lo.Enabled = &b
		case "is_domain":
			b, err := strconv.ParseBool(fmt.Sprintf("%s", v))
			if err != nil {
				return lo, err
			}
			lo.IsDomain = &b
		case "name":
			lo.Name = fmt.Sprintf("%s", v)
		case "parent_id":
			lo.ParentID = fmt.Sprintf("%s", v)
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

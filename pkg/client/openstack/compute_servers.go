package openstack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/diskconfig"
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

func (c *Client) RebootComputeV2Server(params map[string]interface{}) (string, error) {
	var id string
	var rebootType string
	var microversion string

	if found, err := FindCastFromMap(params, "id", &id); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if id == "" {
		return "", fmt.Errorf("server_id is empty")
	}

	if found, err := FindCastFromMap(params, "reboot_type", &rebootType); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if rebootType == "" {
		return "", fmt.Errorf("type is empty")
	}

	if found, err := FindCastFromMap(params, "microversion", &microversion); found && err != nil {
		return "", err
	}

	client, err := openstack.NewComputeV2(c.pClient, gophercloud.EndpointOpts{})
	if err != nil {
		return "", err
	}

	var rebootMethod servers.RebootMethod
	if strings.EqualFold(rebootType, string(servers.HardReboot)) {
		rebootMethod = servers.HardReboot
	} else if strings.EqualFold(rebootType, string(servers.SoftReboot)) {
		rebootMethod = servers.SoftReboot
	} else {
		return "", fmt.Errorf("reboot_type must be 'HARD' or 'SOFT', not %s", rebootType)
	}

	if microversion != "" {
		client.Microversion = microversion
	}

	r := servers.Reboot(client, id, servers.RebootOpts{Type: rebootMethod})
	if r.Err != nil {
		return "", r.Err
	}

	return fmt.Sprintf("successfully requested to reboot a server(%s)", id), nil
}

func (c *Client) ResizeComputeV2Server(params map[string]interface{}) (string, error) {
	var id string
	var flavorRef string
	var diskConfig string
	var microversion string

	if found, err := FindCastFromMap(params, "id", &id); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if id == "" {
		return "", fmt.Errorf("server_id is empty")
	}

	if found, err := FindCastFromMap(params, "resize_flavorRef", &flavorRef); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	} else if flavorRef == "" {
		return "", fmt.Errorf("resize_flavorRef is empty")
	}

	if found, err := FindCastFromMap(params, "resize_diskConfig", &diskConfig); found && err != nil {
		return "", err
	}

	if found, err := FindCastFromMap(params, "microversion", &microversion); found && err != nil {
		return "", err
	}

	client, err := openstack.NewComputeV2(c.pClient, gophercloud.EndpointOpts{})
	if err != nil {
		return "", err
	}

	var opts servers.ResizeOptsBuilder = servers.ResizeOpts{FlavorRef: flavorRef}

	if diskConfig != "" {
		if strings.EqualFold(diskConfig, string(diskconfig.Auto)) {
			diskConfig = string(diskconfig.Auto)
		} else if strings.EqualFold(diskConfig, string(diskconfig.Manual)) {
			diskConfig = string(diskconfig.Manual)
		} else {
			return "", fmt.Errorf("resize_diskConfig must be 'AUTO' or 'MANUAL', not %s", diskConfig)
		}
		opts = diskconfig.ResizeOptsExt{
			ResizeOptsBuilder: opts,
			DiskConfig:        diskconfig.DiskConfig(diskConfig),
		}
	}

	if microversion != "" {
		client.Microversion = microversion
	}

	r := servers.Resize(client, id, opts)
	if r.Err != nil {
		return "", r.Err
	}

	return fmt.Sprintf("successfully requested to resize a server(%s)", id), nil
}

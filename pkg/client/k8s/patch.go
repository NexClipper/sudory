package k8s

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func (c *Client) ResourcePatch(gv schema.GroupVersion, resource string, params map[string]interface{}) (string, error) {
	var result interface{}
	var err error

	var namespace string
	var name string
	var ptype string

	if found, err := FindCastFromMap(params, "namespace", &namespace); found && err != nil {
		return "", err
	}

	if found, err := FindCastFromMap(params, "name", &name); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	}

	if found, err := FindCastFromMap(params, "patch_type", &ptype); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	}

	var pt types.PatchType
	if ptype == "json" {
		pt = types.JSONPatchType
	} else if ptype == "merge" {
		pt = types.MergePatchType
	} else {
		return "", fmt.Errorf("unsupported patch type(%s). must be one of [json merge]", ptype)
	}

	data, err := json.Marshal(params["patch_data"])
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultK8sTimeout)
	defer cancel()

	switch gv.Identifier() {
	case "v1":
		switch resource {
		case "configmaps":
			result, err = c.client.CoreV1().ConfigMaps(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		case "persistentvolumes":
			result, err = c.client.CoreV1().PersistentVolumes().Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		case "persistentvolumeclaims":
			result, err = c.client.CoreV1().PersistentVolumeClaims(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		case "services":
			result, err = c.client.CoreV1().Services(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "apps/v1":
		switch resource {
		case "deployments":
			result, err = c.client.AppsV1().Deployments(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	default:
		err = fmt.Errorf("unsupported group version(%s)", gv.Identifier())
	}

	if err != nil {
		return "", err
	}

	b, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

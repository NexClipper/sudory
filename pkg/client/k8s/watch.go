package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

func (c *Client) ResourceWatch(gv schema.GroupVersion, resource string, params map[string]interface{}) (watch.Interface, error) {
	var watchInf watch.Interface
	var err error

	var namespace string
	var name string

	if found, err := FindCastFromMap(params, "namespace", &namespace); found && err != nil {
		return nil, err
	}

	if found, err := FindCastFromMap(params, "name", &name); found && err != nil {
		return nil, err
	} else if !found {
		return nil, err
	}

	switch gv.Identifier() {
	case "apps/v1":
		switch resource {
		case "deployments":
			watchInf, err = c.client.AppsV1().Deployments(namespace).Watch(context.TODO(), metav1.ListOptions{FieldSelector: fields.OneTermEqualSelector("metadata.name", name).String()})
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
		return nil, err
	}

	if watchInf == nil {
		return nil, fmt.Errorf("watch result is nil")
	}

	return watchInf, nil
}

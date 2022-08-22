package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *Client) ResourceDelete(gv schema.GroupVersion, resource string, params map[string]interface{}) error {
	var err error

	var namespace string
	var name string

	if found, err := FindCastFromMap(params, "namespace", &namespace); found && err != nil {
		return err
	}

	if found, err := FindCastFromMap(params, "name", &name); found && err != nil {
		return err
	} else if !found {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultK8sTimeout)
	defer cancel()

	switch gv.Identifier() {
	case "v1":
		switch resource {
		case "pods":
			err = c.client.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
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
		return err
	}

	return nil
}

package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *Client) ResourceDelete(gv schema.GroupVersion, resource, namespace, name string) error {
	var err error

	switch gv.Identifier() {
	case "v1":
		switch resource {
		case "pods":
			err = c.client.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
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

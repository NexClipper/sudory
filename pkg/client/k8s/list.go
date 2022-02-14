package k8s

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *Client) ResourceList(gv schema.GroupVersion, resource, namespace string, labels map[string]string) (string, error) {
	var result interface{}
	var err error

	labelsString, err := convertMapToLabelSelector(labels)
	if err != nil {
		return "", err
	}

	switch gv.Identifier() {
	case "v1":
		switch resource {
		case "pod":
			result, err = c.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "namespace":
			result, err = c.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("unknown resource(%s)", resource)
		}
	default:
		err = fmt.Errorf("unknown group version(%s)", gv.Identifier())
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

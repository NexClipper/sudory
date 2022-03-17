package k8s

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *Client) ResourceGet(gv schema.GroupVersion, resource, namespace, name string) (string, error) {
	var result interface{}
	var err error

	switch gv.Identifier() {
	case "v1":
		switch resource {
		case "configmaps":
			result, err = c.client.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		case "endpoints":
			result, err = c.client.CoreV1().Endpoints(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		case "events":
			result, err = c.client.CoreV1().Events(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		case "namespaces":
			result, err = c.client.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		case "nodes":
			result, err = c.client.CoreV1().Nodes().Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		case "persistentvolumes":
			result, err = c.client.CoreV1().PersistentVolumes().Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		case "persistentvolumeclaims":
			result, err = c.client.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		case "pods":
			result, err = c.client.CoreV1().Pods(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		case "secrets":
			result, err = c.client.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		case "services":
			result, err = c.client.CoreV1().Services(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "apps/v1":
		switch resource {
		case "daemonsets":
			result, err = c.client.AppsV1().DaemonSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		case "deployments":
			result, err = c.client.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		case "replicasets":
			result, err = c.client.AppsV1().ReplicaSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		case "statefulsets":
			result, err = c.client.AppsV1().StatefulSets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "networking.k8s.io/v1":
		switch resource {
		case "ingresses":
			result, err = c.client.NetworkingV1().Ingresses(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "storage.k8s.io/v1":
		switch resource {
		case "storageclasses":
			result, err = c.client.StorageV1().StorageClasses().Get(context.TODO(), name, metav1.GetOptions{})
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

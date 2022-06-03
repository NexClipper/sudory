package k8s

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *Client) ResourceGet(gv schema.GroupVersion, resource string, params map[string]interface{}) (string, error) {
	var result interface{}
	var err error

	var namespace string
	var name string

	if found, err := FindCastFromMap(params, "namespace", &namespace); found && err != nil {
		return "", err
	}

	if found, err := FindCastFromMap(params, "name", &name); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	}

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
	case "monitoring.coreos.com/v1":
		switch resource {
		case "prometheuses":
			result, err = c.mclient.MonitoringV1().Prometheuses(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		case "prometheusrules":
			result, err = c.mclient.MonitoringV1().PrometheusRules(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		case "servicemonitors":
			result, err = c.mclient.MonitoringV1().ServiceMonitors(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "batch/v1":
		switch resource {
		case "cronjobs":
			result, err = c.client.BatchV1().CronJobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				break
			}
		case "jobs":
			result, err = c.client.BatchV1().Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
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

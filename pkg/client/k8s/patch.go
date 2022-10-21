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
		case "endpoints":
			result, err = c.client.CoreV1().Endpoints(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		case "namespaces":
			result, err = c.client.CoreV1().Namespaces().Patch(ctx, name, pt, data, metav1.PatchOptions{})
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
		case "secrets":
			result, err = c.client.CoreV1().Secrets(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
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
		case "daemonsets":
			result, err = c.client.AppsV1().DaemonSets(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		case "deployments":
			result, err = c.client.AppsV1().Deployments(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		case "statefulsets":
			result, err = c.client.AppsV1().StatefulSets(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "networking.k8s.io/v1":
		switch resource {
		case "ingresses":
			result, err = c.client.NetworkingV1().Ingresses(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "monitoring.coreos.com/v1":
		switch resource {
		case "prometheuses":
			result, err = c.mclient.MonitoringV1().Prometheuses(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		case "prometheusrules":
			result, err = c.mclient.MonitoringV1().PrometheusRules(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		case "servicemonitors":
			result, err = c.mclient.MonitoringV1().ServiceMonitors(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		case "alertmanagers":
			result, err = c.mclient.MonitoringV1().Alertmanagers(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		case "podmonitors":
			result, err = c.mclient.MonitoringV1().PodMonitors(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		case "probes":
			result, err = c.mclient.MonitoringV1().Probes(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		case "thanosrulers":
			result, err = c.mclient.MonitoringV1().ThanosRulers(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "monitoring.coreos.com/v1alpha1":
		switch resource {
		case "alertmanagerconfigs":
			result, err = c.mclient.MonitoringV1alpha1().AlertmanagerConfigs(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "batch/v1":
		switch resource {
		case "cronjobs":
			result, err = c.client.BatchV1().CronJobs(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "autoscaling/v2":
		switch resource {
		case "horizontalpodautoscalers":
			result, err = c.client.AutoscalingV2().HorizontalPodAutoscalers(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "rbac.authorization.k8s.io/v1":
		switch resource {
		case "clusterrolebindings":
			result, err = c.client.RbacV1().ClusterRoleBindings().Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		case "clusterroles":
			result, err = c.client.RbacV1().ClusterRoles().Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		case "rolebindings":
			result, err = c.client.RbacV1().RoleBindings(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
			if err != nil {
				break
			}
		case "roles":
			result, err = c.client.RbacV1().Roles(namespace).Patch(ctx, name, pt, data, metav1.PatchOptions{})
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

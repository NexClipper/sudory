package k8s

import (
	"context"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func (c *Client) ResourceList(gv schema.GroupVersion, resource string, params map[string]interface{}) (string, error) {
	var result interface{}
	var err error

	var namespace string
	var labels = make(map[string]interface{})

	if found, err := FindCastFromMap(params, "namespace", &namespace); found && err != nil {
		return "", err
	}

	if found, err := FindCastFromMap(params, "labels", &labels); found && err != nil {
		return "", err
	}

	labelsString, err := convertMapToLabelSelector(labels)
	if err != nil {
		return "", err
	}

	switch gv.Identifier() {
	case "v1":
		switch resource {
		case "configmaps":
			result, err = c.client.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "endpoints":
			result, err = c.client.CoreV1().Endpoints(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "events":
			result, err = c.client.CoreV1().Events(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "namespaces":
			result, err = c.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "nodes":
			result, err = c.client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "persistentvolumes":
			result, err = c.client.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "persistentvolumeclaims":
			result, err = c.client.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "pods":
			result, err = c.client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "secrets":
			result, err = c.client.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "services":
			result, err = c.client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "limitranges":
			result, err = c.client.CoreV1().LimitRanges(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "podtemplates":
			result, err = c.client.CoreV1().PodTemplates(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "replicationcontrollers":
			result, err = c.client.CoreV1().ReplicationControllers(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "resourcequotas":
			result, err = c.client.CoreV1().ResourceQuotas(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "serviceaccounts":
			result, err = c.client.CoreV1().ServiceAccounts(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "apps/v1":
		switch resource {
		case "daemonsets":
			result, err = c.client.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "deployments":
			result, err = c.client.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "replicasets":
			result, err = c.client.AppsV1().ReplicaSets(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "statefulsets":
			result, err = c.client.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "controllerrevisions":
			result, err = c.client.AppsV1().ControllerRevisions(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "networking.k8s.io/v1":
		switch resource {
		case "ingresses":
			result, err = c.client.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "ingressclasses":
			result, err = c.client.NetworkingV1().IngressClasses().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "networkpolicies":
			result, err = c.client.NetworkingV1().NetworkPolicies(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "storage.k8s.io/v1":
		switch resource {
		case "storageclasses":
			result, err = c.client.StorageV1().StorageClasses().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "csidrivers":
			result, err = c.client.StorageV1().CSIDrivers().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "csinodes":
			result, err = c.client.StorageV1().CSINodes().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "csistoragecapacities":
			result, err = c.client.StorageV1().CSIStorageCapacities(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "volumeattachments":
			result, err = c.client.StorageV1().VolumeAttachments().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "monitoring.coreos.com/v1":
		switch resource {
		case "prometheuses":
			result, err = c.mclient.MonitoringV1().Prometheuses(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "prometheusrules":
			result, err = c.mclient.MonitoringV1().PrometheusRules(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "servicemonitors":
			result, err = c.mclient.MonitoringV1().ServiceMonitors(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "alertmanagers":
			result, err = c.mclient.MonitoringV1().Alertmanagers(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "podmonitors":
			result, err = c.mclient.MonitoringV1().PodMonitors(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "probes":
			result, err = c.mclient.MonitoringV1().Probes(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "thanosrulers":
			result, err = c.mclient.MonitoringV1().ThanosRulers(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "monitoring.coreos.com/v1alpha1":
		switch resource {
		case "alertmanagerconfigs":
			result, err = c.mclient.MonitoringV1alpha1().AlertmanagerConfigs(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "batch/v1":
		switch resource {
		case "cronjobs":
			result, err = c.client.BatchV1().CronJobs(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "jobs":
			result, err = c.client.BatchV1().Jobs(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "admissionregistration.k8s.io/v1":
		switch resource {
		case "mutatingwebhookconfigurations":
			result, err = c.client.AdmissionregistrationV1().MutatingWebhookConfigurations().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "validatingwebhookconfigurations":
			result, err = c.client.AdmissionregistrationV1().ValidatingWebhookConfigurations().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "apiextensions.k8s.io/v1":
		switch resource {
		case "customresourcedefinitions":
			result, err = c.apiextv1client.ApiextensionsV1().CustomResourceDefinitions().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "apiregistration.k8s.io/v1":
		switch resource {
		case "apiservices":
			result, err = c.aggrev1client.ApiregistrationV1().APIServices().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "autoscaling/v2":
		switch resource {
		case "horizontalpodautoscalers":
			result, err = c.client.AutoscalingV2().HorizontalPodAutoscalers(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "certificates.k8s.io/v1":
		switch resource {
		case "certificatesigningrequests":
			result, err = c.client.CertificatesV1().CertificateSigningRequests().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "coordination.k8s.io/v1":
		switch resource {
		case "leases":
			result, err = c.client.CoordinationV1().Leases(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "discovery.k8s.io/v1":
		switch resource {
		case "endpointslices":
			result, err = c.client.DiscoveryV1().EndpointSlices(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "node.k8s.io/v1":
		switch resource {
		case "runtimeclasses":
			result, err = c.client.NodeV1().RuntimeClasses().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "policy/v1":
		switch resource {
		case "poddisruptionbudgets":
			result, err = c.client.PolicyV1().PodDisruptionBudgets(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "rbac.authorization.k8s.io/v1":
		switch resource {
		case "clusterrolebindings":
			result, err = c.client.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "clusterroles":
			result, err = c.client.RbacV1().ClusterRoles().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "rolebindings":
			result, err = c.client.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		case "roles":
			result, err = c.client.RbacV1().Roles(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
			if err != nil {
				break
			}
		default:
			err = fmt.Errorf("group version(%s)'s unsupported resource(%s)", gv.Identifier(), resource)
		}
	case "scheduling.k8s.io/v1":
		switch resource {
		case "priorityclasses":
			result, err = c.client.SchedulingV1().PriorityClasses().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
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

package fetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	watchtools "k8s.io/client-go/tools/watch"

	"github.com/NexClipper/sudory/pkg/client/k8s"
	"github.com/NexClipper/sudory/pkg/client/log"
	"github.com/NexClipper/sudory/pkg/client/service"
)

func (f *Fetcher) UpgradeClient(version service.Version, serviceId string, args map[string]interface{}) (err error) {
	log.Debugf("SudoryClient Upgrade: start")

	t := time.Now()

	// service processing status update
	for {
		up := service.CreateUpdateService(version, serviceId, 1, 0, service.StepStatusProcessing, "", t, time.Time{})
		if err := f.sudoryAPI.UpdateServices(context.Background(), service.ConvertServiceStepUpdateClientToServer(up)); err != nil {
			log.Errorf("SudoryClient Upgrade: failed to update service status(processing): error: %s\n", err.Error())

			// retry and handshake if session expired
			if f.sudoryAPI.IsTokenExpired() {
				f.RetryHandshake()
			}
			continue
		}
		break
	}

	// fetcher polling stop
	f.ticker.Stop()
	defer func() {
		if err != nil {
			log.Errorf("SudoryClient Upgrade: failed to upgrade: %v\n", err)

			for {
				up := service.CreateUpdateService(version, serviceId, 1, 0, service.StepStatusFail, "", t, time.Now())
				if err := f.sudoryAPI.UpdateServices(context.Background(), service.ConvertServiceStepUpdateClientToServer(up)); err != nil {
					log.Errorf(err.Error())

					// retry and handshake if session expired
					if f.sudoryAPI.IsTokenExpired() {
						f.RetryHandshake()
					}
					continue
				}
				break
			}

			// fetcher polling restart because upgrade is failed
			f.ticker.Reset(time.Second * time.Duration(f.pollingInterval))
		}
	}()
	log.Debugf("SudoryClient Upgrade: stop polling")

	// check arguments
	var imageTag string
	var timeout time.Duration
	if args != nil {
		imageTagInf, ok := args["image_tag"]
		if !ok || imageTagInf == nil {
			return fmt.Errorf("failed to find image_tag argument")
		}

		imageTag, ok = imageTagInf.(string)
		if !ok {
			return fmt.Errorf("failed type assertion for image_tag argument: interface to string")
		}

		if imageTag == "" {
			return fmt.Errorf("image_tag argument is empty")
		}

		timeoutInf, ok := args["timeout"]
		if ok && timeoutInf != nil {
			switch v := timeoutInf.(type) {
			case float64: // encoding/json/decode.go:53
				timeout = time.Second * time.Duration(v)
			case string:
				timeoutInt, err := strconv.Atoi(v)
				if err != nil {
					return fmt.Errorf("failed convert string to int for timeout argument")
				}
				timeout = time.Second * time.Duration(timeoutInt)
			default:
				return fmt.Errorf("unsupported timeout argument type(%T): supported type(string, int)", v)
			}
		}
	} else {
		return fmt.Errorf("argument is empty")
	}

	// clean up the remaining services before upgrade(timeout:30s)
	log.Debugf("SudoryClient Upgrade: waiting remain service proccess")
	for cnt := 0; cnt < 10; cnt++ {
		<-time.After(time.Second * 3)
		remainServices := f.RemainServices()
		if len(remainServices) == 0 {
			break
		}

		buf := bytes.Buffer{}
		buf.WriteString("remain services:")
		for uuid, status := range remainServices {
			buf.WriteString(fmt.Sprintf("\n\tuuid: %s, status: %d", uuid, status))
		}
		log.Infof(buf.String() + "\n")
	}
	log.Debugf("SudoryClient Upgrade: end remain service proccess")

	// get namespace
	namespace, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return err
	}
	ns := string(namespace)

	// get pod_name
	podName, err := os.Hostname()
	if err != nil {
		return err
	}
	if podName == "" {
		return fmt.Errorf("pod name is empty")
	}

	log.Debugf("SudoryClient Upgrade: found pod info: {namespace: %s, pod_name: %s}\n", ns, podName)

	// get k8s client
	k8sClient, err := k8s.GetClient()
	if err != nil {
		return err
	}

	// get self-pod -> owner replicaset -> owner deployment
	deploymentObj, err := findDeploymentFromPod(k8sClient, ns, podName)
	if err != nil {
		return err
	}

	prevImage := deploymentObj.Spec.Template.Spec.Containers[0].Image
	upgradeImage := replaceImageTag(prevImage, imageTag)

	// patch deployment's image to upgrade image
	log.Debugf("SudoryClient Upgrade: request to patch deployment's image with %s", upgradeImage)
	patchedObj, err := k8sClient.ResourcePatch(schema.GroupVersion{Group: "apps", Version: "v1"}, "deployments", map[string]interface{}{
		"namespace":  ns,
		"name":       deploymentObj.Name,
		"patch_type": "json",
		"patch_data": []map[string]interface{}{
			{
				"op":    "replace",
				"path":  "/spec/template/spec/containers/0/image",
				"value": upgradeImage,
			},
		},
	})
	if err != nil {
		return err
	}

	origM, patchedM := make(map[string]interface{}), make(map[string]interface{})

	deploymentJson, _ := json.Marshal(deploymentObj)

	json.Unmarshal(deploymentJson, &origM)
	json.Unmarshal([]byte(patchedObj), &patchedM)

	if reflect.DeepEqual(origM, patchedM) {
		return fmt.Errorf("no changes to patch deployment")
	}
	log.Debugf("SudoryClient Upgrade: patched deployment")

	// watch my deployment
	watchInf, err := k8sClient.ResourceWatch(schema.GroupVersion{Group: "apps", Version: "v1"}, "deployments", map[string]interface{}{"namespace": ns, "name": deploymentObj.Name})
	if err != nil {
		return err
	}

	var watchCtx context.Context
	var watchCancel context.CancelFunc
	if timeout == 0 {
		watchCtx, watchCancel = context.WithCancel(context.Background())
	} else {
		watchCtx, watchCancel = context.WithTimeout(context.Background(), timeout)
	}
	defer watchCancel()

	log.Debugf("SudoryClient Upgrade: watch the status of the rollout until it's done")
	if _, watchErr := watchtools.UntilWithoutRetry(watchCtx, watchInf, func(e watch.Event) (bool, error) {
		switch t := e.Type; t {
		case watch.Added, watch.Modified:
			deploymentObj := e.Object.(*appsv1.Deployment)

			if deploymentObj.Generation <= deploymentObj.Status.ObservedGeneration {
				var cond *appsv1.DeploymentCondition
				for i := range deploymentObj.Status.Conditions {
					c := deploymentObj.Status.Conditions[i]
					if c.Type == appsv1.DeploymentProgressing {
						cond = &c
					}
				}

				if cond != nil && cond.Reason == "ProgressDeadlineExceeded" {
					return false, fmt.Errorf("deployment %q exceeded its progress deadline", deploymentObj.Name)
				} else if deploymentObj.Spec.Replicas != nil && deploymentObj.Status.UpdatedReplicas < *deploymentObj.Spec.Replicas {
					log.Debugf("SudoryClient Upgrade: waiting for deployment %q rollout to finish: %d out of %d new replicas have been updated\n", deploymentObj.Name, deploymentObj.Status.UpdatedReplicas, *deploymentObj.Spec.Replicas)
					return false, nil
				} else if deploymentObj.Status.Replicas > deploymentObj.Status.UpdatedReplicas {
					log.Debugf("SudoryClient Upgrade: waiting for deployment %q rollout to finish: %d old replicas are pending termination\n", deploymentObj.Name, deploymentObj.Status.Replicas-deploymentObj.Status.UpdatedReplicas)
					return false, nil
				} else if deploymentObj.Status.AvailableReplicas < deploymentObj.Status.UpdatedReplicas {
					log.Debugf("SudoryClient Upgrade: waiting for deployment %q rollout to finish: %d of %d updated replicas are available\n", deploymentObj.Name, deploymentObj.Status.AvailableReplicas, deploymentObj.Status.UpdatedReplicas)
					return false, nil
				} else {
					log.Debugf("SudoryClient Upgrade: deployment %q successfully rolled out\n", deploymentObj.Name)
					return true, nil
				}
			}
			log.Debugf("SudoryClient Upgrade: waiting for deployment spec update to be observed\n")

			return false, nil

		case watch.Deleted:
			return true, fmt.Errorf("object has been deleted")

		default:
			return true, fmt.Errorf("internal error: unexpected event %#v", e)
		}
	}); watchErr != nil {
		// patch my deployment's prev-image
		if _, err := k8sClient.ResourcePatch(schema.GroupVersion{Group: "apps", Version: "v1"}, "deployments", map[string]interface{}{
			"namespace":  ns,
			"name":       deploymentObj.Name,
			"patch_type": "json",
			"patch_data": []map[string]interface{}{
				{
					"op":    "replace",
					"path":  "/spec/template/spec/containers/0/image",
					"value": prevImage,
				},
			},
		}); err != nil {
			return fmt.Errorf("failed to patch deployment prev-image : watch_error: {%v}, patch_error: {%v}", watchErr, err)
		}

		return watchErr
	}

	// upgrade success
	for {
		up := service.CreateUpdateService(version, serviceId, 1, 0, service.StepStatusSuccess, "", t, time.Now())
		if err := f.sudoryAPI.UpdateServices(context.Background(), service.ConvertServiceStepUpdateClientToServer(up)); err != nil {
			log.Errorf(err.Error())
			// retry and handshake if session expired
			if f.sudoryAPI.IsTokenExpired() {
				f.RetryHandshake()
			}
			continue
		}
		break
	}

	return nil
}

func findDeploymentFromPod(k8sClient *k8s.Client, ns, podName string) (*appsv1.Deployment, error) {
	podJson, err := k8sClient.ResourceGet(schema.GroupVersion{Version: "v1"}, "pods", map[string]interface{}{"namespace": ns, "name": podName})
	if err != nil {
		return nil, err
	}
	podObj := new(corev1.Pod)
	if err := json.Unmarshal([]byte(podJson), podObj); err != nil {
		return nil, err
	}

	// find owner replicaset
	var replicasetName string
	for _, ownerRef := range podObj.OwnerReferences {
		if ownerRef.Kind == "ReplicaSet" {
			replicasetName = ownerRef.Name
		}
	}
	if replicasetName == "" {
		return nil, fmt.Errorf("failed to find replicaset name")
	}

	log.Debugf("SudoryClient Upgrade: found owner replicaset info: {namespace: %s, name: %s}\n", ns, replicasetName)

	replicasetJson, err := k8sClient.ResourceGet(schema.GroupVersion{Group: "apps", Version: "v1"}, "replicasets", map[string]interface{}{"namespace": ns, "name": replicasetName})
	if err != nil {
		return nil, err
	}
	replicasetObj := new(appsv1.ReplicaSet)
	if err := json.Unmarshal([]byte(replicasetJson), replicasetObj); err != nil {
		return nil, err
	}

	// find owner deployment
	var deploymentName string
	for _, ownerRef := range replicasetObj.OwnerReferences {
		if ownerRef.Kind == "Deployment" {
			deploymentName = ownerRef.Name
		}
	}
	if deploymentName == "" {
		return nil, fmt.Errorf("failed to find replicaset name")
	}

	log.Debugf("SudoryClient Upgrade: found owner deployment info: {namespace: %s, name: %s}\n", ns, deploymentName)

	deploymentJson, err := k8sClient.ResourceGet(schema.GroupVersion{Group: "apps", Version: "v1"}, "deployments", map[string]interface{}{"namespace": ns, "name": deploymentName})
	if err != nil {
		return nil, err
	}
	deploymentObj := new(appsv1.Deployment)
	if err := json.Unmarshal([]byte(deploymentJson), deploymentObj); err != nil {
		return nil, err
	}

	return deploymentObj, nil
}

func replaceImageTag(image, tag string) string {
	imageName := image

	index := strings.LastIndex(image, ":")
	if index != -1 {
		imageName = image[:index]
	}

	return fmt.Sprintf("%s:%s", imageName, tag)
}

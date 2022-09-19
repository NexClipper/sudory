package k8s

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/kubectl/pkg/scheme"
)

func searchContainerByName(pod *corev1.Pod, name string) *corev1.Container {
	for i := range pod.Spec.Containers {
		if pod.Spec.Containers[i].Name == name {
			return &pod.Spec.Containers[i]
		}
	}
	for i := range pod.Spec.InitContainers {
		if pod.Spec.InitContainers[i].Name == name {
			return &pod.Spec.InitContainers[i]
		}
	}
	for i := range pod.Spec.EphemeralContainers {
		if pod.Spec.EphemeralContainers[i].Name == name {
			return (*corev1.Container)(&pod.Spec.EphemeralContainers[i].EphemeralContainerCommon)
		}
	}
	return nil
}

func findContainerFrom(pod *corev1.Pod, name string) (*corev1.Container, error) {
	var container *corev1.Container

	if len(name) > 0 {
		container = searchContainerByName(pod, name)
		if container == nil {
			return nil, fmt.Errorf("container %q not found in pod %s", name, pod.Name)
		}
		return container, nil
	}

	if len(pod.Spec.Containers) <= 0 {
		return nil, fmt.Errorf("namespace(%s)'s pod(%s) does not have any containers", pod.Namespace, pod.Name)
	}

	// find default container
	if name := pod.Annotations["kubectl.kubernetes.io/default-container"]; len(name) > 0 {
		if container = searchContainerByName(pod, name); container != nil {
			return container, nil
		}
	}

	return &pod.Spec.Containers[0], nil
}

func AllContainerNames(pod *corev1.Pod) string {
	var containers []string
	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	for _, container := range pod.Spec.EphemeralContainers {
		containers = append(containers, fmt.Sprintf("%s (ephem)", container.Name))
	}
	for _, container := range pod.Spec.InitContainers {
		containers = append(containers, fmt.Sprintf("%s (init)", container.Name))
	}
	return strings.Join(containers, ", ")
}

func (c *Client) ResourceExec(gv schema.GroupVersion, resource string, params map[string]interface{}) (string, error) {
	// var result string
	var result interface{}
	var err error

	var namespace string
	var name string
	var command string
	var containerName string

	if found, err := FindCastFromMap(params, "namespace", &namespace); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	}

	if found, err := FindCastFromMap(params, "name", &name); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	}

	if found, err := FindCastFromMap(params, "command", &command); found && err != nil {
		return "", err
	} else if !found {
		return "", err
	}

	if found, err := FindCastFromMap(params, "container_name", &containerName); found && err != nil {
		return "", err
	}

	if containerName == "" {
		pod, err := c.client.CoreV1().Pods(namespace).Get(context.TODO(), name, v1.GetOptions{})
		if err != nil {
			return "", err
		}

		if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
			return "", fmt.Errorf("cannot exec into a container in a completed pod; current phase is %s", pod.Status.Phase)
		}

		container, err := findContainerFrom(pod, containerName)
		if err != nil {
			return "", err
		}
		containerName = container.Name
	}

	switch gv.Identifier() {
	case "v1":
		switch resource {
		case "pods":
			req := c.client.CoreV1().RESTClient().Post().Resource(resource).Name(name).Namespace(namespace).SubResource("exec").VersionedParams(&corev1.PodExecOptions{
				Container: containerName,
				Command:   strings.Split(command, " "),
				Stdin:     false,
				Stdout:    true,
				Stderr:    true,
				TTY:       false,
			}, scheme.ParameterCodec)

			exec, err := remotecommand.NewSPDYExecutor(c.restconfig, "POST", req.URL())
			if err != nil {
				return "", err
			}

			bufStdout := bytes.NewBuffer([]byte{})
			bufStderr := bytes.NewBuffer([]byte{})
			streamErr := exec.Stream(remotecommand.StreamOptions{
				Stdin:  nil,
				Stdout: bufStdout,
				Stderr: bufStderr,
				Tty:    false,
			})
			cmdRes := &CommandResult{Stdout: strings.TrimSpace(bufStdout.String()), Stderr: strings.TrimSpace(bufStderr.String())}
			if streamErr != nil {
				cmdRes.Err = strings.TrimSpace(streamErr.Error())
			}
			result = cmdRes
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

	// return result, nil
}

type CommandResult struct {
	Stdout string `json:"stdout,omitempty"`
	Stderr string `json:"stderr,omitempty"`
	Err    string `json:"error,omitempty"`
}

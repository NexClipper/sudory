package k8s

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var k8sClient *Client

type Client struct {
	client *kubernetes.Clientset
}

func getClientsetInCluster() (*kubernetes.Clientset, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// creates the clientset
	return kubernetes.NewForConfig(config)
}

func getClientsetOutOfCluster() (*kubernetes.Clientset, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	return kubernetes.NewForConfig(config)
}

func NewClient() (*Client, error) {
	if k8sClient != nil {
		return k8sClient, nil
	}

	clientset, err := getClientsetInCluster()
	if err == nil {
		return &Client{client: clientset}, nil
	}

	// If the getClientsetInCluster() fails to get clientset, use getClientsetOutOfCluster() function.
	clientset, err = getClientsetOutOfCluster()
	if err != nil {
		return nil, err
	}

	k8sClient = &Client{client: clientset}

	return k8sClient, nil
}

func GetClient() (*Client, error) {
	if k8sClient == nil {
		return NewClient()
	}

	return k8sClient, nil
}

func (c *Client) Pod(namespace string) *pods {
	return newPods(c, namespace)
}

func (c *Client) RawRequest() *rawRequest {
	return newRawRequest(c)
}

func (c *Client) ResourceRequest(gv schema.GroupVersion, resource, verb, namespace, name string, labels map[string]string) (string, error) {
	var result interface{}
	var err error

	switch gv.Identifier() {
	case "v1":
		switch resource {
		case "pod":
			switch verb {
			case "get":
				result, err = c.Pod(namespace).Get(context.TODO(), name)
				if err != nil {
					break
				}
			case "list":
				result, err = c.Pod(namespace).List(context.TODO(), labels)
				if err != nil {
					break
				}
			default:
				err = fmt.Errorf("unknown verb(%s)", verb)
			}
		case "namespace":
			switch verb {
			case "get":
				result, err = c.Namespace().Get(context.TODO(), name)
				if err != nil {
					break
				}
			case "list":
				result, err = c.Namespace().List(context.TODO(), labels)
				if err != nil {
					break
				}
			default:
				err = fmt.Errorf("unknown verb(%s)", verb)
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

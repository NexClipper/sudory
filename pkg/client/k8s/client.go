package k8s

import (
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

func (c *Client) RawRequest() *rawRequest {
	return newRawRequest(c)
}

func (c *Client) ResourceRequest(gv schema.GroupVersion, resource, verb, namespace, name string, labels map[string]string) (string, error) {
	var result string
	var err error

	switch verb {
	case "get":
		result, err = c.ResourceGet(gv, resource, namespace, name)
		if err != nil {
			break
		}
	case "list":
		result, err = c.ResourceList(gv, resource, namespace, labels)
		if err != nil {
			break
		}
	default:
		err = fmt.Errorf("unknown verb(%s)", verb)
	}

	if err != nil {
		return "", err
	}

	return result, nil
}

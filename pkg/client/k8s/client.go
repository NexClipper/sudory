package k8s

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

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
	clientset, err := getClientsetInCluster()
	if err == nil {
		return &Client{client: clientset}, nil
	}

	// If the getClientsetInCluster() fails to get clientset, use getClientsetOutOfCluster() function.
	clientset, err = getClientsetOutOfCluster()
	if err != nil {
		return nil, err
	}

	return &Client{client: clientset}, nil
}

func (c *Client) Pod(namespace string) *pods {
	return newPods(c, namespace)
}
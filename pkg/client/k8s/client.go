package k8s

import (
	"flag"
	"fmt"
	"path/filepath"
	"time"

	monitoringclient "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	aggregatorv1 "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
)

const defaultK8sTimeout = 10 * time.Second

var k8sClient *Client

type Client struct {
	client         *kubernetes.Clientset
	mclient        monitoringclient.Interface
	apiextv1client apiextensionsv1.Interface
	aggrev1client  aggregatorv1.Interface
	restconfig     *rest.Config
}

func getClusterConfig() (*rest.Config, error) {
	// creates the in-cluster config
	config, err1 := rest.InClusterConfig()
	if err1 == nil {
		return config, nil
	}

	// creates the out-of-cluster config
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err2 := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err2 != nil {
		return nil, fmt.Errorf("in cluster error: %s, out of cluster error: %s", err1.Error(), err2.Error())
	}

	return config, nil
}

func NewClient() (*Client, error) {
	if k8sClient != nil {
		return k8sClient, nil
	}

	cfg, err := getClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	mclient, err := monitoringclient.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	apiextv1client, err := apiextensionsv1.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	aggrev1client, err := aggregatorv1.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	k8sClient = &Client{client: clientset, mclient: mclient, apiextv1client: apiextv1client, aggrev1client: aggrev1client, restconfig: cfg}

	return k8sClient, nil
}

func GetClient() (*Client, error) {
	if k8sClient == nil {
		return NewClient()
	}

	return k8sClient, nil
}

func (c *Client) GetK8sClientset() *kubernetes.Clientset {
	return c.client
}

func (c *Client) RawRequest() *rawRequest {
	return newRawRequest(c)
}

func (c *Client) ResourceRequest(gv schema.GroupVersion, resource, verb string, args map[string]interface{}) (string, error) {
	var result string
	var err error

	switch verb {
	case "get":
		result, err = c.ResourceGet(gv, resource, args)
		if err != nil {
			break
		}
	case "list":
		result, err = c.ResourceList(gv, resource, args)
		if err != nil {
			break
		}
	case "delete":
		err = c.ResourceDelete(gv, resource, args)
		if err != nil {
			break
		}
	case "patch":
		result, err = c.ResourcePatch(gv, resource, args)
		if err != nil {
			break
		}
	case "exec":
		result, err = c.ResourceExec(gv, resource, args)
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

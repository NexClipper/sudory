package service

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/NexClipper/sudory/pkg/client/k8s"
	"github.com/NexClipper/sudory/pkg/client/p8s"
	"github.com/NexClipper/sudory/pkg/server/macro"
)

type CommandType int

const (
	CommandTypeK8s = iota + 1
	CommandTypeP8s
	CommandTypeHelm
)

func (ct CommandType) String() string {
	if ct == CommandTypeK8s {
		return "kubernetes"
	} else if ct == CommandTypeP8s {
		return "prometheus"
	} else if ct == CommandTypeHelm {
		return "helm"
	}

	return "Unknown CommandType"
}

type Commander interface {
	GetCommandType() CommandType
	Run() (string, error)
}

func NewCommander(command *StepCommand) (Commander, error) {
	mlist := strings.Split(command.Method, ".")
	ctype := mlist[0]

	switch ctype {
	case "kubernetes":
		return NewK8sCommander(command)
	case "prometheus":
		return NewP8sCommander(command)
	case "helm":
		//
	}

	return nil, fmt.Errorf("unknown command method(%s)", command.Method)
}

type K8sCommander struct {
	client    *k8s.Client
	gv        schema.GroupVersion // v1, apps/v1, ...
	resource  string              // pod, namespace, deployment, ...
	namespace string              // "", default, ...
	name      string              // my-pod, my-namespace, ...
	verb      string              // get, list, watch, ...
	labels    map[string]string
}

func NewK8sCommander(command *StepCommand) (Commander, error) {
	client, err := k8s.GetClient()
	if err != nil {
		return nil, err
	}
	cmdr := &K8sCommander{client: client, labels: make(map[string]string)}

	err = cmdr.ParseCommand(command)
	if err != nil {
		return nil, err
	}

	return cmdr, nil
}

func (c *K8sCommander) GetCommandType() CommandType {
	return CommandTypeK8s
}

func (c *K8sCommander) ParseCommand(command *StepCommand) error {
	mlist := strings.SplitN(command.Method, ".", 4)

	if len(mlist) != 4 {
		return fmt.Errorf("there is not enough method(%s) for k8s. want(4) but got(%d)", command.Method, len(mlist))
	}

	gv, err := schema.ParseGroupVersion(mlist[3])
	if err != nil {
		return err
	}

	c.gv = gv
	c.resource = mlist[1]
	c.verb = mlist[2]

	if s, ok := macro.MapString(command.Args, "namespace"); ok {
		c.namespace = s
	}

	if s, ok := macro.MapString(command.Args, "name"); ok {
		c.name = s
	}

	if m, ok := macro.MapMap(command.Args, "labels"); ok {
		for k, v := range m {
			c.labels[k] = fmt.Sprintf("%v", v)
		}
	}

	return nil
}

func (c *K8sCommander) Run() (string, error) {
	return c.client.ResourceRequest(c.gv, c.resource, c.verb, c.namespace, c.name, c.labels)
}

type P8sCommander struct {
	client      *p8s.Client
	apiVersion  string
	api         string
	queryParams map[string]interface{}
}

func NewP8sCommander(command *StepCommand) (Commander, error) {
	cmdr := &P8sCommander{}

	if err := cmdr.ParseCommand(command); err != nil {
		return nil, err
	}

	return cmdr, nil
}

func (c *P8sCommander) GetCommandType() CommandType {
	return CommandTypeP8s
}

func (c *P8sCommander) ParseCommand(command *StepCommand) error {
	mlist := strings.SplitN(command.Method, ".", 3)

	if len(mlist) != 3 {
		return fmt.Errorf("there is not enough method(%s) for p8s. want(3) but got(%d)", command.Method, len(mlist))
	}

	c.api = mlist[1]
	c.apiVersion = mlist[2]
	c.queryParams = command.Args

	url, ok := macro.MapString(command.Args, "url")
	if !ok || len(url) == 0 {
		return fmt.Errorf("prometheus url is empty")
	}

	client, err := p8s.NewClient(url)
	if err != nil {
		return err
	}

	c.client = client

	return nil
}

func (c *P8sCommander) Run() (string, error) {
	return c.client.ApiRequest(c.apiVersion, c.api, c.queryParams)
}

type HelmCommander struct {
	// client *helm.Client
}

func (c *HelmCommander) GetCommandType() CommandType {
	return CommandTypeHelm
}

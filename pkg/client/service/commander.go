package service

import (
	"fmt"
	"strings"

	"github.com/NexClipper/sudory/pkg/client/k8s"
	"github.com/NexClipper/sudory/pkg/server/macro"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
		//
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
	cmdr := &K8sCommander{client: client}

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
		delete(command.Args, "namespace")
	}

	if s, ok := macro.MapString(command.Args, "name"); ok {
		c.name = s
		delete(command.Args, "name")
	}
	// for k, v := range command.Args {
	// 	if k == "namespace" {

	// 		c.namespace = v
	// 		delete(command.Args, k)
	// 	} else if k == "name" {
	// 		c.name = v
	// 		delete(command.Args, k)
	// 	}
	// }
	c.labels = command.Args

	return nil
}

func (c *K8sCommander) Run() (string, error) {
	return c.client.ResourceRequest(c.gv, c.resource, c.verb, c.namespace, c.name, c.labels)
}

type P8sCommander struct {
	// client *p8s.Client
}

func (c *P8sCommander) GetCommandType() CommandType {
	return CommandTypeP8s
}

type HelmCommander struct {
	// client *helm.Client
}

func (c *HelmCommander) GetCommandType() CommandType {
	return CommandTypeHelm
}

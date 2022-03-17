package executor

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/NexClipper/sudory/pkg/client/helm"
	"github.com/NexClipper/sudory/pkg/client/jq"
	"github.com/NexClipper/sudory/pkg/client/k8s"
	"github.com/NexClipper/sudory/pkg/client/p8s"
	"github.com/NexClipper/sudory/pkg/client/service"
	"github.com/NexClipper/sudory/pkg/server/macro"
)

type CommandType int

const (
	CommandTypeK8s = iota + 1
	CommandTypeP8s
	CommandTypeHelm
	CommandTypeJq
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

func NewCommander(command *service.StepCommand) (Commander, error) {
	mlist := strings.Split(command.Method, ".")
	ctype := mlist[0]

	switch ctype {
	case "kubernetes":
		return NewK8sCommander(command)
	case "prometheus":
		return NewP8sCommander(command)
	case "helm":
		return NewHelmCommander(command)
	case "jq":
		return NewJqCommander(command)
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

func NewK8sCommander(command *service.StepCommand) (Commander, error) {
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

func (c *K8sCommander) ParseCommand(command *service.StepCommand) error {
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

func NewP8sCommander(command *service.StepCommand) (Commander, error) {
	cmdr := &P8sCommander{}

	if err := cmdr.ParseCommand(command); err != nil {
		return nil, err
	}

	return cmdr, nil
}

func (c *P8sCommander) GetCommandType() CommandType {
	return CommandTypeP8s
}

func (c *P8sCommander) ParseCommand(command *service.StepCommand) error {
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
	client *helm.Client
	cmd    string
	args   map[string]interface{}
}

func NewHelmCommander(command *service.StepCommand) (Commander, error) {
	cmdr := &HelmCommander{}

	if err := cmdr.ParseCommand(command); err != nil {
		return nil, err
	}

	return cmdr, nil
}

func (c *HelmCommander) GetCommandType() CommandType {
	return CommandTypeHelm
}

func (c *HelmCommander) ParseCommand(command *service.StepCommand) error {
	mlist := strings.SplitN(command.Method, ".", 2)

	if len(mlist) != 2 {
		return fmt.Errorf("there is not enough method(%s) for helm. want(3) but got(%d)", command.Method, len(mlist))
	}

	c.cmd = mlist[1]
	c.args = command.Args

	client, err := helm.NewClient()
	if err != nil {
		return err
	}

	c.client = client

	return nil
}

func (c *HelmCommander) Run() (string, error) {
	return c.client.Request(c.cmd, c.args)
}

type JqCommander struct {
	input  map[string]interface{}
	filter string
}

func NewJqCommander(command *service.StepCommand) (Commander, error) {
	cmdr := &JqCommander{}

	if err := cmdr.ParseCommand(command); err != nil {
		return nil, err
	}

	return cmdr, nil
}

func (c *JqCommander) GetCommandType() CommandType {
	return CommandTypeJq
}

func (c *JqCommander) ParseCommand(command *service.StepCommand) error {
	if m, ok := macro.MapMap(command.Args, "input"); ok {
		c.input = m
	} else {
		return fmt.Errorf("input not found")
	}

	if f, ok := macro.MapString(command.Args, "filter"); ok {
		c.filter = f
	} else {
		return fmt.Errorf("filter not found")
	}

	return nil
}

func (c *JqCommander) Run() (string, error) {
	return jq.Request(c.input, c.filter)
}

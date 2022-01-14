package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) Namespace() *namespaces {
	return newNamespaces(c)
}

type namespaces struct {
	c *Client
}

func newNamespaces(c *Client) *namespaces {
	return &namespaces{c: c}
}

func (ns *namespaces) Get(name string) (*corev1.Namespace, error) {
	return ns.c.client.CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
}

func (ns *namespaces) List(labels map[string]string) (*corev1.NamespaceList, error) {
	labelsString, err := convertMapToLabelSelector(labels)
	if err != nil {
		return nil, err
	}

	return ns.c.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{LabelSelector: labelsString})
}

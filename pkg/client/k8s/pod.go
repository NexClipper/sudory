package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type pods struct {
	c  *Client
	ns string
}

func newPods(c *Client, namespace string) *pods {
	return &pods{
		c:  c,
		ns: namespace,
	}
}

func (c *pods) Get(ctx context.Context, name string) (*corev1.Pod, error) {
	return c.c.client.CoreV1().Pods(c.ns).Get(ctx, name, metav1.GetOptions{})
}

func (c *pods) List(ctx context.Context, labels map[string]string) (*corev1.PodList, error) {
	labelsString, err := convertMapToLabelSelector(labels)
	if err != nil {
		return nil, err
	}

	return c.c.client.CoreV1().Pods(c.ns).List(ctx, metav1.ListOptions{LabelSelector: labelsString})
}

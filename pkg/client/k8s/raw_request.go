package k8s

import (
	"context"
	"fmt"
	"log"
)

type rawRequest struct {
	c *Client
}

func newRawRequest(c *Client) *rawRequest {
	return &rawRequest{
		c: c,
	}
}

func (c *rawRequest) CheckApiServerStatus() error {
	path := "/livez"

	log.Printf("Send request to the endpoint '%s' of the k8s api-server.\n", path)

	result, err := c.c.client.RESTClient().Get().AbsPath(path).DoRaw(context.TODO())
	if err != nil {
		return fmt.Errorf("failed request to the endpoint '%s' of the k8s api-server", path)
	}

	resultStr := string(result)

	log.Printf("Received from the endpoint '%s' of k8s api-server : %s\n", path, resultStr)

	if resultStr != "ok" {
		return fmt.Errorf("k8s api-server's status is bad")
	}

	return nil
}

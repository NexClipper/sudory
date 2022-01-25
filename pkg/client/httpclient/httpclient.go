package httpclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	urlPkg "net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type HttpClient struct {
	url           string
	token         string
	RetryMax      int
	RetryInterval int
}

func NewHttpClient(url, token string, retryMax, retryInterval int) *HttpClient {
	return &HttpClient{url: url, token: token, RetryMax: retryMax, RetryInterval: retryInterval}
}

func (hc *HttpClient) Request(method, path string, params map[string]string, bodyType string, rawBody []byte) ([]byte, error) {
	req, err := retryablehttp.NewRequest(method, hc.url+path, rawBody)
	if err != nil {
		return nil, err
	}

	urlValues := urlPkg.Values{}
	for k, v := range params {
		urlValues.Add(k, v)
	}
	req.URL.RawQuery = urlValues.Encode()

	if len(bodyType) > 0 {
		req.Header.Set("Content-Type", bodyType)

	}

	// TODO: req.Header.Add("token", hc.token)

	client := retryablehttp.NewClient()

	client.RetryMax = hc.RetryMax
	client.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		return time.Millisecond * time.Duration(hc.RetryInterval)
	}

	// client.HTTPClient.Timeout = time.Second

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("received http status error code : %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *HttpClient) Get(path string, params map[string]string) ([]byte, error) {
	return c.Request("GET", path, params, "", nil)
}

func (c *HttpClient) Post(path string, params map[string]string, rawBody []byte) ([]byte, error) {
	return c.Request("POST", path, params, "", rawBody)
}

func (c *HttpClient) PostJson(path string, params map[string]string, rawBody []byte) ([]byte, error) {
	return c.Request("POST", path, params, "application/json", rawBody)
}

func (c *HttpClient) Put(path string, params map[string]string, rawBody []byte) ([]byte, error) {
	return c.Request("PUT", path, params, "", rawBody)
}

func (c *HttpClient) PutJson(path string, params map[string]string, rawBody []byte) ([]byte, error) {
	return c.Request("PUT", path, params, "application/json", rawBody)
}

func (c *HttpClient) Delete(path string, params map[string]string) ([]byte, error) {
	return c.Request("DELETE", path, params, "", nil)
}

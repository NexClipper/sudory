package httpclient

import (
	"io/ioutil"
	"net/http"
	urlPkg "net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/NexClipper/sudory/pkg/client/log"
	"github.com/NexClipper/sudory/pkg/server/macro/jwt"
	sessionv1 "github.com/NexClipper/sudory/pkg/server/model/session/v1"
)

const CustomHeaderClientToken = "x-sudory-client-token"

type HttpClient struct {
	url    string
	token  string
	client *retryablehttp.Client
}

func NewHttpClient(url, token string, retryMax, retryInterval int) *HttpClient {
	client := retryablehttp.NewClient()

	client.HTTPClient.Transport.(*http.Transport).MaxIdleConns = 100
	client.HTTPClient.Transport.(*http.Transport).MaxIdleConnsPerHost = 100
	client.Logger = &log.RetryableHttpLogger{}
	client.RetryMax = retryMax
	client.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		return time.Millisecond * time.Duration(retryInterval)
	}
	client.ErrorHandler = RetryableHttpErrorHandler
	// client.HTTPClient.Timeout = time.Second

	return &HttpClient{url: url, token: token, client: client}
}

func (hc *HttpClient) GetToken() string {
	return hc.token
}

func (hc *HttpClient) IsTokenExpired() bool {
	claims := new(sessionv1.ClientSessionPayload)
	if err := jwt.BindPayload(hc.token, &claims); err != nil {
		log.Warnf("jwt.BindPayload error : %v\n", err)
		return false
	}

	if time.Until(claims.Exp) < 0 {
		return true
	}

	return false
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

	if hc.token != "" {
		req.Header.Add(CustomHeaderClientToken, hc.token)
	}

	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := CheckHttpResponseError(resp); err != nil {
		return nil, err
	}

	if recvToken := resp.Header.Get(CustomHeaderClientToken); recvToken != "" {
		hc.token = recvToken
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

func (c *HttpClient) PostForm(path string, params map[string]string, rawBody []byte) ([]byte, error) {
	return c.Request("POST", path, params, "application/x-www-form-urlencoded", rawBody)
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

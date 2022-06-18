package httpclient

import (
	"bytes"
	"context"
	"io"
	"net/http"
	urlPkg "net/url"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hashicorp/go-retryablehttp"

	"github.com/NexClipper/sudory/pkg/client/log"
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
	jwt_token, _, err := jwt.NewParser().ParseUnverified(hc.token, claims)
	if _, ok := jwt_token.Claims.(*sessionv1.ClientSessionPayload); !ok || err != nil {
		log.Warnf("jwt.ParseUnverified error : %v\n", err)
		return true
	}

	return !claims.VerifyExpiresAt(time.Now().Unix(), true)
}

func (hc *HttpClient) Request(ctx context.Context, method, path string, params map[string]string, bodyType string, rawBody []byte) ([]byte, error) {
	req, err := retryablehttp.NewRequest(method, hc.url+path, rawBody)
	if err != nil {
		return nil, err
	}

	urlValues := urlPkg.Values{}
	for k, v := range params {
		urlValues.Set(k, v)
	}
	req.URL.RawQuery = urlValues.Encode()

	if len(bodyType) > 0 {
		req.Header.Set("Content-Type", bodyType)

	}

	if hc.token != "" {
		req.Header.Set(CustomHeaderClientToken, hc.token)
	}

	req = req.WithContext(ctx)

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

	var result bytes.Buffer
	io.Copy(&result, resp.Body)

	return result.Bytes(), nil
}

func (c *HttpClient) Get(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	return c.Request(ctx, "GET", path, params, "", nil)
}

func (c *HttpClient) GetJson(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	return c.Request(ctx, "GET", path, params, "application/json", nil)
}

func (c *HttpClient) Post(ctx context.Context, path string, params map[string]string, rawBody []byte) ([]byte, error) {
	return c.Request(ctx, "POST", path, params, "", rawBody)
}

func (c *HttpClient) PostJson(ctx context.Context, path string, params map[string]string, rawBody []byte) ([]byte, error) {
	return c.Request(ctx, "POST", path, params, "application/json", rawBody)
}

func (c *HttpClient) PostForm(ctx context.Context, path string, params map[string]string, rawBody []byte) ([]byte, error) {
	return c.Request(ctx, "POST", path, params, "application/x-www-form-urlencoded", rawBody)
}

func (c *HttpClient) Put(ctx context.Context, path string, params map[string]string, rawBody []byte) ([]byte, error) {
	return c.Request(ctx, "PUT", path, params, "", rawBody)
}

func (c *HttpClient) PutJson(ctx context.Context, path string, params map[string]string, rawBody []byte) ([]byte, error) {
	return c.Request(ctx, "PUT", path, params, "application/json", rawBody)
}

func (c *HttpClient) Delete(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	return c.Request(ctx, "DELETE", path, params, "", nil)
}

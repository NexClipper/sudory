package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"

	"github.com/NexClipper/sudory/pkg/client/log"
)

type HttpClient struct {
	root   *url.URL
	client *retryablehttp.Client
}

func NewHttpClient(address string, defaultTLS bool, retryMax, retryInterval int) (*HttpClient, error) {
	defaultUrl, err := DefaultURL(address, defaultTLS)
	if err != nil {
		return nil, err
	}
	client := retryablehttp.NewClient()

	client.HTTPClient.Transport.(*http.Transport).MaxIdleConns = 100
	client.HTTPClient.Transport.(*http.Transport).MaxIdleConnsPerHost = 100

	client.Logger = &log.RetryableHttpLogger{}
	client.RetryMax = retryMax
	client.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		return time.Millisecond * time.Duration(retryInterval)
	}
	client.ErrorHandler = RetryableHttpErrorHandler

	return &HttpClient{root: defaultUrl, client: client}, nil
}

func (c *HttpClient) SetDisableKeepAlives() {
	c.client.HTTPClient.Transport.(*http.Transport).DisableKeepAlives = true
}

func (c *HttpClient) Get(path string) *Request {
	return NewRequest(c, "GET", path)
}

func (c *HttpClient) Post(path string) *Request {
	return NewRequest(c, "POST", path)
}

func (c *HttpClient) Put(path string) *Request {
	return NewRequest(c, "PUT", path)
}

func (c *HttpClient) Delete(path string) *Request {
	return NewRequest(c, "DELETE", path)
}

type Request struct {
	c       *HttpClient
	method  string
	path    string
	params  url.Values
	headers http.Header
	body    interface{}
}

func NewRequest(c *HttpClient, method, path string) *Request {
	r := &Request{
		c:       c,
		method:  method,
		path:    path,
		params:  make(url.Values),
		headers: make(http.Header),
	}

	return r
}

func (r *Request) SetHeader(key, value string) *Request {
	if r.headers == nil {
		r.headers = make(http.Header)
	}
	r.headers.Set(key, value)
	return r
}

func (r *Request) SetParam(key string, values ...string) *Request {
	if r.params == nil {
		r.params = make(url.Values)
	}
	for _, v := range values {
		r.params.Add(key, v)
	}
	return r
}

func (r *Request) SetParamFromQuery(query url.Values) *Request {
	if query == nil || len(query) <= 0 {
		return r
	}

	if r.params == nil {
		r.params = make(url.Values)
	}
	for k, v := range query {
		for _, vv := range v {
			r.params.Add(k, vv)
		}
	}
	return r
}

func (r *Request) SetBody(bodyType string, body interface{}) *Request {
	r.SetHeader("Content-Type", bodyType)
	r.body = body
	return r
}

func (r *Request) URL() *url.URL {
	u := new(url.URL)
	if r.c.root != nil {
		*u = *r.c.root
	}
	u.Path = r.path

	if len(r.params) > 0 {
		u.RawQuery = r.params.Encode()
	}

	return u
}

func (r *Request) Do(ctx context.Context) Result {
	if r.c.client == nil {
		return Result{err: fmt.Errorf("httpclient is nil")}
	}

	url := r.URL().String()
	req, err := retryablehttp.NewRequest(r.method, url, r.body)
	if err != nil {
		return Result{err: err}
	}
	req = req.WithContext(ctx)
	req.Header = r.headers

	resp, err := r.c.client.Do(req)
	defer r.closeBody(resp)

	if err != nil {
		return Result{err: err}
	}

	return r.extractResultFrom(resp)
}

func (r *Request) closeBody(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()

		const maxLimitBytes = int64(4096)
		io.Copy(ioutil.Discard, io.LimitReader(resp.Body, maxLimitBytes))
	}
}

func (r *Request) extractResultFrom(resp *http.Response) Result {
	var result Result

	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()

		var buf bytes.Buffer
		_, result.err = io.Copy(&buf, resp.Body)
		result.body = buf.Bytes()

		if result.err != nil {
			return result
		}
	}

	result.statusCode = resp.StatusCode
	result.headers = resp.Header

	// check response status code
	if result.statusCode < http.StatusOK || result.statusCode >= http.StatusBadRequest {
		if len(result.body) > 0 {
			result.err = fmt.Errorf("%s %s, status code : %d(%s), body : %s", resp.Request.Method, resp.Request.URL.String(), resp.StatusCode, resp.Status, strings.TrimSpace(string(result.body)))
		} else {
			result.err = fmt.Errorf("%s %s, status code : %d(%s)", resp.Request.Method, resp.Request.URL.String(), resp.StatusCode, resp.Status)
		}
		return result
	}

	return result
}

type Result struct {
	body       []byte
	err        error
	statusCode int
	headers    http.Header
}

func (r Result) Raw() ([]byte, error) {
	return r.body, r.err
}

func (r Result) IntoJson(obj interface{}) error {
	if r.err != nil {
		return r.err
	}
	return json.Unmarshal(r.body, obj)
}

func (r Result) Headers() http.Header {
	return r.headers
}

func (r Result) StatusCode() int {
	return r.statusCode
}

func (r Result) Error() error {
	return r.err
}

func (r Result) SetError(err error) Result {
	r.err = err
	return r
}

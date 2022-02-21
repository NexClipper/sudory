package httpclient

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	urlPkg "net/url"
	"strings"
)

func ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("url is empty : got(%s)", url)
	}

	parsedURL, err := urlPkg.Parse(url)
	if err != nil {
		return err
	}

	if parsedURL.Scheme == "" || !(parsedURL.Scheme == "http" || parsedURL.Scheme == "https") {
		return fmt.Errorf("url scheme is empty : want:(http or https), got(%s)", parsedURL.Scheme)
	}

	host, port, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		return err
	}

	if host == "" || port == "" {
		return fmt.Errorf("host or port is empty : host(%s), port(%s)", host, port)
	}

	return nil
}

func RetryableHttpErrorHandler(resp *http.Response, err error, numTries int) (*http.Response, error) {
	if err := CheckHttpResponseError(resp); err != nil {
		return nil, err
	}

	return resp, err
}

func CheckHttpResponseError(resp *http.Response) error {
	if resp != nil {
		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			return fmt.Errorf("%s %s, status code : %d, body : %s", resp.Request.Method, resp.Request.URL.String(), resp.StatusCode, strings.TrimSpace(string(body)))
		}
	}
	return nil
}

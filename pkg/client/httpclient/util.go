package httpclient

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	neturl "net/url"
	"strings"
)

func last(s string, b byte) int {
	i := len(s)
	for i--; i >= 0; i-- {
		if s[i] == b {
			break
		}
	}
	return i
}

func ValidateURL(url string) error {
	if url == "" {
		return fmt.Errorf("url is empty : got(%s)", url)
	}

	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return err
	}

	if parsedURL.Scheme == "" || !(parsedURL.Scheme == "http" || parsedURL.Scheme == "https") {
		return fmt.Errorf("url scheme is empty : want:(http or https), got(%s)", parsedURL.Scheme)
	}

	i := last(parsedURL.Host, ':')

	if i >= 0 {
		host, port, err := net.SplitHostPort(parsedURL.Host)
		if host == "" || port == "" {
			return fmt.Errorf("host or port is empty : host(%s), port(%s)", host, port)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func DefaultURL(url string, defaultTLS bool) (*neturl.URL, error) {
	if url == "" {
		return nil, fmt.Errorf("url is empty : got(%s)", url)
	}

	parsedURL, err := neturl.Parse(url)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		scheme := "http://"
		if defaultTLS {
			scheme = "https://"
		}
		parsedURL, err = neturl.Parse(scheme + url)
		if err != nil {
			return nil, err
		}
		if parsedURL.Path != "" && parsedURL.Path != "/" {
			return nil, fmt.Errorf("default url must not have a path: %s", url)
		}
	}

	return parsedURL, nil
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

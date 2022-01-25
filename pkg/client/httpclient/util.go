package httpclient

import (
	"fmt"
	"net"
	urlPkg "net/url"
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

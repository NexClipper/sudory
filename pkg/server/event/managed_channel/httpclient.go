package managed_channel

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
)

type HttpClient_opt struct {
	Method         string
	Url            string
	RequestHeaders map[string]string
	RequestTimeout time.Duration
}

type part struct {
	Header  textproto.MIMEHeader // textproto.MIMEHeader{"Content-Type": {content_type}}
	Payload []byte
}

func HttpMultipartReq(opt *HttpClient_opt, httpclient *http.Client, parts []part) error {
	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	for _, part := range parts {
		w, err := writer.CreatePart(part.Header)
		if err != nil {
			return errors.Wrapf(err, "create multipart")
		}

		if _, err := w.Write(part.Payload); err != nil {
			return errors.Wrapf(err, "multipart write")
		}
	}
	writer.Close()

	//create http request with timeout context
	requset_timeout := func() time.Duration {
		if opt.RequestTimeout == 0 {
			return 3 * time.Second
		}
		return opt.RequestTimeout
	}

	timeout, cancel := context.WithTimeout(context.Background(), requset_timeout())
	defer cancel()

	req, err := http.NewRequestWithContext(timeout, opt.Method, opt.Url, body)
	if err != nil {
		return errors.Wrapf(err, "make request with context%s",
			logs.KVL(
				"method", opt.Method,
				"url", opt.Url,
			))
	}

	//Header
	for key, val := range opt.RequestHeaders {
		req.Header.Set(key, val) //set http header
	}
	req.Header.Set("Content-Type", "multipart/mixed; boundary="+writer.Boundary())

	//request to host
	rsp, err := httpclient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "request to%s",
			logs.KVL(
				"method", opt.Method,
				"url", opt.Url,
				"headers", opt.RequestHeaders,
			))
	}
	defer rsp.Body.Close()

	//read response
	buffer := bytes.Buffer{}
	if _, err := buffer.ReadFrom(rsp.Body); err != nil {
		return errors.Wrapf(err, "read to response body")
	}

	//check status code
	if rsp.StatusCode/100 != 2 {
		return errors.Errorf("bad response status%s",
			logs.KVL(
				"status", rsp.Status,
				"code", rsp.StatusCode,
				"body", buffer.String(),
			))
	}

	return nil
}

func HttpReq(opt *HttpClient_opt, httpclient *http.Client, content_type string, paylaod []byte) error {
	body := bytes.NewBuffer(paylaod)

	//create http request with timeout context
	requset_timeout := func() time.Duration {
		if opt.RequestTimeout == 0 {
			return 3 * time.Second
		}
		return opt.RequestTimeout
	}

	timeout, cancel := context.WithTimeout(context.Background(), requset_timeout())
	defer cancel()

	req, err := http.NewRequestWithContext(timeout, opt.Method, opt.Url, body)
	if err != nil {
		return errors.Wrapf(err, "make request with context%s",
			logs.KVL(
				"method", opt.Method,
				"url", opt.Url,
			))
	}

	//Header
	for key, val := range opt.RequestHeaders {
		req.Header.Set(key, val) //set http header
	}

	req.Header.Set("Content-Type", content_type)

	//request to host
	rsp, err := httpclient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "request to%s",
			logs.KVL(
				"method", opt.Method,
				"url", opt.Url,
				"headers", opt.RequestHeaders,
			))
	}
	defer rsp.Body.Close()

	//read response
	buffer := bytes.Buffer{}
	if _, err := buffer.ReadFrom(rsp.Body); err != nil {
		return errors.Wrapf(err, "read to response body")
	}

	//check status code
	if rsp.StatusCode/100 != 2 {
		return errors.Errorf("bad response status%s",
			logs.KVL(
				"status", rsp.Status,
				"code", rsp.StatusCode,
				"body", buffer.String(),
			))
	}
	// }
	return nil
}

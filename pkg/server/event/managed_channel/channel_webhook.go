package managed_channel

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"time"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	channelv1 "github.com/NexClipper/sudory/pkg/server/model/channel/v2"
	"github.com/pkg/errors"
)

type ChannelWebhook struct {
	uuid string
	opt  *channelv1.NotifierWebhook_property //config.WebhookNotifierConfig

	httpclient *http.Client //http.Client
}

func NewChannelWebhook(uuid string, opt channelv1.NotifierWebhook_property) *ChannelWebhook {
	if len(opt.Method) == 0 {
		opt.Method = http.MethodGet //set default Method
	}
	opt.Method = strings.ToUpper(opt.Method) //Method to upper

	client := http.DefaultClient
	if transport, ok := client.Transport.(*http.Transport); ok {
		transport.MaxIdleConnsPerHost = 100 //MaxIdleConnsPerHost
	}
	notifier := &ChannelWebhook{}
	notifier.uuid = uuid
	notifier.opt = &opt
	notifier.httpclient = client

	return notifier
}
func (channel ChannelWebhook) Type() fmt.Stringer {
	return channel.opt.Type()
}

func (channel ChannelWebhook) Uuid() string {
	return channel.uuid
}

func (channel ChannelWebhook) Property() map[string]string {
	return map[string]string{
		"type":   channel.opt.Type().String(),
		"uuid":   channel.uuid,
		"method": channel.opt.Method,
		"url":    channel.opt.Url,
	}
}

func (channel *ChannelWebhook) Close() {

}

func (channel ChannelWebhook) OnNotify(factory MarshalFactoryResult) error {

	httpclient := channel.httpclient

	// httpreq := notifier.HttpMultipartReq
	httpreq := channel.HttpReq

	if err := httpreq(channel.opt, httpclient, factory); err != nil {
		return errors.Wrapf(err, "http request")
	}
	return nil
}

func (ChannelWebhook) HttpMultipartReq(opt *channelv1.NotifierWebhook_property, httpclient *http.Client, factory MarshalFactoryResult) error {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	b, err := factory("application/json")
	if err != nil {
		return errors.Wrapf(err, "marshal factory")
	}

	part, err := writer.CreatePart(textproto.MIMEHeader{"Content-Type": {"application/json"}})
	if err != nil {
		return errors.Wrapf(err, "create multipart")
	}
	if _, err := part.Write(b); err != nil {
		return errors.Wrapf(err, "multipart write")
	}

	writer.Close()

	//create http request with timeout context
	requset_timeout := time.Duration(opt.RequestTimeout) * time.Second
	if requset_timeout == 0 {
		requset_timeout = 15 * time.Second
	}

	timeout, cancel := context.WithTimeout(context.Background(), requset_timeout)
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
	for key, val := range opt.RequestHeaders.KeyValue {
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

func (ChannelWebhook) HttpReq(opt *channelv1.NotifierWebhook_property, httpclient *http.Client, factory MarshalFactoryResult) error {

	b, err := factory("application/json")
	if err != nil {
		return errors.Wrapf(err, "marshal factory")
	}

	// for _, b := range b {
	body := bytes.NewBuffer(b)

	//create http request with timeout context
	requset_timeout := func() time.Duration {
		if opt.RequestTimeout == 0 {
			return 15 * time.Second
		}
		return time.Duration(opt.RequestTimeout) * time.Second
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
	for key, val := range opt.RequestHeaders.KeyValue {
		req.Header.Set(key, val) //set http header
	}

	req.Header.Set("Content-Type", "application/json")

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

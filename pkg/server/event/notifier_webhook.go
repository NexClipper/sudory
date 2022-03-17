package event

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
)

type webhookNotifier struct {
	opt WebhookNotifierConfig //config.WebhookNotifierConfig
	sub EventSubscriber       //include.EventSubscriber

	httpclient *http.Client //http.Client
}

func NewWebhookNotifier(opt WebhookNotifierConfig) *webhookNotifier {
	if len(opt.Method) == 0 {
		opt.Method = http.MethodGet //set default Method
	}
	opt.Method = strings.ToUpper(opt.Method) //Method to upper

	if opt.RequestTimeout == 0 {
		opt.RequestTimeout = 15 //set default timeout
	}

	client := http.DefaultClient
	if transport, ok := client.Transport.(*http.Transport); ok {
		transport.MaxIdleConnsPerHost = 100 //MaxIdleConnsPerHost
	}
	notifier := &webhookNotifier{}
	notifier.opt = opt
	notifier.httpclient = client

	return notifier
}
func (notifier webhookNotifier) Type() string {
	return NotifierTypeWebhook.String()
}

func (notifier webhookNotifier) Property() map[string]string {
	return map[string]string{
		"name":   notifier.sub.Config().Name,
		"type":   notifier.Type(),
		"method": notifier.opt.Method,
		"url":    notifier.opt.Url,
	}
}

func (notifier webhookNotifier) PropertyString() string {
	buff := bytes.Buffer{}
	for key, value := range notifier.Property() {
		if 0 < buff.Len() {
			buff.WriteString(" ")
		}
		buff.WriteString(key)
		buff.WriteString("=")
		buff.WriteString(strconv.Quote(value))
	}
	return buff.String()
}

func (notifier *webhookNotifier) Regist(sub EventSubscriber) {
	//Subscribe
	if !(sub == nil && notifier.sub != nil) {
		notifier.sub = sub
		notifier.sub.Notifiers().Add(notifier)
	}
}

func (notifier *webhookNotifier) Close() {
	//Unsubscribe
	if notifier.sub != nil {
		notifier.sub.Notifiers().Remove(notifier)
		notifier.sub = nil
	}
}

func (notifier webhookNotifier) OnNotify(factory MarshalFactory) error {
	opt := notifier.opt
	httpclient := notifier.httpclient

	if err := notifier.HttpMultipartReq(opt, httpclient, factory); err != nil {
		return errors.Wrapf(err, "http multipart request")
	}
	return nil
}

func (notifier webhookNotifier) OnNotifyAsync(factory MarshalFactory) <-chan NotifierFuture {
	future := make(chan NotifierFuture)
	go func() {
		defer close(future)

		future <- NotifierFuture{&notifier, notifier.OnNotify(factory)}
	}()

	return future
}

func (webhookNotifier) HttpMultipartReq(opt WebhookNotifierConfig, httpclient *http.Client, factory MarshalFactory) error {

	timeout, cancel := context.WithTimeout(context.Background(), opt.RequestTimeout)
	defer cancel()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	b, err := factory("application/json")
	if err != nil {
		return errors.Wrapf(err, "marshal factory")
	}
	for _, b := range b {
		part, err := writer.CreatePart(textproto.MIMEHeader{"Content-Type": {"application/json"}})
		if err != nil {
			return errors.Wrapf(err, "create multipart")
		}
		if _, err := part.Write(b); err != nil {
			return errors.Wrapf(err, "multipart write")
		}
	}
	writer.Close()

	//Create request with timeout
	req, err := http.NewRequestWithContext(timeout, opt.Method, opt.Url, body)
	if err != nil {
		return errors.Wrapf(err, "make request with context %s",
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

	//Request to update services's status to server.
	rsp, err := httpclient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "request to %s",
			logs.KVL(
				"method", opt.Method,
				"url", opt.Url,
				"headers", opt.RequestHeaders,
			))
	}
	defer rsp.Body.Close()

	buffer := bytes.Buffer{}
	if _, err := buffer.ReadFrom(rsp.Body); err != nil {
		return errors.Wrapf(err, "read to response body")
	}

	if rsp.StatusCode/100 != 2 {
		err := errors.Errorf("bad response status %s",
			logs.KVL(
				"status", rsp.Status,
				"code", rsp.StatusCode,
				"body", buffer.String(),
			))

		// return errors.Errorf("bad http response status %s", sink.String())

		return err
	}

	return nil
}

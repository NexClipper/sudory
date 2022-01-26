package events

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type WebhookListenOption struct {
	Type        string `yaml:"type,omitempty"`
	Method      string `yaml:"method,omitempty"`
	Url         string `yaml:"url,omitempty"`
	ContentType string `yaml:"content-type,omitempty"`
	Timeout     int32  `yaml:"timeout,omitempty"`
}

type WebhookListen struct {
	mux sync.Mutex

	config EventConfig
	opt    WebhookListenOption
}

//check implementation
var _ ListenerContexter = (*WebhookListen)(nil)

func NewWebhookEventListener(config EventConfig, opt WebhookListenOption) ListenerContexter {

	if len(opt.Method) == 0 {
		opt.Method = http.MethodGet //set default Method
	}
	opt.Method = strings.ToUpper(opt.Method) //Method to upper

	if opt.Timeout == 0 {
		opt.Timeout = 15 //set default timeout
	}

	return &WebhookListen{config: config, opt: opt}
}
func (ctx *WebhookListen) Type() string {
	return ListenerTypeWebhook.String()
}

func (ctx *WebhookListen) Name() string {
	return ctx.config.Name
}

func (ctx *WebhookListen) Summary() string {
	return fmt.Sprintf("%s %s", strings.ToUpper(ctx.opt.Method), ctx.opt.Url)
}

func (ctx *WebhookListen) Raise(v interface{}) error {
	ctx.mux.Lock()
	defer ctx.mux.Unlock()

	return webhook_raise(ctx.opt, v)
}

func webhook_raise(opt WebhookListenOption, v interface{}) error {
	buffer, err := serialize_json(v)
	if err != nil {
		return err
	}
	// Create request
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(opt.Timeout)*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, opt.Method, opt.Url, buffer)
	if err != nil {
		return err
	}
	if 0 < len(opt.ContentType) {
		req.Header.Set("Content-Type", opt.ContentType) //set content-type
	}
	// // Add param
	// q := url.Values{}
	// for k, v := range param {
	// 	q.Add(k, v)
	// }
	// req.URL.RawQuery = q.Encode()

	// Request to update services's status to server.
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	// if rsp.StatusCode < http.StatusOK || rsp.StatusCode >= http.StatusBadRequest {
	// 	return fmt.Errorf("HTTP Error, status code: %d", rsp.StatusCode)
	// }

	if rsp.StatusCode/100 != 2 {
		return fmt.Errorf("HTTP Error, status code: %d", rsp.StatusCode)
	}

	return nil
}

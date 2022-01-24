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
	Type    string `yaml:"type,omitempty"`
	Name    string `yaml:"name,omitempty"`
	Pattern string `yaml:"pattern,omitempty"`
	Option  struct {
		Method      string `yaml:"method,omitempty"`
		Url         string `yaml:"url,omitempty"`
		ContentType string `yaml:"content-type,omitempty"`
		Timeout     int32  `yaml:"timeout,omitempty"`
	} `yaml:"option,omitempty"`
}

type WebhookListen struct {
	opt WebhookListenOption
	mux sync.Mutex
}

//check implementation
var _ ListenerContext = (*WebhookListen)(nil)

func NewWebhookEventListener(opt WebhookListenOption) ListenerContext {

	if len(opt.Option.Method) == 0 {
		opt.Option.Method = http.MethodGet //set default Method
	}
	opt.Option.Method = strings.ToUpper(opt.Option.Method) //Method to upper

	if opt.Option.Timeout == 0 {
		opt.Option.Timeout = 15 //set default timeout
	}

	return &WebhookListen{opt: opt}
}
func (me *WebhookListen) Type() string {
	return "webhook"
}
func (me *WebhookListen) Name() string {
	return me.opt.Name
}
func (me *WebhookListen) Pattern() string {
	return me.opt.Pattern
}
func (me *WebhookListen) Dest() string {
	return me.opt.Option.Url
}

func (me *WebhookListen) Raise(v interface{}) error {
	me.mux.Lock()
	defer me.mux.Unlock()
	return me.raise(v)
}

func (me *WebhookListen) raise(v interface{}) error {
	buffer, err := serialize_json(v)
	if err != nil {
		return err
	}
	// Create request
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(me.opt.Option.Timeout)*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, me.opt.Option.Method, me.opt.Option.Url, buffer)
	if err != nil {
		return err
	}
	if 0 < len(me.opt.Option.ContentType) {
		req.Header.Set("Content-Type", me.opt.Option.ContentType) //set content-type
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

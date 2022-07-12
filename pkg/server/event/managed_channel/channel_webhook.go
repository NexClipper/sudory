package managed_channel

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	channelv2 "github.com/NexClipper/sudory/pkg/server/model/channel/v2"
	"github.com/pkg/errors"
)

type ChannelWebhook struct {
	uuid string
	opt  *channelv2.NotifierWebhook_property

	httpclient *http.Client //http.Client
}

func NewChannelWebhook(uuid string, opt channelv2.NotifierWebhook_property) *ChannelWebhook {
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

func (channel ChannelWebhook) OnNotify(factory *MarshalFactory) error {
	const content_type = "application/json"

	httpclient := channel.httpclient

	b, err := factory.Marshal(content_type)
	if err != nil {
		return errors.Wrapf(err, "marshal factory")
	}

	keyval := make(map[string]string)
	for k, v := range channel.opt.RequestHeaders.KeyValue {
		keyval[k] = v
	}
	opt := HttpClient_opt{
		Method:         channel.opt.Method,
		Url:            channel.opt.Url,
		RequestHeaders: keyval,
		RequestTimeout: time.Duration(channel.opt.RequestTimeout) * time.Second,
	}

	if err := HttpReq(&opt, httpclient, content_type, b); err != nil {
		return errors.Wrapf(err, "http request")
	}
	return nil
}

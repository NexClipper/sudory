package managed_channel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	channelv2 "github.com/NexClipper/sudory/pkg/server/model/channel/v3"
	"github.com/pkg/errors"
)

type ChannelSlackhook struct {
	uuid string
	opt  *channelv2.SlackhookConfig

	httpclient *http.Client //http.Client
}

func NewChannelSlackhook(uuid string, opt *channelv2.SlackhookConfig) *ChannelSlackhook {
	client := http.DefaultClient
	if transport, ok := client.Transport.(*http.Transport); ok {
		transport.MaxIdleConnsPerHost = 100 //MaxIdleConnsPerHost
	}
	notifier := &ChannelSlackhook{}
	notifier.uuid = uuid
	notifier.opt = opt
	notifier.httpclient = client

	return notifier
}
func (channel ChannelSlackhook) Type() fmt.Stringer {
	return channel.opt.Type()
}

func (channel ChannelSlackhook) Uuid() string {
	return channel.uuid
}

func (channel ChannelSlackhook) Property() map[string]string {
	return map[string]string{
		"type": channel.opt.Type().String(),
		"uuid": channel.uuid,
		"url":  channel.opt.Url,
	}
}

func (channel *ChannelSlackhook) Close() {

}

func (channel ChannelSlackhook) OnNotify(factory *MarshalFactory) (err error) {
	const content_type = "application/json"
	httpclient := channel.httpclient

	opt := HttpClient_opt{
		Method:         "POST",
		Url:            channel.opt.Url,
		RequestHeaders: map[string]string{},
		RequestTimeout: time.Duration(channel.opt.RequestTimeout) * time.Second,
	}

	var v interface{}
	if factory.Formatter != nil {
		v, err = factory.Formatter.Format(factory.Value)
	} else {
		v = factory.Value
	}

	format_data := map[string]interface{}{"text": ""}
	switch v := v.(type) {
	case string:
		format_data["text"] = v
	default:
		var b []byte
		b, err = json.Marshal(v)
		err = errors.Wrapf(err, "marshal to json")
		format_data["text"] = string(b)
	}
	if err != nil {
		return err
	}

	payload, err := json.Marshal(format_data)
	if err != nil {
		return errors.Wrapf(err, "convert to slack-incoming-webhook format")
	}

	if err := HttpReq(&opt, httpclient, content_type, payload); err != nil {
		return errors.Wrapf(err, "http request")
	}
	return nil
}

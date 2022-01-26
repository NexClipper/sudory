package events

import (
	"encoding/json"
	"strings"
)

type ListenerContexter interface {
	Raise(v interface{}) error //이벤트 발생
	Type() string              //리스너 타입
	Name() string              //이름
	Summary() string           //요약
}

func listenerContextFactory(config EventConfig, listener ListenerConfig) (ListenerContexter, error) {

	listener_ := (map[string]interface{})(listener)
	if listener_ == nil {
		return nil, ErrorInvaildListenerConfig()
	}

	_, ok := listener_["type"]
	if !ok {
		return nil, ErrorInvaildListenerConfig()
	}

	listener_type, ok := (map[string]interface{})(listener_)["type"].(string)
	if !ok {
		return nil, ErrorInvaildListenerConfig()
	}

	switch strings.ToLower(listener_type) {
	case ListenerTypeWebhook.String(): //webhook
		opt := new(WebhookListenOption)
		err := deepcopy(opt, listener)
		if err != nil {
			return nil, err
		}
		return NewWebhookEventListener(config, *opt), nil
	case ListenerTypeFile.String(): //file
		opt := new(FileListenOption)
		err := deepcopy(opt, listener)
		if err != nil {
			return nil, err
		}
		return NewFileListener(config, *opt)
	default:
		return nil, ErrorUndifinedEventListener(listener_type)
	}
}

// deepcopy
//  by json package
func deepcopy(out, in interface{}) error {

	data, err := json.Marshal(in)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, out)
	if err != nil {
		return err
	}
	return nil
}

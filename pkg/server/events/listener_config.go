package events

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/jinzhu/configor"
)

type Config struct {
	Events []Listener
}

type Listener struct {
	Type    string                 `yaml:"type,omitempty"`
	Name    string                 `yaml:"name,omitempty"`
	Pattern string                 `yaml:"pattern,omitempty"`
	Option  map[string]interface{} `yaml:"option,omitempty"`
}

func New(configPath string) (*Config, error) {

	_, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("file open error:'%w'", err)
	}

	cfg := &Config{}
	if err := configor.Load(cfg, configPath); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (me Config) Regist() error {
	Listeners = make(map[string][]ListenerContext)
	for _, it := range me.Events {

		//미리 정규식 패턴을 컴파일 해본다
		_, err := regexp.Compile(it.Pattern)
		if err != nil {
			return fmt.Errorf("failed to regexp compile pattern='%s' error='%w'", it.Pattern, err) //regexp compile error
		}

		key := it.Pattern
		if _, ok := Listeners[key]; !ok {
			Listeners[key] = make([]ListenerContext, 0, 10)
		}
		ctx, err := factory(it)
		if err != nil {
			return err
		}
		Listeners[key] = append(Listeners[key], ctx)
	}

	return nil
}

func factory(event Listener) (ListenerContext, error) {
	switch strings.ToLower(event.Type) {
	case ListenerTypeWebhook.String(): //webhook
		opt := new(WebhookListenOption)
		err := deepcopy(opt, event)
		if err != nil {
			return nil, err
		}
		return NewWebhookEventListener(*opt), nil
	case ListenerTypeFile.String(): //file
		opt := new(FileListenOption)
		err := deepcopy(opt, event)
		if err != nil {
			return nil, err
		}
		return NewFileListener(*opt), nil
	default:
		return nil, ErrorUndifinedEventListener(event.Type)
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

package events

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/jinzhu/configor"
)

type Config struct {
	Events []EventConfig `yaml:"events,omitempty"`
}

type EventConfig struct {
	Name        string           `yaml:"name,omitempty"`
	Pattern     string           `yaml:"pattern,omitempty"`
	BuzyTimeout *int             `yaml:"buzy-timeout,omitempty"` //타임아웃 (n초) default(10)
	Listeners   []ListenerConfig `yaml:"listeners,omitempty"`
}

type ListenerConfig map[string]interface{}

func NewConfig(configPath string) (*Config, error) {

	_, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("file open error:'%w'", err)
	}

	cfg := Config{}
	if err := configor.Load(&cfg, configPath); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (config Config) MakeEventListener() ([]EventContexter, error) {
	result := make([]EventContexter, 0)
	for _, event_config := range config.Events {

		listeners := make([]ListenerContexter, 0)
		for _, listener := range event_config.Listeners {
			ctx, err := listenerContextFactory(event_config, listener)
			if err != nil {
				return nil, err
			}
			listeners = append(listeners, ctx)
		}

		event_ctx := NewEventContext(event_config, listeners...)
		result = append(result, event_ctx)
	}

	//logging
	var w io.Writer = os.Stdout
	tabwrite := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	defer tabwrite.Flush()
	fmt.Fprintln(w, "*** list event-listeners ***")
	fmt.Fprint(tabwrite, tabwriterArgs("event_no", "listener_no", "name", "pattern", "type", "summary"))
	for event_idx, event := range result {
		for listener_idx, listener := range event.ListenerContexts() {
			fmt.Fprint(tabwrite, tabwriterArgs(event_idx, listener_idx, event.Name(), event.Pattern(), listener.Type(), listener.Summary()))
		}
	}
	return result, nil
}

func tabwriterArgs(v ...interface{}) string {
	s := make([]string, len(v))
	for n := range v {
		s[n] = fmt.Sprintf("%+v", v[n])
	}
	return strings.Join(s, "\t") + "\n"
}

func (config Config) Vaild() error {

	//이름 유효 확인
	vaild_name := func(name string) error {
		//이름 길이 확인
		if len(name) == 0 {
			return ErrorInvalidEventNameZeroLength()
		}
		return nil
	}

	//패턴 유효 확인
	vaild_pattern := func(pattern string) error {
		//패턴 길이 확인
		if len(pattern) == 0 {
			return ErrorInvalidEventPatternZeroLength()
		}
		//미리 정규식 패턴을 컴파일 해본다
		_, err := regexp.Compile(pattern)
		if err != nil {
			return ErrorRegexComplileEventPattern(pattern, err) //regexp compile error
		}
		return nil
	}

	vaild_type := func(config map[string]interface{}) error {

		var ok bool = false

		_, ok = config["type"]
		if !ok {
			return ErrorNotFoundListenerType()
		}
		listener_type, ok := config["type"].(string)
		if !ok {
			return ErrorInvaildListenerType(config["type"])
		}
		//리스너 타입 길이 확인
		if len(listener_type) == 0 {
			return ErrorInvaildListenerTypeZeroLength()
		}
		//타입 리스트에서 확인
		for _, typenames := range ListenerTypeNames() {
			ok = ok || strings.Compare(listener_type, typenames) == 0
		}
		if !ok {
			return ErrorUndifinedEventListener(listener_type)
		}

		return nil
	}

	//리스너 유효 확인
	vaild_listener := func(config ListenerConfig) error {

		config_ := (map[string]interface{})(config)
		if config_ == nil {
			return ErrorInvaildListenerConfig()
		}

		if err := vaild_type(config_); err != nil {
			return err
		}

		return nil
	}

	aggregate_listener := func(configs []ListenerConfig, fn func(cfg ListenerConfig) error) error {
		for _, cfg := range configs {
			if err := fn(cfg); err != nil {
				return err
			}
		}
		return nil
	}

	for _, config := range config.Events {
		var err error
		//name
		err = vaild_name(config.Name)
		if err != nil {
			return err
		}

		//pattern
		err = vaild_pattern(config.Pattern)
		if err != nil {
			return err
		}

		//TODO: vaild buzytimeout (warn)

		//listener Type
		err = aggregate_listener(config.Listeners, vaild_listener)
		if err != nil {
			return err
		}

	}
	return nil
}

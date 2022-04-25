package event

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/NexClipper/sudory/pkg/server/macro/reflected"
	"github.com/NexClipper/sudory/pkg/server/macro/tabwriters"
	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var Invoke func(string, ...interface{})

func init() {
	Invoke = func(s string, i ...interface{}) {}
}

func NewEventConfig(configfile string) (*EventConfig, error) {
	cfgevent := &EventConfig{}
	if err := configor.Load(cfgevent, configfile); err != nil {
		//설정파일 로드 실패
		if b, err_ := ioutil.ReadFile(configfile); err_ != nil {
			//설정 파일 읽기 실패
			return nil, errors.Wrapf(err_, "read config file%s",
				logs.KVL(
					"dest-type-name", reflected.TypeName(cfgevent),
					"config-file-path", configfile,
					"err", err.Error(),
				))
		} else {
			return nil, errors.Wrapf(err, "load config file%s",
				logs.KVL(
					"dest-type-name", reflected.TypeName(cfgevent),
					"config-file-path", configfile,
					"yaml", string(b),
				))
		}
	}

	return cfgevent, nil
}

// // RegistNotifier
// func RegistNotifier(sub EventSubscriber) error {
// 	for i := range sub.Config().NotifierConfigs {
// 		cfgnotifier := sub.Config().NotifierConfigs[i]

// 		//new notifier
// 		notifier, err := NotifierFactory(cfgnotifier)
// 		if err != nil {
// 			return errors.Wrapf(err, "notifier factory%s",
// 				logs.KVL(
// 					"config-event", sub.Config(),
// 					"config-notifier", cfgnotifier,
// 				))
// 		}

// 		//등록
// 		notifier.Regist(sub)
// 	}

// 	return nil
// }

func PrintEventConfiguation(w io.Writer, pub EventPublisher) {
	lst := []func(){
		//subscriber
		func() {
			tabwrite := tabwriters.NewWriter(w, 0, 0, 3, ' ', 0)
			defer tabwrite.Flush()

			fmt.Fprintln(w, "subscriber:")

			// var seq int = 0
			for sub := range pub.Subscribers() {
				tabwrite.PrintKeyValue(
					" ", "-", //empty line
					"event-name", sub.Config().Name,
					"update-interval", sub.Config().UpdateInterval.String(),
				)
				// seq++
			}
		},
		//notifier
		func() {
			tabwrite := tabwriters.NewWriter(w, 0, 0, 3, ' ', 0)
			defer tabwrite.Flush()

			fmt.Fprintln(w, "notifier:")

			// var seq int = 0
			for sub := range pub.Subscribers() {
				for notifier := range sub.Notifiers() {
					tabwrite.PrintKeyValue(
						" ", "-", //empty line
						"event-name", sub.Config().Name,
						"notifier-type", notifier.Type(),
						"notifier-property", notifier.PropertyString(),
					)
					// seq++
				}
			}
		},
	}

	fmt.Fprintln(w, "event configuration:")
	for _, fn := range lst {
		fn()
	}
	fmt.Fprintln(w, strings.Repeat("_", 40))

}

func NotifierFactory(cfgnotifire NotifierConfig) (Notifier, error) {
	if _, ok := cfgnotifire["type"]; !ok {
		return nil, errors.Errorf("not found key int listener config%s",
			logs.KVL(
				"key", "type",
			))
	}

	listener_type_name, ok := cfgnotifire["type"].(string)
	if !ok {
		return nil, errors.Errorf("failed to listener type cast to string%s",
			logs.KVL(
				"type", cfgnotifire["type"],
			))
	}

	notifier_type, err := ParseNotifierType(listener_type_name)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to type parse to ListenerType%s",
			logs.KVL(
				"type", listener_type_name,
			))
	}

	var opt interface{}
	switch notifier_type {
	case NotifierTypeConsole: //console
		opt = &ConsoleNotifierConfig{}
	case NotifierTypeWebhook: //webhook
		opt = &WebhookNotifierConfig{
			RequestHeaders: map[string]string{},
		}
	case NotifierTypeFile: //file
		opt = &FileNotifierConfig{}
	case NotifierTypeRabbitMQ: //rabbitmq
		opt = &RabbitMQNotifierConfig{
			MessageHeaders: map[string]interface{}{},
		}
	default:
		return nil, errors.Errorf("unsupported notifier type%s",
			logs.KVL(
				"notifier_type", notifier_type,
			))
	}

	if err := deepcopy(opt, cfgnotifire); err != nil {
		return nil, errors.Wrapf(err, "deepcopy%s",
			logs.KVL(
				"opt", opt,
			))
	}

	var new_notifier Notifier
	switch opt := opt.(type) {
	case *ConsoleNotifierConfig:
		new_notifier = NewConsoleNotifier()
	case *WebhookNotifierConfig:
		new_notifier = NewWebhookNotifier(*opt)
	case *FileNotifierConfig:
		new_notifier, err = NewFileNotifier(*opt)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create file notifier%s",
				logs.KVL(
					"opt", opt,
				))
		}
	case *RabbitMQNotifierConfig:
		new_notifier, err = NewRabbitMqNotifier(*opt)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create rabbitmq notifier%s",
				logs.KVL(
					"opt", opt,
				))
		}
	default:
		return nil, errors.Errorf("unsupported notifier config%s",
			logs.KVL(
				"opt", opt,
			))
	}

	return new_notifier, nil

}

// deepcopy
//  by yaml package
func deepcopy(dest, src interface{}) error {
	data, err := yaml.Marshal(src)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal yaml%s",
			logs.KVL(
				"src-type-name", reflected.TypeName(src),
				"src", src,
			))
	}

	if err := yaml.Unmarshal(data, dest); err != nil {
		return errors.Wrapf(err, "failed to unmarshal yaml%s",
			logs.KVL(
				"dest-type-name", reflected.TypeName(dest),
				"yaml", data,
			))
	}
	return nil
}

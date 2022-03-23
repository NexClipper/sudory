package event_test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/NexClipper/logger"
	"github.com/NexClipper/sudory/pkg/server/event"
	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/jinzhu/configor"
)

const test_config_filename = "event_test_config.yml"

type Configs struct {
	Events []Config `yaml:"events"`
}

type Config struct {
	Name      string           `yaml:"name"`
	Listeners []ListenerConfig `yaml:"listeners"`
}

type ListenerConfig map[string]interface{}

func TestNewPasudoConfig(t *testing.T) {
	var err error

	cfg := Configs{}
	if err = configor.Load(&cfg, test_config_filename); err != nil {
		t.Fatal(err)
	}
}

func TestNewConfig(t *testing.T) {

	os.Setenv("CONFIGOR_DEBUG_MODE", "true")

	//에러 핸들러 등록
	errorHandlers := event.HashsetErrorHandlers{}
	errorHandlers.Add(func(err error) {

		var stack string
		logs.CauseIter(err, func(err error) {
			logs.StackIter(err, func(s string) {
				stack = s
			})
		})

		logger.Error(fmt.Errorf("event notify: %w %s", err,
			logs.KVL(
				"stack", stack,
			)))
	})

	cfgevent, err := event.NewEventConfig(test_config_filename)
	if err != nil {
		t.Fatal(err)
	}

	pub := event.NewEventPublish()

	for i := range cfgevent.EventSubscribeConfigs {
		cfgsub := cfgevent.EventSubscribeConfigs[i]

		sub := event.NewEventSubscribe(cfgsub, errorHandlers)

		if err := event.RegistNotifier(sub); err != nil {
			t.Fatal(err)
		}

		sub.Regist(pub)

		// print(sub)
	}
	event.PrintEventConfiguation(os.Stdout, pub)

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			event.Invoke("test",
				struct {
					I       int
					J       int
					Message string
				}{I: i, J: j, Message: fmt.Sprintf("hello %v %v", i, j)})
		}
	}

	t.Log("closing")
	pub.Close()

	t.Log("done")
}

func TestJsonMashalWithDiffrentTypesSlice(t *testing.T) {

	v := make([]interface{}, 0) //make diffrent types slice
	v = append(v, 1)            //append int type item
	v = append(v, "string")     //append string type item
	v = append(v, struct {
		Foo int
		Bar string
	}{Foo: 123, Bar: "bar"})
	v = append(v, struct {
		Baz    string
		Foobar int
	}{Foobar: 123, Baz: "bar"})

	b, err := json.Marshal(v)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%s", b)

}

package events_test

import (
	"log"
	"os"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/events"
)

const test_config_filename = "test-event-config.yml"

func TestConfigRegist(t *testing.T) {

	_, err := os.Open(test_config_filename)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := events.New(test_config_filename)
	if err != nil {
		log.Fatal(err)
	}

	for n, it := range cfg.Events {
		t.Log(n)
		t.Log(it)
	}

	err = cfg.Regist()
	if err != nil {
		log.Fatal(err)
	}

}

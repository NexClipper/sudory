package service

import (
	"encoding/json"
	"testing"

	v3 "github.com/NexClipper/sudory/pkg/server/model/service/v3"
)

func TestPollingV3(t *testing.T) {

	var v3_polling = v3.HttpRsp_ClientServicePolling{}

	v3_polling.Service = v3.Service{
		Name: "v3service",
	}
	v3_polling.Steps = []v3.ServiceStep{
		{Name: "v3 step 1"},
		{Name: "v3 step 2"},
	}

	var polling HttpRsp_ClientServicePolling

	polling.Version = "v3"

	polling.V3 = v3_polling

	b, err := json.Marshal(polling)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(b))

	err = json.Unmarshal(b, &polling)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(polling.Version)
}

func TestPollingV4(t *testing.T) {

	var v4_polling = HttpRsp_ClientServicePolling_multistep{Name: "v4service"}

	var polling HttpRsp_ClientServicePolling

	polling.Version = "v4"

	polling.V4 = v4_polling

	b, err := json.Marshal(polling)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(b))

	err = json.Unmarshal(b, &polling)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(polling.Version)
}

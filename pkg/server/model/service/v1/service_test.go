package v1

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	cryptov1 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v1"
)

func TestServiceJson(t *testing.T) {

	m := NewService()

	data, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		t.Error(err)
	}
	println(string(data))

	cmd_ := new(Service)

	err = json.Unmarshal(data, cmd_)
	if err != nil {
		t.Error(err)
	}
}

func TestDbSchemaServiceJson(t *testing.T) {

	m := NewService()

	data, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		t.Error(err)
	}
	println(string(data))

	cmd_ := new(Service)

	err = json.Unmarshal(data, cmd_)
	if err != nil {
		t.Error(err)
	}
}

func TestHttpHttpRspService(t *testing.T) {

	m := HttpRspService{Service: NewService()}

	data, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		t.Error(err)
	}
	println(string(data))

	cmd_ := new(Service)

	err = json.Unmarshal(data, cmd_)
	if err != nil {
		t.Error(err)
	}
}

func TestHttpReqClientSideServiceJson(t *testing.T) {

	m := HttpReqService_ClientSide{Service: NewService()}

	data, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		t.Error(err)
	}
	println(string(data))

	cmd_ := new(Service)

	err = json.Unmarshal(data, cmd_)
	if err != nil {
		t.Error(err)
	}
}

const ServiceUuid = "cda6498a235d4f7eae19661d41bc154c"
const ClusterUuid = "cda6498a235d4f7eae19661d41bc154c"

func NewService() Service {

	out := Service{}

	out.Id = 11112222333344445555
	out.Created = newist.Time(time.Now())
	out.Updated = newist.Time(time.Now())
	out.Deleted = nil
	out.Uuid = "00001111222233334444555566667777"
	out.Name = "test-name"
	out.Summary = newist.String("test: ...")
	// out.ApiVersion = newist.String("v1")
	out.ClusterUuid = ClusterUuid
	out.StepCount = newist.Int32(0)
	out.StepPosition = newist.Int32(0)
	// out.Type = newist.Int32(0)
	// out.Epoch = newist.Int32(0)
	// out.Interval = newist.Int32(0)
	out.Status = newist.Int32(0)
	out.Result = func(s string) *cryptov1.String { return (*cryptov1.String)(&s) }("success")

	return out
}

func EmptyServiceService() Service {
	return Service{}
}

func TestJsonBind(t *testing.T) {

	s := `
	{
		"cluster_uuid": "3161ff9142c54c96b274361d62fe6e21",
		"name": "kps helm installation",
		"template_uuid": "20000000000000000000000000000001",
		"summary": "kps helm installation",
		"subscribed_channel": "",
		"steps": [
			{
				"args":{
					"name": "kps",
					"chart_name": "kube-prometheus-stack",
					"repo_url": "https://prometheus-community.github.io/helm-charts",
					"namespace": "monitor",
					"chart_version": "35.0.3",
					"values": {}
				}
			}
		]
	}
	
	`
	service_create := HttpReqService_Create{}

	if err := json.Unmarshal([]byte(s), &service_create); err != nil {

		t.Fatal(err)
	}

}

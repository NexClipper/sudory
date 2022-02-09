package jwt_test

import (
	"log"
	"testing"

	"github.com/NexClipper/sudory/pkg/server/macro/jwt"
)

const _secret_ = "your-256-bit-secret"
const _jwt_ = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjIsIm5hbWUiOiJKb2huIERvZSIsInN1YiI6IjEyMzQ1Njc4OTAifQ.fdOPQ05ZfRhkST2-rIWgUpbqUsVhkkNVNcuG7Ki0s-8"

func TestNew(t *testing.T) {

	payload := map[string]interface{}{
		"sub":  "1234567890",
		"name": "John Doe",
		"iat":  1516239022,
	}

	jwt, err := jwt.New(payload, []byte(_secret_))

	if err != nil {
		log.Fatal(err)
	}

	if jwt != _jwt_ {
		t.Errorf("\nexpect:'%s' \nactual: '%s'", _jwt_, jwt)
	}

	println(jwt)
}
func TestVerify(t *testing.T) {

	err := jwt.Verify(_jwt_, []byte(_secret_))

	if err != nil {
		t.Error(err)
	}

}

func TestGetPayload(t *testing.T) {

	payload, err := jwt.GetPayload(_jwt_)

	if err != nil {
		t.Error(err)
	}

	if payload == nil {
		t.Error("empty")
	}

	if payload["name"].(string) != "John Doe" {
		t.Errorf("expect: '%v', actual: '%v'\n", "John Doe", payload["name"])
	}
	if payload["sub"].(string) != "1234567890" {
		t.Errorf("expect: '%v', actual: '%v'\n", "1234567890", payload["sub"])
	}
	if payload["iat"].(float64) != 1516239022 {
		t.Errorf("expect: '%v', actual: '%v'\n", 1516239022, payload["iat"])
	}

}

func TestBindPayload(t *testing.T) {

	var err error
	var payload interface{}

	check := func() {

		if err != nil {
			t.Error(err)
		}

		if payload == nil {
			t.Error("empty")
		}
	}

	//bind map[string]interface{}
	payload = make(map[string]interface{})
	err = jwt.BindPayload(_jwt_, &payload)
	check()

	map_payload := payload.(map[string]interface{})

	if map_payload["name"].(string) != "John Doe" {
		t.Errorf("expect: '%v', actual: '%v'\n", "John Doe", map_payload["name"])
	}
	if map_payload["sub"].(string) != "1234567890" {
		t.Errorf("expect: '%v', actual: '%v'\n", "1234567890", map_payload["sub"])
	}
	if map_payload["iat"].(float64) != 1516239022 {
		t.Errorf("expect: '%v', actual: '%v'\n", 1516239022, map_payload["iat"])
	}

	//bind payload_value
	payload = &payload_value{}
	err = jwt.BindPayload(_jwt_, &payload)
	check()

	st_payload := payload.(*payload_value)

	if st_payload.Name != "John Doe" {
		t.Errorf("expect: '%v', actual: '%v'\n", "John Doe", st_payload.Name)
	}
	if st_payload.Sub != "1234567890" {
		t.Errorf("expect: '%v', actual: '%v'\n", "1234567890", st_payload.Sub)
	}
	if st_payload.Iat != 1516239022 {
		t.Errorf("expect: '%v', actual: '%v'\n", 1516239022, st_payload.Iat)
	}

	t.Logf("%v", payload)
}

type payload_value struct {
	Name string
	Sub  string
	Iat  int
}

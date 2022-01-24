package macro

import (
	"log"
	"testing"
)

const _secret_ = "your-256-bit-secret"
const _jwt_ = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjIsIm5hbWUiOiJKb2huIERvZSIsInN1YiI6IjEyMzQ1Njc4OTAifQ.fdOPQ05ZfRhkST2-rIWgUpbqUsVhkkNVNcuG7Ki0s-8"

func TestNew(t *testing.T) {

	payload := map[string]interface{}{
		"sub":  "1234567890",
		"name": "John Doe",
		"iat":  1516239022,
	}

	jwt, err := New(payload, []byte(_secret_))

	if err != nil {
		log.Fatal(err)
	}

	if jwt != _jwt_ {
		t.Errorf("\nexpect:'%s' \nactual: '%s'", _jwt_, jwt)
	}

	println(jwt)
}
func TestVerify(t *testing.T) {

	err := Verify(_jwt_, []byte(_secret_))

	if err != nil {
		t.Error(err)
	}

}

func TestGetPayload(t *testing.T) {

	payload, err := GetPayload(_jwt_)

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

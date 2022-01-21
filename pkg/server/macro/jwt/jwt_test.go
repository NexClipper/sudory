package macro

import (
	"testing"
)

func TestMakeJwt(t *testing.T) {
	const secret = "your-256-bit-secret"

	body := map[string]interface{}{
		"sub":  "1234567890",
		"name": "John Doe",
		"iat":  1516239022,
	}

	const jwt_ = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjIsIm5hbWUiOiJKb2huIERvZSIsInN1YiI6IjEyMzQ1Njc4OTAifQ.fdOPQ05ZfRhkST2-rIWgUpbqUsVhkkNVNcuG7Ki0s-8"

	jwt := New(body, secret)

	if jwt != jwt_ {
		t.Errorf("\nexpect:'%s' \nactual: '%s'", jwt_, jwt)
	}

	println(jwt)
}
func TestVerifyJWT(t *testing.T) {
	const secret = "your-256-bit-secret"

	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjIsIm5hbWUiOiJKb2huIERvZSIsInN1YiI6IjEyMzQ1Njc4OTAifQ.fdOPQ05ZfRhkST2-rIWgUpbqUsVhkkNVNcuG7Ki0s-8"

	err := Verify(jwt, secret)

	if err != nil {
		t.Error(err)
	}

}

func TestGetJwtPayload(t *testing.T) {
	const secret = "your-256-bit-secret"

	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE1MTYyMzkwMjIsIm5hbWUiOiJKb2huIERvZSIsInN1YiI6IjEyMzQ1Njc4OTAifQ.fdOPQ05ZfRhkST2-rIWgUpbqUsVhkkNVNcuG7Ki0s-8"

	payload, err := GetToken(jwt, secret)

	if err != nil {
		t.Error(err)
	}

	if payload == nil {
		t.Error("empty")
	}

}

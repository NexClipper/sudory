package tenants_test

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestSha1(t *testing.T) {

	s := ""

	h := sha1.Sum([]byte(s))
	t.Log(h)

	b64 := base64.StdEncoding.EncodeToString(h[0:])
	t.Log(b64)

	x := hex.EncodeToString(h[0:])
	t.Log(x)
}

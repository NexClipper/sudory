package echoutil

import (
	"net/http"
	"strings"
)

func SeHttpHeader(header http.Header, key string, value ...string) {
	for i, value := range value {
		switch i {
		case 0:
			header.Set(key, value)
		default:
			header.Add(key, value)
		}
	}
}

func ParseAuthorizationHeader(header http.Header) (auth_scheme string, auth_value string, ok bool) {
	value := header.Get("Authorization")
	if len(value) == 0 {
		return
	}

	tokens := strings.Split(value, " ")
	ok = len(tokens) == 2
	if !ok {
		return
	}

	auth_scheme = tokens[0]
	auth_value = tokens[1]
	return
}

const (
	AuthSchemaBearer = "Bearer"
)

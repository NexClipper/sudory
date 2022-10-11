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

func GetAuthorizationHeader(header http.Header) string {
	return header.Get(HTTP_HEAD_AUTHORIZATION)
}

func ParseAuthorizationHeader(header http.Header) (auth_scheme string, auth_value string, ok bool) {
	value := header.Get(HTTP_HEAD_AUTHORIZATION)
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
	HTTP_AUTH_SCHEMA_BEARER = "Bearer"
	HTTP_HEAD_AUTHORIZATION = "Authorization"
)

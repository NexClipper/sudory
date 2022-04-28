//go:generate go run github.com/abice/go-enum --file=environment_keys.go --names --nocase
package env

import (
	"strconv"
	"syscall"
	"time"

	"github.com/NexClipper/sudory/pkg/server/macro/newist"
)

/* ENUM(
bearer-token-signature-secret
bearer-token-expiration-time

client-session-signature-secret
client-session-expiration-time

client-config-poll-interval
client-config-loglevel
)
*/
type Env int

func atoi(s string, d int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return d
	}
	return i
}

func getString(key Env, d string) string {
	s, ok := syscall.Getenv(key.String())

	if !ok {
		s = d
	}
	return s
}

// o: o[0]=default, o[1]=min, o[2]=max
func getInt(key Env, o ...int) int {

	var def int
	var min *int
	var max *int
	if 0 < len(o) {
		def = o[0] //default
	}
	if 1 < len(o) {
		min = newist.Int(o[1]) //min
	}
	if 2 < len(o) {
		max = newist.Int(o[2]) //max
	}

	key_ := key.String()
	s, ok := syscall.Getenv(key_)
	if !ok {
		return def
	}
	val := atoi(s, def)
	if min != nil && val < *min {
		return *min
	}
	if max != nil && *max < val {
		return *max
	}

	return val
}

// // ClusterTokenSignatureSecret
// //  클러스터 토큰 시그니처 시크릿
// func ClusterTokenSignatureSecret() string {
// 	return getString(EnvClusterTokenSignatureSecret, "")
// }

// BearerTokenSignatureSecret
//  bearer-토큰 시그니처 시크릿
func BearerTokenSignatureSecret() string {
	return getString(EnvBearerTokenSignatureSecret, "")
}

// EnvBearerTokenExpirationTime
//  bearer-토큰 만료 시간 (month)
func BearerTokenExpirationTime(t time.Time) time.Time {

	month := getInt(EnvBearerTokenExpirationTime, 12, 1)

	if true {
		const day = 1 * 60 * 60 * 24
		//만료일 다음날 0시 정각
		usec := (t.Unix() / day) * day
		t = time.Unix(usec, 0)
		return t.AddDate(0, month, 1)

	} else {
		//만료일 현재 시간
		return t.AddDate(0, month, 0)
	}
}

// ClientSessionSignatureSecret
//  클라이언트 세션 시그니처 시크릿
func ClientSessionSignatureSecret() string {
	return getString(EnvClientSessionSignatureSecret, "")
}

// EnvClientSessionExpirationTime
//  클라이언트 세션 만료 시간 (초)
func ClientSessionExpirationTime(t time.Time) time.Time {
	value := getInt(EnvClientSessionExpirationTime, 60, 1)
	return t.Add(time.Duration(value) * time.Second)
}

// ClientConfigPollInterval
//  클라이언트 폴 주기
func ClientConfigPollInterval() int {
	return getInt(EnvClientConfigPollInterval, 15, 1)
}

// ClientConfigLoglevel
//  클라이언트 로그 레벨
func ClientConfigLoglevel() string {
	return getString(EnvClientConfigLoglevel, "")
}

//go:generate go run github.com/abice/go-enum --file=keys.go --names --nocase
package globvar

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
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
type Key int

type StoreManager struct {
	Store map[Key]func(s string) error
}

func (gvset *StoreManager) Setter(gv Key, fn func(s string) error) {
	if gvset.Store == nil {
		gvset.Store = map[Key]func(s string) error{}
	}
	gvset.Store[gv] = fn
}

func (gvset *StoreManager) Call(gv Key, s string) error {
	fn, ok := gvset.Store[gv]
	if !ok {
		return errors.Errorf("not found setter%v",
			logs.KVL(
				"key", gv,
			))
	}

	return fn(s)
}

var storeManager *StoreManager

func init() {
	storeManager = &StoreManager{}
	for k, v := range defaultValueSet {
		storeManager.Setter(k, v.Setter)
	}
}

// BearerTokenSignatureSecret
//  bearer-토큰 시그니처 시크릿
func BearerTokenSignatureSecret() string {
	return bearerTokenSignatureSecret
}

var bearerTokenSignatureSecret string = ""

func BearerTokenExpirationTime(t time.Time) time.Time {
	if true {
		const day = 1 * 60 * 60 * 24
		//만료일 다음날 0시 정각
		usec := (t.Unix() / day) * day
		t = time.Unix(usec, 0)
		return t.AddDate(0, bearerTokenExpirationTime_month, 1)

	} else {
		//만료일 현재 시간
		return t.AddDate(0, bearerTokenExpirationTime_month, 0)
	}
}

var bearerTokenExpirationTime_month int = 1

// ClientSessionSignatureSecret
//  클라이언트 세션 시그니처 시크릿
func ClientSessionSignatureSecret() string {
	return clientSessionSignatureSecret
}

var clientSessionSignatureSecret string = ""

// EnvClientSessionExpirationTime
//  클라이언트 세션 만료 시간 (초)
func ClientSessionExpirationTime(t time.Time) time.Time {
	return t.Add(time.Duration(clientSessionExpirationTime_sec) * time.Second)
}

var clientSessionExpirationTime_sec int

// ClientConfigPollInterval
//  클라이언트 폴 주기
func ClientConfigPollInterval() int {
	return clientConfigPollInterval_sec
}

var clientConfigPollInterval_sec int = 1

// ClientConfigLoglevel
//  클라이언트 로그 레벨
func ClientConfigLoglevel() string {
	return clientConfigLoglevel
}

var clientConfigLoglevel string = "debug"

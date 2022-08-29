package globvar

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/macro/logs"
	"github.com/pkg/errors"
)

//go:generate go run github.com/abice/go-enum --file=keys.go --names --nocase

/* ENUM(
bearer-token-signature-secret
bearer-token-expiration-time

client-session-signature-secret
client-session-expiration-time

client-config-poll-interval
client-config-loglevel

event-notifier-status-rotate-limit

service-session-signature-secret
service-session-expiration-time
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

type bearerToken struct {
	// bearer-토큰 시그니처 시크릿
	signatureSecret string
	// bearer-토큰 만료 시간
	expirationTime int // (month)
}

func (value bearerToken) SignatureSecret() string {
	return value.signatureSecret
}

func (value bearerToken) ExpirationTime(t time.Time) time.Time {
	if true {
		return TrimDay(t).AddDate(0, value.expirationTime, 1)
	}
	return t.AddDate(0, value.expirationTime, 0)
}

func TrimDay(t time.Time) time.Time {
	const day = 1 * 60 * 60 * 24
	usec := (t.Unix() / day) * day
	return time.Unix(usec, 0)
}

type clientSession struct {
	// 클라이언트 세션 토큰 시그니처 시크릿
	signatureSecret string
	// 클라이언트 세션 만료 시간
	expirationTime int
}

func (value clientSession) SignatureSecret() string {
	return value.signatureSecret
}

func (value clientSession) ExpirationTime(t time.Time) time.Time {
	return t.Add(time.Duration(value.expirationTime) * time.Second)
}

type clientConfig struct {
	// 클라이언트 로그 레벨
	loglevel string
	// 클라이언트 폴 주기
	pollInterval int // (second)
}

func (value clientConfig) PollInterval() int {
	return value.pollInterval
}

func (value clientConfig) Loglevel() string {
	return value.loglevel
}

type event struct {
	// 이벤트 알림 상태 rotate limit
	nofitierStatusRotateLimit uint
}

func (value event) NofitierStatusRotateLimit() uint {
	return value.nofitierStatusRotateLimit
}

type serviceSession struct {
	// 서비스 세션 시그니처 시크릿
	signatureSecret string
	// 서비스 세션 만료 시간 (month)
	expirationTime int // (month)
}

func (value serviceSession) SignatureSecret() string {
	return value.signatureSecret
}

func (value serviceSession) ExpirationTime(t time.Time) time.Time {
	if true {
		return TrimDay(t).AddDate(0, value.expirationTime, 1)
	}
	return t.AddDate(0, value.expirationTime, 0)
}

var (
	BearerToken = bearerToken{
		signatureSecret: "",
		expirationTime:  6, // month(6)
	}

	ClientSession = clientSession{
		signatureSecret: "",
		expirationTime:  60, // second(60)
	}

	ClientConfig = clientConfig{
		loglevel:     "debug",
		pollInterval: 15, // second(15)
	}

	Event = event{
		nofitierStatusRotateLimit: 20,
	}

	ServiceSession = serviceSession{
		signatureSecret: "",
		expirationTime:  12, // month(12)
	}
)

package globvar

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	globvarv2 "github.com/NexClipper/sudory/pkg/server/model/global_variables/v2"
)

type defaultValue struct {
	Uuid    string
	Value   string
	Summary string
	Setter  func(s string) error
}

var defaultValueSet = map[Key]defaultValue{
	KeyBearerTokenSignatureSecret: {Uuid: "e2db6f6b08e94cb58bc6a35e244aaa29", Value: BearerToken.signatureSecret, Summary: "bearer-토큰 시그니처 생성 시크릿", Setter: func(s string) error {
		BearerToken.signatureSecret = s
		return nil
	}}, //(사용안함)
	KeyBearerTokenExpirationTime: {Uuid: "0f5658f37f2b45d881f19c7f56ea2e23", Value: fmt.Sprintf("%v", BearerToken.expirationTime), Summary: "bearer-토큰 만료 시간 (month)", Setter: func(s string) error {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		BearerToken.expirationTime = i
		return nil
	}},

	KeyClientSessionSignatureSecret: {Uuid: "77f7b2aeb0aa4254ad073ae7743291ab", Value: ClientSession.signatureSecret, Summary: "클라이언트 세션 시그니처 생성 시크릿", Setter: func(s string) error {
		ClientSession.signatureSecret = s
		return nil
	}},
	KeyClientSessionExpirationTime: {Uuid: "af9a14a58b254d13ae69c065a27811b6", Value: fmt.Sprintf("%v", ClientSession.expirationTime), Summary: "클라이언트 세션 만료 시간 (초)", Setter: func(s string) error {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		ClientSession.expirationTime = i
		return nil
	}},

	KeyClientConfigPollInterval: {Uuid: "75531e760ee6423cb3a050ddcc83e275", Value: fmt.Sprintf("%v", ClientConfig.pollInterval), Summary: "클라이언트 poll interval (초)", Setter: func(s string) error {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		ClientConfig.pollInterval = i
		return nil
	}},
	KeyClientConfigLoglevel: {Uuid: "4e55651f63814b648f7284ab9113cf67", Value: ClientConfig.loglevel, Summary: "클라이언트 log level ['debug', 'info', 'warn', 'error', 'fatal']", Setter: func(s string) error {
		ClientConfig.loglevel = s
		return nil
	}},
	KeyClientConfigServiceValidTimeLimit: {Uuid: "bc2cd0f95b6d4db68870d30862523a04", Value: fmt.Sprintf("%v", ClientConfig.serviceValidTimeLimit), Summary: "Service Valid Time Limit (minute)", Setter: func(s string) error {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		ClientConfig.serviceValidTimeLimit = i
		return nil
	}},
	KeyEventNotifierStatusRotateLimit: {Uuid: "997c1631c9dd47f9a0c75448fb557ab0", Value: fmt.Sprintf("%v", Event.nofitierStatusRotateLimit), Summary: "이벤트 알림 상태 rotate limit", Setter: func(s string) error {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		if i == 0 {
			i = 20
		} else if i < 0 {
			i = 20
		} else if math.MaxUint8 < i {
			i = math.MaxUint8
		}
		Event.nofitierStatusRotateLimit = uint(i)

		return nil
	}},

	KeyServiceSessionSignatureSecret: {Uuid: "ae806b0bc0ef4d0093a87d2635761aa3", Value: ServiceSession.signatureSecret, Summary: "서비스 세션 시그니처 생성 시크릿", Setter: func(s string) error {
		ServiceSession.signatureSecret = s
		return nil
	}},
	KeyServiceSessionExpirationTime: {Uuid: "bb281e5e882047e1a79dfdd011ee7543", Value: fmt.Sprintf("%v", ServiceSession.expirationTime), Summary: "서비스 세션 만료 시간 (month)", Setter: func(s string) error {
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		ServiceSession.expirationTime = i
		return nil
	}},
}

func GetDefaultValue(key Key) (defaultValue, bool) {
	value, ok := defaultValueSet[key]
	return value, ok
}

func GetDefaultGlobalVariable(key Key, t time.Time) (globvarv2.GlobalVariables, []string, bool) {
	on_dupe_update_columns := []string{
		"summary",
		"value",
		"updated",
	}

	value, ok := GetDefaultValue(key)

	globvar := globvarv2.GlobalVariables{}
	globvar.Uuid = value.Uuid
	globvar.Name = key.String()
	globvar.Summary = *vanilla.NewNullString(fmt.Sprintf("%s (default='%s')", value.Summary, value.Value))
	globvar.Value = *vanilla.NewNullString(value.Value)
	globvar.Created = t
	globvar.Updated = *vanilla.NewNullTime(t)

	return globvar, on_dupe_update_columns, ok
}

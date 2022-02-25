package env

import (
	"fmt"

	"github.com/NexClipper/sudory/pkg/server/macro/newist"
	envv1 "github.com/NexClipper/sudory/pkg/server/model/environment/v1"
)

type Default struct {
	Uuid    string
	Value   string
	Summary string
}

var DefaultEnvironmanets = map[Env]Default{
	// EnvClusterTokenSignatureSecret: {Uuid: "cc6eeb942b9a4a9ca34dc4dfabc54275", Value: "", Summary: "클러스터 토큰 시그니처 생성 시크릿"},

	EnvBearerTokenSignatureSecret: {Uuid: "e2db6f6b08e94cb58bc6a35e244aaa29", Value: "", Summary: "bearer-토큰 시그니처 생성 시크릿"}, //(사용안함)
	EnvBearerTokenExpirationTime:  {Uuid: "0f5658f37f2b45d881f19c7f56ea2e23", Value: "6", Summary: "bearer-토큰 만료 시간 (month)"},

	EnvClientSessionSignatureSecret: {Uuid: "77f7b2aeb0aa4254ad073ae7743291ab", Value: "", Summary: "클라이언트 세션 시그니처 생성 시크릿"},
	EnvClientSessionExpirationTime:  {Uuid: "af9a14a58b254d13ae69c065a27811b6", Value: "60", Summary: "클라이언트 세션 만료 시간 (초)"},

	EnvClientConfigPollInterval: {Uuid: "75531e760ee6423cb3a050ddcc83e275", Value: "15", Summary: "클라이언트 poll interval (초)"},
	EnvClientConfigLoglevel:     {Uuid: "4e55651f63814b648f7284ab9113cf67", Value: "debug", Summary: "클라이언트 log level ['debug', 'info', 'warn', 'error', 'fatal']"},
}

func Convert(key Env, value Default) envv1.Environment {

	const ApiVersion = "v1"

	out := envv1.Environment{}

	out.Uuid = value.Uuid
	out.ApiVersion = newist.String(ApiVersion)
	out.Name = newist.String(key.String())
	out.Summary = newist.String(fmt.Sprintf("%s default='%s'", value.Summary, value.Value))
	out.Value = newist.String(value.Value)

	return out
}

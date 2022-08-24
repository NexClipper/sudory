package v3

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	crypto "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v2"
)

type serviceResultTableName struct{}

func (serviceResultTableName) TableName() string {
	return "service_result"
}

type ServiceResult_create struct {
	serviceResultTableName `json:"-"`

	PK         pkService           `json:",inline"`
	ResultType ResultType          `column:"result_type" json:"result_type,omitempty"`
	Result     crypto.CryptoString `column:"result"      json:"result,omitempty"`
	Created    time.Time           `column:"created"     json:"created,omitempty"`
}

type ServiceResult_update struct {
	serviceResultTableName `json:"-"`

	ResultType ResultType
	Result     crypto.CryptoString
	Updated    vanilla.NullTime
}

type ServiceResult struct {
	serviceResultTableName `json:"-"`

	PK         pkService           `json:",inline"`
	ResultType ResultType          `column:"result_type" json:"result_type,omitempty"`
	Result     crypto.CryptoString `column:"result"      json:"result,omitempty"`
	Created    time.Time           `column:"created"     json:"created,omitempty"`
	Updated    vanilla.NullTime    `column:"updated"     json:"updated,omitempty"`
}

func (ServiceResult) TableName() string {
	return "service_result"
}

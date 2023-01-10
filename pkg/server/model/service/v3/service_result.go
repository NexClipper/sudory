package service

import (
	crypto "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v2"
)

type serviceResultTableName struct{}

func (serviceResultTableName) TableName() string {
	return "service_result"
}

// type ServiceResult_create struct {
// 	serviceResultTableName `json:"-"`

// 	pkService      `json:",inline"`
// 	ResultSaveType ResultSaveType      `column:"result_type" json:"result_type,omitempty"`
// 	Result         crypto.CryptoString `column:"result"      json:"result,omitempty"`
// }

// type ServiceResult_update struct {
// 	serviceResultTableName `json:"-"`

// 	ResultSaveType ResultSaveType
// 	Result         crypto.CryptoString
// 	Timestamp      time.Time
// }

type ServiceResult struct {
	serviceResultTableName `json:"-"`

	pkService      `json:",inline"`
	ResultSaveType ResultSaveType      `column:"result_type" json:"result_type,omitempty"`
	Result         crypto.CryptoString `column:"result"      json:"result,omitempty"`
}

func (ServiceResult) TableName() string {
	return "service_result"
}

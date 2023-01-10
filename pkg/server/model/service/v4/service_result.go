package service

import (
	"time"

	crypto "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v2"
)

type serviceResultTableName struct{}

func (serviceResultTableName) TableName() string {
	return "service_result_v2"
}

type pkServiceResult struct {
	PartitionDate time.Time `column:"pdate"        json:"partition_date"` // pk date
	ClusterUuid   string    `column:"cluster_uuid" json:"cluster_uuid"`   // pk char(32) cluster.uuid
	Uuid          string    `column:"uuid"         json:"uuid"`           // pk char(32) service.uuid
	// Created       time.Time `column:"created"      json:"created"`        // pk datetime(6)
}

// type ServiceResult_create struct {
// 	serviceResultTableName `json:"-"`

// 	pkServiceResult `json:",inline"`
// 	ResultSaveType  ResultSaveType      `column:"result_type" json:"result_type,omitempty"`
// 	Result          crypto.CryptoString `column:"result"      json:"result,omitempty"`
// }

// type ServiceResult_update struct {
// 	serviceResultTableName `json:"-"`

// 	ResultSaveType ResultSaveType
// 	Result         crypto.CryptoString
// 	Timestamp      time.Time
// }

type ServiceResult struct {
	serviceResultTableName `json:"-"`

	pkServiceResult `json:",inline"`
	ResultSaveType  ResultSaveType      `column:"result_type" json:"result_type,omitempty"`
	Result          crypto.CryptoString `column:"result"      json:"result,omitempty"`
	Created         time.Time           `column:"created"     json:"created"`
}

// func (ServiceResult) TableName() string {
// 	return "service_result"
// }

package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

//Cluster Property
type ClusterProperty struct {
	PollingOption map[string]interface{} `json:"polling_option,omitempty" xorm:"text null 'polling_option' comment('polling option')"`
}

//Cluster
type Cluster struct {
	metav1.UuidMeta  `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	ClusterProperty  `json:",inline" xorm:"extends"` //inline property
}

//DATABASE SCHEMA: Cluster
type DbSchema struct {
	metav1.DbMeta `xorm:"extends"`
	Cluster       `xorm:"extends"`
}

func (DbSchema) TableName() string {
	return "cluster"
}

//HTTP REQUEST BODY: Cluster
type HttpReqCluster struct {
	Cluster `json:",inline"`
}

//HTTP RESPONSE BODY: Cluster
type HttpRspCluster struct {
	DbSchema `json:",inline"`
}

// //변환 DbSchema -> Cluster
// func TransFormDbSchema(s []DbSchema) []Cluster {
// 	var out = make([]Cluster, len(s))
// 	for n, it := range s {
// 		out[n] = it.Cluster
// 	}
// 	return out
// }

//변환 Cluster -> HttpRsp
func TransToHttpRsp(s []DbSchema) []HttpRspCluster {
	var out = make([]HttpRspCluster, len(s))
	for n, it := range s {
		out[n].DbSchema = it
	}
	return out
}

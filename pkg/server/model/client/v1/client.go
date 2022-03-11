package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

//Client Property
type ClientProperty struct {
	//ClusterUuid
	ClusterUuid string `json:"cluster_uuid" xorm:"char(32) notnull index 'cluster_uuid' comment('cluster_uuid')"`
}

//Client
type Client struct {
	metav1.UuidMeta  `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	ClientProperty   `json:",inline" xorm:"extends"` //inline property
}

//DATABASE SCHEMA: Client
type DbSchema struct {
	metav1.DbMeta `xorm:"extends"`
	Client        `xorm:"extends"`
}

func (DbSchema) TableName() string {
	return "client"
}

//HTTP REQUEST BODY: Client
type HttpReqClient struct {
	Client `json:",inline"`
}

//HTTP RESPONSE BODY: Client
type HttpRspClient struct {
	DbSchema `json:",inline"`
}

//변환 DbSchema -> Client
func TransFormDbSchema(s []DbSchema) []Client {
	var out = make([]Client, len(s))
	for n, it := range s {
		out[n] = it.Client
	}
	return out
}

//변환 Client -> HttpRsp
func TransToHttpRsp(s []DbSchema) []HttpRspClient {
	var out = make([]HttpRspClient, len(s))
	for n, it := range s {
		out[n].DbSchema = it
	}
	return out
}

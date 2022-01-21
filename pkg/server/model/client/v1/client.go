package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/model/orm"
)

//Client Property
type ClientProperty struct {
	//ClusterUuid
	ClusterUuid string `json:"cluster_uuid" xorm:"char(32) notnull index 'cluster_uuid' comment('cluster_uuid')"`
}

//Client
type Client struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	ClientProperty   `json:",inline" xorm:"extends"` //inline property
}

//DATABASE SCHEMA: Client
type DbSchemaClient struct {
	metav1.DbMeta `xorm:"extends"`
	Client        `xorm:"extends"`
}

var _ orm.TableName = (*DbSchemaClient)(nil)

func (DbSchemaClient) TableName() string {
	return "client"
}

//HTTP REQUEST BODY: Client
type HttpReqClient struct {
	Client `json:",inline"`
}

//HTTP RESPONSE BODY: Client
type HttpRspClient struct {
	Client `json:",inline"`
}

//변환 DbSchema -> Client
func TransFormDbSchema(s []DbSchemaClient) []Client {
	var out = make([]Client, len(s))
	for n, it := range s {
		out[n] = it.Client
	}
	return out
}

//변환 Client -> HttpRsp
func TransToHttpRsp(s []Client) []HttpRspClient {
	var out = make([]HttpRspClient, len(s))
	for n, it := range s {
		out[n].Client = it
	}
	return out
}

package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

//Client Property
type ClientProperty struct {
	ClusterUuid string `json:"cluster_uuid" xorm:"'cluster_uuid' char(32) notnull index comment('cluster_uuid')"`
}

//Client
type Client struct {
	metav1.DbMeta    `json:",inline" xorm:"extends"` //inline dbmeta
	metav1.UuidMeta  `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	ClientProperty   `json:",inline" xorm:"extends"` //inline property
}

func (Client) TableName() string {
	return "client"
}

package v1

import (
	"time"

	cryptov1 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v1"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

//ClusterToken Property
type ClusterTokenProperty struct {
	ClusterUuid    string          `json:"cluster_uuid"    xorm:"'cluster_uuid'    char(32)     notnull index  comment('cluster_uuid')"`
	Token          cryptov1.String `json:"token"           xorm:"'token'           varchar(255) notnull unique comment('token')"`
	IssuedAtTime   time.Time       `json:"issued_at_time"  xorm:"'issued_at_time'  varchar(255) notnull        comment('issued at time')"`
	ExpirationTime time.Time       `json:"expiration_time" xorm:"'expiration_time' varchar(255) notnull        comment('expiration time')"`
}

//DATABASE SCHEMA: ClusterToken
type ClusterToken struct {
	metav1.DbMeta        `json:",inline" xorm:"extends"`
	metav1.UuidMeta      `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta     `json:",inline" xorm:"extends"` //inline labelmeta
	ClusterTokenProperty `json:",inline" xorm:"extends"` //inline property
}

func (ClusterToken) TableName() string {
	return "cluster_token"
}

type HttpReqClusterToken_Create struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	ClusterUuid      string                          `json:"cluster_uuid" `
}

type HttpReqClusterToken_UpdateLabel struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
}

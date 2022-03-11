package v1

import (
	"time"

	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/model/orm"
)

//Token Property
type TokenProperty struct {
	UserKind       string    `json:"user_kind" xorm:"varchar(255) notnull index 'user_kind' comment('user_kind')"`
	UserUuid       string    `json:"user_uuid" xorm:"char(32) notnull index 'user_uuid' comment('user_uuid')"`
	Token          string    `json:"token" xorm:"varchar(255) notnull unique 'token' comment('token')"`
	IssuedAtTime   time.Time `json:"issued_at_time" xorm:"varchar(255) notnull 'issued_at_time' comment('issued at time')"`
	ExpirationTime time.Time `json:"expiration_time" xorm:"varchar(255) notnull 'expiration_time' comment('expiration time')"`
}

//Token
type Token struct {
	metav1.UuidMeta  `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	TokenProperty    `json:",inline" xorm:"extends"` //inline property
}

//DATABASE SCHEMA: Token
type DbSchema struct {
	metav1.DbMeta `xorm:"extends"`
	Token         `xorm:"extends"`
}

var _ orm.TableName = (*DbSchema)(nil)

func (DbSchema) TableName() string {
	return "token"
}

//HTTP REQUEST BODY: Token
type HttpReqToken struct {
	Token `json:",inline"`
}

//HTTP RESPONSE BODY: Token
type HttpRspToken struct {
	DbSchema `json:",inline"`
}

//변환 Token -> HttpRsp
func TransToHttpRsp(s []DbSchema) []HttpRspToken {
	var out = make([]HttpRspToken, len(s))
	for n, it := range s {
		out[n].DbSchema = it
	}
	return out
}

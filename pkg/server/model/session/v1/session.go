package v1

import (
	"time"

	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
	"github.com/NexClipper/sudory/pkg/server/model/orm"
)

//Session Property
type SessionProperty struct {
	UserKind       string    `json:"user_kind" xorm:"varchar(255) notnull index 'user_kind' comment('user_kind')"`
	UserUuid       string    `json:"user_uuid" xorm:"char(32) notnull index 'user_uuid' comment('user_uuid')"`
	Token          string    `json:"token" xorm:"text notnull 'token' comment('token')"`
	IssuedAtTime   time.Time `json:"issued_at_time" xorm:"varchar(255) notnull 'issued_at_time' comment('issued at time')"`
	ExpirationTime time.Time `json:"expiration_time" xorm:"varchar(255) notnull 'expiration_time' comment('expiration time')"`
}

//Session
type Session struct {
	metav1.UuidMeta `json:",inline" xorm:"extends"` //inline uuidmeta
	SessionProperty `json:",inline" xorm:"extends"` //inline property
}

//DATABASE SCHEMA: Session
type DbSchema struct {
	metav1.DbMeta `xorm:"extends"`
	Session       `xorm:"extends"`
}

var _ orm.TableName = (*DbSchema)(nil)

func (DbSchema) TableName() string {
	return "session"
}

//HTTP REQUEST BODY: Session
type HttpReqSession struct {
	Session `json:",inline"`
}

//HTTP RESPONSE BODY: Session
type HttpRspSession struct {
	DbSchema `json:",inline"`
}

//변환 Session -> HttpRsp
func TransToHttpRsp(s []DbSchema) []HttpRspSession {
	var out = make([]HttpRspSession, len(s))
	for n, it := range s {
		out[n].DbSchema = it
	}
	return out
}

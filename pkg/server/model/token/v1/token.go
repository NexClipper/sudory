//go:generate go-enum --file=token.go --names --nocase=true
package v1

import (
	"time"

	cryptov1 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v1"
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

/* ENUM (
	cluster
)
*/
type TokenUserKind int32

//Token Property
type TokenProperty struct {
	UserKind       string          `json:"user_kind"       xorm:"'user_kind'       varchar(255) notnull index  comment('user_kind')"`
	UserUuid       string          `json:"user_uuid"       xorm:"'user_uuid'       char(32)     notnull index  comment('user_uuid')"`
	Token          cryptov1.String `json:"token"           xorm:"'token'           varchar(255) notnull unique comment('token')"`
	IssuedAtTime   time.Time       `json:"issued_at_time"  xorm:"'issued_at_time'  varchar(255) notnull        comment('issued at time')"`
	ExpirationTime time.Time       `json:"expiration_time" xorm:"'expiration_time' varchar(255) notnull        comment('expiration time')"`
}

//DATABASE SCHEMA: Token
type Token struct {
	metav1.DbMeta    `json:",inline" xorm:"extends"`
	metav1.UuidMeta  `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	TokenProperty    `json:",inline" xorm:"extends"` //inline property
}

func (Token) TableName() string {
	return "token"
}

type HttpReqToken_CreateClusterToken struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	UserUuid         string                          `json:"user_uuid" `
}

type HttpReqToken_UpdateLabel struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
}

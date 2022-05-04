package v1

import (
	"time"

	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

//Session Property
type SessionProperty struct {
	ClusterUuid    string    `json:"cluster_uuid"    xorm:"'cluster_uuid'    char(32)     notnull index comment('user_uuid')"`
	Token          string    `json:"token"           xorm:"'token'           text         notnull       comment('token')"`
	IssuedAtTime   time.Time `json:"issued_at_time"  xorm:"'issued_at_time'  varchar(255) notnull       comment('issued at time')"`
	ExpirationTime time.Time `json:"expiration_time" xorm:"'expiration_time' varchar(255) notnull       comment('expiration time')"`
}

//DATABASE SCHEMA: Session
type Session struct {
	metav1.DbMeta   `json:",inline" xorm:"extends"`
	metav1.UuidMeta `json:",inline" xorm:"extends"` //inline uuidmeta
	SessionProperty `json:",inline" xorm:"extends"` //inline property
}

func (Session) TableName() string {
	return "session"
}

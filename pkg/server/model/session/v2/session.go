package v2

import (
	"time"

	noxorm "github.com/NexClipper/sudory/pkg/server/model/noxorm/v2"
)

type Session_essential struct {
	Name           string            `column:"name"            json:"name,omitempty"`
	Summary        noxorm.NullString `column:"summary"         json:"summary,omitempty"`
	ClusterUuid    string            `column:"cluster_uuid"    json:"cluster_uuid"`
	Token          string            `column:"token"           json:"token"`
	IssuedAtTime   noxorm.NullTime   `column:"issued_at_time"  json:"issued_at_time"`
	ExpirationTime noxorm.NullTime   `column:"expiration_time" json:"expiration_time"`
}

func (Session_essential) TableName() string {
	return "session"
}

type Session struct {
	ID   int64  `column:"id"   json:"id,omitempty"`   // pk
	Uuid string `column:"uuid" json:"uuid,omitempty"` // uuid

	Session_essential `json:",inline"` //inline property

	Created time.Time       `column:"created" json:"created,omitempty"`
	Updated noxorm.NullTime `column:"updated" json:"updated,omitempty"`
	Deleted noxorm.NullTime `column:"deleted" json:"deleted,omitempty"`
}

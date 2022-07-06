package v2

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type Session_essential struct {
	Name           string             `column:"name"            json:"name,omitempty"`
	Summary        vanilla.NullString `column:"summary"         json:"summary,omitempty"`
	ClusterUuid    string             `column:"cluster_uuid"    json:"cluster_uuid"`
	Token          string             `column:"token"           json:"token"`
	IssuedAtTime   vanilla.NullTime   `column:"issued_at_time"  json:"issued_at_time"`
	ExpirationTime vanilla.NullTime   `column:"expiration_time" json:"expiration_time"`
}

func (Session_essential) TableName() string {
	return "session"
}

type Session struct {
	ID   int64  `column:"id"   json:"id,omitempty"`   // pk
	Uuid string `column:"uuid" json:"uuid,omitempty"` // uuid

	Session_essential `json:",inline"` //inline property

	Created time.Time        `column:"created" json:"created,omitempty"`
	Updated vanilla.NullTime `column:"updated" json:"updated,omitempty"`
	Deleted vanilla.NullTime `column:"deleted" json:"deleted,omitempty"`
}

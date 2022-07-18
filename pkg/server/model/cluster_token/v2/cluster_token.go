package v2

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	cryptov2 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v2"
)

type ClusterToken_essential struct {
	Name        string             `column:"name"         json:"name,omitempty"`
	Summary     vanilla.NullString `column:"summary"      json:"summary,omitempty"      swaggertype:"string"`
	ClusterUuid string             `column:"cluster_uuid" json:"cluster_uuid,omitempty"`
}

//ClusterToken Property
type ClusterToken_property struct {
	ClusterToken_essential `json:",inline"`

	Token          cryptov2.CryptoString `column:"token"           json:"token,omitempty"`
	IssuedAtTime   time.Time             `column:"issued_at_time"  json:"issued_at_time,omitempty"`
	ExpirationTime time.Time             `column:"expiration_time" json:"expiration_time,omitempty"`
}

func (ClusterToken_property) TableName() string {
	return "cluster_token"
}

//DATABASE SCHEMA: ClusterToken
type ClusterToken struct {
	ClusterToken_property `json:",inline"` //inline property

	ID      int64            `column:"id"      json:"id,omitempty"`   // pk
	Uuid    string           `column:"uuid"    json:"uuid,omitempty"` // uuid
	Created time.Time        `column:"created" json:"created,omitempty"`
	Updated vanilla.NullTime `column:"updated" json:"updated,omitempty" swaggertype:"string"`
	Deleted vanilla.NullTime `column:"deleted" json:"deleted,omitempty" swaggertype:"string"`
}

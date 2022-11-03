package cluster_token

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
	cryptov2 "github.com/NexClipper/sudory/pkg/server/model/default_crypto_types/v2"
)

type ClusterToken_create struct {
	Uuid        string             `json:"uuid,omitempty"` //optional
	Name        string             `json:"name,omitempty"`
	Summary     vanilla.NullString `json:"summary,omitempty"      swaggertype:"string"`
	ClusterUuid string             `json:"cluster_uuid,omitempty"`
}

//ClusterToken Property
type ClusterToken_update struct {
	Name    string             `json:"name,omitempty"`
	Summary vanilla.NullString `json:"summary,omitempty" swaggertype:"string"`
}

// ClusterToken
type ClusterToken struct {
	ID             int64                 `column:"id"              json:"id,omitempty"`   // pk
	Uuid           string                `column:"uuid"            json:"uuid,omitempty"` // uuid
	Name           string                `column:"name"            json:"name,omitempty"`
	Summary        vanilla.NullString    `column:"summary"         json:"summary,omitempty"         swaggertype:"string"`
	ClusterUuid    string                `column:"cluster_uuid"    json:"cluster_uuid,omitempty"`
	Token          cryptov2.CryptoString `column:"token"           json:"token,omitempty"`
	IssuedAtTime   time.Time             `column:"issued_at_time"  json:"issued_at_time,omitempty"`
	ExpirationTime time.Time             `column:"expiration_time" json:"expiration_time,omitempty"`
	Created        time.Time             `column:"created"         json:"created,omitempty"`
	Updated        vanilla.NullTime      `column:"updated"         json:"updated,omitempty"         swaggertype:"string"`
	Deleted        vanilla.NullTime      `column:"deleted"         json:"deleted,omitempty"         swaggertype:"string"`
}

func (ClusterToken) TableName() string {
	return "cluster_token"
}

package v3

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

// type Session_essential struct {
// 	ClusterUuid    string           `column:"cluster_uuid"    json:"cluster_uuid"`
// 	Token          string           `column:"token"           json:"token"`
// 	IssuedAtTime   vanilla.NullTime `column:"issued_at_time"  json:"issued_at_time"    swaggertype:"object"`
// 	ExpirationTime vanilla.NullTime `column:"expiration_time" json:"expiration_time"   swaggertype:"object"`
// }

type Session struct {
	ID             int64            `column:"id"              json:"id,omitempty"`   // pk
	Uuid           string           `column:"uuid"            json:"uuid,omitempty"` // uuid
	ClusterUuid    string           `column:"cluster_uuid"    json:"cluster_uuid"`
	Token          string           `column:"token"           json:"token"`
	IssuedAtTime   vanilla.NullTime `column:"issued_at_time"  json:"issued_at_time"    swaggertype:"object"`
	ExpirationTime vanilla.NullTime `column:"expiration_time" json:"expiration_time"   swaggertype:"object"`
	Created        time.Time        `column:"created"         json:"created,omitempty"`
	Updated        vanilla.NullTime `column:"updated"         json:"updated,omitempty" swaggertype:"object"`
	Deleted        vanilla.NullTime `column:"deleted"         json:"deleted,omitempty" swaggertype:"object"`
}

func (Session) TableName() string {
	return "session"
}

package v3

import (
	"fmt"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type Tenant struct {
	ID      int64              `column:"id"      json:"id,omitempty"`      // pk
	Hash    string             `column:"hash"    json:"hash,omitempty"`    // key pattern->sha1+hex_string
	Pattern string             `column:"pattern" json:"pattern,omitempty"` // tenant pattern
	Name    string             `column:"name"    json:"name,omitempty"`
	Summary vanilla.NullString `column:"summary" json:"summary,omitempty" swaggertype:"string"`
	Created time.Time          `column:"created" json:"created,omitempty"`
	Updated vanilla.NullTime   `column:"updated" json:"updated,omitempty" swaggertype:"string"`
	Deleted vanilla.NullTime   `column:"deleted" json:"deleted,omitempty" swaggertype:"string"`
}

func NewTenant(hash, pattern string, t time.Time) *Tenant {
	tenant := new(Tenant)
	tenant.Hash = hash
	tenant.Pattern = pattern
	tenant.Name = fmt.Sprintf("tenant.%v", pattern)
	tenant.Summary = *vanilla.NewNullString(fmt.Sprintf("hash='%v' pattern='%v'", tenant.Hash, pattern))
	tenant.Created = t
	tenant.Updated = *vanilla.NewNullTime(t)
	return tenant
}

func (Tenant) TableName() string {
	return "tenant"
}

// TenantClusters
// Tenant(1): Cluster(1)
type TenantClusters struct {
	ClusterId int64 `column:"cluster_id" json:"cluster_id,omitempty"` // pk FK(cluster.id)
	TenantId  int64 `column:"tenant_id"  json:"tenant_id,omitempty"`  //    FK(tenant.id)
}

func (tc TenantClusters) TableName() string {
	return "tenant_clusters"
}

// func (TenantClusters) TableNameWithTenant(tenant_hash string, cluster_uuid string) string {
// 	return fmt.Sprintf("( SELECT tenant_clusters.* FROM tenant_clusters, tenant WHERE tenant_clusters.tenant_id = tenant.id AND tenant.deleted IS NULL AND tenant.hash = %v ) x",
// 		tenant_hash,
// 	)
// }

// TenantChannels
// Tenant(1): Channel(1)
type TenantChannels struct {
	ChannelUuid string `column:"channel_uuid" json:"channel_uuid,omitempty"` // pk FK(managed_channel.uuid)
	TenantId    int64  `column:"tenant_id"    json:"tenant_id,omitempty"`    //    FK(tenant.id)
}

func (TenantChannels) TableName() string {
	return "tenant_channels"
}

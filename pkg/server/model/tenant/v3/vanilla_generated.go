// Code generated by Ice-cream-maker DO NOT EDIT.
package v3
 
func (Tenant) ColumnNames() []string {
	return []string{
 		"id",
 		"hash",
 		"pattern",
 		"name",
 		"summary",
 		"created",
 		"updated",
 		"deleted",
	}
}
 
func (TenantClusters) ColumnNames() []string {
	return []string{
 		"cluster_id",
 		"tenant_id",
	}
}
 
func (TenantChannels) ColumnNames() []string {
	return []string{
 		"channel_uuid",
 		"tenant_id",
	}
}
 
func (row Tenant) Values() []interface{} {
	return []interface{}{
		row.ID,
		row.Hash,
		row.Pattern,
		row.Name,
		row.Summary,
		row.Created,
		row.Updated,
		row.Deleted,
	}
}
 
func (row TenantClusters) Values() []interface{} {
	return []interface{}{
		row.ClusterId,
		row.TenantId,
	}
}
 
func (row TenantChannels) Values() []interface{} {
	return []interface{}{
		row.ChannelUuid,
		row.TenantId,
	}
}

type Scanner interface {
	Scan(dest ...interface{}) error
}
 
func (row *Tenant) Scan(scanner Scanner) error {
	return scanner.Scan(
		&row.ID,
		&row.Hash,
		&row.Pattern,
		&row.Name,
		&row.Summary,
		&row.Created,
		&row.Updated,
		&row.Deleted,
	)
}
 
func (row *TenantClusters) Scan(scanner Scanner) error {
	return scanner.Scan(
		&row.ClusterId,
		&row.TenantId,
	)
}
 
func (row *TenantChannels) Scan(scanner Scanner) error {
	return scanner.Scan(
		&row.ChannelUuid,
		&row.TenantId,
	)
}
 
func (row *Tenant) Ptrs() []interface{} {
	return []interface{}{
		&row.ID,
		&row.Hash,
		&row.Pattern,
		&row.Name,
		&row.Summary,
		&row.Created,
		&row.Updated,
		&row.Deleted,
	}
}
 
func (row *TenantClusters) Ptrs() []interface{} {
	return []interface{}{
		&row.ClusterId,
		&row.TenantId,
	}
}
 
func (row *TenantChannels) Ptrs() []interface{} {
	return []interface{}{
		&row.ChannelUuid,
		&row.TenantId,
	}
}

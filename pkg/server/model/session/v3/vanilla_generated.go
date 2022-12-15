// Code generated by Ice-cream-maker DO NOT EDIT.
package sessions
 
func (Session) ColumnNames() []string {
	return []string{
 		"id",
 		"uuid",
 		"cluster_uuid",
 		"token",
 		"issued_at_time",
 		"expiration_time",
 		"created",
 		"updated",
 		"deleted",
	}
}
 
func (Session) ColumnNamesWithAlias() []string {
	return []string{
 		"id",
 		"uuid",
 		"cluster_uuid",
 		"token",
 		"issued_at_time",
 		"expiration_time",
 		"created",
 		"updated",
 		"deleted",
	}
}
 
func (row Session) Values() []interface{} {
	return []interface{}{
		row.ID,
		row.Uuid,
		row.ClusterUuid,
		row.Token,
		row.IssuedAtTime,
		row.ExpirationTime,
		row.Created,
		row.Updated,
		row.Deleted,
	}
}

type Scanner interface {
	Scan(dest ...interface{}) error
}
 
func (row *Session) Scan(scanner Scanner) error {
	return scanner.Scan(
		&row.ID,
		&row.Uuid,
		&row.ClusterUuid,
		&row.Token,
		&row.IssuedAtTime,
		&row.ExpirationTime,
		&row.Created,
		&row.Updated,
		&row.Deleted,
	)
}
 
func (row *Session) Ptrs() []interface{} {
	return []interface{}{
		&row.ID,
		&row.Uuid,
		&row.ClusterUuid,
		&row.Token,
		&row.IssuedAtTime,
		&row.ExpirationTime,
		&row.Created,
		&row.Updated,
		&row.Deleted,
	}
}

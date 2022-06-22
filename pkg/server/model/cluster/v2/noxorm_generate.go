package v2

func (Cluster_essential) ColumnNames() []string {
	return []string{"name", "summary", "polling_option", "polling_limit"}
}

func (Cluster) ColumnNames() []string {
	return []string{"id", "uuid", "name", "summary", "polling_option", "polling_limit", "created", "updated", "deleted"}
}

func (row Cluster_essential) Values() []interface{} {
	return []interface{}{row.Name, row.Summary, row.PollingOption, row.PoliingLimit}
}

func (row Cluster) Values() []interface{} {
	return []interface{}{row.ID, row.Uuid, row.Cluster_essential.Name, row.Cluster_essential.Summary, row.Cluster_essential.PollingOption, row.Cluster_essential.PoliingLimit, row.Created, row.Updated, row.Deleted}
}

type Scanner interface {
	Scan(dest ...interface{}) error
}

func (row *Cluster_essential) Scan(scanner Scanner) error {
	return scanner.Scan(&row.Name, &row.Summary, &row.PollingOption, &row.PoliingLimit)
}

func (row *Cluster) Scan(scanner Scanner) error {
	return scanner.Scan(&row.ID, &row.Uuid, &row.Cluster_essential.Name, &row.Cluster_essential.Summary, &row.Cluster_essential.PollingOption, &row.Cluster_essential.PoliingLimit, &row.Created, &row.Updated, &row.Deleted)
}

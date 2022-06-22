package v2

func (Template_essential) ColumnNames() []string {
	return []string{"name", "summary", "origin"}
}

func (Template) ColumnNames() []string {
	return []string{"uuid", "name", "summary", "origin", "created", "updated", "deleted"}
}

func (TemplateCommand_essential) ColumnNames() []string {
	return []string{"name", "summary", "template_uuid", "sequence", "method", "args", "result_filter"}
}

func (TemplateCommand) ColumnNames() []string {
	return []string{"uuid", "name", "summary", "template_uuid", "sequence", "method", "args", "result_filter", "created", "updated", "deleted"}
}

func (row Template_essential) Values() []interface{} {
	return []interface{}{row.Name, row.Summary, row.Origin}
}

func (row Template) Values() []interface{} {
	return []interface{}{row.Uuid, row.Template_essential.Name, row.Template_essential.Summary, row.Template_essential.Origin, row.Created, row.Updated, row.Deleted}
}

func (row TemplateCommand_essential) Values() []interface{} {
	return []interface{}{row.Name, row.Summary, row.TemplateUuid, row.Sequence, row.Method, row.Args, row.ResultFilter}
}

func (row TemplateCommand) Values() []interface{} {
	return []interface{}{row.Uuid, row.TemplateCommand_essential.Name, row.TemplateCommand_essential.Summary, row.TemplateCommand_essential.TemplateUuid, row.TemplateCommand_essential.Sequence, row.TemplateCommand_essential.Method, row.TemplateCommand_essential.Args, row.TemplateCommand_essential.ResultFilter, row.Created, row.Updated, row.Deleted}
}

type Scanner interface {
	Scan(dest ...interface{}) error
}

func (row *Template_essential) Scan(scanner Scanner) error {
	return scanner.Scan(&row.Name, &row.Summary, &row.Origin)
}

func (row *Template) Scan(scanner Scanner) error {
	return scanner.Scan(&row.Uuid, &row.Template_essential.Name, &row.Template_essential.Summary, &row.Template_essential.Origin, &row.Created, &row.Updated, &row.Deleted)
}

func (row *TemplateCommand_essential) Scan(scanner Scanner) error {
	return scanner.Scan(&row.Name, &row.Summary, &row.TemplateUuid, &row.Sequence, &row.Method, &row.Args, &row.ResultFilter)
}

func (row *TemplateCommand) Scan(scanner Scanner) error {
	return scanner.Scan(&row.Uuid, &row.TemplateCommand_essential.Name, &row.TemplateCommand_essential.Summary, &row.TemplateCommand_essential.TemplateUuid, &row.TemplateCommand_essential.Sequence, &row.TemplateCommand_essential.Method, &row.TemplateCommand_essential.Args, &row.TemplateCommand_essential.ResultFilter, &row.Created, &row.Updated, &row.Deleted)
}

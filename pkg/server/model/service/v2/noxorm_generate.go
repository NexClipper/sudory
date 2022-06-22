package v2

func (Service_essential) ColumnNames() []string {
	return []string{"name", "summary", "cluster_uuid", "template_uuid", "step_count", "subscribed_channel", "on_completion"}
}

func (Service) ColumnNames() []string {
	return []string{"uuid", "created", "name", "summary", "cluster_uuid", "template_uuid", "step_count", "subscribed_channel", "on_completion"}
}

func (ServiceStatus_essential) ColumnNames() []string {
	return []string{"assigned_client_uuid", "step_position", "status", "message"}
}

func (ServiceStatus) ColumnNames() []string {
	return []string{"uuid", "created", "assigned_client_uuid", "step_position", "status", "message"}
}

func (ServiceResults_essential) ColumnNames() []string {
	return []string{"result_type", "result"}
}

func (ServiceResult) ColumnNames() []string {
	return []string{"uuid", "created", "result_type", "result"}
}

func (Service_tangled) ColumnNames() []string {
	return []string{"uuid", "created", "name", "summary", "cluster_uuid", "template_uuid", "step_count", "subscribed_channel", "on_completion", "assigned_client_uuid", "step_position", "status", "message", "result_type", "result", "updated"}
}

func (Service_status) ColumnNames() []string {
	return []string{"uuid", "created", "name", "summary", "cluster_uuid", "template_uuid", "step_count", "subscribed_channel", "on_completion", "assigned_client_uuid", "step_position", "status", "message", "updated"}
}

func (ServiceStep_essential) ColumnNames() []string {
	return []string{"name", "summary", "method", "args", "result_filter"}
}

func (ServiceStep) ColumnNames() []string {
	return []string{"uuid", "sequence", "created", "name", "summary", "method", "args", "result_filter"}
}

func (ServiceStepStatus_essential) ColumnNames() []string {
	return []string{"status", "started", "ended"}
}

func (ServiceStepStatus) ColumnNames() []string {
	return []string{"uuid", "sequence", "created", "status", "started", "ended"}
}

func (ServiceStep_tangled) ColumnNames() []string {
	return []string{"uuid", "sequence", "created", "name", "summary", "method", "args", "result_filter", "status", "started", "ended", "updated"}
}

func (row Service_essential) Values() []interface{} {
	return []interface{}{row.Name, row.Summary, row.ClusterUuid, row.TemplateUuid, row.StepCount, row.SubscribedChannel, row.OnCompletion}
}

func (row Service) Values() []interface{} {
	return []interface{}{row.Uuid, row.Created, row.Service_essential.Name, row.Service_essential.Summary, row.Service_essential.ClusterUuid, row.Service_essential.TemplateUuid, row.Service_essential.StepCount, row.Service_essential.SubscribedChannel, row.Service_essential.OnCompletion}
}

func (row ServiceStatus_essential) Values() []interface{} {
	return []interface{}{row.AssignedClientUuid, row.StepPosition, row.Status, row.Message}
}

func (row ServiceStatus) Values() []interface{} {
	return []interface{}{row.Uuid, row.Created, row.ServiceStatus_essential.AssignedClientUuid, row.ServiceStatus_essential.StepPosition, row.ServiceStatus_essential.Status, row.ServiceStatus_essential.Message}
}

func (row ServiceResults_essential) Values() []interface{} {
	return []interface{}{row.ResultType, row.Result}
}

func (row ServiceResult) Values() []interface{} {
	return []interface{}{row.Uuid, row.Created, row.ServiceResults_essential.ResultType, row.ServiceResults_essential.Result}
}

func (row Service_tangled) Values() []interface{} {
	return []interface{}{row.Service.Uuid, row.Service.Created, row.Service.Service_essential.Name, row.Service.Service_essential.Summary, row.Service.Service_essential.ClusterUuid, row.Service.Service_essential.TemplateUuid, row.Service.Service_essential.StepCount, row.Service.Service_essential.SubscribedChannel, row.Service.Service_essential.OnCompletion, row.ServiceStatus_essential.AssignedClientUuid, row.ServiceStatus_essential.StepPosition, row.ServiceStatus_essential.Status, row.ServiceStatus_essential.Message, row.ServiceResults_essential.ResultType, row.ServiceResults_essential.Result, row.Updated}
}

func (row Service_status) Values() []interface{} {
	return []interface{}{row.Service.Uuid, row.Service.Created, row.Service.Service_essential.Name, row.Service.Service_essential.Summary, row.Service.Service_essential.ClusterUuid, row.Service.Service_essential.TemplateUuid, row.Service.Service_essential.StepCount, row.Service.Service_essential.SubscribedChannel, row.Service.Service_essential.OnCompletion, row.ServiceStatus_essential.AssignedClientUuid, row.ServiceStatus_essential.StepPosition, row.ServiceStatus_essential.Status, row.ServiceStatus_essential.Message, row.Updated}
}

func (row ServiceStep_essential) Values() []interface{} {
	return []interface{}{row.Name, row.Summary, row.Method, row.Args, row.ResultFilter}
}

func (row ServiceStep) Values() []interface{} {
	return []interface{}{row.Uuid, row.Sequence, row.Created, row.ServiceStep_essential.Name, row.ServiceStep_essential.Summary, row.ServiceStep_essential.Method, row.ServiceStep_essential.Args, row.ServiceStep_essential.ResultFilter}
}

func (row ServiceStepStatus_essential) Values() []interface{} {
	return []interface{}{row.Status, row.Started, row.Ended}
}

func (row ServiceStepStatus) Values() []interface{} {
	return []interface{}{row.Uuid, row.Sequence, row.Created, row.ServiceStepStatus_essential.Status, row.ServiceStepStatus_essential.Started, row.ServiceStepStatus_essential.Ended}
}

func (row ServiceStep_tangled) Values() []interface{} {
	return []interface{}{row.ServiceStep.Uuid, row.ServiceStep.Sequence, row.ServiceStep.Created, row.ServiceStep.ServiceStep_essential.Name, row.ServiceStep.ServiceStep_essential.Summary, row.ServiceStep.ServiceStep_essential.Method, row.ServiceStep.ServiceStep_essential.Args, row.ServiceStep.ServiceStep_essential.ResultFilter, row.ServiceStepStatus_essential.Status, row.ServiceStepStatus_essential.Started, row.ServiceStepStatus_essential.Ended, row.Updated}
}

type Scanner interface {
	Scan(dest ...interface{}) error
}

func (row *Service_essential) Scan(scanner Scanner) error {
	return scanner.Scan(&row.Name, &row.Summary, &row.ClusterUuid, &row.TemplateUuid, &row.StepCount, &row.SubscribedChannel, &row.OnCompletion)
}

func (row *Service) Scan(scanner Scanner) error {
	return scanner.Scan(&row.Uuid, &row.Created, &row.Service_essential.Name, &row.Service_essential.Summary, &row.Service_essential.ClusterUuid, &row.Service_essential.TemplateUuid, &row.Service_essential.StepCount, &row.Service_essential.SubscribedChannel, &row.Service_essential.OnCompletion)
}

func (row *ServiceStatus_essential) Scan(scanner Scanner) error {
	return scanner.Scan(&row.AssignedClientUuid, &row.StepPosition, &row.Status, &row.Message)
}

func (row *ServiceStatus) Scan(scanner Scanner) error {
	return scanner.Scan(&row.Uuid, &row.Created, &row.ServiceStatus_essential.AssignedClientUuid, &row.ServiceStatus_essential.StepPosition, &row.ServiceStatus_essential.Status, &row.ServiceStatus_essential.Message)
}

func (row *ServiceResults_essential) Scan(scanner Scanner) error {
	return scanner.Scan(&row.ResultType, &row.Result)
}

func (row *ServiceResult) Scan(scanner Scanner) error {
	return scanner.Scan(&row.Uuid, &row.Created, &row.ServiceResults_essential.ResultType, &row.ServiceResults_essential.Result)
}

func (row *Service_tangled) Scan(scanner Scanner) error {
	return scanner.Scan(&row.Service.Uuid, &row.Service.Created, &row.Service.Service_essential.Name, &row.Service.Service_essential.Summary, &row.Service.Service_essential.ClusterUuid, &row.Service.Service_essential.TemplateUuid, &row.Service.Service_essential.StepCount, &row.Service.Service_essential.SubscribedChannel, &row.Service.Service_essential.OnCompletion, &row.ServiceStatus_essential.AssignedClientUuid, &row.ServiceStatus_essential.StepPosition, &row.ServiceStatus_essential.Status, &row.ServiceStatus_essential.Message, &row.ServiceResults_essential.ResultType, &row.ServiceResults_essential.Result, &row.Updated)
}

func (row *Service_status) Scan(scanner Scanner) error {
	return scanner.Scan(&row.Service.Uuid, &row.Service.Created, &row.Service.Service_essential.Name, &row.Service.Service_essential.Summary, &row.Service.Service_essential.ClusterUuid, &row.Service.Service_essential.TemplateUuid, &row.Service.Service_essential.StepCount, &row.Service.Service_essential.SubscribedChannel, &row.Service.Service_essential.OnCompletion, &row.ServiceStatus_essential.AssignedClientUuid, &row.ServiceStatus_essential.StepPosition, &row.ServiceStatus_essential.Status, &row.ServiceStatus_essential.Message, &row.Updated)
}

func (row *ServiceStep_essential) Scan(scanner Scanner) error {
	return scanner.Scan(&row.Name, &row.Summary, &row.Method, &row.Args, &row.ResultFilter)
}

func (row *ServiceStep) Scan(scanner Scanner) error {
	return scanner.Scan(&row.Uuid, &row.Sequence, &row.Created, &row.ServiceStep_essential.Name, &row.ServiceStep_essential.Summary, &row.ServiceStep_essential.Method, &row.ServiceStep_essential.Args, &row.ServiceStep_essential.ResultFilter)
}

func (row *ServiceStepStatus_essential) Scan(scanner Scanner) error {
	return scanner.Scan(&row.Status, &row.Started, &row.Ended)
}

func (row *ServiceStepStatus) Scan(scanner Scanner) error {
	return scanner.Scan(&row.Uuid, &row.Sequence, &row.Created, &row.ServiceStepStatus_essential.Status, &row.ServiceStepStatus_essential.Started, &row.ServiceStepStatus_essential.Ended)
}

func (row *ServiceStep_tangled) Scan(scanner Scanner) error {
	return scanner.Scan(&row.ServiceStep.Uuid, &row.ServiceStep.Sequence, &row.ServiceStep.Created, &row.ServiceStep.ServiceStep_essential.Name, &row.ServiceStep.ServiceStep_essential.Summary, &row.ServiceStep.ServiceStep_essential.Method, &row.ServiceStep.ServiceStep_essential.Args, &row.ServiceStep.ServiceStep_essential.ResultFilter, &row.ServiceStepStatus_essential.Status, &row.ServiceStepStatus_essential.Started, &row.ServiceStepStatus_essential.Ended, &row.Updated)
}

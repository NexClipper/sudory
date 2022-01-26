package control

// http 요청에서 받는 parameter 이름 정의
const (
	__BODY__          = "__body__" //request body 를 넘겨야 할 떄 사용하는 이름 (이벤트에서 json object KEY로 사용)
	__UUID__          = "uuid"
	__CLUSTER_UUID__  = "cluster_uuid"
	__CLIENT_UUID__   = "client_uuid"
	__SERVICE_UUID__  = "service_uuid"
	__TEMPLATE_UUID__ = "template_uuid"
)

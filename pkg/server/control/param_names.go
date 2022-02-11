package control

// http 요청에서 받는 parameter 이름 정의
const (
	//HTTP HEAD
	__HTTP_HEADER_X_SUDORY_CLIENT_TOKEN__ = "x-sudory-client-token"

	//AUTH jwt-bearer
	__ASSERTION__  = "assertion"
	__GRANT_TYPE__ = "grant_type"
	// // @Param grant_type   formData string true "grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer"
	// __GRANT_TYPE_VALUE__ = "urn:ietf:params:oauth:grant-type:jwt-bearer"

	__BODY__    = "__body__" //request body 를 넘겨야 할 떄 사용하는 이름 (이벤트에서 json object KEY로 사용)
	__UUID__    = "uuid"
	__NAME__    = "name"
	__SUMMARY__ = "summary"

	__CLUSTER_UUID__  = "cluster_uuid"
	__CLIENT_UUID__   = "client_uuid"
	__TOKEN__         = "token"
	__VALUE__         = "value"
	__SERVICE_UUID__  = "service_uuid"
	__TEMPLATE_UUID__ = "template_uuid"
)

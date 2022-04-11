package v1

//Auth Property
type AuthProperty struct {
	ClusterUuid   string `json:"cluster_uuid,omitempty"`   //cluster uuid
	ClientUuid    string `json:"client_uuid,omitempty"`    //client uuid
	Assertion     string `json:"assertion,omitempty"`      //<bearer-token>
	ClientVersion string `json:"client_version,omitempty"` //client version
	// GrantType   string `json:"grant_type,omitempty" default:"urn:ietf:params:oauth:grant-type:jwt-bearer"` //grant_type
}

//HttpReqAuth
type HttpReqAuth struct {
	AuthProperty `json:",inline"` //inline property
}

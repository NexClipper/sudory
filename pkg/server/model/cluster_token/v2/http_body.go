package v2

type HttpReq_ClusterToken_create struct {
	Uuid        string  `json:"uuid,omitempty"` //optional
	Name        string  `json:"name,omitempty"`
	Summary     *string `json:"summary,omitempty"`
	ClusterUuid string  `json:"cluster_uuid,omitempty"`
}

type HttpReq_ClusterToken_update struct {
	Name    string  `json:"name,omitempty"`
	Summary *string `json:"summary,omitempty"`
}

type HttpRsp_ClusterToken struct {
	ClusterToken `json:",inline"`
}

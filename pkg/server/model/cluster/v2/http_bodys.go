package v2

type HttpRsp_Cluster struct {
	Cluster `json:",inline"`
}

type HttpReq_Cluster_create struct {
	Cluster_essential `json:",inline"`

	Uuid string `json:"uuid,omitempty"` // uuid
}

type HttpReq_Cluster_update struct {
	Cluster_essential `json:",inline"`
}

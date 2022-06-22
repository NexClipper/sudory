package v2

type HttpRsp_Cluster struct {
	Cluster `json:",inline"`
}

type HttpReq_ClusterCreate struct {
	Cluster_essential `json:",inline"`
}

type HttpReq_ClusterUpdate struct {
	Cluster_essential `json:",inline"`
}

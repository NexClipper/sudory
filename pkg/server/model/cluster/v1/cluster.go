package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

//Cluster Property
type ClusterProperty struct {
	PollingOption map[string]interface{} `json:"polling_option,omitempty" xorm:"text null 'polling_option' comment('polling option')"`
}

//Cluster
type Cluster struct {
	metav1.DbMeta    `json:",inline" xorm:"extends"` //inline dbmeta
	metav1.UuidMeta  `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	ClusterProperty  `json:",inline" xorm:"extends"` //inline property
}

func (Cluster) TableName() string {
	return "cluster"
}

//HTTP REQUEST BODY: Create Cluster
type HttpReqCluster_Create struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	ClusterProperty  `json:",inline" xorm:"extends"` //inline property
}

//HTTP REQUEST BODY: Update Cluster
type HttpReqCluster_Update struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	ClusterProperty  `json:",inline" xorm:"extends"` //inline property
}

package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

type ChannelProperty struct {
	ClusterUuid string `json:"cluster_uuid,omitempty" xorm:"'cluster_uuid' char(32)     notnull index"` //
	// Pattern     string `json:"pattern,omitempty"      xorm:"'pattern'      varchar(255) notnull index"` //
}

//DATABASE SCHEMA: EVENT
type Channel struct {
	metav1.DbMeta    `json:",inline" xorm:"extends"` //inline dbmeta
	metav1.UuidMeta  `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	ChannelProperty  `json:",inline" xorm:"extends"` //inline property
}

func (Channel) TableName() string {
	return "channel"
}

type ChannelWithEdges struct {
	Channel       `json:",inline" xorm:"extends"`
	NotifierEdges []NotifierEdge `json:"notifier_edges,omitempty"`
}

type NotifierEdge struct {
	NotifierType string `json:"notifier_type,omitempty"` //
	NotifierUuid string `json:"notifier_uuid,omitempty"` //
}

type Channel_create struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	ChannelProperty  `json:",inline" xorm:"extends"` //inline property
	NotifierEdges    []NotifierEdge                  `json:"notifier_edges,omitempty"`
}

type Channel_update struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	ChannelProperty  `json:",inline" xorm:"extends"` //inline property
}

// type NotifierEdges struct {
// 	NotifierEdges []NotifierEdge `json:"notifier_edges,omitempty"`
// }

// type EventNotifier map[string]interface{}

// type EventNotifiers []EventNotifier

package v1

import (
	metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"
)

type EventProperty struct {
	ClusterUuid string `json:"cluster_uuid,omitempty" xorm:"'cluster_uuid' char(32)     notnull index"` //
	Pattern     string `json:"pattern,omitempty"      xorm:"'pattern'      varchar(255) notnull index"` //
}

//DATABASE SCHEMA: EVENT
type Event struct {
	metav1.DbMeta    `json:",inline" xorm:"extends"` //inline dbmeta
	metav1.UuidMeta  `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	EventProperty    `json:",inline" xorm:"extends"` //inline property
}

func (Event) TableName() string {
	return "event"
}

type EventWithEdges struct {
	Event         `json:",inline" xorm:"extends"`
	NotifierEdges []NotifierEdge `json:"notifier_edges,omitempty"`
}

type NotifierEdge struct {
	NotifierType string `json:"notifier_type,omitempty"` //
	NotifierUuid string `json:"notifier_uuid,omitempty"` //
}

type Event_create struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	EventProperty    `json:",inline" xorm:"extends"` //inline property
	NotifierEdges    []NotifierEdge                  `json:"notifier_edges,omitempty"`
}

type Event_update struct {
	metav1.LabelMeta `json:",inline" xorm:"extends"` //inline labelmeta
	EventProperty    `json:",inline" xorm:"extends"` //inline property
}

type NotifierEdges struct {
	NotifierEdges []NotifierEdge `json:"notifier_edges,omitempty"`
}

type EventNotifier map[string]interface{}

type EventNotifiers []EventNotifier

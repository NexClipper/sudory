package v1

import metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"

//go:generate go run github.com/abice/go-enum --file=event_notifier_edge.go --names --nocase=true

/* ENUM(
console
webhook
rabbitmq
)
*/
type EventNotifierType int

type EventNotifierEdgeProperty struct {
	NotifierType string `json:"notifier_type,omitempty" xorm:"'notifier_type' varchar(255) notnull index"` //
	NotifierUuid string `json:"notifier_uuid,omitempty" xorm:"'notifier_uuid' char(32)     notnull index"` //
	EventUuid    string `json:"event_uuid,omitempty"    xorm:"'event_uuid'    char(32)     notnull index"` //
}

type EventNotifierEdge struct {
	metav1.DbMeta             `json:",inline" xorm:"extends"`
	EventNotifierEdgeProperty `json:",inline" xorm:"extends"` //inline property
}

func (EventNotifierEdge) TableName() string {
	return "event_notifier_edge"
}

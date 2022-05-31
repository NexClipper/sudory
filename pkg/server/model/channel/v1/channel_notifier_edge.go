package v1

import metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"

//go:generate go run github.com/abice/go-enum --file=channel_notifier_edge.go --names --nocase=true

/* ENUM(
console
webhook
rabbitmq
)
*/
type NotifierType int

type ChannelNotifierEdgeProperty struct {
	ChannelUuid  string `json:"channel_uuid,omitempty"  xorm:"'channel_uuid'  char(32)     notnull index"` //
	NotifierType string `json:"notifier_type,omitempty" xorm:"'notifier_type' varchar(255) notnull index"` //
	NotifierUuid string `json:"notifier_uuid,omitempty" xorm:"'notifier_uuid' char(32)     notnull index"` //
}

type ChannelNotifierEdge struct {
	metav1.DbMeta               `json:",inline" xorm:"extends"`
	ChannelNotifierEdgeProperty `json:",inline" xorm:"extends"` //inline property
}

func (ChannelNotifierEdge) TableName() string {
	return "channel_notifier_edge"
}

// type ChannelNotifierEdge_update struct {
// 	ChannelNotifierEdgeProperty `json:",inline"`
// }

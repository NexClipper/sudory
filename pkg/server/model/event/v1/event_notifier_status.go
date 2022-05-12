package v1

import metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"

type EventNotifierStatusProperty struct {
	NotifierType string `json:"notifier_type,omitempty" xorm:"'notifier_type' varchar(255) notnull index"` //
	NotifierUuid string `json:"notifier_uuid,omitempty" xorm:"'notifier_uuid' char(32)     notnull index"` //
	Error        string `json:"error,omitempty"         xorm:"'error'         TEXT         null"`          //
}

type EventNotifierStatus struct {
	metav1.DbMeta               `json:",inline" xorm:"extends"`
	metav1.UuidMeta             `json:",inline" xorm:"extends"` //inline uuidmeta
	EventNotifierStatusProperty `json:",inline" xorm:"extends"` //inline property
}

func (EventNotifierStatus) TableName() string {
	return "event_notifier_status"
}

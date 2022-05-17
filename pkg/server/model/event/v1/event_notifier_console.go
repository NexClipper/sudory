package v1

import metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"

type EventNotifierConsoleProperty struct{}

func (EventNotifierConsoleProperty) Type() EventNotifierType {
	return EventNotifierTypeConsole
}

type EventNotifierConsole struct {
	metav1.DbMeta                `json:",inline" xorm:"extends"` //inline dbmeta
	metav1.UuidMeta              `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta             `json:",inline" xorm:"extends"` //inline labelmeta
	EventNotifierConsoleProperty `json:",inline" xorm:"extends"` //inline property
	MIME                         `json:",inline" xorm:"extends"` //inline MIME
}

func (EventNotifierConsole) TableName() string {
	return "event_notifier_console"
}

type EventNotifierConsole_create struct {
	metav1.LabelMeta             `json:",inline" xorm:"extends"` //inline labelmeta
	EventNotifierConsoleProperty `json:",inline" xorm:"extends"` //inline property
	MIME                         `json:",inline" xorm:"extends"` //inline MIME
}

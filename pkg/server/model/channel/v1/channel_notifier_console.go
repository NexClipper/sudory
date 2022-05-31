package v1

import metav1 "github.com/NexClipper/sudory/pkg/server/model/meta/v1"

type NotifierConsoleProperty struct{}

func (NotifierConsoleProperty) Type() NotifierType {
	return NotifierTypeConsole
}

type NotifierConsole struct {
	metav1.DbMeta           `json:",inline" xorm:"extends"` //inline dbmeta
	metav1.UuidMeta         `json:",inline" xorm:"extends"` //inline uuidmeta
	metav1.LabelMeta        `json:",inline" xorm:"extends"` //inline labelmeta
	NotifierConsoleProperty `json:",inline" xorm:"extends"` //inline property
	MIME                    `json:",inline" xorm:"extends"` //inline MIME
}

func (NotifierConsole) TableName() string {
	return "channel_notifier_console"
}

type NotifierConsole_create struct {
	metav1.LabelMeta        `json:",inline" xorm:"extends"` //inline labelmeta
	NotifierConsoleProperty `json:",inline" xorm:"extends"` //inline property
	MIME                    `json:",inline" xorm:"extends"` //inline MIME
}

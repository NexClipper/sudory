package v3

import (
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type ManagedChannel_create struct {
	Uuid    string             `json:"uuid,omitempty"` // optional
	Name    string             `json:"name,omitempty"`
	Summary vanilla.NullString `json:"summary,omitempty"        swaggertype:"string"`
	// enums:"NaV(0), nonspecified(1), client-auth(2), client-polling-out(3), client-polling-in(4)"
	EventCategory EventCategory `json:"event_category,omitempty" enums:"0,1,2,3,4"`
}

type ManagedChannel_update = ManagedChannel_property

type ManagedChannel_property struct {
	Name    string             `column:"name"                       json:"name,omitempty"`
	Summary vanilla.NullString `column:"summary"                    json:"summary,omitempty"        swaggertype:"string"`
	// enums:"NaV(0), nonspecified(1), client-auth(2), client-polling-out(3), client-polling-in(4)"
	EventCategory EventCategory `column:"event_category,default(0)" json:"event_category,omitempty" enums:"0,1,2,3,4"`
}

type ManagedChannel struct {
	Uuid string `column:"uuid"                          json:"uuid,omitempty"` // pk

	ManagedChannel_property `json:",inline"`

	Created time.Time        `column:"created"                   json:"created,omitempty"`
	Updated vanilla.NullTime `column:"updated"                   json:"updated,omitempty"        swaggertype:"string"`
	Deleted vanilla.NullTime `column:"deleted"                   json:"deleted,omitempty"        swaggertype:"string"`
}

func (ManagedChannel) TableName() string {
	return "managed_channel"
}

// type ManagedChannel_option struct {
// 	ManagedChannel `alias:"CN," json:",inline"`

// 	StatusOption ChannelStatusOption_property `alias:"SO,status_option" json:"status_option,omitempty"`
// 	Format       Format_property              `alias:"FM,format"        json:"format,omitempty"`
// 	Edge         NotifierEdge_property        `alias:"EG,edge"          json:"edge,omitempty"`
// }

// func (record ManagedChannel_option) TableName() string {
// 	q := `(
// 		SELECT %v /**columns**/
// 		  FROM %v CN /**managed_channel**/
// 		  LEFT JOIN %v SO /**managed_channel_status_option**/
// 				 ON CN.uuid = SO.uuid
// 		  LEFT JOIN %v FM /**managed_channel_format**/
// 				 ON CN.uuid = FM.uuid
// 		  LEFT JOIN %v EG /**managed_channel_notifier_edge**/
// 				 ON CN.uuid = EG.uuid
// 		) X`

// 	CN := record.ManagedChannel.TableName()
// 	SO := record.StatusOption.TableName()
// 	FM := record.Format.TableName()
// 	EG := record.Edge.TableName()
// 	return fmt.Sprintf(q, strings.Join(record.ColumnNamesWithAlias(), ", "), CN, SO, FM, EG)
// }

// func (record ManagedChannel_tangled) TableName() string {
// 	q := `(
// 		SELECT %v /**columns**/
// 		  FROM %v CN /**managed_channel**/
// 		  LEFT JOIN %v SO /**managed_channel_status_option**/
// 				 ON CN.uuid = SO.uuid
// 		  LEFT JOIN %v FM /**managed_channel_format**/
// 				 ON CN.uuid = FM.uuid
// 		  LEFT JOIN %v EG /**managed_channel_notifier_edge**/
// 				 ON CN.uuid = EG.uuid
// 		  LEFT JOIN %v CS /**managed_channel_notifier_console**/
// 			     ON CN.uuid = CS.uuid
// 		  LEFT JOIN %v WH /**managed_channel_notifier_webhook**/
// 			     ON CN.uuid = WH.uuid
// 		  LEFT JOIN %v RQ /**managed_channel_notifier_rabbitmq**/
// 				 ON CN.uuid = RQ.uuid
// 		  LEFT JOIN %v SH /**managed_channel_notifier_slackhook**/
// 				 ON CN.uuid = SH.uuid
// 		) X`

// 	CN := record.ManagedChannel.TableName()
// 	SO := record.StatusOption.TableName()
// 	FM := record.Format.TableName()
// 	EG := record.Edge.TableName()
// 	CS := record.Notifiers.Console.TableName()
// 	WH := record.Notifiers.Webhook.TableName()
// 	RQ := record.Notifiers.RabbitMq.TableName()
// 	SH := record.Notifiers.Slackhook.TableName()
// 	return fmt.Sprintf(q, strings.Join(record.ColumnNamesWithAlias(), ", "), CN, SO, FM, EG, CS, WH, RQ, SH)
// }

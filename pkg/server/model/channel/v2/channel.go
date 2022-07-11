package v2

import (
	"fmt"
	"strings"
	"time"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type ManagedChannel_property struct {
	Name    string             `column:"name"           json:"name,omitempty"`
	Summary vanilla.NullString `column:"summary"        json:"summary,omitempty"        swaggertype:"string"`
	// enums:"NaV(0), nonspecified(1), client-auth(2), client-polling-out(3), client-polling-in(4)"
	EventCategory EventCategory `column:"event_category,default(0)" json:"event_category,omitempty" enums:"0,1,2,3,4"`
}

func (ManagedChannel_property) TableName() string {
	return "managed_channel"
}

type ManagedChannel struct {
	ManagedChannel_property `json:",inline"`

	Created time.Time        `column:"created" json:"created,omitempty"`
	Updated vanilla.NullTime `column:"updated" json:"updated,omitempty" swaggertype:"string"`
	Deleted vanilla.NullTime `column:"deleted" json:"deleted,omitempty" swaggertype:"string"`
	Uuid    string           `column:"uuid"    json:"uuid,omitempty"` // pk
}

type ManagedChannel_option struct {
	ManagedChannel `alias:"CN," json:",inline"`

	StatusOption ChannelStatusOption_property `alias:"SO,status_option" json:"status_option,omitempty"`
	Format       Format_property              `alias:"FM,format"        json:"format,omitempty"`
	Edge         NotifierEdge_property        `alias:"EG,edge"          json:"edge,omitempty"`
}

func (record ManagedChannel_option) TableName() string {
	q := `(
		SELECT %v /**columns**/
		  FROM %v CN /**managed_channel CN**/
		  LEFT JOIN %v SO /**managed_channel_status_option SO**/
				 ON CN.uuid = SO.uuid 		  
		  LEFT JOIN %v FM /**managed_channel_format FM**/
				 ON CN.uuid = FM.uuid 
		  LEFT JOIN %v EG /**managed_channel_notifier_edge EG**/
				 ON CN.uuid = EG.uuid 
		) X`

	CN := record.ManagedChannel.TableName()
	SO := record.StatusOption.TableName()
	FM := record.Format.TableName()
	EG := record.Edge.TableName()
	return fmt.Sprintf(q, strings.Join(record.ColumnNamesWithAlias(), ", "), CN, SO, FM, EG)
}

type ManagedChannel_tangled struct {
	ManagedChannel `alias:"CN," json:",inline"`

	StatusOption ChannelStatusOption_property `alias:"SO,status_option" json:"status_option,omitempty"`
	Format       Format_property              `alias:"FM,format"        json:"format,omitempty"`

	Notifier struct {
		NotifierEdge_property `alias:"EG,edge" json:",inline"`

		Console  NotifierConsole_property  `alias:"CS,console"  json:"console,omitempty"`
		Webhook  NotifierWebhook_property  `alias:"WH,webhook"  json:"webhook,omitempty"`
		RabbitMq NotifierRabbitMq_property `alias:"RQ,rabbitmq" json:"rabbitmq,omitempty"`
	} `json:"notifier,omitempty"`
}

func (record ManagedChannel_tangled) TableName() string {
	q := `(
		SELECT %v /**columns**/
		  FROM %v CN /**managed_channel CN**/
		  LEFT JOIN %v SO /**managed_channel_status_option SO**/
				 ON CN.uuid = SO.uuid 		  
		  LEFT JOIN %v FM /**managed_channel_format FM**/
				 ON CN.uuid = FM.uuid 
		  LEFT JOIN %v EG /**managed_channel_notifier_edge EG**/
				 ON CN.uuid = EG.uuid 
		  LEFT JOIN %v CS /**managed_channel_notifier_console B**/
			     ON CN.uuid = CS.uuid 
		  LEFT JOIN %v WH /**managed_channel_notifier_webhook C**/
			     ON CN.uuid = WH.uuid 		  
		  LEFT JOIN %v RQ /**managed_channel_notifier_rabbitmq D**/
				 ON CN.uuid = RQ.uuid 
		) X`

	CN := record.ManagedChannel.TableName()
	SO := record.StatusOption.TableName()
	FM := record.Format.TableName()
	EG := record.Notifier.NotifierEdge_property.TableName()
	CS := record.Notifier.Console.TableName()
	WH := record.Notifier.Webhook.TableName()
	RQ := record.Notifier.RabbitMq.TableName()
	return fmt.Sprintf(q, strings.Join(record.ColumnNamesWithAlias(), ", "), CN, SO, FM, EG, CS, WH, RQ)
}

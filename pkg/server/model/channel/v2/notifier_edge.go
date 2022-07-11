package v2

import (
	"fmt"
	"strings"

	"github.com/NexClipper/sudory/pkg/server/database/vanilla"
)

type NotifierEdge_essential struct {
	// enums:"NaV(0), console(1), webhook(2), rabbitmq(3)"
	NotifierType NotifierType `column:"notifier_type,default(0)" json:"notifier_type,omitempty" swaggertype:"string" enums:"0,1,2,3"`
}

type NotifierEdge_property struct {
	NotifierEdge_essential `json:",inline"`

	Created vanilla.NullTime `column:"created" json:"created,omitempty" swaggertype:"string"`
	Updated vanilla.NullTime `column:"updated" json:"updated,omitempty" swaggertype:"string"`
}

func (NotifierEdge_property) TableName() string {
	return "managed_channel_notifier_edge"
}

type NotifierEdge struct {
	NotifierEdge_property `json:",inline"`

	Uuid string `column:"uuid" json:"uuid,omitempty"` // pk
}

type NotifierEdge_option struct {
	NotifierEdge `alias:"EG," json:",inline"`

	Console  NotifierConsole_property  `alias:"CS,console"  json:"console,omitempty"`
	Webhook  NotifierWebhook_property  `alias:"WH,webhook"  json:"webhook,omitempty"`
	RabbitMq NotifierRabbitMq_property `alias:"RQ,rabbitmq" json:"rabbitmq,omitempty"`
}

func (record NotifierEdge_option) TableName() string {
	q := `(
		SELECT %v /**columns**/
	      FROM %v EG /**managed_channel_notifier_edge EG**/
		  LEFT JOIN %v CS /**managed_channel_notifier_console B**/
			     ON EG.uuid = CS.uuid 
		  LEFT JOIN %v WH /**managed_channel_notifier_webhook C**/
			     ON EG.uuid = WH.uuid 		  
		  LEFT JOIN %v RQ /**managed_channel_notifier_rabbitmq D**/
				 ON EG.uuid = RQ.uuid 
	) X
`

	EG := record.NotifierEdge_property.TableName()
	CS := record.Console.TableName()
	WH := record.Webhook.TableName()
	RQ := record.RabbitMq.TableName()
	return fmt.Sprintf(q, strings.Join(record.ColumnNamesWithAlias(), ", "), EG, CS, WH, RQ)
}
